package tg

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func send(b []byte, method string, apiToken string) error {
	body, err := sendGetBody(b, method, apiToken)
	if err != nil {
		return err
	}
	body.Close()
	return nil
}

func sendGetBody(b []byte, method string, apiToken string) (io.ReadCloser, error) {
	client := http.Client{}
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.telegram.org/bot"+apiToken+"/"+method,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		return nil, fmt.Errorf("telegram response status code %d. error %s", res.StatusCode, string(body))
	}
	return res.Body, nil
}
