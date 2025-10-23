package entity

type CustomField struct {
	FieldID   int    `json:"field_id"`
	FieldName string `json:"field_name"`
	FieldCode string `json:"field_code"`
	FieldType string `json:"field_type"`
	Values    []struct {
		Value    string `json:"value"`
		EnumID   int    `json:"enum_id"`
		EnumCode string `json:"enum_code"`
	} `json:"values"`
}
