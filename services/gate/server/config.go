package server

type KafkaConfig struct {
	Topic         string
	GroupID       string
	Addr          string
	NumPartitions int
}

type CommunicationConfig struct {
	SelfAddress         string
	AuthAddress         string
	LogAddress          string
	NotificationAddress string
}
