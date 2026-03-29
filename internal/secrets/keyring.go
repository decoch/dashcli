package secrets

import "github.com/99designs/keyring"

func Set(profile, apiKey string) error {
	ring, err := open()
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:  profile,
		Data: []byte(apiKey),
	})
}

func Get(profile string) (string, error) {
	ring, err := open()
	if err != nil {
		return "", err
	}
	item, err := ring.Get(profile)
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func Delete(profile string) error {
	ring, err := open()
	if err != nil {
		return err
	}
	return ring.Remove(profile)
}

func open() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: "dashcli",
	})
}
