package emails

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Success bool   `json:"success"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Code    string `json:"code"`
	Error   string `json:"error,omitempty"`
}

func GetStreamlabsCode(email string) (string, error) {
	url := fmt.Sprintf("http://mail.elevatehosting.net:8000/api/streamlabs.com?email=%s&type=1", email)

	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("failed to make GET request: %v", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("failed to read response body: %v", err)
		}

		var response Response
		err = json.Unmarshal(body, &response)
		if err != nil {
			return "", fmt.Errorf("failed to parse JSON response: %v", err)
		}

		if response.Success {
			return response.Code, nil
		}

		time.Sleep(2 * time.Second)
	}

	return "", fmt.Errorf("failed to retrieve code for email %s after 10 attempts", email)
}

func GetCustomEmail() string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	usernameLength := 8 + rand.Intn(8)
	var username strings.Builder
	for i := 0; i < usernameLength; i++ {
		username.WriteByte(letters[rand.Intn(len(letters))])
	}

	email := fmt.Sprintf("%s@elevatehosting.net", username.String())
	return email
}
