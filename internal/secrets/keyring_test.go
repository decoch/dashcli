package secrets

import (
	"errors"
	"testing"

	"github.com/99designs/keyring"
)

// mockKeyring is an in-memory keyring used in tests.
type mockKeyring struct {
	data map[string][]byte
}

func newMockKeyring() *mockKeyring {
	return &mockKeyring{data: make(map[string][]byte)}
}

func (m *mockKeyring) Get(key string) (keyring.Item, error) {
	data, ok := m.data[key]
	if !ok {
		return keyring.Item{}, keyring.ErrKeyNotFound
	}
	return keyring.Item{Key: key, Data: data}, nil
}

func (m *mockKeyring) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

func (m *mockKeyring) Set(item keyring.Item) error {
	m.data[item.Key] = item.Data
	return nil
}

func (m *mockKeyring) Remove(key string) error {
	if _, ok := m.data[key]; !ok {
		return keyring.ErrKeyNotFound
	}
	delete(m.data, key)
	return nil
}

func (m *mockKeyring) Keys() ([]string, error) {
	keys := make([]string, 0, len(m.data))
	for k := range m.data {
		keys = append(keys, k)
	}
	return keys, nil
}

func stubOpenKeyring(t *testing.T) *mockKeyring {
	t.Helper()
	mock := newMockKeyring()
	prev := openKeyring
	openKeyring = func() (keyring.Keyring, error) { return mock, nil }
	t.Cleanup(func() { openKeyring = prev })
	return mock
}

func TestSetAndGetAPIKey(t *testing.T) {
	stubOpenKeyring(t)

	if err := SetAPIKey("test-api-key"); err != nil {
		t.Fatalf("SetAPIKey() error = %v", err)
	}
	got, err := GetAPIKey()
	if err != nil {
		t.Fatalf("GetAPIKey() error = %v", err)
	}
	if got != "test-api-key" {
		t.Fatalf("GetAPIKey() = %q, want %q", got, "test-api-key")
	}
}

func TestSetAndGetBaseURL(t *testing.T) {
	stubOpenKeyring(t)

	if err := SetBaseURL("https://redash.example.com"); err != nil {
		t.Fatalf("SetBaseURL() error = %v", err)
	}
	got, err := GetBaseURL()
	if err != nil {
		t.Fatalf("GetBaseURL() error = %v", err)
	}
	if got != "https://redash.example.com" {
		t.Fatalf("GetBaseURL() = %q, want %q", got, "https://redash.example.com")
	}
}

func TestDeleteAPIKey(t *testing.T) {
	stubOpenKeyring(t)

	if err := SetAPIKey("key"); err != nil {
		t.Fatalf("SetAPIKey() error = %v", err)
	}
	if err := DeleteAPIKey(); err != nil {
		t.Fatalf("DeleteAPIKey() error = %v", err)
	}
	_, err := GetAPIKey()
	if !errors.Is(err, keyring.ErrKeyNotFound) {
		t.Fatalf("GetAPIKey() after delete error = %v, want ErrKeyNotFound", err)
	}
}

func TestDeleteBaseURL(t *testing.T) {
	stubOpenKeyring(t)

	if err := SetBaseURL("https://redash.example.com"); err != nil {
		t.Fatalf("SetBaseURL() error = %v", err)
	}
	if err := DeleteBaseURL(); err != nil {
		t.Fatalf("DeleteBaseURL() error = %v", err)
	}
	_, err := GetBaseURL()
	if !errors.Is(err, keyring.ErrKeyNotFound) {
		t.Fatalf("GetBaseURL() after delete error = %v, want ErrKeyNotFound", err)
	}
}

func TestGetAPIKey_NotFound(t *testing.T) {
	stubOpenKeyring(t)

	_, err := GetAPIKey()
	if !errors.Is(err, keyring.ErrKeyNotFound) {
		t.Fatalf("GetAPIKey() error = %v, want ErrKeyNotFound", err)
	}
}

func TestOpenKeyring_Error(t *testing.T) {
	t.Parallel()

	prev := openKeyring
	openKeyring = func() (keyring.Keyring, error) { return nil, errors.New("backend unavailable") }
	t.Cleanup(func() { openKeyring = prev })

	if err := SetAPIKey("key"); err == nil {
		t.Fatal("SetAPIKey() error = nil, want error when keyring unavailable")
	}
	if _, err := GetAPIKey(); err == nil {
		t.Fatal("GetAPIKey() error = nil, want error when keyring unavailable")
	}
	if err := DeleteAPIKey(); err == nil {
		t.Fatal("DeleteAPIKey() error = nil, want error when keyring unavailable")
	}
}
