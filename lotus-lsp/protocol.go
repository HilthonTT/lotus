package main

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync   int                `json:"textDocumentSync"`
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`
	HoverProvider      bool               `json:"hoverProvider,omitempty"`
}

type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label         string         `json:"label"`
	Kind          *int           `json:"kind,omitempty"`
	Detail        *string        `json:"detail,omitempty"`
	Documentation *MarkupContent `json:"documentation,omitempty"`
}

type MarkupContent struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

type HoverResult struct {
	Contents MarkupContent `json:"contents"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type TextDocumentItem struct {
	URI  string `json:"uri"`
	Text string `json:"text"`
}

type DidOpenParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type ContentChangeEvent struct {
	Text string `json:"text"`
}

type DidChangeParams struct {
	TextDocument   TextDocumentIdentifier `json:"textDocument"`
	ContentChanges []ContentChangeEvent   `json:"contentChanges"`
}

type CompletionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type HoverParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}
