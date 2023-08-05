package gate

type KafkaConfig struct {
	Topic         string
	GroupID       string
	Addr          string
	NumPartitions int
}
