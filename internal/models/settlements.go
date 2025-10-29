package models

type SettlementType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Settlement struct {
	ID               string          `json:"id"`
	PostalCode       string          `json:"postalCode"`
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	SettlementTypeID string          `json:"settlementTypeId,omitempty"`
	SettlementType   *SettlementType `json:"settlementType"`
	MunicipalityID   string          `json:"municipalityId,omitempty"`
	Municipality     *Municipality   `json:"municipality"`
	CityID           string          `json:"cityId,omitempty"`
	City             *City           `json:"city"`
	StateID          string          `json:"stateId,omitempty"`
	State            *State          `json:"state"`
	OfficePostalCode string          `json:"officePostalCode"`
	Zone             string          `json:"zone"`
	Latitude         *float64        `json:"latitude,omitempty"`
	Longitude        *float64        `json:"longitude,omitempty"`
}
