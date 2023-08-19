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

type Config struct {
	SelfAddress       string
	PrometheusAddress string
	Kafka             KafkaConfig
	DB                DBConfig
}
