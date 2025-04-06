package main

import (
	"aleron/src/alerr"
	"aleron/src/files"
	"aleron/src/lang_implementation/tokenizer"
)

func main() {
	file, _ := files.NewFile("test.al")
	l := tokenizer.Lexer{}
	l.Init(file)
	tokens, hadErr := l.Tokenize()
	if hadErr {
		println("Lexer had errors")
		alerr.PrintErrStream(l.Errors)
	}
	tokenizer.PrintTokStream(tokens)
}
