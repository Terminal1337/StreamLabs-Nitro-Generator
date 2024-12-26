package twitter

import http "github.com/bogdanfinn/fhttp"

var (
	TWITTER_HEADERS = http.Header{
		"Accept":                    {"application/json, text/plain, */*"},
		"Accept-Encoding":           {"gzip, deflate, br, zstd"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Connection":                {"keep-alive"},
		"Host":                      {"api.x.com"},
		"Referer":                   {"https://api.twitter.com/"},
		"Sec-Ch-Ua":                 {`"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`},
		"sec-ch-ua-mobile":          {"?0"},
		"Sec-Ch-Ua-Platform":        {`"Windows"`},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"cross-site"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"},
	}
)
