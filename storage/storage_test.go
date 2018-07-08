package storage_test

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"

	"github.com/dwburke/prefs/storage"
)

func TestStorage(t *testing.T) {
	viper.SetConfigType("yaml")
	var yamlExample = []byte(`
prefs:
  port: 4441
  https: true
  search:
  - "{context}.someapp.{key}"
  - "{context}.someapp.{customer_id}.{key}"
storage:
  type: "mysql"
  dsn: "addict:abc123@/prefs?charset=utf8"
`)

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	st, err := storage.New()

	if err != nil {
		t.Errorf("Error creating storage: %s", err)
	}

	err = st.Set("test.key", "test.value")
	if err != nil {
		t.Errorf("Error setting key: %s", err)
	}

	value, err := st.Get("test.key")
	if err != nil {
		t.Errorf("Error setting key: %s", err)
	}

	if value != "test.value" {
		t.Errorf("value not as expected, wanted 'test.value', got '%s'", value)
	}

}
