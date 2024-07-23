package main

func ping (args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

var Handler = map[string]func([] Value) Value{
	"Ping", ping
}