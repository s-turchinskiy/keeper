package models

import (
	"fmt"
	"regexp"
	"strings"
)

type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url,omitempty"`
}

func (d LoginData) Validate() error {
	username := strings.TrimSpace(d.Username)
	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) > MaxUsernameLength {
		return fmt.Errorf("username too long: %d characters (max: %d)",
			len(username), MaxUsernameLength)
	}

	password := strings.TrimSpace(d.Password)
	if password == "" {
		return fmt.Errorf("password is required")
	}
	if len(password) > MaxPasswordLength {
		return fmt.Errorf("password too long: %d characters (max: %d)",
			len(password), MaxPasswordLength)
	}

	if d.URL != "" {
		url := strings.TrimSpace(d.URL)
		if len(url) > MaxURLLength {
			return fmt.Errorf("URL too long: %d characters (max: %d)",
				len(url), MaxURLLength)
		}
	}

	return nil
}

type TextData struct {
	Content string `json:"content"`
}

func (d TextData) Validate() error {
	if strings.TrimSpace(d.Content) == "" {
		return fmt.Errorf("content is required")
	}
	if len(d.Content) > MaxTextSize {
		return fmt.Errorf("text content too long: %d characters (max: %d)",
			len(d.Content), MaxTextSize)
	}
	return nil
}

type FileData struct {
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Content  string `json:"content"`
}

func (d FileData) Validate() error {
	if strings.TrimSpace(d.FileName) == "" {
		return fmt.Errorf("file name is required")
	}
	if strings.TrimSpace(d.Content) == "" {
		return fmt.Errorf("content is required")
	}
	if d.FileSize > MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum %d bytes",
			d.FileSize, MaxFileSize)
	}
	return nil
}

type CardData struct {
	Number string `json:"number"`
	Holder string `json:"holder"`
	Expiry string `json:"expiry"`
	CVV    string `json:"cvv,omitempty"`
}

func (d CardData) Validate() error {
	cardRegex := regexp.MustCompile(`^\d{13,19}$`)
	if !cardRegex.MatchString(strings.ReplaceAll(d.Number, " ", "")) {
		return fmt.Errorf("invalid card number format")
	}

	holder := strings.TrimSpace(d.Holder)
	if holder == "" {
		return fmt.Errorf("card holder name is required")
	}
	if len(holder) > MaxCardHolderLength {
		return fmt.Errorf("card holder name too long: %d characters (max: %d)",
			len(holder), MaxCardHolderLength)
	}

	expiryRegex := regexp.MustCompile(`^(0[1-9]|1[0-2])\/(\d{2}|\d{4})$`)
	if !expiryRegex.MatchString(d.Expiry) {
		return fmt.Errorf("invalid expiry format, use MM/YY or MM/YYYY")
	}

	if strings.TrimSpace(d.CVV) == "" {
		return fmt.Errorf("CVV is required")
	}

	cvvRegex := regexp.MustCompile(`^\d+$`)
	if !cvvRegex.MatchString(d.CVV) {
		return fmt.Errorf("CVV must contain only digits")
	}

	if len(d.CVV) != 3 && len(d.CVV) != 4 {
		return fmt.Errorf("CVV must be 3 or 4 digits, got %d", len(d.CVV))
	}

	return nil
}
