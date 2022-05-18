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

type ConfigState struct {
	Sensors   []models.Sensor   `json:"sensors"`
	Executors []models.Executor `json:"executors"`
}

type ConfigStore struct {
	*ConfigState
}

func UseConfig() ConfigStore {
	return ConfigStore{
		config,
	}
}

func (self *ConfigStore) Load(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Cannot read file %s", path)
	}

	err = json.Unmarshal(file, self)

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
	err = encoder.Encode(self)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Error writing file")
	}

	return nil
}

func init() {
	config = &ConfigState{}
}
