package model

type Config struct {
	Port      string
	Host      string
	FromKafka Kafka
	ToKafka   Kafka
	DBconf
}

type Kafka struct {
	Brokers []string
	Topic   string
	GroupID string
}
type DBconf struct {
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

type MessageFromKafka struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type MessageToKafka struct {
	MessageFromKafka []byte `json:"message_from_kafka"`
	Error            error  `json:"error,omitempty"`
}

type MessageToDB struct {
	UID        string
	Name       string
	Surname    string
	Patronymic string
	Age        int
	Gender     string
	Nation     string
	Error      string
}
