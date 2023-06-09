package queue

import "encoding/json"

type QueueDTO struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	ID       int    `json:"id"`
}

func (dto *QueueDTO) Marshal() ([]byte, error) {
	return json.Marshal(dto)
}

func (dto *QueueDTO) Unmarshal(data []byte) error {
	return json.Unmarshal(data, dto)
}
