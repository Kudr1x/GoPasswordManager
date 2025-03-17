package models

type Password struct {
	Service     string `json:"service"`
	Login       string `json:"login"`
	Version     string `json:"version"`
	Description string `json:"description"`
}
