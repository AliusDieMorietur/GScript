package main

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	u "github.com/core/utils"
)

const (
	// Single-character tokens
	LeftBracket      = "leftBracket"
	RightBracket     = "rightBracket"
	LeftBrace        = "leftBrace"
	RightBrace       = "rightBrace"
	LeftCurlyBracket = "leftCurlyBracket"
	RightCurlyBrace  = "rightCurlyBrace"
	Comma            = "comma"
	Dot              = "dot"
	Minus            = "minus"
	Plus             = "plus"
	Semicolon        = "semicolon"
	Slash            = "slash"
	Star             = "star"
	Colon            = "colon"
	Question         = "question"

	// One or two character tokens
	Bang         = "bang"
	BangEqual    = "bangEqual"
	Equal        = "equal"
	EqualEqual   = "equalEqual"
	Greater      = "greater"
	GreaterEqual = "greaterEqual"
	Less         = "less"
	LessEqual    = "lessEqual"
	Or           = "or"
	And          = "and"

	// Literals
	Identifier = "identifier"
	String     = "string"
	Number     = "number"

	// Keywords

	Struct = "struct"
	Else   = "else"
	True   = "true"
	False  = "false"
	Fn     = "fn"
	For    = "for"
	If     = "if"
	Null   = "null"
	Print  = "print"
	Return = "return"
	Super  = "super"
	This   = "this"
	Let    = "let"
	While  = "while"
	Eof    = "eof"
)

var keywords = map[string]string{
	"struct": Struct,
	"else":   Else,
	"true":   True,
	"false":  False,
	"for":    For,
	"fn":     Fn,
	"if":     If,
	"null":   Null,
	"print":  Print,
	"return": Return,
	"super":  Super,
	"this":   This,
	"let":    Let,
	"while":  While,
}

type Token struct {
	tokenType string
	lexeme    string
	literal   any
	line      uint
}

func NewToken(tokenType string, lexeme string, literal any, line uint) Token {
	return Token{
		tokenType,
		lexeme,
		literal,
		line,
	}
}

func (t Token) ToString() string {
	return fmt.Sprintf("%s %d", t.tokenType, t.line)
}

type Scanner struct {
	source  string
	start   uint
	current uint
	line    uint
	tokens  []Token
	onError func()
}

func NewScanner(source string, onError func()) Scanner {
	return Scanner{
		source,
		0,
		0,
		1,
		[]Token{},
		onError,
	}
}

func (s Scanner) isAtEnd() bool {
	return s.current >= uint(len(s.source))
}

func (s *Scanner) addToken(tokenType string, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) advance() byte {
	s.current += 1
	return s.source[s.current-1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return byte(0)
	}
	return s.source[s.current]
}

func (s Scanner) peekNext() byte {
	if s.current+1 >= uint(len(s.source)) {
		return byte(0)
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		u.Error(s.line, "Unterminated string.")
		s.onError()
		return
	}
	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(String, value)
}

func (s *Scanner) isDigit(c byte) bool {
	return unicode.IsDigit(rune(c))
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	value := s.source[s.start:s.current]
	float, _ := strconv.ParseFloat(value, 64)
	s.addToken(Number, float)
}

func (s Scanner) isAlpha(char byte) bool {
	re := regexp.MustCompile("[A-Za-z_]")
	return re.MatchString(string(char))
}

func (s Scanner) isAlphaNumeric(char byte) bool {
	return s.isAlpha(char) || s.isDigit(char)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tokenType, exists := keywords[text]
	if !exists {
		tokenType = Identifier
	}
	s.addToken(tokenType, "")
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '&':
		if s.match('&') {
			s.addToken(And, "")
		} else {
			u.Error(s.line, "Unterminted &")
			s.onError()
		}
	case '|':
		if s.match('|') {
			s.addToken(Or, "")
		} else {
			u.Error(s.line, "Unterminted |")
			s.onError()
		}
	case '(':
		s.addToken(LeftBrace, "")
	case ')':
		s.addToken(RightBrace, "")
	case '{':
		s.addToken(LeftCurlyBracket, "")
	case '}':
		s.addToken(RightCurlyBrace, "")
	case '[':
		s.addToken(LeftBracket, "")
	case ']':
		s.addToken(RightBracket, "")
	case ',':
		s.addToken(Comma, "")
	case '.':
		s.addToken(Dot, "")
	case '-':
		s.addToken(Minus, "")
	case '+':
		s.addToken(Plus, "")
	case ';':
		s.addToken(Semicolon, "")
	case '*':
		s.addToken(Star, "")
	case '?':
		s.addToken(Question, "")
	case ':':
		s.addToken(Colon, "")
	case '!':
		s.addToken(u.Ternary(s.match('='), BangEqual, Bang), "")
	case '=':
		s.addToken(u.Ternary(s.match('='), EqualEqual, Equal), "")
	case '<':
		s.addToken(u.Ternary(s.match('='), LessEqual, Less), "")
	case '>':
		s.addToken(u.Ternary(s.match('='), GreaterEqual, Greater), "")
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(Slash, "")
		}
	case ' ':
	case '\r':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			u.Error(s.line, fmt.Sprintf("Unexpected character \"%s\"", string(c)))
			s.onError()
		}
	}
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(Eof, "", "", s.line))
	return s.tokens
}
