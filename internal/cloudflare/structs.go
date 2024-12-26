package cloudflare

type Response struct {
	Cookies []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"cookies"`
}
