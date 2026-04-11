package main

import (
	"context"
	"fmt"
	"os"

	"github.com/sourcegraph/jsonrpc2"
)

const lsName = "lotus-ls"

func main() {
	f, _ := os.OpenFile("/mnt/c/Users/Public/lotus-lsp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	fmt.Fprintln(f, "LSP started")
	fmt.Fprintln(os.Stderr, "LSP Running")

	stream := jsonrpc2.NewBufferedStream(stdrwc{}, jsonrpc2.VSCodeObjectCodec{})
	h := NewLotusHandler()
	conn := jsonrpc2.NewConn(context.Background(), stream, jsonrpc2.HandlerWithError(h.handle))
	<-conn.DisconnectNotify()

	fmt.Fprintln(f, "LSP stopped")
	f.Close()
}

type stdrwc struct{}

func (stdrwc) Read(p []byte) (int, error) {
	return os.Stdin.Read(p)
}

func (stdrwc) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (stdrwc) Close() error {
	return nil
}
