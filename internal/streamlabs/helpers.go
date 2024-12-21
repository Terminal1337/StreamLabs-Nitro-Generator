package streamlabs

import (
	"fmt"
	"net/url"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func CreateClient() (tls_client.HttpClient, error) {

	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Chrome_131),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar),
	}
	Client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return Client, nil

}

func GetXSRFTokenFromJar(Client tls_client.HttpClient) (string, error) {
	// Parse the URL string into a *url.URL
	parsedURL, err := url.Parse("https://streamlabs.com")
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	// Retrieve cookies for the parsed URL
	cookies := Client.GetCookieJar().Cookies(parsedURL)

	// Loop through the cookies to find the X-XSRF-TOKEN
	for _, cookie := range cookies {
		if cookie.Name == "X-XSRF-TOKEN" {
			return cookie.Value, nil
		}
	}

	// If the cookie is not found, return an error
	return "", fmt.Errorf("X-XSRF-TOKEN cookie not found")
}
