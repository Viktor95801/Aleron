package alerr

import (
	"fmt"

	"github.com/fatih/color"
)

type AlParsErr struct {
	ErrType ErrType
	Msg     string

	File   string
	Line   int
	Column int
}

func (e AlParsErr) Wrap(msg string) AlParsErr {
	return AlParsErr{
		File:   e.File,
		Line:   e.Line,
		Column: e.Column,

		ErrType: e.ErrType,
		Msg:     e.Msg + " (internal: " + msg + ")",
	}
}

func (e AlParsErr) Fmt() string {
	return color.New(color.FgRed).Sprintf("[ERROR] ") + fmt.Sprintf("at %s:%d:%d during parsing", e.File, e.Line, e.Column) + ": " + e.Msg
}

func (e AlParsErr) Print() {
	fmt.Println(e.Fmt())
}

func PrintErrStream(errors []AlParsErr) {
	for _, err := range errors {
		err.Print()
	}
}

func (e AlParsErr) Ok() bool {
	return e.ErrType == ErrOK
}

func AlParsErrOK() AlParsErr {
	return AlParsErr{
		ErrType: ErrOK,
	}
}
