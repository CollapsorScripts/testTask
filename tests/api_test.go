package tests

import (
	"auth/pkg/logger"
	"auth/tests/suite"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"testing"
)

type tokenPair struct {
	AccessToken  string
	RefreshToken string
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "   ")
	if err != nil {
		return in
	}
	return out.String()
}

func toBytesJSON(object any) []byte {
	jsonByte, err := json.Marshal(object)
	if err != nil {
		logger.Error("Ошибка при получении JSON: ", err.Error())
	}
	n := len(jsonByte)             //Find the length of the byte array
	result := string(jsonByte[:n]) //convert to string

	return []byte(jsonPrettyPrint(result))
}

func Test_CreateToken(t *testing.T) {
	_, st := suite.New(t)

	guid := uuid.New().String()

	endpoint := fmt.Sprintf("%s/%s?id=%s", st.BaseURL, "token", guid)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		t.Errorf("Ошибка создания запроса: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := st.Client.Do(req)
	if err != nil {
		t.Errorf("Ошибка выполнения запроса: %v", err)
		return
	}
	defer response.Body.Close()

	token := new(tokenPair)

	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		t.Errorf("Ошибка при анмаршлинге тела ответа: %v", err)
		return
	}

	t.Logf("Статус: %s", response.Status)
	t.Logf("Ответ: %s", toBytesJSON(token))
}

func Test_RefreshToken(t *testing.T) {
	_, st := suite.New(t)

	token := tokenPair{
		AccessToken:  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ4MTg3ODMsImlhdCI6MTc0NDgxNjk4MywiaXAiOiIxMjcuMC4wLjEiLCJzZXNzaW9uIjoiMmIyZWI3ODgtMzNlYS00ZDRiLWFkOTItYjgyY2I5ZDEwMTZkIiwic3ViIjoiYmI5ZjAwNzgtY2U5Ny00OThhLTkzZDYtMGYyMzU5ZjIyMzc5In0.9jWkzb89wiKD0LdfyLbHxpQasbEBauV7XkOx3vBNzhJjN8uC0mCwMt9jMJA7NZtJ6dmvVB3PW5RhQ2Zly-ZYlw",
		RefreshToken: "5BRC6cGiZFRmeBlhMmnnQVcI3sRHhNr0i+6xQV+5lU4=",
	}

	endpoint := fmt.Sprintf("%s/%s", st.BaseURL, "refresh")

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(toBytesJSON(token)))
	if err != nil {
		t.Errorf("Ошибка создания запроса: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := st.Client.Do(req)
	if err != nil {
		t.Errorf("Ошибка выполнения запроса: %v", err)
		return
	}
	defer response.Body.Close()

	newToken := new(tokenPair)

	if err := json.NewDecoder(response.Body).Decode(&newToken); err != nil {
		t.Errorf("Ошибка при анмаршлинге тела ответа: %v", err)
		return
	}

	t.Logf("Статус: %s", response.Status)
	t.Logf("Ответ: %v", newToken)
}
