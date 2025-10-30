package models

type ListMetadata struct {
	Total int `json:"total"`
	Limit int `json:"limit,omitempty"`
}

type PostalCodeDetails struct {
	PostalCode     string `json:"postalCode"`
	Settlement     string `json:"settlement"`
	SettlementType string `json:"settlementType"`
	Municipality   string `json:"municipality"`
	State          string `json:"state"`
	City           string `json:"city"`
}

type SettlementSearchResult struct {
	Name           string `json:"name"`
	PostalCode     string `json:"postalCode"`
	SettlementType string `json:"settlementType"`
	Municipality   string `json:"municipality"`
	State          string `json:"state"`
	City           string `json:"city"`
}
