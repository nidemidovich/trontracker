package tron

import "encoding/json"

// /wallet/getblock
type GetBlockResponseDto struct {
	BlockID      string        `json:"blockID"`
	BlockHeader  BlockHeader   `json:"block_header"`
	Transactions []Transaction `json:"transactions"`
}

type BlockHeader struct {
	RawData          RawData `json:"raw_data"`
	WitnessSignature string  `json:"witness_signature"`
}

type RawData struct {
	Number         int64  `json:"number"`
	TxTrieRoot     string `json:"txTrieRoot"`
	WitnessAddress string `json:"witness_address"`
	ParentHash     string `json:"parentHash"`
	Version        int64  `json:"version"`
	Timestamp      int64  `json:"timestamp"`
}

type Transaction struct {
	Ret        []interface{}      `json:"ret"`
	Signature  []string           `json:"signature"`
	TxID       string             `json:"txID"`
	RawData    TransactionRawData `json:"raw_data"`
	RawDataHex string             `json:"raw_data_hex"`
}

type TransactionRawData struct {
	Contract      []Contract `json:"contract"`
	RefBlockBytes string     `json:"ref_block_bytes"`
	RefBlockHash  string     `json:"ref_block_hash"`
	Expiration    int64      `json:"expiration"`
	FeeLimit      int64      `json:"fee_limit"`
	Timestamp     int64      `json:"timestamp"`
}

type Contract struct {
	Type      string          `json:"type"`
	Parameter json.RawMessage `json:"parameter"`
}

type TriggerSmartContractParameter struct {
	Value   Value  `json:"value"`
	TypeUrl string `json:"type_url"`
}

type Value struct {
	Data            string `json:"data"`
	OwnerAddress    string `json:"owner_address"`
	ContractAddress string `json:"contract_address"`
	CallValue       *int64 `json:"call_value,omitempty"`
}

// /wallet/gettransactioninfobyblocknum
type GetTransactionInfoByBlockNumResponseDto []TransactionInfo

type TransactionInfo struct {
	Result         *string    `json:"result,omitempty"`
	ContractResult []string   `json:"contractResult"`
	Log            []LogEntry `json:"log,omitempty"`
}

type LogEntry struct {
	Address string   `json:"address"`
	Data    string   `json:"data"`
	Topics  []string `json:"topics"`
}

// /wallet/triggerconstantcontract
type GetTokenSymbolResponseDto struct {
	ConstantResult []string `json:"constant_result"`
}
