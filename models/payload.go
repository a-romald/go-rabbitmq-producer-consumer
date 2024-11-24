package models

type Payload struct {
	Word   string `json:name`
	Action string `json:"action"`
}
