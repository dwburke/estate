package key_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/dwburke/prefs/api/key"
)

var TestKeys = []struct {
	Successful bool
	Test       string
	Expect     string
}{
	{true, "{context}.{key}", "dev.foo"},
	{true, "{context}.someapp.{customer_id}.{key}", "dev.someapp.12345678.foo"},
	{false, "{context}.region.{region}.{key}", "dev.someapp.12345678.foo"},
}

func TestTranslateKey(t *testing.T) {

	p := gin.Params{
		{Key: "context", Value: "dev"},
		{Key: "key", Value: "foo"},
		{Key: "customer_id", Value: "12345678"},
	}

	for _, test_data := range TestKeys {
		result, err := key.TranslateKey(test_data.Test, &p)

		if test_data.Successful == true {
			if err != nil {
				t.Errorf("Error translating string: %s", err)
			}

			if result != test_data.Expect {
				t.Errorf("expected 'dev.foo', got '%s'", result)
			}
		} else {
			if err == nil {
				t.Errorf("Expected an error translating the string, did not get one")
			}
		}
	}
}

func TestKey(t *testing.T) {
	viper.SetConfigType("yaml")
	var yamlExample = []byte(`
prefs:
  port: 4441
  https: true
  search:
  - "{context}.someapp.{customer_id}.{key}"
  - "{context}.someapp.{key}"
  storage:
    type: "memory"
`)
	//dsn: "prefs:abc123@/prefs?charset=utf8"

	viper.ReadConfig(bytes.NewBuffer(yamlExample))

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/prefs/:context/:key", key.GetKey)
	r.POST("/prefs/:context/:key", key.SetKey)

	params := map[string]interface{}{
		"key":         "test.foo",
		"value":       "test.bar",
		"customer_id": "123456",
	}
	jsonstr, err := json.Marshal(params)
	if err != nil {
		t.Errorf("Error marshalling params to json: %s", err)
	}
	req, _ := http.NewRequest("POST", "/prefs/dev/foo", bytes.NewReader([]byte(jsonstr)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response code should be 200, was: %d", w.Code)
	}
	bodyAsString := w.Body.String()

	fmt.Println(bodyAsString)

	type SetData struct {
		Ok int `json:"ok"`
	}

	var set_data SetData
	err = json.Unmarshal([]byte(bodyAsString), &set_data)

	if err != nil {
		t.Errorf("Error unmarshalling json: %s", err)
	}

	if set_data.Ok != 1 {
		t.Errorf("expected %d got %d", 1, set_data.Ok)
	}

	req, _ = http.NewRequest("GET", "/prefs/dev/foo", nil)

	q := req.URL.Query()
	q.Add("customer_id", "123456")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Response code should be 200, was: %d; %s", w.Code, w.Body)
	}

	bodyAsString = w.Body.String()

	type GetData struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	var get_data GetData
	err = json.Unmarshal([]byte(bodyAsString), &get_data)

	if err != nil {
		t.Errorf("Error unmarshalling json: %s", err)
	}

	if get_data.Key != "dev.someapp.foo" || get_data.Value != "bar" {
		t.Errorf("expected 'dev.someapp.foo' = bar, got %s = %s", get_data.Key, get_data.Value)
	}
}
