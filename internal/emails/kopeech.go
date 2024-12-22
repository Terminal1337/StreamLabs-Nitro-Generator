package emails

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ItsMerkz/Kopeechka_Wrapper/kopeechka"
)

var (
	client   *kopeechka.EmailClient
	domains_ = []string{"mail.com"} //, "outlook.com", "yahoo.com", "hotmail.com", "mail.com", "email.com"
)

func init() {
	client = &kopeechka.EmailClient{} // Assuming NewEmailClient() is a constructor function
	rand.Seed(time.Now().UnixNano())  // Seed the random number generator
}

func GetRandomDomain() string {
	// Choose a random index from the domains slice
	randomIndex := rand.Intn(len(domains_))
	// Return the domain at the random index
	return domains_[randomIndex]
}

func GetEmailKopeechka(key string) (string, string) {
	for i := 0; i < 5; i++ {
		domain := GetRandomDomain() // Get a random domain
		MailID, Email := client.BuyEmail(key, domain, "streamlabs.com")
		if MailID != "" && Email != "" {
			return MailID, Email
		}
		// If MailID or Email is empty, retry after some time
		time.Sleep(2 * time.Second) // Adjust the retry interval as needed
	}
	return "", ""
}

// func GetCodeKopeechka(emailKey, mailID string) string {
// 	for i := 0; i < 30; i++ {
// 		MailID, err := strconv.Atoi(mailID) // Changes MailId from string to int
// 		if err != nil {
// 			log.Println(err.Error())
// 		}
// 		value := client.GetLetter(emailKey, MailID)
// 		fmt.Println(value)
// 		if !strings.Contains(value, `WAIT`) {
// 			return value

// 		}
// 		time.Sleep(2 * time.Second) // Adjust the retry interval as needed

//		}
//		DeleteEmail(emailKey, mailID)
//		return ""
//	}
func GetCodeKopeechka(emailKey, mailID string) string {
	// Construct the API URL
	url := fmt.Sprintf("https://api.kopeechka.store/mailbox-get-message?full=1&id=%s&token=%s&type=JSON&api=2.0", mailID, emailKey)

	// Make the API request
	for i := 0; i < 30; i++ {
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Error fetching data:", err)
			return ""
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body:", err)
			return ""
		}

		// Parse the JSON response
		var kopeechkaResp KopeechkaResponse
		err = json.Unmarshal(body, &kopeechkaResp)
		if err != nil {
			log.Println("Error unmarshalling response:", err)
			return ""
		}

		// Check if the status is OK and there's a message
		if kopeechkaResp.Status == "OK" && kopeechkaResp.FullMessage != "" {
			// fmt.Println(kopeechkaResp.FullMessage)
			code := extractCodeFromHTML(kopeechkaResp.FullMessage)
			if code != "" {
				return code
			}
		}
		// Retry if the response doesn't contain the code
		time.Sleep(2 * time.Second)
	}

	return ""
}

func extractCodeFromHTML(html string) string {
	var code string
	if strings.Contains(html, `<div style="margin-bottom: 24.0px;background-color: rgb(227,232,235);padding: 16.0px;text-align: center;font-size: 30.0px;line-height: 1.5;color: rgb(0,0,0);">`) {
		code = strings.Split(html, ` <div style="margin-bottom: 24.0px;background-color: rgb(227,232,235);padding: 16.0px;text-align: center;font-size: 30.0px;line-height: 1.5;color: rgb(0,0,0);">`)[1]
	}
	code = strings.TrimSpace(strings.Split(code, `<`)[0])
	return code
}

func DeleteEmail(key, mailid string) {
	MailID, err := strconv.Atoi(mailid)
	if err != nil {
		log.Fatal(err)
	}
	response, _ := client.DeleteMail(key, MailID)
	if response == "OK" {
		log.Println("Cancelled Email > ", mailid)
	}

}
