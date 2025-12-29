package models

import (
	"encoding/json"
	"fmt"
)

func parseSecretData(secretType string, inData []byte) (SecretData, error) {
	switch secretType {
	case SecretTypePassword:
		var outData LoginData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case SecretTypeText:
		var outData TextData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case SecretTypeBinary:
		var outData FileData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	case SecretTypeCard:
		var outData CardData
		err := json.Unmarshal(inData, &outData)
		return outData, err
	default:
		return nil, fmt.Errorf("unknown secret type: %s", secretType)
	}
}
