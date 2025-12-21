package models

import "encoding/json"

type SecretDataContainer struct {
	Type       string     `json:"type"`
	Name       string     `json:"name"`
	SecretData SecretData `json:"-"`
}

func (c *SecretDataContainer) MarshalJSON() ([]byte, error) {
	type Alias SecretDataContainer
	return json.Marshal(&struct {
		*Alias
		SecretData json.RawMessage `json:"secret_data"`
	}{
		Alias: (*Alias)(c),
		SecretData: func() json.RawMessage {
			data, _ := json.Marshal(c.SecretData)
			return data
		}(),
	})
}

func (c *SecretDataContainer) UnmarshalJSON(data []byte) error {
	type Alias SecretDataContainer
	aux := &struct {
		*Alias
		SecretData json.RawMessage `json:"secret_data"`
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	parsedData, err := parseSecretData(c.Type, aux.SecretData)
	if err != nil {
		return err
	}

	c.SecretData = parsedData

	return nil
}
