package streamlabs

import (
	"fmt"
	"net/url"
	"strings"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func CreateClient() (tls_client.HttpClient, error) {

	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(60),
		tls_client.WithClientProfile(profiles.Chrome_131),
		tls_client.WithCookieJar(jar),
		tls_client.WithProxyUrl("http://" + Proxies.Next()),
	}
	Client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return nil, err
	}

	return Client, nil

}

func GetXSRFTokenFromJar(Client tls_client.HttpClient) (string, error) {
	parsedURL, err := url.Parse("https://streamlabs.com")
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %v", err)
	}

	cookies := Client.GetCookieJar().Cookies(parsedURL)

	for _, cookie := range cookies {
		if cookie.Name == "XSRF-TOKEN" {
			// Replace '%3D' with '=' manually
			decodedValue := strings.Replace(cookie.Value, "%3D", "=", -1)
			return decodedValue, nil
		}
	}

	return "", fmt.Errorf("XSRF-TOKEN cookie not found")
}
