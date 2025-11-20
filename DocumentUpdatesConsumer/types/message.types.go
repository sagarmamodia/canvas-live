package types

type Message struct {
	DocumentID string `json:"documentId"`
	UserID     string `json:"userId"`
	Username   string `json:"username"`
	Type       int    `json:"type"`
	Body       string `json:"body"`
}
