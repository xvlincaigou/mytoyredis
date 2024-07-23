package main

import (
	"bufio"
	"io"
	"strconv"
	"fmt"
)

/* 
Different between bulk string and simple string:
bulk string : $<length>\r\n<data>\r\n
e.g. $6\r\nfoobar\r\n
simple string : +<data>\r\n
e.g. +OK\r\n
*/
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

/*
Package bufio implements buffered I/O. 
It wraps an io.Reader or io.Writer object, creating another object (Reader or Writer) that 
also implements the interface but provides buffering and some help for textual I/O.
*/
/*
strings.NewReader(s string) *strings.Reader
func NewReader(s string) *Reader
NewReader returns a new Reader reading from s. It is similar to bytes.NewBufferString but more efficient and non-writable.
*/

type Resp struct {
	reader *bufio.Reader
}

// It is like factory method!
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readline() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		line = append(line, b)
		n += 1
		if len(line) >= 2 && line[n-2] == '\r' {
			return line[:n-2], n, nil
		}
	}
	return line, n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readline()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// Read the command from the connection
// r *Resp is a pointer to the Resp object, it is a receiver. This read() method is binded to the Resp object.
// param: none
// return: Value, error
func (r *Resp) read() (Val Value, err error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {
	// In our toy redis, the command is always an array of strings.
	// So we only need to consider these two cases.
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Val Value, err error){
	Val = Value{}
	Val.typ = "array"

	size, _, err := r.readInteger()
	if err != nil {
		return Val, err
	}

	for size > 0 {
		val, err := r.read()
		if err != nil {
			return Val, err
		}

		Val.array = append(Val.array, val)
		size -= 1
	}

	return Val, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	r.readline()

	return v, nil
}

func (v Value) Marshal() ([]byte) {
	switch v.typ {
		case "array":
			return v.marshalArray()
		case "bulk":
			return v.marshalBulk()
		case "string":
			return v.marshalString()
		case "null":
			return v.marshallNull()
		case "error":
			return v.marshallError()
		default:
			return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}