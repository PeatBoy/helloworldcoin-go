package model

/*
 @author king 409060350@qq.com
*/

import "helloworld-blockchain-go/dto"

type Payer struct {
	PrivateKey             string
	TransactionHash        string
	TransactionOutputIndex uint64
	Value                  uint64
	Address                string
}
type Payee struct {
	Address string
	Value   uint64
}
type AutoBuildTransactionRequest struct {
	NonChangePayees []*Payee
}
type AutoBuildTransactionResponse struct {
	BuildTransactionSuccess bool
	Message                 string
	TransactionHash         string
	Fee                     uint64
	Payers                  []*Payer
	NonChangePayees         []*Payee
	ChangePayee             *Payee
	Payees                  []*Payee
	Transaction             *dto.TransactionDto
}

//PayAlert
const PAYEE_CAN_NOT_EMPTY = "PAYEE_CAN_NOT_EMPTY"
const PAYEE_ADDRESS_CAN_NOT_EMPTY = "PAYEE_ADDRESS_CAN_NOT_EMPTY"
const PAYEE_VALUE_CAN_NOT_LESS_EQUAL_THAN_ZERO = "PAYEE_VALUE_CAN_NOT_LESS_EQUAL_THAN_ZERO"
const NOT_ENOUGH_MONEY_TO_PAY = "NOT_ENOUGH_MONEY_TO_PAY"
const BUILD_TRANSACTION_SUCCESS = "BUILD_TRANSACTION_SUCCESS"
