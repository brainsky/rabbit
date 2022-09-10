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
	//l.readChar()
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
			//多字节的情况，该字节高位为1，获取高位1的个数
			var mask1 byte = 0x55
			var mask2 byte = 0x33
			var mask3 byte = 0x0f
			oneLayout := (firstCh & mask1) + ((firstCh >> 1) & mask1)
			twoLayout := (oneLayout & mask2) + ((oneLayout >> 2) & mask2)
			charLen := (twoLayout & mask3) + ((twoLayout >> 4) & mask3)
			//计算移动多少位后
			var charLenInt = int(charLen)
			tempLen := l.readPosition + charLenInt
			if tempLen >= len(l.input) {
				//不正确的字符
				l.utf8Char = "0"
			} else {
				l.utf8Char = l.input[l.readPosition:tempLen]
				l.position = l.readPosition
				l.readPosition += charLenInt
			}
		}

	}

}

func (l *Lexer) NextTokenUTF8() token.Token {
	var tok token.Token
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
	}
	l.readUTF8Char()
	return tok

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
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
	}
	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newTokenUTF8(tokenType token.TokenType, utf8Char string) token.Token {
	return token.Token{Type: tokenType, Literal: utf8Char}
}
