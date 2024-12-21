package emails

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/ItsMerkz/Kopeechka_Wrapper/kopeechka"
)

var (
	client   *kopeechka.EmailClient
	domains_ = []string{"outlook.com"} //, "outlook.com", "yahoo.com", "hotmail.com", "mail.com", "email.com"
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

func GetCodeKopeechka(emailKey, mailID string) string {
	for i := 0; i < 20; i++ {
		MailID, err := strconv.Atoi(mailID) // Changes MailId from string to int
		if err != nil {
			log.Println(err.Error())
		}
		value := client.GetLetter(emailKey, MailID)
		if !strings.Contains(value, `WAIT`) {
			return value

		}
		time.Sleep(2 * time.Second) // Adjust the retry interval as needed

	}
	DeleteEmail(emailKey, mailID)
	return ""
}

func DeleteEmail(key, mailid string) {
	MailID, err := strconv.Atoi(mailid) // Changes MailId from string to int
	if err != nil {
		log.Fatal(err)
	}
	response, _ := client.DeleteMail(key, MailID)
	if response == "OK" {
		log.Println("Cancelled Email > ", mailid)
	}

}
