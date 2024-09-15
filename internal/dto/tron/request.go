package tron

type GetBlockRequestDto struct {
	IdOrNum *string `json:"id_or_num,omitempty"`
	Detail  bool    `json:"detail"`
}

type GetTransactionInfoByBlockNumRequestDto struct {
	Num int64 `json:"num"`
}

type GetTokenAddressRequestDto struct {
	ContractAddress  string `json:"contract_address"`
	FunctionSelector string `json:"function_selector"`
	OwnerAddress     string `json:"owner_address"`
	Visible          bool   `json:"visible"`
}
