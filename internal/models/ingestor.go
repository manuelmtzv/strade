package models

type RawDataRecord struct {
	PostalCode         string `json:"d_codigo"`
	Settlement         string `json:"d_asenta"`
	SettlementType     string `json:"d_tipo_asenta"`
	Municipality       string `json:"D_mnpio"`
	State              string `json:"d_estado"`
	City               string `json:"d_ciudad"`
	AdminPostalCode    string `json:"d_CP"`
	StateCode          string `json:"c_estado"`
	OfficePostalCode   string `json:"c_oficina"`
	EmptyField         string `json:"c_CP"`
	SettlementTypeCode string `json:"c_tipo_asenta"`
	MunicipalityCode   string `json:"c_mnpio"`
	SettlementID       string `json:"id_asenta_cpcons"`
	Zone               string `json:"d_zona"`
	CityCode           string `json:"c_eve_ciudad"`
}
