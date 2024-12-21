package streamlabs

type Payload struct {
	Email            string `json:"email"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Agree            bool   `json:"agree"`
	AgreePromotional bool   `json:"agreePromotional"`
	Dob              string `json:"dob"`
	CaptchaToken     string `json:"captcha_token"`
	Locale           string `json:"locale"`
}

type Verify struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
