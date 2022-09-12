package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

/**
 * 定义一些固定tokentype
 */
const (
	ILEGAL = "ILLEGAL"
	EOF    = "EOF"

	IDENT = "IDENT" //变量
	INT   = "INT"   // 整数

	//操作符
	ASSIGN = "="
	PLUS   = "+"

	//分隔符
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	//关键词
	FUNCTION = "FUNCTION"
	LET      = "LET"

	//添加中文Token
	VAR        = "变量"
	FUNCTIONCN = "函数"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

var keywordsUTF8 = map[string]TokenType{
	"变量": VAR,
	"函数": FUNCTIONCN,
}

func LookupUTFIdent(ident string) TokenType {

	if tok, ok := keywordsUTF8[ident]; ok {
		//如果是中文TOKEN则返回该Token
		return tok
	}
	//否则返回为标识符
	return IDENT
}
