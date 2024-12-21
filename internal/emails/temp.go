package emails

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var (
	emailDomain string
	waitTime    int
	domains     []string
)

func Log(tag, message string, writeToLogFile bool) {
	logMessage := fmt.Sprintf("[%s] %s", tag, message)
	fmt.Println(logMessage)
}

func RequestEmailT() string {
	rand.Seed(time.Now().UnixNano())
	for {
		json_data := map[string]int{
			"min_name_length": 2,
			"max_name_length": 5,
		}
		jsonDataBytes, err := json.Marshal(json_data)
		if err != nil {
			Log("EMAIL", fmt.Sprintf("Error marshalling JSON: %v", err), true)
			return ""
		}

		// currentDomain := domains[rand.Intn(len(domains))]
		response, err := http.Post("https://api.internal.temp-mail.io/api/v3/email/new", "application/json", bytes.NewReader(jsonDataBytes))
		if err != nil {
			Log("EMAIL", fmt.Sprintf("Unknown error occurred while getting email: %v", err), true)
			return ""
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusOK {
			var data map[string]interface{}
			if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
				Log("EMAIL", fmt.Sprintf("Error decoding JSON: %v", err), true)
				return ""
			}
			email, ok := data["email"].(string)
			if !ok {
				Log("EMAIL", "Email not found in response", true)
				return ""
			}
			return email
		} else {
			Log("EMAIL", fmt.Sprintf("Unknown error occurred while getting email: %v", response.Status), true)
			return ""
		}
	}
}

func GetCodeT(mail string) string {
	tries := 0
	for tries < 10 {
		tries++
		response, err := http.Get(fmt.Sprintf("https://api.internal.temp-mail.io/api/v3/email/%s/messages", mail))
		if err != nil {
			Log("EMAIL", fmt.Sprintf("Error getting messages: %v", err), true)
			time.Sleep(3 * time.Second)
			continue
		}
		defer response.Body.Close()

		var messages []map[string]interface{}
		if err := json.NewDecoder(response.Body).Decode(&messages); err != nil {
			Log("EMAIL", fmt.Sprintf("Error decoding JSON: %v", err), true)
			time.Sleep(3 * time.Second)
			continue
		}

		if len(messages) > 0 {
			subject, ok := messages[0]["subject"].(string)
			if !ok {
				Log("EMAIL", "Error parsing subject", true)
				time.Sleep(3 * time.Second)
				continue
			}
			parts := strings.Split(subject, " ")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			} else {
				Log("EMAIL", "Confirmation code not found in subject", true)
			}
		}

		time.Sleep(3 * time.Second)
	}

	Log("EMAIL", "No code received", true)
	return ""
}
