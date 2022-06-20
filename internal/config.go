package trove

// trove configuration
type Config struct {
	// Broker Addresses for Kafka
	Broker []string
	// Address for Server
	Address string
	// Certificate Path
	Certificate string
	// Key Path
	Key string
}

func getConfig(config ...Config) Config {
	if len(config) != 0 {
		return config[0]
	} else {
		return Config{
			Broker: []string{"127.0.0.1:9092"},
		}
	}
}
