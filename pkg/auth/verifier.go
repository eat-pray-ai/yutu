// Copyright 2026 eat-pray-ai & OpenWaygate
// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	sdkauth "github.com/modelcontextprotocol/go-sdk/auth"
)

const defaultTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo"

// NewGoogleTokenVerifier returns a TokenVerifier that validates Google OAuth
// access tokens by calling the given tokeninfo endpoint. Use
// GoogleTokenVerifier for production (calls Google's endpoint directly).
func NewGoogleTokenVerifier(tokenInfoURL string) sdkauth.TokenVerifier {
	return func(ctx context.Context, token string, req *http.Request) (*sdkauth.TokenInfo, error) {
		url := fmt.Sprintf("%s?access_token=%s", tokenInfoURL, token)
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", sdkauth.ErrInvalidToken, err)
		}

		resp, err := http.DefaultClient.Do(httpReq)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", sdkauth.ErrInvalidToken, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf(
				"%w: tokeninfo returned %d", sdkauth.ErrInvalidToken, resp.StatusCode,
			)
		}

		var result struct {
			ExpiresIn int64  `json:"expires_in,string"`
			Scope     string `json:"scope"`
			Sub       string `json:"sub"`
			Error     string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("%w: %w", sdkauth.ErrInvalidToken, err)
		}
		if result.Error != "" {
			return nil, fmt.Errorf("%w: %s", sdkauth.ErrInvalidToken, result.Error)
		}

		return &sdkauth.TokenInfo{
			Scopes:     strings.Fields(result.Scope),
			Expiration: time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
			UserID:     result.Sub,
			Extra:      map[string]any{"access_token": token},
		}, nil
	}
}

// GoogleTokenVerifier is a production TokenVerifier that validates tokens
// against Google's tokeninfo endpoint.
var GoogleTokenVerifier = NewGoogleTokenVerifier(defaultTokenInfoURL)
