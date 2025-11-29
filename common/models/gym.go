package models

type ErrResponse struct {
	Ok               bool        `json:"ok"`
	Code             string      `json:"code"`
	Message          string      `json:"msg"`
	Type             string      `json:"type"`
	ValidationErrors interface{} `json:"validation_errors,omitempty"`
	Details          interface{} `json:"details,omitempty"`
}

type SuccessResponse struct {
	Ok      bool   `json:"ok"`
	Code    string `json:"code"`
	Message string `json:"msg"`
}

type ValidationError struct {
	FailedField string
	Tag         string
	Value       string
}
