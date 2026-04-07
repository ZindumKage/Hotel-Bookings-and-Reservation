package audit_logs

import "log"

type LogPublisher struct{}

func NewLogPublisher() *LogPublisher {
	return &LogPublisher{}
}

func (l *LogPublisher) Publish(topic string, data []byte) error {

	log.Printf("AUDIT EVENT [%s]: %s\n", topic, string(data))

	return nil
}