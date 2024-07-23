package main

import (
	"bufio"
	"os"
	"sync"
	"time"
	"io"
)

type Aof struct {
	file *os.File
	rd *bufio.Reader
	mu sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	// Open the file in append mode. If the file does not exist, create it.
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	// Define a new Aof object.
	aof := &Aof {
		file: f,
		rd: bufio.NewReader(f),
	}

	// This is a goroutine that synchronizes the file every second.
	// How graceful!
	// go is a keyword that creates a new goroutine.
	// func is a keyword that declares a function.
	// () means call the function immediately.
	go func () {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

// Close the file.
func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

// Write the file.
func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	return err
} 

func (aof *Aof) Read(fn func(value Value)) error {
    aof.mu.Lock()
    defer aof.mu.Unlock()

    // 将文件指针移到文件开头
    _, err := aof.file.Seek(0, 0)
    if err != nil {
        return err
    }

    // 创建一个新的 Resp 对象来读取和解析数据
    resp := NewResp(aof.rd)

    for {
        // 读取一个 Value
        value, err := resp.read()
        if err != nil {
            // 如果到达文件末尾，就退出循环
            if err == io.EOF {
                break
            }
            return err
        }

        // 调用传入的函数处理这个 Value
        fn(value)
    }

    return nil
}