package models

type Municipality struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	StateID string `json:"stateId"`
	State   *State `json:"state"`
}
