package helpful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

const jsonPathDelimiter = "."

func NewJsonCfg(fileName string) (Config, error) {
	return newJsonCfg(fileName)
}

func newJsonCfg(fileName string) (*jsonConfig, error) {
	cfg := &jsonConfig{
		cfg: make(map[string]interface{}),
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("cant read config: %v", err)
	}

	err = json.Unmarshal(data, &cfg.cfg)
	if err != nil {
		return nil, fmt.Errorf("cant unmarshal config: %v", err)
	}

	return cfg, nil
}

func NewJsonCfgWithGenerator(fileName string) (ConfigGenerator, error) {
	cfg, err := newJsonCfg(fileName)
	if err != nil {
		return nil, err
	}
	cfg.generator = NewJsonConfigGenerator()
	return cfg, nil
}

type jsonConfig struct {
	cfg        map[string]interface{}
	pathPrefix string
	generator  *jsonConfigGenerator
}

func (j *jsonConfig) Child(path string) Config {
	if j.generator != nil {
		j.generator.Child(path)
	}
	if j.pathPrefix != "" {
		path = j.pathPrefix + jsonPathDelimiter + path
	}
	return &jsonConfig{
		cfg:        j.cfg,
		pathPrefix: path,
	}
}

func (j *jsonConfig) GetArray(path string) ([]Config, error) {
	if j.generator != nil {
		j.generator.GetArray(path)
	}
	val, err := j.getValByPath(path)
	if err != nil {
		return nil, err
	}

	switch arr := val.(type) {
	case []interface{}:
		res := make([]Config, 0)
		for i, v := range arr {
			switch m := v.(type) {
			case map[string]interface{}:
				res = append(res, &jsonConfig{
					cfg:        m,
					pathPrefix: "",
				})
			default:
				return nil, fmt.Errorf("element no %v value of array by path %v is not a json object", i, path)
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf("value by path %v is not an array", path)
	}

}

func (j *jsonConfig) GetInterface(path string) (interface{}, error) {
	return j.getValByPath(path)
}

func (j *jsonConfig) getValByPath(path string) (interface{}, error) {
	if j.pathPrefix != "" {
		path = j.pathPrefix + jsonPathDelimiter + path
	}
	names := strings.Split(path, jsonPathDelimiter)
	var v interface{} = j.cfg
	for _, name := range names {
		switch m := v.(type) {
		case map[string]interface{}:
			var exists bool
			v, exists = m[name]
			if !exists {
				return nil, fmt.Errorf("cant get value for %v, element %v doesn't exist", path, name)
			}
		default:
			return nil, fmt.Errorf("cant get value for %v, element %v doesn't exist", path, name)
		}
	}
	return v, nil
}

func (j *jsonConfig) GetInt(path string) (int, error) {
	if j.generator != nil {
		j.generator.GetInt(path)
	}
	val, err := j.getValByPath(path)
	if err != nil {
		return 0, err
	}

	if res, ok := val.(int); ok {
		return res, nil
	}
	switch k := reflect.ValueOf(val).Kind(); k {
	case reflect.Float32:
		return int(val.(float32)), nil
	case reflect.Float64:
		return int(val.(float64)), nil
	case reflect.String:
		val, err := strconv.Atoi(val.(string))
		if err != nil {
			return 0, err
		}
		return val, nil
	default:
		return 0, fmt.Errorf("value of %v is %v type and has unconvertable value %v", path, reflect.TypeOf(val), val)
	}
}

func (j *jsonConfig) GetString(path string) (string, error) {
	if j.generator != nil {
		j.generator.GetString(path)
	}
	val, err := j.getValByPath(path)
	if err != nil {
		return "", err
	}
	switch t := val.(type) {
	case string:
		return val.(string), nil
	default:
		return "", fmt.Errorf("value %v is %v type", path, t)
	}
}

func (j *jsonConfig) Generate() ([]byte, error) {
	if j.generator != nil {
		return j.generator.Generate()
	}
	return nil, fmt.Errorf("nil generator in config, something went wrong")
}
