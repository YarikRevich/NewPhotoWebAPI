package models

const (
	AUTH_ERROR = iota
	NOT_THIS_TIME_ERROR
)

type ERRORAuthModel struct {
	Service struct {
		Error int `json:"error"`
	} `json:"service"`
}
