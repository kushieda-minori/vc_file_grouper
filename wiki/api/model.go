package api

import "encoding/json"

// {"batchcomplete":"","query":{"tokens":{"logintoken":"fcca4b8a65305fb470b5620c8396e29a5fd97a00+\\"}}}
type tokenResponse struct {
	BatchComplete string          `json:"batchcomplete"`
	Query         tokenQuery      `json:"query"`
	Warnings      json.RawMessage `json:"warnings"`
}

type tokenQuery struct {
	Tokens tokenInfo `json:"tokens"`
}

type tokenInfo struct {
	LoginToken string `json:"logintoken"`
	CSRFToken  string `json:"csrftoken"`
}

//{"login":{"result":"Failed","reason":"The supplied credentials could not be authenticated."}}
type loginResponse struct {
	Login result `json:"login"`
}

type result struct {
	Result string `json:"result"`
	Reason string `json:"reason"`
}
