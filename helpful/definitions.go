package helpful

// Configuration interface
// provides some values by given path
// path delimter must be dot (".")
type Config interface {
	GetString(path string) (string, error)
	GetInt(path string) (int, error)
	GetInterface(path string) (interface{}, error)
	GetArray(path string) ([]Config, error)
	GetInterfaceArray(path string) ([]interface{}, error)
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

type Printer interface {
	Printf(format string, a ...interface{}) (n int, err error)
}

type LogsCache interface {
	GetLastLogs() []string
}
