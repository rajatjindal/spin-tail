package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type retryTransport struct {
	envName string
	token   *Token
	base    http.RoundTripper
}

func (r *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token.Token))
	res, err := r.base.RoundTrip(req)
	if err != nil {
		return res, err
	}

	// only retry if unauthorized
	if res.StatusCode != http.StatusUnauthorized {
		return res, err
	}

	// refresh token and try once more
	newToken, err := r.refreshToken()
	if err != nil {
		return res, err
	}
	r.token = newToken

	err = saveNewToken(r.envName, r.token)
	if err != nil {
		fmt.Println("WARN: failed to save updated token file %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.token.Token))
	return r.base.RoundTrip(req)
}

func (r *retryTransport) refreshToken() (*Token, error) {
	refreshRequest := map[string]string{
		"token":        r.token.Token,
		"refreshToken": r.token.RefreshToken,
	}

	raw, err := json.Marshal(refreshRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/auth-tokens/refresh", CloudBase), bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	resp, err := r.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	token := &Token{}
	err = json.Unmarshal(raw, token)
	if err != nil {
		return nil, err
	}

	token.URL = CloudBase
	return token, err
}
