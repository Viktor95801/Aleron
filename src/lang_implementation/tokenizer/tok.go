package tokenizer

import "fmt"

type TokType int

type Tok struct {
	// data
	Type   TokType
	Lit    string
	Line   int
	Column int

	// debug
	Pos int
}

// ======================================================================
// Tok Types
// ======================================================================
const (
	//======================================================================
	// Especial
	//======================================================================
	INVALID = iota
	EOF
	COMMENT

	//======================================================================
	// Literals
	//======================================================================
	IDENT
	STRING
	INTEGER
	NUMBER

	//======================================================================
	// Operators
	//======================================================================
	EQUAL
	NOT_EQUAL
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL

	PLUS
	MINUS
	STAR
	SLASH
	MOD

	//======================================================================
	// Logical Operators
	//======================================================================
	AND
	OR
	NOT
	XOR

	//======================================================================
	// Bitwise Operators
	//======================================================================
	AMPERSAND
	PIPE
	CARET
	TILDE

	//======================================================================
	// Assignment Operators
	//======================================================================
	ASSIGN
	INFER_ASSIGN

	PLUS_ASSIGN
	MINUS_ASSIGN
	STAR_ASSIGN
	SLASH_ASSIGN
	MOD_ASSIGN

	AMP_ASSIGN
	PIPE_ASSIGN
	CARET_ASSIGN
	TILDE_ASSIGN

	//======================================================================
	// Delimiters
	//======================================================================
	COMMA
	SEMICOLON
	COLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LCURLY
	RCURLY

	DOT

	//======================================================================
	// Keywords
	//======================================================================

	// Modules
	PACKAGE
	IMPORT
	EXPORT
	AS

	// Declarations
	LET
	CONST
	VAR

	FUNC
	RETURN // technically not a declaration, but it fits

	PROC

	// Statements
	IF
	MATCH
	CASE
	ELSE
	DEFER

	FOR
	BREAK
	CONTINUE

	WITH

	// Making types
	TYPE
	STRUCT
	ENUM
	UNION

	// Primitive Types
	INT8_T
	UINT8_T

	INT16_T
	UINT16_T

	INT32_T
	UINT32_T

	INT64_T
	UINT64_T

	INT_T // size of a pointer (64 bit on 64 bit archs, 32 bit on 32 bit archs)
	UINT_T

	FLOAT32_T
	FLOAT64_T

	STRING_T

	BOOLEAN_T
	// not a type, but it fits
	TRUE
	FALSE

	FUNCTION_T
	PROCEDURE_T
)

// ======================================================================
// Keyword lookup
// ======================================================================
var Keywords = map[string]TokType{
	"package": PACKAGE,
	"import":  IMPORT,
	"export":  EXPORT,
	"as":      AS,

	"let":   LET,
	"const": CONST,
	"var":   VAR,

	"func":   FUNC,
	"return": RETURN,

	"proc": PROC,

	"if":    IF,
	"match": MATCH,
	"case":  CASE,
	"else":  ELSE,
	"defer": DEFER,

	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,

	"with": WITH,

	"type":   TYPE,
	"struct": STRUCT,
	"enum":   ENUM,
	"union":  UNION,

	"true":  TRUE,
	"false": FALSE,
}

//======================================================================
// Type lookups
//======================================================================
var Types = map[string]TokType{
	"int8":  INT8_T,
	"uint8": UINT8_T,

	"int16":  INT16_T,
	"uint16": UINT16_T,

	"int32":  INT32_T,
	"uint32": UINT32_T,

	"int64":  INT64_T,
	"uint64": UINT64_T,

	"int":  INT_T,
	"uint": UINT_T,

	"float32": FLOAT32_T,
	"float64": FLOAT64_T,

	"string": STRING_T,

	"bool": BOOLEAN_T,

	// Aliases
	"byte": UINT8_T,
	"str":  STRING_T,

	// Shorthand
	"i8":  INT8_T,
	"u8":  UINT8_T,
	"i16": INT16_T,
	"u16": UINT16_T,
	"i32": INT32_T,
	"u32": UINT32_T,
	"i64": INT64_T,
	"u64": UINT64_T,
	"f32": FLOAT32_T,
	"f64": FLOAT64_T,

	"function":  FUNCTION_T,
	"procedure": PROCEDURE_T,
}

//======================================================================
// Implementations
//======================================================================

func (t Tok) String() string {
	var kind string
	if t.OneOf(IDENT, STRING, INTEGER, NUMBER, COMMENT) {
		switch t.Type {
		case IDENT:
			kind = "IDENT"
		case STRING:
			kind = "STRING"
		case INTEGER:
			kind = "INTEGER"
		case NUMBER:
			kind = "NUMBER"
		case COMMENT:
			kind = "COMMENT"
		default:
			panic("WTF?????????? HOW DID YOU GET HERE? HACKER!!!!!!!!!!!!!!!!!!! BURN THE HACKER!!!!")
		}
		return fmt.Sprintf("%s: %s", kind, t.Lit)
	}

	switch t.Type {
	case INVALID:
		return "INVALID"
	case EOF:
		return "EOF"
	case INFER_ASSIGN:
		return "INFER_ASSIGN"

	case EQUAL:
		return "EQUAL"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"

	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case MOD:
		return "MOD"

	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case XOR:
		return "XOR"

	case AMPERSAND:
		return "AMPERSAND"
	case PIPE:
		return "PIPE"
	case CARET:
		return "CARET"
	case TILDE:
		return "TILDE"

	case ASSIGN:
		return "ASSIGN"

	case PLUS_ASSIGN:
		return "PLUS_ASSIGN"
	case MINUS_ASSIGN:
		return "MINUS_ASSIGN"
	case STAR_ASSIGN:
		return "STAR_ASSIGN"
	case SLASH_ASSIGN:
		return "SLASH_ASSIGN"
	case MOD_ASSIGN:
		return "MOD_ASSIGN"

	case AMP_ASSIGN:
		return "AMP_ASSIGN"
	case PIPE_ASSIGN:
		return "PIPE_ASSIGN"
	case CARET_ASSIGN:
		return "CARET_ASSIGN"
	case TILDE_ASSIGN:
		return "TILDE_ASSIGN"

	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case COLON:
		return "COLON"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LCURLY:
		return "LCURLY"
	case RCURLY:
		return "RCURLY"

	case DOT:
		return "DOT"

	case PACKAGE:
		return "PACKAGE"
	case IMPORT:
		return "IMPORT"
	case EXPORT:
		return "EXPORT"
	case AS:
		return "AS"

	case LET:
		return "LET"
	case CONST:
		return "CONST"
	case VAR:
		return "VAR"

	case FUNC:
		return "FUNC"
	case RETURN:
		return "RETURN"

	case PROC:
		return "PROC"

	case IF:
		return "IF"
	case MATCH:
		return "MATCH"
	case CASE:
		return "CASE"
	case ELSE:
		return "ELSE"

	case FOR:
		return "FOR"
	case BREAK:
		return "BREAK"
	case CONTINUE:
		return "CONTINUE"

	case WITH:
		return "WITH"

	case TYPE:
		return "TYPE"
	case STRUCT:
		return "STRUCT"
	case ENUM:
		return "ENUM"
	case UNION:
		return "UNION"

	case INT8_T:
		return "INT8_T"
	case UINT8_T:
		return "UINT8_T"

	case INT16_T:
		return "INT16_T"
	case UINT16_T:
		return "UINT16_T"

	case INT32_T:
		return "INT32_T"
	case UINT32_T:
		return "UINT32_T"

	case INT64_T:
		return "INT64_T"
	case UINT64_T:
		return "UINT64_T"

	case INT_T:
		return "INT_T"
	case UINT_T:
		return "UINT_T"

	case FLOAT32_T:
		return "FLOAT32_T"
	case FLOAT64_T:
		return "FLOAT64_T"

	case STRING_T:
		return "STRING_T"

	case BOOLEAN_T:
		return "BOOLEAN_T"

	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"

	case FUNCTION_T:
		return "FUNCTION_T"
	case PROCEDURE_T:
		return "PROCEDURE_T"

	default:
		return "UNKNOWN"
	}
}

func (t Tok) OneOf(types ...TokType) bool {
	for _, typ := range types {
		if t.Type == typ {
			return true
		}
	}
	return false
}

func PrintTokStream(stream []Tok) {
	for _, tok := range stream {
		fmt.Printf("%d:%d %s\n", tok.Line, tok.Column, tok.String())
	}
}
