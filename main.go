package main

import "fmt"
import "net"
import "strings"
import "bufio"

func main() {
	// The Listen function creates servers.
	/*e.g.
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
	*/
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Accept waits for and returns the next connection to the listener.
	// conn is a generic network connection.
	// When l.Accept() is called, it blocks until a connection is made. Once a connection is made, it returns a net.Conn object.
	// The net.Conn object is a generic network connection. It has methods for reading and writing to the connection.
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening on port :6379")
	// The defer statement defers the execution of a function until the surrounding function returns.
	// The defer statement is often used to close a file, so that the file is closed as soon as the function returns.
	defer conn.Close()

	// We have a connection, now we need to read and write to it!\
	// This loop will keep reading from the connection and writing to it.
	for {
		// conn is a generic network connection which has methods for reading and writing to the connection.
		// So NewResp(conn) is creating a new Resp object which wraps the connection.
		resp := NewResp(conn)
		// resp.read() is reading the command from the connection.
		value, err := resp.read()
		if err != nil {
			fmt.Println(err)
			return
		}

		// If the command is not an array, then it is invalid.
		if value.typ != "array" {
			fmt.Println("Invalid command")
			return
		}

		// If the command is an empty array, then it is invalid.
		if len(value.array) == 0 {
			fmt.Println("Invalid command")
			return
		}

		// The first element of the array is the command. The rest of the elements are the arguments.
		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		// NewWriter wraps the connection and has a method Write which writes to the connection.
		writer := NewWriter(conn)

		// Handler is a map of commands to functions. We get the coresponding function for the command.
		// If the command is not in the map, then it is invalid. ok is false in that case.
		handler, ok := Handler[cmd]
		if !ok {
			fmt.Println("Invalid command", cmd)
			writer.Write(Value{typ: "string", str: "ERR unknown command '" + cmd + "'"})
			continue
		}

		// We call the function with the arguments and get the result.
		result := handler(args)
		writer.Write(result)
	}
}