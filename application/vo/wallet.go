package vo

/*
 @author king 409060350@qq.com
*/

import "helloworld-blockchain-go/netcore-dto/dto"

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
const PAYEE_CAN_NOT_EMPTY = "payee_can_not_empty"
const PAYEE_ADDRESS_CAN_NOT_EMPTY = "payee_address_can_not_empty"
const PAYEE_VALUE_CAN_NOT_LESS_EQUAL_THAN_ZERO = "payee_value_can_not_less_equal_than_zero"
const NOT_ENOUGH_MONEY_TO_PAY = "not_enough_money_to_pay"
const BUILD_TRANSACTION_SUCCESS = "build_transaction_success"
