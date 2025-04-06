package alerr

type ErrType int

const (
	ErrOK ErrType = iota

	ErrLEX_InvalidTok
	ErrLEX_InvalidNumber
	ErrLEX_InvalidByteLiteral

	ErrLEX_UnterminatedString
	ErrLEX_UnterminatedComment
)
