package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
)

type JSONRPCDemoServer struct {
	subprocess *exec.Cmd
	stdin      io.WriteCloser
	stdout     io.ReadCloser
	client     *jrpc2.Client
}

func NewJSONRPCDemoServer(name string, args ...string) *JSONRPCDemoServer {
	subprocess := exec.Command(name, args...)
	subprocess.Stderr = os.Stderr
	stdin, err := subprocess.StdinPipe()
	if err != nil {
		log.Fatalln(err)
	}
	stdout, err := subprocess.StdoutPipe()
	if err != nil {
		log.Fatalln(err)
	}
	if err = subprocess.Start(); err != nil {
		log.Fatalln(err)
	}
	client := jrpc2.NewClient(channel.Line(stdout, stdin), nil)
	time.Sleep(time.Second)
	return &JSONRPCDemoServer{
		subprocess: subprocess,
		stdin:      stdin,
		stdout:     stdout,
		client:     client,
	}
}

func (s *JSONRPCDemoServer) Close() (err error) {
	if err = s.client.Close(); err != nil {
		return
	}
	if err = s.stdout.Close(); err != nil {
		return
	}
	if err = s.stdin.Close(); err != nil {
		return
	}
	return
}

func (s *JSONRPCDemoServer) Hello(ctx context.Context, name string) (string, error) {
	var result string
	if err := s.client.CallResult(ctx, "pythonic_hello", []string{name}, &result); err != nil {
		return "", err
	}
	return result, nil
}

func main() {
	server := NewJSONRPCDemoServer("./server")
	defer server.Close()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Go> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln(err)
		}
		if strings.TrimSpace(line) == "" {
			break
		}

		resp, err := server.Hello(context.Background(), line)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(resp)
	}
}
