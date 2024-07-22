package main

import (
	"bufio"
	"io"
	"strconv"
	"fmt"
)

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

type Resp struct {
	reader *bufio.Reader
}

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

func (r *Resp) read() (Val Value, err error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {
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