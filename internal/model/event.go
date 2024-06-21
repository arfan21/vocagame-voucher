package model

type EventName string

func (e EventName) MarshalBinary() ([]byte, error) {
	return []byte(e), nil
}

func (e *EventName) UnmarshalBinary(data []byte) error {
	*e = EventName(data)
	return nil
}

const (
	EventTransactionNotification EventName = "transaction_notification"
)

const (
	StreamNotification = "notification"
)

type Event struct {
	Name          EventName `json:"name"`
	TransactionID string    `json:"transaction_id"`
	Email         string    `json:"email"`
	Url           string    `json:"url"`
	RetryCount    int       `json:"retry_count"`
}
