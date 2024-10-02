package nodetreemodel

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func (c *ntmConfig) getConfigFile() string {
	if c.configFile == "" {
		return "datadog.yaml"
	}
	return c.configFile
}

// ReadInConfig wraps Viper for concurrent access
func (c *ntmConfig) ReadInConfig() error {
	c.Lock()
	defer c.Unlock()
	// ReadInConfig reset configuration with the main config file
	err := c.readInConfig()
	if err != nil {
		return err
	}

	type extraConf struct {
		path    string
		content []byte
	}

	// Read extra config files
	// TODO: handle c.extraConfigFilePaths, read and merge files
	return nil
}

func (c *ntmConfig) readInConfig() error {
	filename := c.getConfigFile()
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	root, err := c.readConfigurationContent(content)
	if err != nil {
		return err
	}
	c.root = root
	return nil
}

// ReadConfig wraps Viper for concurrent access
func (c *ntmConfig) ReadConfig(in io.Reader) error {
	c.Lock()
	defer c.Unlock()
	content, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	root, err := c.readConfigurationContent(content)
	if err != nil {
		return err
	}
	c.root = root
	return nil
}

func (c *ntmConfig) readConfigurationContent(content []byte) (Node, error) {
	var obj map[string]interface{}
	if err := yaml.Unmarshal(content, &obj); err != nil {
		return nil, err
	}
	root, err := NewNode(obj)
	if err != nil {
		return nil, err
	}
	return root, nil
}
