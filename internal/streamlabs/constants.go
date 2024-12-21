package streamlabs

import (
	http "github.com/bogdanfinn/fhttp"
)

var (
	BASE_HEADERS = http.Header{
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Encoding":           {"gzip, deflate, br, zstd"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Connection":                {"keep-alive"},
		"Host":                      {"streamlabs.com"},
		"Sec-Ch-Ua":                 {`"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`},
		"Sec-Ch-Ua-Mobile":          {"?0"},
		"Sec-Ch-Ua-Platform":        {`"Windows"`},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"},
		http.HeaderOrderKey: {
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Connection",
			"Host",
			"Sec-Ch-Ua",
			"Sec-Ch-Ua-Mobile",
			"Sec-Ch-Ua-Platform",
			"Sec-Fetch-Dest",
			"Sec-Fetch-Mode",
			"Sec-Fetch-Site",
			"Sec-Fetch-User",
			"Upgrade-Insecure-Requests",
			"User-Agent",
		}}

	REGISTER_HEADERS = http.Header{
		"Accept":             {"application/json, text/plain, */*"},
		"Accept-Encoding":    {"gzip, deflate, br, zstd"},
		"Accept-Language":    {"en-US,en;q=0.9"},
		"Client-ID":          {"419049641753968640"},
		"Connection":         {"keep-alive"},
		"Content-Length":     {"1017"},
		"Content-Type":       {"application/json"},
		"Cookie":             {`__cf_bm=VYFzPXzz9ke0uVz1RINP886BSmyG00lf2Gp9cYv9xB0-1734815975-1.0.1.1-XklHSC5WfMDCZe8d9rVE3jeg9StmT5wfCRkMQD.GQHohHXAyp75Qjm36O7PJ5B_IpT6HHMZMu94DAY2h.1aihQ; XSRF-TOKEN=eyJpdiI6IkhDMG9VOGQ1YkY4dE0xK2l5YWFrZ2c9PSIsInZhbHVlIjoiWmFqTUwraGVocys5bGgwTDZhcE9IMUtVYTNvQzdjdWdqMWJjK3l5YklrV21jRWNvcm15UUlKWXQ4ZXBJS3V4TyIsIm1hYyI6Ijk5Mzc4NDU2NmU4NDJkYThmMWI5MTk5MzI3NTk3YTllM2E4NGY1ZGU5YTYzMTBmZGUyYTVjOTk3NGQ0OTAwYWQiLCJ0YWciOiIifQ%3D%3D; slsid=eyJpdiI6IjBMemlONUdaNUVINW1OdkYwVXQvSEE9PSIsInZhbHVl...`},
		"Host":               {"api-id.streamlabs.com"},
		"Origin":             {"https://streamlabs.com"},
		"Referer":            {"https://streamlabs.com/"},
		"Sec-Ch-Ua":          {`"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {`"Windows"`},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"same-site"},
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"},
		"X-XSRF-TOKEN":       {`eyJpdiI6IkhDMG9VOGQ1YkY4dE0xK2l5YWFrZ2c9PSIsInZhbHVlIjoiWmFqTUwraGVocys5bGgwTDZhcE9IMUtVYTNvQzdjdWdqMWJjK3l5YklrV21jRWNvcm15UUlKWXQ4ZXBJS3V4TyIsIm1hYyI6Ijk5Mzc4NDU2NmU4NDJkYThmMWI5MTk5MzI3NTk3YTllM2E4NGY1ZGU5YTYzMTBmZGUyYTVjOTk3NGQ0OTAwYWQiLCJ0YWciOiIifQ==`},
		http.HeaderOrderKey: {"Accept", "Accept-Encoding", "Accept-Language", "Client-ID", "Connection", "Content-Length",
			"Content-Type", "Cookie", "Host", "Origin", "Referer", "Sec-Ch-Ua", "Sec-Ch-Ua-Mobile", "Sec-Ch-Ua-Platform",
			"Sec-Fetch-Dest", "Sec-Fetch-Mode", "Sec-Fetch-Site", "User-Agent", "X-XSRF-TOKEN"},
	}
)
