package helpful

import "fmt"

func NewMapConfig(m map[string]interface{}) (Config, error){
	if m == nil{
		return nil, fmt.Errorf("must be not-nul map")
	}
	return &jsonConfig{
		cfg:        m,
	}, nil
}
