package hydrogen

/*
A few methods:

conn.GetInt("var_name") -> Looks for any non-specific type of int and converts to golang type `int`
conn.GetString("var_name") -> Looks for a string

conn.Return(data... []byte) -> Will give the client back some data


server := hydrogen.NewServer()
	-> returns server type (Must be initiated before using)
server.Start(true, "localhost:3030")
	-> Will start on localhost:3030 and will log all interactions (Print to terminal)


function type:

func PaymentHandler(conn hydrogen.Connection) {
	sender := conn.GetString("address")
	amount := conn.GetInt("amount")
}

server.AddEndpoint("/SendPayment", PaymentHandler)
*/
