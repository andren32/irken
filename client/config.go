// used to config different settings like standard nick
package client

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Config type stores cfg values which can be saved and loaded from a file
type Config struct {
	cfgValues map[string]string
	fileName  string
}

// NewConfig takes a filename, open or creates the file depending on if it
// exists or not and return a pointer to a Config type.
func NewConfig(fileName string) (c *Config, err error) {
	c = &Config{make(map[string]string), fileName}
	err = c.Load()
	if err != nil {
		c = &Config{make(map[string]string), fileName}
		return
	}
	return
}

// Adds a config value. If it already exists it is overwritten with the new value
func (c *Config) AddCfgValue(key, value string) {
	c.cfgValues[key] = value
}

// Removes the config value with the given key
func (c *Config) RemoveCfgValue(key string) {
	delete(c.cfgValues, key)
}

// Removes the config value with the given key
func (c *Config) RemoveAllCfgValues() {
	c.cfgValues = make(map[string]string)
}

// Returns pointer to the cfg values.
func (c *Config) GetCfgValues() map[string]string {
	return c.cfgValues
}

// Loads the values from the config file
func (c *Config) Load() error {
	file, err := os.Open(c.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	c.RemoveAllCfgValues()
	reader := bufio.NewReader(file)

	for {
		ln, err := reader.ReadString('\n')
		switch err {
		case nil:
			// MIGHT BE BEST TO CHANGE HOW WE HANDLE THE DATA HERE IF FILE ISN'T AS WE EXPECT
			s := strings.SplitN(strings.Replace(ln, "\n", "", -1), " ", 2)
			c.AddCfgValue(s[0], s[1])
		case io.EOF:
			return nil
		default:
			return err
		}
	}

	return nil
}

// Saves the values in the config file
func (c *Config) Save() error {
	file, err := os.Create(c.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for k, v := range c.cfgValues {
		s := k + " " + v + "\n"

		_, err := writer.WriteString(s)
		if err != nil {
			return err
		}
	}
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
