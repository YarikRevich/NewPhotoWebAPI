package models

const (
	AUTH_ERROR = iota
)

type ERRORAuthModel struct {
	Service struct {
		Error int `json:"error"`
	} `json:"service"`
}