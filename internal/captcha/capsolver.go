package captcha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type capSolverResponse struct {
	ErrorId          int32          `json:"errorId"`
	ErrorCode        string         `json:"errorCode"`
	ErrorDescription string         `json:"errorDescription"`
	TaskId           string         `json:"taskId"`
	Status           string         `json:"status"`
	Solution         map[string]any `json:"solution"`
}

var (
	apikey = "CAP-F19DE1CD3C46ADED7C93C246552170DD"
)

func capSolver(ctx context.Context, apiKey string, taskData map[string]any) (*capSolverResponse, error) {
	uri := "https://api.capsolver.com/createTask"
	res, err := request(ctx, uri, map[string]any{
		"clientKey": apiKey,
		"task":      taskData,
	})
	if err != nil {
		return nil, err
	}
	if res.ErrorId == 1 {
		return nil, errors.New(res.ErrorDescription)
	}

	uri = "https://api.capsolver.com/getTaskResult"
	for {
		select {
		case <-ctx.Done():
			return res, errors.New("solve timeout")
		case <-time.After(time.Second):
			break
		}
		res, err = request(ctx, uri, map[string]any{
			"clientKey": apiKey,
			"taskId":    res.TaskId,
		})
		if err != nil {
			return nil, err
		}
		if res.ErrorId == 1 {
			return nil, errors.New(res.ErrorDescription)
		}
		if res.Status == "ready" {
			return res, err
		}
	}
}

func request(ctx context.Context, uri string, payload interface{}) (*capSolverResponse, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", uri, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	capResponse := &capSolverResponse{}
	err = json.Unmarshal(responseData, capResponse)
	if err != nil {
		return nil, err
	}
	return capResponse, nil
}

func CapSolve() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	res, err := capSolver(ctx, apikey, map[string]any{
		"type":       "AntiTurnstileTaskProxyLess",
		"websiteURL": "https://streamlabs.com",
		"websiteKey": "0x4AAAAAAACELUBpqiwktdQ9",
	})
	if err != nil {
		return "", err
	}
	return res.Solution["token"].(string), nil
}

func CapGetClearance(proxy string) (string, string, string, error) {
	fmt.Println(proxy)
	proxy = strings.Replace(proxy, "http://", "", -1)
	fmt.Println(proxy)

	parts := strings.Split(proxy, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid proxy format")
	}
	auth := strings.Split(parts[0], ":")
	if len(auth) != 2 {
		return "", "", "", fmt.Errorf("invalid authentication format")
	}
	ipPort := parts[1]
	proxyFormatted := fmt.Sprintf("%s:%s:%s", ipPort, auth[0], auth[1])

	maxRetries := 3
	var res *capSolverResponse
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
		defer cancel()

		res, err = capSolver(ctx, apikey, map[string]any{
			"type":       "AntiCloudflareTask",
			"websiteURL": "https://streamlabs.com/discord/nitro",
			"proxy":      proxyFormatted,
		})
		if err == nil {
			break
		}
		if attempt < maxRetries {
			fmt.Printf("Attempt %d/%d failed. Retrying...\n", attempt, maxRetries)
		}
	}

	if err != nil {
		return "", "", "", fmt.Errorf("failed after %d attempts: %w", maxRetries, err)
	}
	fmt.Println(res)
	token := res.Solution["token"].(string)
	secChUa := res.Solution["headers"].(map[string]any)["sec-ch-ua"].(string)
	userAgent := res.Solution["headers"].(map[string]any)["User-Agent"].(string)

	return token, secChUa, userAgent, nil
}
