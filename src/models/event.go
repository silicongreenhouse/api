package models

type Event struct {
	Id string `json:"id"`
	Executor string `json:"executor"`
	State string `json:"state"`
	Value string `json:"value,omitemtpy"`
	Above float32 `json:"above,omitemtpy"`
	Equal float32 `json:"equal,omitemtpy"`
	Below float32 `json:"below,omitemtpy"`
}

