package twitter

import (
	"fmt"
	"io/ioutil"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func GetCSRF(Client tls_client.HttpClient, headers http.Header) (string, error) {

	req, err := http.NewRequest(http.MethodGet, "https://streamlabs.com/dashboard", nil)
	if err != nil {
		return "", err
	}
	req.Header = headers
	fmt.Println(req.Header)
	resp, err := Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))
	fmt.Println(resp.StatusCode)

	return "", nil
}
