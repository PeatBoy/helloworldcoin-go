package core

import (
	"fmt"
	"helloworldcoin-go/core/model"
	"helloworldcoin-go/core/model/TransactionType"
	"helloworldcoin-go/core/tool/BlockTool"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/core/tool/ScriptTool"
	"helloworldcoin-go/core/tool/TransactionTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/crypto/ByteUtil"
	"helloworldcoin-go/setting/GenesisBlockSetting"
	"helloworldcoin-go/util/JsonUtil"
	"helloworldcoin-go/util/ThreadUtil"
	"helloworldcoin-go/util/TimeUtil"
)

type Miner struct {
	CoreConfiguration              *CoreConfiguration
	Wallet                         *Wallet
	BlockchainDatabase             *BlockchainDatabase
	UnconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
}

func (i *Miner) Start() {
	for {
		ThreadUtil.MillisecondSleep(10)
		if !i.isActive() {
			continue
		}
		minerAccount := i.Wallet.CreateAccount()
		block := i.buildMiningBlock(i.BlockchainDatabase, i.UnconfirmedTransactionDatabase, minerAccount)
		startTimestamp := TimeUtil.MillisecondTimestamp()
		for {
			if !i.isActive() {
				break
			}
			//在挖矿的期间，可能收集到新的交易。每隔一定的时间，重新组装挖矿中的区块，这样新收集到交易就可以被放进挖矿中的区块了。
			if TimeUtil.MillisecondTimestamp()-startTimestamp > i.CoreConfiguration.GetMinerMineTimeInterval() {
				break
			}
			//随机数
			block.Nonce = ByteUtil.BytesToHexString(ByteUtil.Random32Bytes())
			block.Hash = BlockTool.CalculateBlockHash(block)
			//挖矿成功
			if i.BlockchainDatabase.Consensus.CheckConsensus(i.BlockchainDatabase, block) {
				i.Wallet.SaveAccount(minerAccount)
				blockDto := Model2DtoTool.Block2BlockDto(block)
				fmt.Println(JsonUtil.ToJson(blockDto))
				isAddBlockToBlockchainSuccess := i.BlockchainDatabase.AddBlockDto(blockDto)
				if !isAddBlockToBlockchainSuccess {
					//LogUtil.debug("挖矿成功，但是区块放入区块链失败。")
				}
				break
			}
		}
	}
}
func (i *Miner) isActive() bool {
	//TODO
	return true
}

func (i *Miner) buildMiningBlock(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase, minerAccount *AccountUtil.Account) *model.Block {
	timestamp := TimeUtil.MillisecondTimestamp()

	tailBlock := blockchainDatabase.QueryTailBlock()
	var nonNonceBlock model.Block
	nonNonceBlock.Timestamp = timestamp

	if tailBlock == nil {
		nonNonceBlock.Height = GenesisBlockSetting.HEIGHT + uint64(1)
		nonNonceBlock.PreviousHash = GenesisBlockSetting.HASH
	} else {
		nonNonceBlock.Height = tailBlock.Height + uint64(1)
		nonNonceBlock.PreviousHash = tailBlock.Hash
	}
	packingTransactions := i.packingTransactions(blockchainDatabase, unconfirmedTransactionDatabase)

	incentive := blockchainDatabase.Incentive
	incentiveValue := incentive.IncentiveValue(blockchainDatabase, &nonNonceBlock)

	mineAwardTransaction := i.buildIncentiveTransaction(minerAccount.Address, incentiveValue)
	var mineAwardTransactions []model.Transaction
	mineAwardTransactions = append(mineAwardTransactions, *mineAwardTransaction)

	packingTransactions = append(mineAwardTransactions, packingTransactions...)
	nonNonceBlock.Transactions = packingTransactions

	fmt.Println(JsonUtil.ToJsonStringBlock(&nonNonceBlock))
	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(&nonNonceBlock)
	nonNonceBlock.MerkleTreeRoot = merkleTreeRoot
	fmt.Println(JsonUtil.ToJsonStringBlock(&nonNonceBlock))

	//计算挖矿难度
	nonNonceBlock.Difficulty = blockchainDatabase.Consensus.CalculateDifficult(blockchainDatabase, &nonNonceBlock)
	return &nonNonceBlock
}

func (i *Miner) buildIncentiveTransaction(address string, incentiveValue uint64) *model.Transaction {
	var transaction model.Transaction
	transaction.TransactionType = TransactionType.GENESIS_TRANSACTION

	var outputs []model.TransactionOutput
	var output model.TransactionOutput
	output.Address = address
	output.Value = incentiveValue
	fmt.Println("address:" + address)
	output.OutputScript = ScriptTool.CreatePayToPublicKeyHashOutputScript(address)
	outputs = append(outputs, output)

	transaction.Outputs = outputs
	transaction.TransactionHash = TransactionTool.CalculateTransactionHash(transaction)
	return &transaction
}
func (i *Miner) packingTransactions(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase) []model.Transaction {
	forMineBlockTransactionDtos := unconfirmedTransactionDatabase.SelectTransactions(uint64(1), uint64(10000))

	transactions := []model.Transaction{}
	if forMineBlockTransactionDtos != nil {
		for _, transactionDto := range forMineBlockTransactionDtos {
			//TODO exception
			transaction := TransactionDto2Transaction(blockchainDatabase, &transactionDto)
			transactions = append(transactions, *transaction)
		}
	}
	return transactions
}
