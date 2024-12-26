package streamlabs

import (
	"fmt"
	"net/url"
	"os"
	"streamlabs/internal/emails"
	"streamlabs/internal/helpers"
	"streamlabs/internal/twitter"
	"streamlabs/pkg/logging"

	"github.com/Terminal1337/GoCycle"
)

var (
	Proxies   *GoCycle.Cycle
	TwTs      *GoCycle.Cycle
	kopeechka = false
)

func init() {
	var err error

	Proxies, err = GoCycle.NewFromFile("data/input/proxies.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	TwTs, err = GoCycle.NewFromFile("data/input/twitter.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

func Create() {

	var (
		Email   string
		EmailID string
		Code    string
	)
	Client, err := CreateClient()
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error creating client")
		return
	}

	if err := GetCookies(Client); err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error getting cookies")
		return
	}
	if kopeechka {
		EmailID, Email := emails.GetEmailKopeechka("47e48acd01f988fdd2627acf3f3ae3da")
		if EmailID == "" || Email == "" {
			logging.Logger.Error().Msg("No email received from Kopeechka")
			return
		}
	} else {
		Email = emails.GetCustomEmail()
	}

	logging.Logger.Info().Str("Email", Email).Str("ID", EmailID).Msg("Received Email")

	Password := helpers.GeneratePassword(12)
	resp, err := Signup(Client, Email, Password)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		if resp != nil {
			return
		}
		logging.Logger.Error().Str("msg", "Error signing up or unexpected status code").Int("status_code", resp.StatusCode).Msg("Signup Request")
		return
	}
	defer resp.Body.Close()
	if kopeechka {
		Code = emails.GetCodeKopeechka("47e48acd01f988fdd2627acf3f3ae3da", EmailID)
		if Code == "" {
			logging.Logger.Error().Msg("No verification code received from Kopeechka")
			return
		}
	} else {
		Code, err = emails.GetStreamlabsCode(Email)
		if err != nil {
			logging.Logger.Error().Msg("No verification code received from Kopeechka")
		}
	}

	logging.Logger.Info().Str("Code", Code).Msg("Received Code")

	resp, err = EmailVerify(Client, Email, Code)
	if err != nil || resp.StatusCode != 204 {
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

	csrf, err := GetOauth(Client, REGISTER_HEADERS.Clone())
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Msg("GetOauth")
	}

	twt := TwTs.Next()
	err = twitter.Merge(twt, Client, REGISTER_HEADERS.Clone(), csrf)
	if err != nil {
		logging.Logger.Error().Str("msg", err.Error()).Str("twt", twt).Msg("Failed to Link Twitter")

	} else {
		logging.Logger.Info().Str("twt", twt).Str("email", Email).Msg("Linked Twitter")

	}
}
