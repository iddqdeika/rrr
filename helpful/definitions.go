package helpful


// Configuration interface
// provides some values by given path
// path delimter must be dot (".")
type Config interface {
	GetString(path string) (string, error)
	GetInt(path string) (int, error)
	GetArray(path string) ([]Config, error)
	Child(path string) Config
}

// Config interface.
// Also can provide data for Configuration source template (json footprint, etc)
type ConfigGenerator interface {
	Generate() ([]byte, error)
	Config
}

// Logging interface
type Logger interface {
	Errorf(format string, a ...interface{})
	Infof(format string, a ...interface{})
}
