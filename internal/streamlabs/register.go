package streamlabs

import (
	"net/url"
	"streamlabs/internal/emails"
	"streamlabs/internal/helpers"
	"streamlabs/pkg/logging"
)

func Create() {
	Client, err := CreateClient()
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error creating client")
		return
	}

	if err := GetCookies(Client); err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error getting cookies")
		return
	}

	EmailID, Email := emails.GetEmailKopeechka("47e48acd01f988fdd2627acf3f3ae3da")
	if EmailID == "" || Email == "" {
		logging.Logger.Error().Msg("No email received from Kopeechka")
		return
	}

	logging.Logger.Info().Str("Email", Email).Str("ID", EmailID).Msg("Received Email")

	Password := helpers.GeneratePassword(12)
	resp, err := Signup(Client, Email, Password)
	if err != nil || resp.StatusCode != 200 {
		if resp != nil {
			return
		}
		logging.Logger.Error().Str("msg", "Error signing up or unexpected status code").Int("status_code", resp.StatusCode).Msg("Signup Request")
		return
	}
	defer resp.Body.Close()

	Code := emails.GetCodeKopeechka("47e48acd01f988fdd2627acf3f3ae3da", EmailID)
	if Code == "" {
		logging.Logger.Error().Msg("No verification code received from Kopeechka")
		return
	}

	logging.Logger.Info().Str("Code", Code).Msg("Received Code")

	resp, err = EmailVerify(Client, Email, Code)
	if err != nil || resp.StatusCode == 204 {
		if resp != nil {
			return
		}
		logging.Logger.Error().Str("msg", "Error verifying email or unexpected status code").Int("status_code", resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	logging.Logger.Info().Str("email", Email).Str("password", Password).Msg("Account successfully created")

	t, _ := url.Parse("https://streamlabs.com")
	err = helpers.WriteCookiesToFile(Email, Password, Client.GetCookieJar().Cookies(t))
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Saving File")
	}
}
