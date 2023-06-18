package queue

type MockQueueConnection struct {
	q []*QueueDTO
}

func (m *MockQueueConnection) Publish(b []byte) error {
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
