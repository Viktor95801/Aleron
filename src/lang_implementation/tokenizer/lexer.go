package tokenizer

import (
	"aleron/src/alerr"
	"aleron/src/files"
	"fmt"
	"strings"
)

type Lexer struct {
	File   files.File
	source []byte
	tokens []Tok

	pos    int
	Line   int
	Column int

	Errors   []alerr.AlParsErr
	hadError bool
}

func (l *Lexer) Init(file files.File) {
	l.File = file
	l.source = append(file.Content, ' ')
	l.tokens = make([]Tok, 0)

	l.pos = 0
	l.Line = 1
	l.Column = 1

	l.Errors = make([]alerr.AlParsErr, 0)
	l.hadError = false
}

func (l *Lexer) newErr(typ alerr.ErrType, msg string) alerr.AlParsErr {
	return alerr.AlParsErr{ErrType: typ, Msg: msg, File: l.File.Name, Line: l.Line, Column: l.Column}
}

func (l *Lexer) newTok(tokType TokType, lit string) Tok {
	return Tok{Type: tokType, Lit: lit, Line: l.Line, Column: l.Column, Pos: l.pos}
}

func (l *Lexer) updatePos(c byte) {
	if c == '\n' {
		l.Line++
		l.Column = 1
	} else {
		l.Column++
	}

	l.pos++
}
func (l *Lexer) deUpdatePos() {
	l.pos--
	c := l.source[l.pos]

	if c == '\n' {
		l.Line--
	} else {
		l.Column--
	}
}

func (l *Lexer) peek() byte {
	return l.source[l.pos+1]
}

func (l *Lexer) advance() byte {
	l.updatePos(l.source[l.pos])
	return l.source[l.pos]
}

func (l *Lexer) handleComment(c byte, multiLine bool) Tok {
	buf := strings.Builder{}
	tok := Tok{}

	// TODO: single line comments for some reason cannot be sequentially used
	if !multiLine {
		for c != '\n' && c != 0 {
			buf.WriteByte(c)
			l.updatePos(c)
			c = l.source[l.pos]
		}
		tok = l.newTok(COMMENT, buf.String())
	} else {
		for c != '*' || l.peek() != '/' {
			buf.WriteByte(c)
			l.updatePos(c)
			if l.pos >= len(l.source) {
				l.hadError = true
				l.Errors = append(l.Errors, l.newErr(alerr.ErrLEX_UnterminatedComment,
					"Comment was never terminated").Wrap("Unexpected EOF"))
			}
			c = l.source[l.pos]
		}
		buf.WriteString("*/")
		tok = l.newTok(COMMENT, buf.String())
		l.updatePos(c)
	}
	return tok
}

func (l *Lexer) handleIdent(c byte) Tok {
	buf := strings.Builder{}
	for isAlphaNum(c) {
		buf.WriteByte(c)
		l.updatePos(c)
		c = l.source[l.pos]
	}
	l.deUpdatePos()

	keyword, ok := Keywords[buf.String()]
	if ok {
		return l.newTok(keyword, buf.String())
	}

	primitiveType, ok := Types[buf.String()]
	if ok {
		return l.newTok(primitiveType, buf.String())
	}

	return l.newTok(IDENT, buf.String())
}

func (l *Lexer) handleString(c byte) Tok {
	if c == '\'' {
		c = l.advance()
		if l.peek() != '\'' {
			l.hadError = true
			l.Errors = append(l.Errors, l.newErr(alerr.ErrLEX_InvalidByteLiteral,
				"Invalid byte literal"))
		}

		return l.newTok(UINT8_T, fmt.Sprintf("%d", c))
	} else {
		buf := strings.Builder{}
		c = l.advance()
		for c != '"' {
			buf.WriteByte(c)
			l.updatePos(c)
			if l.pos >= len(l.source) {
				l.hadError = true
				l.Errors = append(l.Errors, l.newErr(alerr.ErrLEX_UnterminatedString,
					"String was never terminated").Wrap("Unexpected EOF"))
				break
			}
			c = l.source[l.pos]
		}

		return l.newTok(STRING, buf.String())
	}
}

func (l *Lexer) handleNumber(c byte) Tok {
	buf := strings.Builder{}
	isFloat := false

	if c == '.' {
		if !isDigit(l.peek()) {
			return l.newTok(DOT, ".")
		}
		isFloat = true
		buf.Write([]byte{'0', '.'})
		l.updatePos(c)
		c = l.source[l.pos]
	}

	for isDigit(c) || c == '.' {
		if c == '.' {
			if isFloat {
				l.hadError = true
				l.Errors = append(l.Errors, l.newErr(alerr.ErrLEX_InvalidNumber,
					"Invalid number literal"))
				break
			}
			isFloat = true
		}
		buf.WriteByte(c)
		l.updatePos(c)
		c = l.source[l.pos]
	}
	l.deUpdatePos()

	if isFloat {
		return l.newTok(NUMBER, buf.String())
	}

	return l.newTok(INTEGER, buf.String())
}

func (l *Lexer) nextTok() Tok {
	c := l.source[l.pos]
	tok := l.newTok(INVALID, string(c))

	for isSpace(c) {
		l.updatePos(c)
		if l.pos >= len(l.source) {
			return l.newTok(EOF, "")
		}
		c = l.source[l.pos]
	}

	//======================================================================
	// Literals/Keywords/Identifiers and Comments
	//======================================================================
	if c == '/' && (l.peek() == '*' || l.peek() == '/') {
		if l.peek() == '*' {
			tok = l.handleComment(c, true)
		} else if l.peek() == '/' {
			tok = l.handleComment(c, false)
		}
	} else if c == '"' || c == '\'' {
		tok = l.handleString(c)
	} else if isAlpha(c) {
		tok = l.handleIdent(c)
	} else if isDigit(c) || c == '.' {
		tok = l.handleNumber(c)
	} else {
		//======================================================================
		// Operators and Delimiters
		//======================================================================
		switch c {
		case '=':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(EQUAL, "==")
			default:
				tok = l.newTok(ASSIGN, "=")
			}
		case '!':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(NOT_EQUAL, "!=")
			default:
				tok = l.newTok(NOT, "!")
			}
		case '<':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(LESS_EQUAL, "<=")
			default:
				tok = l.newTok(LESS, "<")
			}
		case '>':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(GREATER_EQUAL, ">=")
			default:
				tok = l.newTok(GREATER, ">")
			}

		case '+':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(PLUS_ASSIGN, "+=")
			default:
				tok = l.newTok(PLUS, "+")
			}
		case '-':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(MINUS_ASSIGN, "-=")
			default:
				tok = l.newTok(MINUS, "-")
			}
		case '*':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(STAR_ASSIGN, "*=")
			default:
				tok = l.newTok(STAR, "*")
			}
		case '/':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(SLASH_ASSIGN, "/=")
			default:
				tok = l.newTok(SLASH, "/")
			}
		case '%':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(MOD_ASSIGN, "%=")
			default:
				tok = l.newTok(MOD, "%")
			}

		case '&':
			switch l.peek() {
			case '&':
				l.updatePos(c)
				tok = l.newTok(AND, "&&")
			case '=':
				l.updatePos(c)
				tok = l.newTok(AMP_ASSIGN, "&=")
			default:
				tok = l.newTok(AMPERSAND, "&")
			}
		case '|':
			switch l.peek() {
			case '|':
				l.updatePos(c)
				tok = l.newTok(OR, "||")
			case '=':
				l.updatePos(c)
				tok = l.newTok(PIPE_ASSIGN, "|=")
			default:
				tok = l.newTok(PIPE, "|")
			}
		case '^':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(CARET_ASSIGN, "^=")
			case '^':
				l.updatePos(c)
				tok = l.newTok(XOR, "^^")
			default:
				tok = l.newTok(CARET, "^")
			}
		case '~':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(TILDE_ASSIGN, "~=")
			default:
				tok = l.newTok(TILDE, "~")
			}

		case '(':
			tok = l.newTok(LPAREN, "(")
		case ')':
			tok = l.newTok(RPAREN, ")")
		case '{':
			tok = l.newTok(LBRACE, "{")
		case '}':
			tok = l.newTok(RBRACE, "}")
		case '[':
			tok = l.newTok(LBRACE, "[")
		case ']':
			tok = l.newTok(RBRACE, "]")
		case ',':
			tok = l.newTok(COMMA, ",")
		case ';':
			tok = l.newTok(SEMICOLON, ";")
		case ':':
			switch l.peek() {
			case '=':
				l.updatePos(c)
				tok = l.newTok(INFER_ASSIGN, ":=")
			default:
				tok = l.newTok(COLON, ":")
			}
		}
	}

	l.updatePos(c)

	return tok
}

func (l *Lexer) Tokenize() ([]Tok, bool) {
	for l.pos < len(l.source) {
		tok := l.nextTok()
		l.tokens = append(l.tokens, tok)
	}

	return l.tokens, l.hadError
}
