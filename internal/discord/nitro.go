package discord

import (
	"fmt"
	"streamlabsuwu/internal/cloudflare"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func ClaimNitro(Client tls_client.HttpClient, Headers http.Header) (string, error) {

	req, err := http.NewRequest(http.MethodGet, "https://streamlabs.com/discord/nitro", nil)
	if err != nil {
		return "", err
	}

	req.Header = Headers
	req.Header.Set("Origin", "https://streamlabs.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"129\", \"Not=A?Brand\";v=\"8\"")

	// Get Clearance
	clearance, err := cloudflare.GetClearance(Client.GetProxy())
	if err != nil || clearance == "" {
		return "", err
	}

	req.AddCookie(&http.Cookie{
		Name:  "cf_clearance",
		Value: clearance,
	})
	Client.SetFollowRedirect(false)

	resp, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 302 {
		return "", fmt.Errorf("Nitro Link Not Found")
	}
	return resp.Header.Get("location"), nil

}
