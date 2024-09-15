package tronscan

type GetTokenResponseDto struct {
	Trc20Tokens []Trc20Token `json:"trc20_tokens"`
}

type Trc20Token struct {
	Symbol string `json:"symbol"`
}
