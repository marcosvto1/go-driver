package queue

import "errors"

type MockQueueConfig struct {
	PublishWillReturnErr bool
}

type MockQueueConnection struct {
	q           []*QueueDTO
	mockOptions MockQueueConfig
}

func (m *MockQueueConnection) Publish(b []byte) error {

	if m.mockOptions.PublishWillReturnErr {
		return errors.New("publish err")
	}

	dto := new(QueueDTO)
	dto.Unmarshal(b)

	m.q = append(m.q, dto)

	return nil
}

func (m *MockQueueConnection) Consume(ch chan<- QueueDTO) error {
	data := QueueDTO{
		Filename: "any",
		ID:       1,
	}

	ch <- data

	return nil
}
