package models

type Patient struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	AuthId   string `json:"authid"`
	Authcode string `json:"authcode"`
}
