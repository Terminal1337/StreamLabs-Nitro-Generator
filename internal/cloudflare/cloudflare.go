package cloudflare

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetClearance(proxy string) (string, error) {

	proxy = strings.Replace(proxy, "http://", "", -1)
	parts := strings.Split(proxy, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid proxy format, expected username:password@host:port")
	}

	credentials := strings.Split(parts[0], ":")
	if len(credentials) != 2 {
		return "", fmt.Errorf("invalid proxy credentials format, expected username:password")
	}

	address := strings.Split(parts[1], ":")
	if len(address) != 2 {
		return "", fmt.Errorf("invalid proxy address format, expected host:port")
	}

	username := credentials[0]
	password := credentials[1]
	host := address[0]

	port, err := strconv.Atoi(address[1])
	if err != nil {
		return "", fmt.Errorf("invalid port: %v", err)
	}

	url := "https://solver.proxiflare.com/cf-clearance-scraper"
	method := "POST"

	payload := fmt.Sprintf(`{
		"url": "https://streamlabs.com/discord/nitro",
		"mode": "waf-session",
		"proxy": {
			"host": "%s",
			"port": %d,
			"username": "%s",
			"password": "%s"
		}
	}`, host, port, username, password)

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %v", err)
	}

	for _, cookie := range response.Cookies {
		if cookie.Name == "cf_clearance" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("cf_clearance cookie not found")
}
