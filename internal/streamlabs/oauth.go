package streamlabs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func GetOauth(client tls_client.HttpClient, Headers http.Header) (string, error) {
	payload := `{
		"origin": "https://streamlabs.com",
		"intent": "connect",
		"state": ""
	}`

	req, err := http.NewRequest(http.MethodPost, "https://api-id.streamlabs.com/v1/identity/clients/419049641753968640/oauth2", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header = Headers.Clone()

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(b, &responseData)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %v", err)
	}

	redirectURL, ok := responseData["redirect_url"].(string)
	if !ok || redirectURL == "" {
		return "", fmt.Errorf("no redirect_url found in response")
	}

	req, err = http.NewRequest(http.MethodGet, redirectURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GET request for redirect: %v", err)
	}
	req.Header = Headers.Clone()

	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send GET request to redirect URL: %v", err)
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from GET request: %v", err)
	}

	if !strings.Contains(string(b), `var redirectUrl = '`) {
		return "", fmt.Errorf("redirectUrl Var Not Present")
	}

	redirectURL = strings.Split(strings.Split(string(b), `var redirectUrl = '`)[1], `';`)[0]
	req, err = http.NewRequest(http.MethodGet, redirectURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create GET request for redirect: %v", err)
	}
	req.Header = Headers.Clone()
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send GET request to redirect URL: %v", err)
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from GET request: %v", err)
	}

	var csrf string
	if !strings.Contains(string(b), `"csrf-token" content="`) {
		return "", fmt.Errorf("CSRF token not found")
	}
	csrf = strings.Split(strings.Split(string(b), `"csrf-token" content="`)[1], `"`)[0]

	return csrf, nil
}
