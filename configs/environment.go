package configs

// Environment additional in this service
type Environment struct {
	
}

var env Environment

// GetEnv get global additional environment
func GetEnv() Environment {
	return env
}

func loadAdditionalEnv() {
	
}
