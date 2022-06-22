package trove

// trove configuration
type Config struct {
	// Brokers for Kafka
	Brokers []string
	// Address for Server
	Address string
	// Certificate Path
	Certificate string
	// Key Path
	Key string
	// Consumer Group ID
	GroupID string
	// MongoDB URI
	MongoURI string
}

func getConfig(config ...Config) Config {
	if len(config) != 0 {
		return config[0]
	} else {
		return Config{
			Brokers:     []string{"127.0.0.1:9092"},
			Address:     ":443",
			Certificate: "",
			Key:         "",
			GroupID:     "unassigned-group",
			MongoURI:    "",
		}
	}
}
