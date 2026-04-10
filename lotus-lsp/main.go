package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

const lsName = "lotus-ls"

func main() {
	fmt.Println("LSP Running")

	stream := jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{})
	h := NewLotusHandler()
	conn := jsonrpc2.NewConn(context.Background(), stream, jsonrpc2.HandlerWithError(h.handle))
	<-conn.DisconnectNotify()
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}
	return os.Stdout.Close()
}
