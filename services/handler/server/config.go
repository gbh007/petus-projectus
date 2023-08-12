package server

type KafkaConfig struct {
	Topic         string
	GroupID       string
	Addr          string
	NumPartitions int
}

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type RabbitMQConfig struct {
	Username, Password, Addr, QueueName string
}
