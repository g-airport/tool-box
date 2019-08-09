package entity

const IsValid = 1

type EmailInfo struct {
	Email     string `json:"address"`
	SrcStatus string `json:"-"`
	RetStatus bool   `json:"deliverable"`
	Extra     string `json:"-"`
	Err       error  `json:"-"`
}
