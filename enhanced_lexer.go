package parser

import "strings"

// An enhanced lexer that supports negative numbers
type EnhancedLexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

// NewEnhancedLexer creates a new enhanced lexer that supports negative numbers
func NewEnhancedLexer(input string) *EnhancedLexer {
	l := &EnhancedLexer{input: input}
	l.readChar()
	return l
}

func (l *EnhancedLexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *EnhancedLexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// peekTwoChars returns the character two positions ahead
func peekTwoChars(l *EnhancedLexer) byte {
	if l.readPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition+1]
}

func (l *EnhancedLexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(EQ, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: NE, Literal: literal}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: LE, Literal: literal}
		} else {
			tok = newToken(LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = Token{Type: GE, Literal: literal}
		} else {
			tok = newToken(GT, l.ch)
		}
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		// Check if this comma is between digits, which would make it part of a number
		if isDigit(l.peekChar()) && l.position > 0 && isDigit(l.input[l.position-1]) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			// Include the comma in the literal
			tok.Literal = "," + tok.Literal
			return tok
		} else {
			tok = newToken(COMMA, l.ch)
		}
	case '\'':
		tok.Type = STRING
		tok.Literal = l.readString()
		// Check for unclosed string marker
		if strings.HasPrefix(tok.Literal, "_UNCLOSED_STRING_") {
			tok.Type = ILLEGAL
			tok.Literal = "unclosed string: " + strings.TrimPrefix(tok.Literal, "_UNCLOSED_STRING_")
		}
		return tok
	case '-':
		// Check if it's a negative number
		if isDigit(l.peekChar()) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case 0: // EOF
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			typeFromIdent := LookupIdentifier(tok.Literal)
			tok.Type = typeFromIdent
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *EnhancedLexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *EnhancedLexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '.' || l.ch == '_' || l.ch == '-' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *EnhancedLexer) readNumber() string {
	position := l.position
	// Handle negative numbers
	if l.ch == '-' {
		l.readChar()
	}

	// Read integer part, allowing commas for readability (e.g., 1,000,000)
	hasDigits := false
	for isDigit(l.ch) || l.ch == ',' {
		if isDigit(l.ch) {
			hasDigits = true
		}
		l.readChar()
	}

	if !hasDigits {
		// Ensure we have at least one digit
		return l.input[position:l.position]
	}

	// Read decimal part if present
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // Read '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	// Read scientific notation if present (e.g., 1.23e45, 1e10)
	if (l.ch == 'e' || l.ch == 'E') && (isDigit(l.peekChar()) ||
		((l.peekChar() == '+' || l.peekChar() == '-') && isDigit(peekTwoChars(l)))) {

		l.readChar() // Read 'e' or 'E'

		// Read sign if present
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}

		// Read exponent
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.position]
}

func (l *EnhancedLexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
		// Handle escape sequences
		if l.ch == '\\' && l.peekChar() == '\'' {
			l.readChar() // Skip the backslash and include the quote
		}
	}
	literal := l.input[position:l.position]
	// Check if the string was unclosed (ended with EOF)
	isUnclosed := l.ch == 0
	if l.ch != 0 { // Only advance if we're not at EOF
		l.readChar()
	}
	if isUnclosed {
		// Return a special marker that can be checked by the token consumer
		return "_UNCLOSED_STRING_" + literal
	}
	return literal
}
