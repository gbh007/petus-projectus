package server

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type RabbitMQConfig struct {
	Username, Password, Addr, QueueName string
}

type Config struct {
	NotificationAddress string
	PrometheusAddress   string
	DB                  DBConfig
	RabbitMQ            RabbitMQConfig
}
