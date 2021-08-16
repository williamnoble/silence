package silence

// Configuration holds information regarding the server specifics.
type Configuration struct {
	Host string
	Port string
}

// NewConfiguration returns sensible defaults for hosting the server.
func NewConfiguration() (*Configuration, error) {
	return &Configuration{
		Host: "127.0.0.1",
		Port: "8000",
	}, nil
}
