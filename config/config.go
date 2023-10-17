// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Config package save and load local config files.
package config

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

const (
	configBaseDir  = "teonet"
)

// config stucture and methods to store teocrypt config to local host.
type config[T any] struct {
	Data     *T     `json:"data"`
	fileName string `json:"-"`
}

// New creates config object.
func New[T any](appShortName, configName string, data ...*T) (cfg *config[T], err error) {
	cfg = new(config[T])
	if len(data) > 0 {
		cfg.Data = data[0]
	} else {
		cfg.Data = new(T)
	}
	if cfg.fileName, err = cfg.createFolder(appShortName, configName); err != nil {
		return
	}
	return
}

// loadConfig creates config object and load config from local host file.
func Load[T any](appShortName, configName string, data ...*T) (cfg *config[T], err error) {
	cfg, err = New[T](appShortName, configName, data...)
	cfg.load()
	return
}

// save config to local host file.
func (c config[T]) Save() (err error) {

	f, err := os.Create(c.fileName)
	if err != nil {
		return
	}

	data, err := c.Marshal()
	if err != nil {
		return
	}

	_, err = f.Write(data)

	return
}

// Marshal config.
func (c config[T]) Marshal() (data []byte, err error) {
	data, err = json.MarshalIndent(c, "", " ")
	if err != nil {
		return
	}
	return
}

// Unmarshal config.
func (c *config[T]) Unmarshal(data []byte) error {
	return json.Unmarshal(data, c)
}

// load config file.
func (c *config[T]) load() (err error) {

	// Open file
	f, err := os.Open(c.fileName)
	if err != nil {
		return
	}

	// Read file data
	data, err := io.ReadAll(f)
	if err != nil {
		return
	}

	// Unmarshal config data
	err = c.Unmarshal(data)
	if err != nil {
		return
	}

	return
}

// createFolder creates config folder and return config fileName.
func (c *config[T]) createFolder(appShortName, configFileName string) (fileName string, err error) {

	fileName, err = c.getFileName(configBaseDir, appShortName, configFileName+".cfg")
	if err != nil {
		return
	}

	err = os.MkdirAll(path.Dir(fileName), os.ModePerm)
	if err != nil {
		return
	}

	return
}

// getFileName gets full path and file name of config file in os config folder.
func (c *config[T]) getFileName(configDir, appShortName, fileName string) (
	res string, err error) {

	res, err = os.UserConfigDir()
	if err != nil {
		return
	}

	res += "/" + configDir + "/" + appShortName + "/" + fileName
	return
}
