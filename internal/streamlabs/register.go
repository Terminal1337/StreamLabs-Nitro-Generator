package streamlabs

import (
	"fmt"
	"streamlabs/internal/emails"
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

	logging.Logger.Info().Str("Email", Email).Str("ID", EmailID).Msg("Recieved Email")

	resp, err := Signup(Client, Email)
	if err != nil || resp.StatusCode == 200 {
		if resp != nil {
			resp.Body.Close()
		}
		logging.Logger.Error().Str("msg", "Error signing up or unexpected status code").Int("status_code", resp.StatusCode)
		return
	}

	Code := emails.GetCodeKopeechka("47e48acd01f988fdd2627acf3f3ae3da", EmailID)
	if Code == "" {
		logging.Logger.Error().Msg("No verification code received from Kopeechka")
		return
	}
	logging.Logger.Info().Str("Code", Code).Msg("Recieved Code")

	resp, err = EmailVerify(Client, Email, Code)
	if err != nil || resp.StatusCode == 204 {
		if resp != nil {
			resp.Body.Close()
		}
		logging.Logger.Error().Str("msg", "Error verifying email or unexpected status code").Int("status_code", resp.StatusCode)
		return
	}

	fmt.Println("Account created successfully")
	logging.Logger.Info().Msg("Account successfully created")
}
