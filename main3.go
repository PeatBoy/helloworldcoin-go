package main

import (
	"fmt"
	"helloworldcoin-go/core"
	"helloworldcoin-go/dto"
	"helloworldcoin-go/util/JsonUtil"
)

func main() {

	transactionOutputDto := dto.TransactionOutputDto{OutputScript: []string{"01", "02", "00", "82f46bdbb4550d3c552f1764b53fd0005c81ad3d", "03", "04"}, Value: uint64(5000000000)}
	var outputs []dto.TransactionOutputDto
	outputs = append(outputs, transactionOutputDto)

	transactionDto := dto.TransactionDto{Inputs: nil, Outputs: outputs}
	var transactions []dto.TransactionDto
	transactions = append(transactions, transactionDto)

	blockDto := &dto.BlockDto{Timestamp: uint64(1623645017179), PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000", Transactions: transactions, Nonce: "c21afb034d3be7f2f233b72aa4136dfcc4ee258af213b91a616bcca1ab780f5b"}

	consensus := &core.Consensus{}
	incentive := &core.Incentive{}
	coreConfiguration := &core.CoreConfiguration{CorePath: "d:"}
	blockchainDatabase := core.BlockchainDatabase{Consensus: consensus, Incentive: incentive, CoreConfiguration: coreConfiguration}
	blockchainDatabase.AddBlockDto(blockDto)

	block := blockchainDatabase.QueryBlockByBlockHeight(uint64(1))
	fmt.Println(JsonUtil.ToJsonStringBlock(block))
}
