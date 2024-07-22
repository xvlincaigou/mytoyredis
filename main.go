package main

import "fmt"
import "net"

func main() {
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening on port :6379")
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		conn.Write([]byte("+OK\r\n"))
	}
}