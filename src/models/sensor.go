package models

type Sensor struct {
	Id string `json:"id"`
	Name string `json:"name"`
	ShortName string `json:"short_name"`
	Events []Event `json:"events,omitempty"`
}

