package sms

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/faozimipa/gomux/api/models"
)

func SendMessage(message string, phone string) (*http.Response, error) {
	data := models.Payload{
		From: models.From{
			Type:   "sms",
			Number: "Nexmo",
		},
		To: models.To{
			Type:   "sms",
			Number: phone,
		},
		Message: models.Message{
			Content: models.Content{
				Type: "text",
				Text: message,
			},
		},
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.nexmo.com/v0.1/messages", body)
	if err != nil {
		return nil, err
	}
	//Ensure headers
	api_key := os.Getenv("NEXMO_API_KEY")
	secret_key := os.Getenv("NEXMO_API_SECRET")
	req.SetBasicAuth(api_key, secret_key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
