package tg

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func send(b []byte, method string, apiToken string) error {
	client := http.Client{}
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.telegram.org/bot"+apiToken+"/"+method,
		bytes.NewBuffer(b),
	)
	println("https://api.telegram.org/bot" + apiToken + "/" + method)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("telegram response status code %d. error %s", res.StatusCode, string(body))
	}
	return nil
}
