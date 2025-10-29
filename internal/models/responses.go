package models

type PostalCodeDetails struct {
	PostalCode     string `json:"postal_code"`
	City           string `json:"city"`
	State          string `json:"state"`
	SettlementName string `json:"settlement_name"`
	SettlementType string `json:"settlement_type"`
}
