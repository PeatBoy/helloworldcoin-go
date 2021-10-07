package model

/*
 @author king 409060350@qq.com
*/

import "helloworld-blockchain-go/netcore-dto/dto"

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
const PAYEE_CAN_NOT_EMPTY = "payee_can_not_empty"
const PAYEE_ADDRESS_CAN_NOT_EMPTY = "payee_address_can_not_empty"
const PAYEE_VALUE_CAN_NOT_LESS_EQUAL_THAN_ZERO = "payee_value_can_not_less_equal_than_zero"
const NOT_ENOUGH_MONEY_TO_PAY = "not_enough_money_to_pay"
const BUILD_TRANSACTION_SUCCESS = "build_transaction_success"
