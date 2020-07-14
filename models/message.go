package models

// Response ... our response structure
type Response struct {
	Success      bool   `json:"status"`
	ResponseCode int    `json:"responseCode"`
	Message      string `json:"message"`
}

// TokenResponse ... our response struct with token string
type TokenResponse struct {
	Success      bool   `json:"status"`
	ResponseCode int    `json:"responseCode"`
	Message      string `json:"message"`
	Token        string `json:"token"`
}
