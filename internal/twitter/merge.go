package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"streamlabsuwu/internal/helpers"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func Merge(auth_token string, Client tls_client.HttpClient, Headers http.Header, CSRF string) error {
	params := url.Values{"r": []string{"/dashboard#/settings/account-settings/platforms"}}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://streamlabs.com/api/v5/user/accounts/merge/twitter_account?%s", params.Encode()), nil)
	if err != nil {
		return err
	}
	req.Header = Headers
	req.Header.Set("Referer", "https://streamlabs.com/dashboard")
	req.Header.Set("X-CSRF-TOKEN", CSRF)

	// Set IPv4 For X Requests

	Client.SetProxy("http://" + helpers.IPv4.Next())
	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to get OAuth URL: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return fmt.Errorf("failed to parse response body: %v", err)
	}

	token, ok := responseData["redirect_url"].(string)
	if !ok || token == "" {
		return fmt.Errorf("no token found in response")
	}

	token = strings.Split(token, `oauth_token=`)[1]
	req, err = http.NewRequest(http.MethodGet, "https://api.x.com/oauth/authenticate?oauth_token="+token, nil)
	if err != nil {
		return err
	}

	req.Header = TWITTER_HEADERS.Clone()
	req.Header.Set("cookie", fmt.Sprintf("auth_token=%s", auth_token))

	resp, err = Client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to Do Oauth Request: %d", resp.StatusCode)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(body), `"authenticity_token" type="hidden" value="`) {
		return fmt.Errorf("authenticity token not found")
	}

	authenticity_token := strings.Split(strings.Split(string(body), `"authenticity_token" type="hidden" value="`)[1], `"`)[0]

	payload := fmt.Sprintf(`authenticity_token=%s&oauth_token=%s`, authenticity_token, token)

	req, err = http.NewRequest(http.MethodPost, "https://x.com/oauth/authorize", strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header = TWITTER_HEADERS.Clone()
	req.Header.Set("cookie", fmt.Sprintf("auth_token=%s", auth_token))
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("referer", "https://api.x.com/")
	req.Header.Set("origin", "https://api.x.com/")

	resp, err = Client.Do(req)
	if err != nil {
		return err
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if !strings.Contains(string(body), `If your browser doesn't redirect you please`) {
		return fmt.Errorf("Could Not Find Callback Url")
	}

	redirect := strings.Split(strings.Split(string(body), `<a class="maintain-context" href="`)[1], `">`)[0]
	redirect = strings.Replace(redirect, "amp;", "", -1)

	req, err = http.NewRequest(http.MethodGet, redirect, nil)
	if err != nil {
		return err
	}
	req.Header = Headers
	req.Header.Set("Referer", "https://x.com/")
	Client.SetFollowRedirect(false)

	resp, err = Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 302 {
		return nil
	}

	return fmt.Errorf("Failed to link Twitter account: %d", resp.StatusCode)
}
