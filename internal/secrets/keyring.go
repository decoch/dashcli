package secrets

import "github.com/99designs/keyring"

const (
	keyAPIKey  = "api_key"
	keyBaseURL = "base_url"
)

var openKeyring = func() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName:     "dashcli",
		AllowedBackends: keyring.AvailableBackends(),
	})
}

func SetAPIKey(apiKey string) error {
	return setValue(keyAPIKey, apiKey)
}

func GetAPIKey() (string, error) {
	return getValue(keyAPIKey)
}

func DeleteAPIKey() error {
	return deleteValue(keyAPIKey)
}

func SetBaseURL(baseURL string) error {
	return setValue(keyBaseURL, baseURL)
}

func GetBaseURL() (string, error) {
	return getValue(keyBaseURL)
}

func DeleteBaseURL() error {
	return deleteValue(keyBaseURL)
}

func setValue(key, value string) error {
	ring, err := openKeyring()
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:  key,
		Data: []byte(value),
	})
}

func getValue(key string) (string, error) {
	ring, err := openKeyring()
	if err != nil {
		return "", err
	}
	item, err := ring.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func deleteValue(key string) error {
	ring, err := openKeyring()
	if err != nil {
		return err
	}
	return ring.Remove(key)
}
