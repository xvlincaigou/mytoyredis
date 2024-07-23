package main

import (
	"bufio"
	"os"
	"sync"
)

type Aof struct {
	file os.File,
	rd bufio.Reader,
	mu sync.Mutex,
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		return nil, err
	}

	aof := &Aof {
		file: f,
		rd: bufio.NewReader(f),
	}

	go func () {
		for {
			aof.mu.Lock()
			aof.file.sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}

	return aof, nil
}

func (aof *Aof) Close() error {
	aof..mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.mu.lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	return err
} 

