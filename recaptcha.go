package captcha

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type recaptcha struct {
	url       string
	secretKey string
	threshold float32
}

type result struct {
	Success    bool     `json:"success"`
	Timestamp  string   `json:"challenge_ts"`
	Hostname   string   `json:"hostname"`
	ErrorCodes []string `json:"error-codes"`
	Score      float32  `json:"score"`
}

func (s *recaptcha) Verify(response string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.post(ctx, response)
}

func (s *recaptcha) VerifyWithContext(ctx context.Context, response string) error {
	return s.post(ctx, response)
}

func (s *recaptcha) post(ctx context.Context, response string) error {
	data := url.Values{"secret": {s.secretKey}, "response": {response}}

	req, err := http.NewRequest("POST", s.url, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	var r result
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	if !r.Success {
		for _, code := range r.ErrorCodes {
			err = errors.Join(err, errors.New(code))
		}
		return err
	}

	if r.Score < s.threshold {
		return fmt.Errorf("score below threshold: %.1f", r.Score)
	}

	return nil
}

func Recaptcha(secretKey string, threshold float32, urls ...string) Service {
	url := "https://www.google.com/recaptcha/api/siteverify"
	if len(urls) != 0 {
		url = urls[0]
	}

	return &recaptcha{
		url:       url,
		secretKey: secretKey,
		threshold: threshold,
	}
}
