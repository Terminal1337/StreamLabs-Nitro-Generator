package streamlabs

import (
	"fmt"
	"streamlabsuwu/pkg/logging"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func GetCookies(Client tls_client.HttpClient) error {
	req, err := http.NewRequest(http.MethodGet, "https://streamlabs.com/signup", nil)
	req.Header = BASE_HEADERS.Clone()
	if err != nil {
		logging.Logger.Error().
			Str("msg", err.Error()).
			Msg("First Cookie Request")
		return err
	}
	resp, err := Client.Do(req)
	if err != nil {
		logging.Logger.Error().
			Str("msg", err.Error()).
			Msg("First Cookie Response")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Status Code MisMatch : ", resp.StatusCode)
	}

	req2, err := http.NewRequest(http.MethodGet, "https://streamlabs.com/api/v5/user/basic-information", nil)
	req2.Header = req.Header
	if err != nil {
		logging.Logger.Error().
			Str("msg", err.Error()).
			Msg("Second Cookie Request")
		return err
	}

	resp2, err := Client.Do(req2)
	if err != nil {
		logging.Logger.Error().
			Str("msg", err.Error()).
			Msg("First Cookie Response")
		return err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 401 {
		return fmt.Errorf("Status Code MisMatch : ", resp.StatusCode)
	}
	return nil
}
