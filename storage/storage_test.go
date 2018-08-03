package storage_test

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"

	"github.com/dwburke/lode/storage"
)

func TestStorage(t *testing.T) {
	viper.SetConfigType("yaml")
	var yamlExample = []byte(`
lode:
  port: 4441
  https: true
  search:
  - "{context}.someapp.{customer_id}.{key}"
  - "{context}.someapp.{key}"
  storage:
    type: "memory"
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	st, err := storage.New()

	if err != nil {
		t.Errorf("Error creating storage: %s", err)
	}

	err = st.Set("test.key", []byte("test.value"))
	if err != nil {
		t.Errorf("Error setting key: %s", err)
	}

	value, err := st.Get("test.key")
	if err != nil {
		t.Errorf("Error setting key: %s", err)
	}

	if string(value) != "test.value" {
		t.Errorf("value not as expected, wanted 'test.value', got '%s'", value)
	}

}
