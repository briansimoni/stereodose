package config

import "testing"

func TestVerify(t *testing.T) {
	c := &Config{
		ClientID:      "aspoidfjpoijadsf",
		ClientSecret:  "apiosjdfposijdf",
		AuthKey:       "apiodsjfpojasf",
		RedirectURL:   "asjkdfpoija",
		EncryptionKey: "apsoifjdfjds",
	}

	err := c.Verify()
	if err != nil {
		t.Errorf("Expected verification to pass. Got: %s", err.Error())
	}

	c.ClientID = ""
	err = c.Verify()
	if err == nil {
		t.Error("ClientID was empty string, Expected non-nil error")
	}
}
