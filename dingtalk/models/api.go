package models

type TopResult[T any] struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
	Result       T      `json:"result"`
	RequestID    string `json:"request_id"`
}
