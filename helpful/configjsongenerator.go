package helpful

import (
	"encoding/json"
	"fmt"
	"strings"
)

type jsonConfigGenerator struct {
	Root   map[string]interface{}
	prefix string
}

func (j *jsonConfigGenerator) GetInterfaceArray(path string) ([]interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (j *jsonConfigGenerator) GetInterface(path string) (interface{}, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (j *jsonConfigGenerator) GetArray(path string) ([]Config, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (j *jsonConfigGenerator) GetString(path string) (string, error) {
	if j.prefix != "" {
		path = j.prefix + jsonPathDelimiter + path
	}
	paths := strings.Split(path, jsonPathDelimiter)
	m := j.Root
	for _, p := range paths[:len(paths)-1] {
		_, ok := m[p]
		if !ok {
			m[p] = make(map[string]interface{})
		}
		m = m[p].(map[string]interface{})
	}
	s := "string"
	m[paths[len(paths)-1]] = &s
	return s, nil
}

func (j *jsonConfigGenerator) GetInt(path string) (int, error) {
	if j.prefix != "" {
		path = j.prefix + jsonPathDelimiter + path
	}
	paths := strings.Split(path, jsonPathDelimiter)
	m := j.Root
	for _, p := range paths[:len(paths)-1] {
		_, ok := m[p]
		if !ok {
			m[p] = make(map[string]interface{})
		}
		m = m[p].(map[string]interface{})
	}
	i := 1
	m[paths[len(paths)-1]] = &i
	return i, nil
}

func (j *jsonConfigGenerator) Child(path string) Config {
	if j.prefix != "" {
		path = j.prefix + jsonPathDelimiter + path
	}
	return &jsonConfigGenerator{
		Root:   j.Root,
		prefix: path,
	}
}

func (j *jsonConfigGenerator) Generate() ([]byte, error) {
	return json.MarshalIndent(j.Root, "", "\t")
}
