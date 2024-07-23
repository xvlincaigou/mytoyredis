package main

import (
	"sync"
)

// Declare a map of strings to strings. And it is empty now.
var SETs = map[string]string{}
// Declare a mutex. It is used to lock the SETs map and ensure that only one goroutine can access it at a time.
var SETsMu = sync.RWMutex{}

// To set the value of a key.
func set(args []Value) Value {
	// If the number of arguments is not 2, then return an error.
	if len(args) != 2 {
		return Value{typ: "string", str: "ERR wrong number of arguments for 'set' command"}
	}

	// The first argument is the key and the second argument is the value.
	key := args[0].bulk
	value := args[1].bulk

	// The code between SETsMu.Lock() and SETsMu.Unlock() is a critical section. It is protected by the mutex.
	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	// Return OK if the value is set successfully.
	return Value{typ: "string", str: "OK"}
}
 
// To get the value of a key.
func get(args []Value) Value {
	// If the number of arguments is not 1, then return an error.
	if len(args) != 1 {
		return Value{typ: "string", str: "ERR wrong number of arguments for 'get' command"}
	}

	// The first argument is the key.
	key := args[0].bulk

	// The code between SETsMu.RLock() and SETsMu.RUnlock() is a critical section. It is protected by the mutex.
	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	// We need to handle this after releasing the lock, or the program will deadlock.
	if !ok {
		return Value{typ: "string", str: "nil"}
	}
	return Value{typ: "string", str: value}
}

func ping (args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

// Declare a map of strings to functions. The functions take a slice of Value and return a Value. And the {} is the body of the map.
var Handler = map[string]func([] Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}