package vo

/*
 @author king 409060350@qq.com
*/

import "helloworld-blockchain-go/dto"

type PayerVo struct {
	PrivateKey             string `json:"privateKey"`
	TransactionHash        string `json:"transactionHash"`
	TransactionOutputIndex uint64 `json:"transactionOutputIndex"`
	Value                  uint64 `json:"value"`
	Address                string `json:"address"`
}
type PayeeVo struct {
	Address string `json:"address"`
	Value   uint64 `json:"value"`
}

type AutomaticBuildTransactionRequest struct {
	NonChangePayees []*PayeeVo `json:"nonChangePayees"`
}
type AutomaticBuildTransactionResponse struct {
	BuildTransactionSuccess bool                `json:"buildTransactionSuccess"`
	Message                 string              `json:"message"`
	TransactionHash         string              `json:"transactionHash"`
	Fee                     uint64              `json:"fee"`
	Payers                  []*PayerVo          `json:"payers"`
	NonChangePayees         []*PayeeVo          `json:"nonChangePayees"`
	ChangePayee             *PayeeVo            `json:"changePayee"`
	Payees                  []*PayeeVo          `json:"payees"`
	Transaction             *dto.TransactionDto `json:"transaction"`
}

//PayAlert
const PAYEE_CAN_NOT_EMPTY = "PAYEE_CAN_NOT_EMPTY"
const PAYEE_ADDRESS_CAN_NOT_EMPTY = "PAYEE_ADDRESS_CAN_NOT_EMPTY"
const PAYEE_VALUE_CAN_NOT_LESS_EQUAL_THAN_ZERO = "PAYEE_VALUE_CAN_NOT_LESS_EQUAL_THAN_ZERO"
const NOT_ENOUGH_MONEY_TO_PAY = "NOT_ENOUGH_MONEY_TO_PAY"
const BUILD_TRANSACTION_SUCCESS = "BUILD_TRANSACTION_SUCCESS"
