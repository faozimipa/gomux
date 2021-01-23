package models

type Payload struct {
	From    From    `json:"from"`
	To      To      `json:"to"`
	Message Message `json:"message"`
}

type From struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type To struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Message struct {
	Content Content `json:"content"`
}
