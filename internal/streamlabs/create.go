package streamlabs

import (
	"bytes"
	"encoding/json"
	"streamlabs/internal/captcha"
	"streamlabs/pkg/logging"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func Signup(Client tls_client.HttpClient, Email string, Password string) (*http.Response, error) {
	Captcha, err := captcha.CapSolve()
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Captcha Solve Failed")
		return nil, err
	}

	logging.Logger.Info().Str("captcha_token", Captcha[:100]).Msg("Captcha Solved")

	Data := Payload{
		Email:            Email,
		Password:         Password,
		CaptchaToken:     Captcha,
		Locale:           "en-US",
		Agree:            true,
		AgreePromotional: false,
	}

	jsonData, err := json.Marshal(Data)
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error marshaling JSON")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api-id.streamlabs.com/v1/auth/register", bytes.NewBuffer(jsonData))
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Request creation failed")
		return nil, err
	}
	req.Header = REGISTER_HEADERS.Clone()

	xsrf, err := GetXSRFTokenFromJar(Client)
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error Getting XSRF")
		return nil, err
	}

	req.Header.Set("XSRF-TOKEN", xsrf)

	resp, err := Client.Do(req)
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Request failed")
		return nil, err
	}

	return resp, nil
}

func EmailVerify(Client tls_client.HttpClient, Email, Code string) (*http.Response, error) {
	Data := Verify{Email: Email, Code: Code}
	jsonData, err := json.Marshal(Data)
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error marshaling JSON")
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://api-id.streamlabs.com/v1/users/@me/email/verification/confirm", bytes.NewBuffer(jsonData))
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Request creation failed")
		return nil, err
	}

	resp, err := Client.Do(req)
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Request failed")
		return nil, err
	}

	return resp, nil
}
