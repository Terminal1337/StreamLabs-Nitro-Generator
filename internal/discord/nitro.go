package discord

import (
	"fmt"
	"streamlabsuwu/internal/captcha"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func ClaimNitro(Client tls_client.HttpClient, Headers http.Header) (string, error) {

	req, err := http.NewRequest(http.MethodGet, "https://streamlabs.com/discord/nitro", nil)
	if err != nil {
		return "", err
	}
	clearance, _, ua, err := captcha.CapGetClearance(Client.GetProxy())
	if err != nil || clearance == "" {
		fmt.Println(err.Error())
		return "", err
	}
	req.Header = Headers
	req.Header.Set("Origin", "https://streamlabs.com")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", ua)
	req.Header.Del("Sec-Ch-Ua")
	// req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"129\", \"Not=A?Brand\";v=\"8\"")

	// Get Clearance
	// clearance, err := cloudflare.GetClearance(Client.GetProxy())
	// if err != nil || clearance == "" {
	// 	return "", err
	// }

	req.AddCookie(&http.Cookie{
		Name:  "cf_clearance",
		Value: clearance,
	})
	Client.SetFollowRedirect(false)

	resp, err := Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	// defer resp.Body.Close()
	fmt.Println("Nitro : ", resp.StatusCode)
	if resp.StatusCode != 302 {
		return "", fmt.Errorf("Nitro Link Not Found")
	}
	return resp.Header.Get("location"), nil

}
