package config

import "errors"

// Config contains all of the variables that
// are needed for Stereodose to function
type Config struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AuthKey       string
	EncryptionKey string
}

// Verify can be used to conveniently check
// if the configuration was passed in correctly
func (c *Config) Verify() error {
	if c.ClientID == "" {
		return errors.New("STEREODOSE_CLIENT_ID was empty string")
	}
	if c.ClientSecret == "" {
		return errors.New("STEREODOSE_CLIENT_SECRET was empty string")
	}
	if c.RedirectURL == "" {
		return errors.New("STEREODOSE_REDIRECT_URL was empty string")
	}
	if c.AuthKey == "" {
		return errors.New("STEREODOSE_AUTH_KEY was empty string")
	}
	if c.EncryptionKey == "" {
		return errors.New("STEREODOSE_ENCRYPTION_KEY was empty string")
	}
	return nil
}
