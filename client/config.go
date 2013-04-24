// used to config different settings like standard nick
package client

import (
	"os"
	"strings"
)

// Config type stores cfg values which can be saved and loaded from a file
type Config struct {
	cfgValues *map[string]string
	file      *File
}

// NewConfig takes a filename, open or creates the file depending on if it
// exists or not and return a pointer to a Config type.
func NewConfig(fileName string) (*Config, err) {
	file * File
	cfg * Config

	file, err = os.OpenFile(fileName, os.O_CREAT|os.O_RDWR|os.O_TRUNC, 0666) // open file or create if it doesn't exist.

	if err != nil {
		return nil, err
	}

	return &Config{make(map[string]string), file}, nil
}

// Adds a config value. If it already exists it is overwritten with the new value
func (c *Config) AddCfgValue(key, value string) {
	cfgValues[key] = value
}

// Removes the config value with the given key
func (c *Config) RemoveCfgValue(key string) {
	delete(c.cfgValues, key)
}

// Loads the values from the config file
func (c *Config) Load() err {
	reader := bufio.NewReader(c.file)

	for {
		ln, err := reader.ReadString('\n')
		switch err {
		case nil:
			// MIGHT BE BEST TO CHANGE HOW WE HANDLE THE DATA HERE IF FILE ISN'T AS WE EXPECT
			s[] := ln.Split(" ")
			cfgValues[s[0]] = s[1]
		case io.EOF:
			return
		case default:
			return err
		}
	}
}

// Saves the values in the config file
func (c *Config) Save() err {
	writer := bufio.NewWriter(c.file)
	for k, v := range cfgValues {
		s := k + " " + v + "\n"

		_, err := writer.WriteString(s)
		if err != nil {
			return err
		}
	}
	err := writer.Flush()
	if err != nil {
		return err
	}
}

// Closes the config file, should be called when done with it
func (c *Config) Close() err {
	return file.Close()
}
