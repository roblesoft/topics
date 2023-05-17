package models

type Message struct {
	Content  string `json:"content,omitempty"`
	Channel  string `json:"channel,omitempty"`
	Username string `json:"username,omitempty"`
	Command  int    `json:"command,omitempty"`
	Err      string `json:"err,omitempty"`
}
