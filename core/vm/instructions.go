package vm

type Instruction byte

// Ethereum-like arithmetic instructions
const (
	STOP       Instruction = 0x00
	ADD        Instruction = 0x01
	MUL        Instruction = 0x02
	SUB        Instruction = 0x03
	DIV        Instruction = 0x04
	SDIV       Instruction = 0x05
	MOD        Instruction = 0x06
	SMOD       Instruction = 0x07
	ADDMOD     Instruction = 0x08
	MULMOD     Instruction = 0x09
	EXP        Instruction = 0x0a
	SIGNEXTEND Instruction = 0x0b
)

// Comparison and bitwise logic instructions
const (
	LT      Instruction = 0x10
	GT      Instruction = 0x11
	SLT     Instruction = 0x12
	SGT     Instruction = 0x13
	EQ      Instruction = 0x14
	ISZERO  Instruction = 0x15
	AND     Instruction = 0x16
	OR      Instruction = 0x17
	XOR     Instruction = 0x18
	NOT     Instruction = 0x19
	BYTE    Instruction = 0x1a
	SHL     Instruction = 0x1b
	SHR     Instruction = 0x1c
	SAR     Instruction = 0x1d
	BITTEST Instruction = 0x1e // Test if specific bit is set
	BSWAP   Instruction = 0x1f // Byte swap (endianness conversion)
)

// Cryptographic and hashing instructions
const (
	KECCAK256     Instruction = 0x20
	SHA256        Instruction = 0x21
	RIPEMD160     Instruction = 0x22
	ECDSA_RECOVER Instruction = 0x23 // Recover public key from signature
	HASH160       Instruction = 0x24 // Bitcoin's hash160 (SHA256 followed by RIPEMD160)
	HASH256       Instruction = 0x25 // Double SHA256
)

// XRPL-like ledger and payment operations
const (
	PAYMENT      Instruction = 0x30
	ESCROWCREATE Instruction = 0x31
	ESCROWCANCEL Instruction = 0x32
	ESCROWFINISH Instruction = 0x33
	BLACKLIST    Instruction = 0x34
	OFFERCREATE  Instruction = 0x35
	OFFERCANCEL  Instruction = 0x36
	TRUSTSET     Instruction = 0x37
	ACCOUNTSET   Instruction = 0x38 // Modify account properties
	DELEGATE     Instruction = 0x39 // Delegate authority (like staking or voting)

)

// Ethereum-like storage and execution instructions
const (
	SLOAD        Instruction = 0x50
	SSTORE       Instruction = 0x51
	JUMP         Instruction = 0x52
	JUMPI        Instruction = 0x53
	PC           Instruction = 0x54
	MSIZE        Instruction = 0x55
	GAS          Instruction = 0x56
	JUMPDEST     Instruction = 0x57
	CREATE2      Instruction = 0x58 // Deterministic contract creation
	STATICCALL   Instruction = 0x59 // Call another contract without modifying state
	DELEGATECALL Instruction = 0x5a // Delegate call to another contract
)

// Additional utility instructions
const (
	PUSH1        Instruction = 0x60
	PUSH32       Instruction = 0x7f
	DUP1         Instruction = 0x80
	DUP16        Instruction = 0x8f
	SWAP1        Instruction = 0x90
	SWAP16       Instruction = 0x9f
	LOG0         Instruction = 0xa0
	LOG4         Instruction = 0xa4
	CREATE       Instruction = 0xf0
	CALL         Instruction = 0xf1
	RETURN       Instruction = 0xf2
	REVERT       Instruction = 0xf3
	SELFDESTRUCT Instruction = 0xff
	SERIALIZE    Instruction = 0xa5 // Serialize data structure into bytes
	DESERIALIZE  Instruction = 0xa6 // Deserialize bytes into data structure
	MERKLEPROOF  Instruction = 0xa7 // Verify Merkle proof
	RANGEPROOF   Instruction = 0xa8 // Verify range proof (useful in privacy-preserving protocols)
)

// Advanced blockchain and consensus-related instructions
const (
	VERIFYBLOCK  Instruction = 0xb0 // Verify block header
	VERIFYTX     Instruction = 0xb1 // Verify transaction format and signatures
	GETBLOCKHASH Instruction = 0xb2 // Get hash of a block in the blockchain
	GETTXSTATUS  Instruction = 0xb3 // Get status of a transaction
	CHAINID      Instruction = 0xb4 // Get current chain ID
	TIMESTAMP    Instruction = 0xb5 // Get current block timestamp
	BLOCKNUMBER  Instruction = 0xb6 // Get current block number
	DIFFICULTY   Instruction = 0xb7 // Get current block difficulty
)

// Define additional opcodes as needed to cover more complex operations
// or specific features from both Ethereum and XRPL.
