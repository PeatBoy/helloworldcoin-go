package core

import (
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/model/TransactionType"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/core/tool/ScriptTool"
	"helloworld-blockchain-go/core/tool/TransactionTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
	"helloworld-blockchain-go/util/TimeUtil"
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
		if !i.IsActive() {
			continue
		}

		blockChainHeight := i.BlockchainDatabase.QueryBlockchainHeight()
		//'当前区块链的高度'是否大于'矿工最大被允许的挖矿高度'
		if blockChainHeight >= i.CoreConfiguration.getMaxBlockHeight() {
			continue
		}

		minerAccount := i.Wallet.CreateAccount()
		block := i.buildMiningBlock(i.BlockchainDatabase, i.UnconfirmedTransactionDatabase, minerAccount)
		startTimestamp := TimeUtil.MillisecondTimestamp()
		for {
			if !i.IsActive() {
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
				LogUtil.Debug("祝贺您！挖矿成功！！！区块高度:" + StringUtil.ValueOfUint64(block.Height) + ",区块哈希:" + block.Hash)
				blockDto := Model2DtoTool.Block2BlockDto(block)
				isAddBlockToBlockchainSuccess := i.BlockchainDatabase.AddBlockDto(blockDto)
				if !isAddBlockToBlockchainSuccess {
					LogUtil.Debug("挖矿成功，但是区块放入区块链失败。")
				}
				break
			}
		}
	}
}
func (i *Miner) IsActive() bool {
	return i.CoreConfiguration.IsMinerActive()
}
func (i *Miner) Deactive() {
	i.CoreConfiguration.deactiveMiner()
}
func (i *Miner) Active() {
	i.CoreConfiguration.activeMiner()
}

func (i *Miner) SetMaxBlockHeight(maxHeight uint64) {
	i.CoreConfiguration.setMaxBlockHeight(maxHeight)
}

func (i *Miner) GetMaxBlockHeight() uint64 {
	return i.CoreConfiguration.getMaxBlockHeight()
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
	var mineAwardTransactions []*model.Transaction
	mineAwardTransactions = append(mineAwardTransactions, mineAwardTransaction)

	packingTransactions = append(mineAwardTransactions, packingTransactions...)
	nonNonceBlock.Transactions = packingTransactions

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(&nonNonceBlock)
	nonNonceBlock.MerkleTreeRoot = merkleTreeRoot

	//计算挖矿难度
	nonNonceBlock.Difficulty = blockchainDatabase.Consensus.CalculateDifficult(blockchainDatabase, &nonNonceBlock)
	return &nonNonceBlock
}

func (i *Miner) buildIncentiveTransaction(address string, incentiveValue uint64) *model.Transaction {
	var transaction model.Transaction
	transaction.TransactionType = TransactionType.GENESIS_TRANSACTION

	var outputs []*model.TransactionOutput
	var output model.TransactionOutput
	output.Address = address
	output.Value = incentiveValue
	output.OutputScript = ScriptTool.CreatePayToPublicKeyHashOutputScript(address)
	outputs = append(outputs, &output)

	transaction.Outputs = outputs
	transaction.TransactionHash = TransactionTool.CalculateTransactionHash(transaction)
	return &transaction
}
func (i *Miner) packingTransactions(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase) []*model.Transaction {
	forMineBlockTransactionDtos := unconfirmedTransactionDatabase.SelectTransactions(uint64(1), uint64(10000))

	var transactions []*model.Transaction
	if forMineBlockTransactionDtos != nil {
		for _, transactionDto := range forMineBlockTransactionDtos {
			//TODO exception
			transaction := TransactionDto2Transaction(blockchainDatabase, transactionDto)
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}
