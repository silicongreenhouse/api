package sensors

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/silicongreenhouse/api/src/models"
)

func createNewId(sensorId string, executorId string) (string, error) {
	newId := ""
	sensors := config.Sensors
	var events []models.Event

	// Fill the events variable with the current sensor events
	for _, sensor := range sensors {
		if sensor.Id == sensorId {
			events = sensor.Events
		}
	}

	// Creating new id
	// If there are no events then the id will be the sensorid + i1
	if !(len(events) > 0) {
		newId = fmt.Sprintf("%s%s", executorId, "i1")
		// If there are events then get the index of the last event and increment 1
	} else {
		regex := regexp.MustCompile(`\d+$`)
		lastEvent := events[len(events)-1]
		lastIndex := string(regex.Find([]byte(lastEvent.Id)))

		convertedIndex, err := strconv.Atoi(lastIndex)
		if err != nil {
			return "", fmt.Errorf("Error converting index")
		}

		newIndex := fmt.Sprint(convertedIndex + 1)
		newId = fmt.Sprintf("%si%s", executorId, newIndex)
	}

	return newId, nil
}

func deleteElement[T any](slice *[]T, index int) {
	*slice = append((*slice)[:index], (*slice)[index + 1:]...)
}
