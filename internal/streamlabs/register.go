package streamlabs

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"streamlabsuwu/internal/discord"
	"streamlabsuwu/internal/emails"
	"streamlabsuwu/internal/helpers"
	"streamlabsuwu/internal/twitter"
	"streamlabsuwu/pkg/logging"
	"syscall"
	"unsafe"

	"github.com/Terminal1337/GoCycle"
	http "github.com/bogdanfinn/fhttp"
)

var (
	Proxies   *GoCycle.Cycle // IPV6
	TwTs      *GoCycle.Cycle
	kopeechka = false
	created   int
	errors    int
	nitros    int
	connected int
)

func init() {
	var err error

	Proxies, err = GoCycle.NewFromFile("data/input/proxies_ipv6.txt") // v6 proxy
	if err != nil || Proxies.ListLength() == 0 {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Failed to load proxies")
		os.Exit(1)
	}

	TwTs, err = GoCycle.NewFromFile("data/input/twitter.txt")
	if err != nil || TwTs.ListLength() == 0 {
		logging.Logger.Error().Str("msg", err.Error()).Msg("Failed to load Twitter accounts")
		os.Exit(1)
	}
	go SetConsoleTitle()
}

func Create() {
	defer func() {
		if r := recover(); r != nil {
			logging.Logger.Error().Interface("recover", r).Msg("Recovered from panic in Create")
		}
	}()

	var (
		Email   string
		EmailID string
		Code    string
	)

	Client, err := CreateClient()
	if err != nil {
		errors += 1
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error creating client")
		return
	}

	if err := GetCookies(Client); err != nil {
		errors += 1

		logging.Logger.Error().Str("msg", err.Error()).Msg("Error getting cookies")
		return
	}

	if kopeechka {
		EmailID, Email = emails.GetEmailKopeechka("47e48acd01f988fdd2627acf3f3ae3da")
		if EmailID == "" || Email == "" {
			errors += 1

			logging.Logger.Error().Msg("No email received from Kopeechka")
			return
		}
	} else {
		Email = emails.GetCustomEmail()
		EmailID = "0"
		if Email == "" {
			errors += 1

			logging.Logger.Error().Msg("Failed to get a custom email")
			return
		}
	}

	logging.Logger.Info().Str("Email", Email).Str("ID", EmailID).Msg("Received Email")

	Password := helpers.GeneratePassword(12)
	resp, err := Signup(Client, Email, Password)
	if err != nil || resp == nil {
		errors += 1

		logging.Logger.Error().Str("msg", "Error signing up or unexpected status code").Int("status_code", getStatusCode(resp)).Msg("Signup Request")
		return
	}
	defer safeClose(resp.Body)

	if kopeechka {
		Code = emails.GetCodeKopeechka("47e48acd01f988fdd2627acf3f3ae3da", EmailID)
		if Code == "" {
			errors += 1

			logging.Logger.Error().Msg("No verification code received from Kopeechka")
			return
		}
	} else {
		Code, err = emails.GetStreamlabsCode(Email)
		if err != nil || Code == "" {
			errors += 1

			logging.Logger.Error().Msg("Failed to get verification code from CustomMail")
			return
		}
	}

	logging.Logger.Info().Str("Code", Code).Str("Email", Email).Msg("Received Code")

	resp, err = EmailVerify(Client, Email, Code)
	if err != nil || resp == nil || resp.StatusCode != 204 {
		errors += 1
		logging.Logger.Error().Str("msg", "Error verifying email or unexpected status code").Int("status_code", getStatusCode(resp)).Msg("Email Verification")
		return
	}
	defer safeClose(resp.Body)

	logging.Logger.Info().Str("email", Email).Str("password", Password).Msg("Account successfully created")
	created += 1
	t, _ := url.Parse("https://streamlabs.com")
	err = helpers.WriteCookiesToFile(Email, Password, Client.GetCookieJar().Cookies(t))
	if err != nil {
		errors += 1
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error saving cookies to file")
		return
	}

	csrf, err := GetOauth(Client, REGISTER_HEADERS.Clone())
	if err != nil {
		errors += 1
		logging.Logger.Error().Str("msg", err.Error()).Msg("Error fetching OAuth token")
		return
	}

	twt := TwTs.Next()
	TwTs.Remove(twt)
	if TwTs.ListLength() <= 1 {
		logging.Logger.Error().Str("msg", "Twitter Tokens Out of Stock...").Msg("Please Refill")
		os.Exit(1)
	}

	err = twitter.Merge(twt, Client, REGISTER_HEADERS.Clone(), csrf)
	if err != nil {
		errors += 1
		logging.Logger.Error().Str("msg", err.Error()).Str("twt", twt).Msg("Failed to link Twitter")
		return
	} else {
		errors += 1
		logging.Logger.Info().Str("twt", twt).Str("email", Email).Msg("Linked Twitter successfully")
		connected += 1
	}
	// Set IPv6 Proxy Back
	Client.SetProxy("http://" + Proxies.Next())

	nitro, err := discord.ClaimNitro(Client, REGISTER_HEADERS.Clone())
	if err != nil {
		errors += 1
		logging.Logger.Error().Interface("msg", err.Error()).Msg("Discord Nitro Err")
		return
	}

	logging.Logger.Info().Str("Email", Email).Str("twt", twt).Str("nitro", nitro[:100]).Msg("Nitro Retrieved Successfully!")
	nitros += 1
	helpers.SaveNitro(nitro)

}

func getStatusCode(resp *http.Response) int {
	if resp != nil {
		return resp.StatusCode
	}
	return 0
}

func safeClose(closer io.Closer) {
	if closer != nil {
		_ = closer.Close()
	}
}

func SetConsoleTitle() {
	for {
		titlestr := fmt.Sprintf("StreamLabs Nitro Creator | Created : %d | Linked : %d | Emails : %d | Nitros : %d | Errors : %d", created, connected, nitros, errors)
		kernel32, err := syscall.LoadLibrary("kernel32.dll")
		if err != nil {
			continue
		}
		defer syscall.FreeLibrary(kernel32)

		proc, err := syscall.GetProcAddress(kernel32, "SetConsoleTitleW")
		if err != nil {
			continue
		}

		titlePtr, err := syscall.UTF16PtrFromString(titlestr)
		if err != nil {
			continue
		}

		_, _, err = syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(titlePtr)), 0, 0)
		if err != nil {
			continue
		}
	}
}
