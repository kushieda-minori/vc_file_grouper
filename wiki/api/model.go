package api

import "encoding/json"

// {"batchcomplete":"","query":{"tokens":{"logintoken":"fcca4b8a65305fb470b5620c8396e29a5fd97a00+\\"}}}
type tokenResponse struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Tokens struct {
			LoginToken string `json:"logintoken"`
			CSRFToken  string `json:"csrftoken"`
		} `json:"tokens"`
	} `json:"query"`
	Warnings json.RawMessage `json:"warnings"`
}

//{"login":{"result":"Failed","reason":"The supplied credentials could not be authenticated."}}
type loginResponse struct {
	Login struct {
		Result string `json:"result"`
		Reason string `json:"reason"`
	} `json:"login"`
}

type editResponse struct {
	Edit struct {
		Result       string `json:"result"`
		PageID       int    `json:"pageid"`
		Title        string `json:"title"`
		ContentModel string `json:"contentmodel"`
		OldRevID     int    `json:"oldrevid"`
		NewRevID     int    `json:"newrevid"`
		NewTimestamp string `json:"newtimestamp"`
	} `json:"edit"`
	Error json.RawMessage `json:"error"`
}
