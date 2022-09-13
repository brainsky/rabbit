package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int    //输入当前位置
	readPosition int    //读取位置
	ch           byte   //当前字符
	utf8Char     string //用string去存储字符
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}
func NewUTF(input string) *Lexer {
	l := &Lexer{input: input}
	l.readUTF8Char()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readUTF8Char() {
	if l.readPosition >= len(l.input) {
		l.utf8Char = "0"
	} else {
		//判断是否是UTF8字符串
		firstCh := l.input[l.readPosition]

		if firstCh&0x80 == 0x00 {
			l.utf8Char = string(l.input[l.readPosition])
			l.position = l.readPosition
			l.readPosition += 1
		} else {
			//多字节的情况，该字节高位为1，获取高位1的个数，这个只计算1不符合我们的需求
			// var mask1 byte = 0x55
			// var mask2 byte = 0x33
			// var mask3 byte = 0x0f
			// oneLayout := (firstCh & mask1) + ((firstCh >> 1) & mask1)
			// twoLayout := (oneLayout & mask2) + ((oneLayout >> 2) & mask2)
			// charLen := (twoLayout & mask3) + ((twoLayout >> 4) & mask3)
			charLen := caculateUTF8Length(firstCh)
			//计算移动多少位后
			tempLen := l.readPosition + charLen
			if tempLen >= len(l.input) {
				//不正确的字符
				l.utf8Char = "Error"
			} else {
				l.utf8Char = l.input[l.readPosition:tempLen]
				l.position = l.readPosition
				l.readPosition += charLen
			}
		}

	}

}

func (l *Lexer) NextTokenUTF8() token.Token {
	var tok token.Token
	l.skipWhitespaceString()
	switch l.utf8Char {
	case "=":
		tok = newTokenUTF8(token.ASSIGN, l.utf8Char)
	case ";":
		tok = newTokenUTF8(token.SEMICOLON, l.utf8Char)
	case "(":
		tok = newTokenUTF8(token.LPAREN, l.utf8Char)
	case ")":
		tok = newTokenUTF8(token.RPAREN, l.utf8Char)
	case ",":
		tok = newTokenUTF8(token.COMMA, l.utf8Char)
	case "+":
		tok = newTokenUTF8(token.PLUS, l.utf8Char)
	case "{":
		tok = newTokenUTF8(token.LBRACE, l.utf8Char)
	case "}":
		tok = newTokenUTF8(token.RBRACE, l.utf8Char)
	case "0":
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		//先判断是否合法
		if l.utf8Char == "Error" {
			tok = newToken(token.ILEGAL, l.ch)
		} else if isDigitString(l.utf8Char) {
			//判断是否是数字
			tok.Type = token.INT
			tok.Literal = l.readNumeberUTF8()
			return tok
		} else {
			//可能是中文字符，英文字符
			//都需要循环读取
			tok.Literal = l.readUTF8Ident()
			tok.Type = token.LookupUTFIdent(tok.Literal)
			return tok
		}
	}
	l.readUTF8Char()
	return tok

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)

	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILEGAL, l.ch)
		}

	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newTokenUTF8(tokenType token.TokenType, utf8Char string) token.Token {
	return token.Token{Type: tokenType, Literal: utf8Char}
}

/**
 * 引用csdb博主的代码 https://blog.csdn.net/zzqhost/article/details/7613716
 */
func caculateUTF8Length(firstCh byte) int {
	len := 0
	if firstCh >= 0xFC && firstCh <= 0xFD {
		len = 6
	} else if firstCh >= 0xF8 {
		len = 5
	} else if firstCh >= 0xF0 {
		len = 4
	} else if firstCh >= 0xE0 {
		len = 3
	} else if firstCh >= 0xC0 {
		len = 2
	}
	return len

}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipWhitespaceString() {
	for l.utf8Char == string(' ') || l.utf8Char == string('\t') || l.utf8Char == string('\n') || l.utf8Char == string('\r') {
		l.readUTF8Char()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumeberUTF8() string {
	position := l.position

	for isDigitString(l.utf8Char) {
		l.readUTF8Char()
	}

	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isDigitString(utf8Ch string) bool {
	return "0" <= utf8Ch && utf8Ch <= "9"
}

func (l *Lexer) readUTF8Ident() string {
	position := l.position

	//如果中文或者是字母则继续读取
	for isLetterString(l.utf8Char) || isChineseChar(l.utf8Char) {
		l.readUTF8Char()
	}

	return l.input[position:l.position]
}

func isLetterString(utf8Char string) bool {
	return string('a') <= utf8Char && utf8Char <= string('z') || string('A') <= utf8Char && utf8Char <= string('Z') || utf8Char == string('_')
}

/**
 * 大于127的字符，看作是中文字符
 */
func isChineseChar(utf8Char string) bool {
	firstCh := utf8Char[0] //获取第一个字符
	result := firstCh & 0x80
	return result != 0x00
}
