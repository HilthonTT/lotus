package main

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/sourcegraph/jsonrpc2"
)

type LotusHandler struct {
	documents map[string]string
	analyzer  *Analyzer
}

func NewLotusHandler() *LotusHandler {
	return &LotusHandler{
		documents: make(map[string]string),
		analyzer:  NewAnalyzer(),
	}
}

func (h *LotusHandler) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
	switch req.Method {
	case "initialize":
		return h.initialize(), nil

	case "initialized", "$/setTrace", "$/cancelRequest":
		return nil, nil

	case "shutdown":
		return nil, nil

	case "textDocument/didOpen":
		var p DidOpenParams
		if err := json.Unmarshal(*req.Params, &p); err != nil {
			return nil, err
		}
		h.documents[p.TextDocument.URI] = p.TextDocument.Text
		return nil, nil

	case "textDocument/didChange":
		var p DidChangeParams
		if err := json.Unmarshal(*req.Params, &p); err != nil {
			return nil, err
		}
		if len(p.ContentChanges) > 0 {
			h.documents[p.TextDocument.URI] = p.ContentChanges[len(p.ContentChanges)-1].Text
		}
		return nil, nil

	case "textDocument/completion":
		var p CompletionParams
		if err := json.Unmarshal(*req.Params, &p); err != nil {
			return nil, err
		}
		return CompletionList{
			IsIncomplete: false,
			Items:        h.complete(p),
		}, nil

	case "textDocument/hover":
		var p HoverParams
		if err := json.Unmarshal(*req.Params, &p); err != nil {
			return nil, err
		}
		return h.hover(p), nil
	}

	return nil, nil
}

func (h *LotusHandler) initialize() InitializeResult {
	syncFull := 1
	return InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: syncFull,
			CompletionProvider: &CompletionOptions{
				TriggerCharacters: []string{"."},
			},
			HoverProvider: true,
		},
		ServerInfo: &ServerInfo{Name: lsName, Version: "0.1.0"},
	}
}

func (h *LotusHandler) complete(p CompletionParams) []CompletionItem {
	content, ok := h.documents[p.TextDocument.URI]
	if !ok {
		return nil
	}
	line := getLine(content, p.Position.Line)
	col := p.Position.Character
	if col > len(line) {
		col = len(line)
	}

	prefix := ""
	receiver := ""

	if col > 0 {
		trimmed := strings.TrimRight(line[:col], " \t")
		if strings.HasSuffix(trimmed, ".") {
			receiver = dotReceiver(trimmed)
		} else {
			prefix = lastWord(line[:col])
		}
	}

	return h.analyzer.Complete(content, prefix, receiver)
}

func (h *LotusHandler) hover(p HoverParams) *HoverResult {
	content, ok := h.documents[p.TextDocument.URI]
	if !ok {
		return nil
	}
	line := getLine(content, p.Position.Line)
	word := wordAt(line, p.Position.Character)
	if word == "" {
		return nil
	}
	doc := h.analyzer.HoverDoc(word)
	if doc == "" {
		return nil
	}
	return &HoverResult{
		Contents: MarkupContent{Kind: "markdown", Value: doc},
	}
}
