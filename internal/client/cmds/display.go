package cmds

import (
	"fmt"
	"strings"

	"github.com/s-turchinskiy/keeper/internal/client/models"
)

const timeFormat = "2006-01-02 15:04:05"

func displaySecrets(responses []*models.LocalSecret) {
	if len(responses) == 0 {
		fmt.Println("No secrets found")
		return
	}

	fmt.Printf("%-12s %-12s %s\n", "Name", "Type", "Last Modified")
	fmt.Println(strings.Repeat("-", 82))
	for _, resp := range responses {
		fmt.Printf("%-12s %-12s %s\n",
			resp.Name,
			resp.Type,
			resp.LastModified.Local().Format(timeFormat))
	}
}

func displaySecret(secret *models.LocalSecret, full bool) error {
	fmt.Printf("Name: %s\n", secret.Name)
	fmt.Printf("Type: %s\n", secret.Type)
	fmt.Printf("Last Modified: %s\n", secret.LastModified.Local().Format(timeFormat))
	if secret.Metadata != "" {
		fmt.Printf("Metadata: %s\n", secret.Metadata)
	}
	fmt.Println()

	data, err := secret.ParseData()
	if err != nil {
		return err
	}

	switch data := data.(type) {
	case models.LoginData:
		fmt.Printf("Username: %s\n", data.Username)
		if full {
			fmt.Printf("Password: %s\n", data.Password)
		} else {
			fmt.Printf("Password: ********\n")
		}
		if data.URL != "" {
			fmt.Printf("URL: %s\n", data.URL)
		}

	case models.TextData:
		fmt.Printf("Content: %s\n", data.Content)

	case models.FileData:
		fmt.Printf("File Name: %s\n", data.FileName)
		fmt.Printf("File Size: %d bytes\n", data.FileSize)
		fmt.Printf("Content Size: %d bytes (base64)\n", len(data.Content))

	case models.CardData:
		fmt.Printf("Card Number: %s\n", data.Number)
		fmt.Printf("Card Holder: %s\n", data.Holder)
		fmt.Printf("Expiry: %s\n", data.Expiry)
		if full {
			fmt.Printf("CVV: %s\n", data.CVV)
		} else {
			fmt.Printf("CVV: ***\n")
		}

	default:
		return fmt.Errorf("unknown data type: %T", data)
	}

	return nil
}
