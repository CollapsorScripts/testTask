package config

import (
	"bytes"
	"encoding/json"
	"testing"
)

const configPathLocal = "../../config/local.yaml"

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "   ")
	if err != nil {
		return in
	}
	return out.String()
}

// ToJSON - конвертирует объект в JSON строку
func ToJSON(object any) string {
	jsonByte, err := json.Marshal(object)
	if err != nil {
		panic(any(err))
	}
	n := len(jsonByte)
	result := string(jsonByte[:n])

	return jsonPrettyPrint(result)
}

func Test_MustLoadByPath_Happy(t *testing.T) {
	cfg := MustLoadByPath(configPathLocal)

	t.Logf("Конфигурация local: %s", ToJSON(cfg))
}

func Test_MustLoadByPath_Bad(t *testing.T) {
	cfg := MustLoadByPath("local.yaml")

	t.Logf("Конфигурация: %+v", cfg)
}
