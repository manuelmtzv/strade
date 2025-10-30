package models

type City struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	StateID string `json:"stateId,omitempty"`
	State   *State `json:"state,omitempty"`
}
