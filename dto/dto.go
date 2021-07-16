package dto

type BlockDto struct {
	Timestamp    uint64           `json:"timestamp"`
	PreviousHash string           `json:"previousHash"`
	Transactions []TransactionDto `json:"transactions"`
	Nonce        string           `json:"nonce"`
}
type TransactionDto struct {
	Inputs  []TransactionInputDto  `json:"inputs"`
	Outputs []TransactionOutputDto `json:"outputs"`
}
type TransactionInputDto struct {
	TransactionHash        string   `json:"transactionHash"`
	TransactionOutputIndex uint64   `json:"transactionOutputIndex"`
	InputScript            []string `json:"inputScript"`
}
type TransactionOutputDto struct {
	OutputScript []string `json:"outputScript"`
	Value        uint64   `json:"value"`
}
type ScriptDto = []string
type InputScriptDto = ScriptDto
type OutputScriptDto = ScriptDto
