package stores

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/silicongreenhouse/api/src/models"
)

var config *ConfigState

type Signal bool
var signalchannel = make(chan Signal)

var ClientConnected = false
var ClientControllerConnected = false
var RaspberryConnected = false
var RaspberryControllerConnected = false

type ConfigState struct {
	Sensors   []models.Sensor   `json:"sensors"`
	Executors []models.Executor `json:"executors"`
}

type ConfigStore struct {
	State *ConfigState
	SignalChannel chan Signal
}

func UseConfig() ConfigStore {
	return ConfigStore{
		State: config,
		SignalChannel: signalchannel,
	}
}

func (self *ConfigStore) Load(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Cannot read file %s", path)
	}

	err = json.Unmarshal(file, &self.State)

	return err
}

func (self *ConfigStore) Write(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Error opening file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(&self.State)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Error writing file")
	}

	if RaspberryConnected {
		self.SignalChannel <- true
	}

	return nil
}

func init() {
	config = &ConfigState{}
}
