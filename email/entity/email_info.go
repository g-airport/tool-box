package entity

const IsValid = 1

type EmailInfo struct {
	Email     string
	SrcStatus string
	RetStatus bool
	Extra     string
	Err       error
}