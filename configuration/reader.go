package configuration

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Reader interface {
	Read(filename string) (*Configuration, error)
}

func NewReader() Reader {
	return &YamlReader{}
}

type YamlReader struct{}

func (yr *YamlReader) Read(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var c Configuration
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
