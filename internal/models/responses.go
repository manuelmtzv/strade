package models

type ListMetadata struct {
	Total int `json:"total"`
}

type PostalCodeDetails struct {
	PostalCode     string `json:"postalCode"`
	Settlement     string `json:"settlement"`
	SettlementType string `json:"settlementType"`
	Municipality   string `json:"municipality"`
	State          string `json:"state"`
	City           string `json:"city"`
}
