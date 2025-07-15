package types

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(content string) *Message {
	return &Message{
		ID:        randomId(),
		Content:   content,
		Timestamp: time.Now(),
	}
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func MessageFromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// talvez mudar isso para um uuid?
func randomId() string {
	// format -> ano-mes-dia-hora-minuto-segundos : sem o hífen
	return time.Now().Format("20060102150405") + "-" + time.Now().Format("000000")
}
