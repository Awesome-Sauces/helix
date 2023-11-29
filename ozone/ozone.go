package ozone

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"errors"
	"os"
	"reflect"
	"sync"
	"time"
)

func init() {
	gob.Register(DataItem{})
}

// Error variables to represent specific errors
var (
	ErrOzoneClosed = errors.New("database closed")
	ErrKeyNotFound = errors.New("key not found")
)

// DataItem represents the data to be stored with additional metadata
type DataItem struct {
	Key       string      // Unique identifier
	Value     interface{} // Actual data to be stored
	Timestamp string      // ISO timestamp indicating when the data was added
	Size      float64     // Size of the data in KB
}

// item represents the unit that will be stored in the cache
type item struct {
	key   string   // Unique identifier
	value DataItem // Data item containing actual data and metadata
}

// Database represents the main database structure
type Database struct {
	db        map[string]DataItem      // Main database storage
	cache     map[string]*list.Element // LRU cache for quick access to recent data
	lru       *list.List               // Doubly linked list to implement LRU caching
	cacheSize int                      // Max size of the cache
	mutex     sync.RWMutex             // Mutex for handling concurrent access
	closed    bool                     // Flag indicating whether the database is closed
}

// New initializes and returns a new Database
func New(cacheSize int) *Database {
	return &Database{
		db:        make(map[string]DataItem),
		cache:     make(map[string]*list.Element),
		lru:       list.New(),
		cacheSize: cacheSize,
	}
}

// Close safely closes the database
func (db *Database) Close() error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrOzoneClosed
	}

	db.closed = true
	return nil
}

func (db *Database) Contains(key string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if db.closed {
		return false
	}

	if _, found := db.cache[key]; found {
		return true
	}

	_, exists := db.db[key]

	return exists
}

// Get retrieves a data item from the database
func (db *Database) Get(key string) (DataItem, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if db.closed {
		return DataItem{}, ErrOzoneClosed
	}

	if elem, found := db.cache[key]; found {
		db.lru.MoveToFront(elem)
		return elem.Value.(*item).value, nil
	}

	value, exists := db.db[key]
	if !exists {
		return DataItem{}, ErrKeyNotFound
	}

	db.updateCache(key, value)
	return value, nil
}

// Set adds or updates a data item in the database
func (db *Database) Set(key string, value interface{}) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrOzoneClosed
	}

	data := DataItem{
		Key:       key,
		Value:     value,
		Timestamp: time.Now().Format(time.RFC3339),
		Size:      calculateSize(value),
	}

	db.db[key] = data
	db.updateCache(key, data)
	return nil
}

// Delete removes a data item from the database
func (db *Database) Delete(key string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrOzoneClosed
	}

	if _, contains := db.db[key]; !contains {
		return ErrKeyNotFound
	}

	delete(db.db, key)
	if elem, found := db.cache[key]; found {
		db.lru.Remove(elem)
		delete(db.cache, key)
	}
	return nil
}

// updateCache updates the cache when a data item is added or accessed
func (db *Database) updateCache(key string, value DataItem) {
	if elem, found := db.cache[key]; found {
		db.lru.MoveToFront(elem)
		elem.Value.(*item).value = value
		return
	}

	if db.lru.Len() >= db.cacheSize {
		toRemove := db.lru.Back()
		if toRemove != nil {
			db.lru.Remove(toRemove)
			delete(db.cache, toRemove.Value.(*item).key)
		}
	}

	elem := db.lru.PushFront(&item{key: key, value: value})
	db.cache[key] = elem
}

// IterateFunc is a callback function to process key-value pairs during iteration.
type IterateFunc func(key string, value DataItem) error

// Iterate iterates over the database and applies the provided callback function to each key-value pair.
func (db *Database) Iterate(callback IterateFunc) error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if db.closed {
		return ErrOzoneClosed
	}

	// Create a channel to send key-value pairs
	ch := make(chan struct {
		Key   string
		Value DataItem
	})

	// Start a goroutine to send key-value pairs to the channel
	go func() {
		defer close(ch)
		for key, value := range db.db {
			ch <- struct {
				Key   string
				Value DataItem
			}{Key: key, Value: value}
		}
	}()

	// Process key-value pairs using the provided callback
	for pair := range ch {
		err := callback(pair.Key, pair.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// calculateSize estimates the size of the data item in KB
func calculateSize(value interface{}) float64 {
	v := reflect.ValueOf(value)
	sizeInBytes := 0

	// Logic for calculating the size of different types of values
	// and converting size from bytes to KB
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		sizeInBytes = 8 // assuming max 8 bytes for these types
	case reflect.Bool:
		sizeInBytes = 1
	case reflect.String:
		sizeInBytes = len(value.(string))
	case reflect.Slice, reflect.Array:
		sizeInBytes = v.Len() * 8 // rough estimation
	case reflect.Map, reflect.Struct:
		sizeInBytes = v.NumField() * 8 // rough estimation
	}

	// Converting size in bytes to kilobytes (KB)
	sizeInKB := float64(sizeInBytes) / 1024.0
	return sizeInKB
}

// Encodes database to GOB
func (db *Database) Save(filename string) error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if db.closed {
		return ErrOzoneClosed
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(db.db)
	if err != nil {
		return err
	}

	return nil
}

// Loads GOB Encoded Database
func (db *Database) Load(filename string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrOzoneClosed
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&db.db)
	if err != nil {
		return err
	}

	// You might want to also update the cache based on the newly loaded data.
	// Clearing the cache as an example.
	db.cache = make(map[string]*list.Element)
	db.lru = list.New()

	return nil
}

// SaveToString encodes the database to a string and returns it
func (db *Database) SaveString() (string, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if db.closed {
		return "", ErrOzoneClosed
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(db.db)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// LoadFromString decodes the database from a string
func (db *Database) LoadString(data string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if db.closed {
		return ErrOzoneClosed
	}

	buf := bytes.NewBufferString(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&db.db)
	if err != nil {
		return err
	}

	// You might want to also update the cache based on the newly loaded data.
	// Clearing the cache as an example.
	db.cache = make(map[string]*list.Element)
	db.lru = list.New()

	return nil
}
