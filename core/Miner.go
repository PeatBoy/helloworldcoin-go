package core

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/model/TransactionType"
	"helloworld-blockchain-go/core/tool/BlockTool"
	"helloworld-blockchain-go/core/tool/Model2DtoTool"
	"helloworld-blockchain-go/core/tool/ScriptTool"
	"helloworld-blockchain-go/core/tool/SizeTool"
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/core/tool/TransactionTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/crypto/ByteUtil"
	"helloworld-blockchain-go/setting/BlockSetting"
	"helloworld-blockchain-go/setting/GenesisBlockSetting"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/StringUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
	"helloworld-blockchain-go/util/TimeUtil"
)

type Miner struct {
	coreConfiguration              *CoreConfiguration
	wallet                         *Wallet
	blockchainDatabase             *BlockchainDatabase
	unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase
}

func NewMiner(coreConfiguration *CoreConfiguration, wallet *Wallet, blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase) *Miner {
	var miner Miner
	miner.coreConfiguration = coreConfiguration
	miner.wallet = wallet
	miner.blockchainDatabase = blockchainDatabase
	miner.unconfirmedTransactionDatabase = unconfirmedTransactionDatabase
	return &miner
}

func (i *Miner) Start() {
	for {
		ThreadUtil.MillisecondSleep(10)
		if !i.IsActive() {
			continue
		}

		blockChainHeight := i.blockchainDatabase.QueryBlockchainHeight()
		//'当前区块链的高度'是否大于'矿工最大被允许的挖矿高度'
		if blockChainHeight >= i.coreConfiguration.getMaxBlockHeight() {
			continue
		}

		minerAccount := i.wallet.CreateAccount()
		block := i.buildMiningBlock(i.blockchainDatabase, i.unconfirmedTransactionDatabase, minerAccount)
		startTimestamp := TimeUtil.MillisecondTimestamp()
		for {
			if !i.IsActive() {
				break
			}
			//在挖矿的期间，可能收集到新的交易。每隔一定的时间，重新组装挖矿中的区块，这样新收集到交易就可以被放进挖矿中的区块了。
			if TimeUtil.MillisecondTimestamp()-startTimestamp > i.coreConfiguration.GetMinerMineTimeInterval() {
				break
			}
			//随机数
			block.Nonce = ByteUtil.BytesToHexString(ByteUtil.Random32Bytes())
			block.Hash = BlockTool.CalculateBlockHash(block)
			//挖矿成功
			if i.blockchainDatabase.consensus.CheckConsensus(i.blockchainDatabase, block) {
				i.wallet.SaveAccount(minerAccount)
				LogUtil.Debug("祝贺您！挖矿成功！！！区块高度:" + StringUtil.ValueOfUint64(block.Height) + ",区块哈希:" + block.Hash)
				blockDto := Model2DtoTool.Block2BlockDto(block)
				isAddBlockToBlockchainSuccess := i.blockchainDatabase.AddBlockDto(blockDto)
				if !isAddBlockToBlockchainSuccess {
					LogUtil.Debug("挖矿成功，但是区块放入区块链失败。")
				}
				break
			}
		}
	}
}
func (i *Miner) IsActive() bool {
	return i.coreConfiguration.IsMinerActive()
}
func (i *Miner) Deactive() {
	i.coreConfiguration.deactiveMiner()
}
func (i *Miner) Active() {
	i.coreConfiguration.activeMiner()
}

func (i *Miner) SetMaxBlockHeight(maxHeight uint64) {
	i.coreConfiguration.setMaxBlockHeight(maxHeight)
}

func (i *Miner) GetMaxBlockHeight() uint64 {
	return i.coreConfiguration.getMaxBlockHeight()
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

	incentive := blockchainDatabase.incentive
	incentiveValue := incentive.IncentiveValue(blockchainDatabase, &nonNonceBlock)

	mineAwardTransaction := i.buildIncentiveTransaction(minerAccount.Address, incentiveValue)
	var mineAwardTransactions []*model.Transaction
	mineAwardTransactions = append(mineAwardTransactions, mineAwardTransaction)

	packingTransactions = append(mineAwardTransactions, packingTransactions...)
	nonNonceBlock.Transactions = packingTransactions

	merkleTreeRoot := BlockTool.CalculateBlockMerkleTreeRoot(&nonNonceBlock)
	nonNonceBlock.MerkleTreeRoot = merkleTreeRoot

	//计算挖矿难度
	nonNonceBlock.Difficulty = blockchainDatabase.consensus.CalculateDifficult(blockchainDatabase, &nonNonceBlock)
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
	transaction.TransactionHash = TransactionTool.CalculateTransactionHash(&transaction)
	return &transaction
}
func (i *Miner) packingTransactions(blockchainDatabase *BlockchainDatabase, unconfirmedTransactionDatabase *UnconfirmedTransactionDatabase) []*model.Transaction {
	forMineBlockTransactionDtos := unconfirmedTransactionDatabase.SelectTransactions(uint64(1), uint64(10000))

	var transactions []*model.Transaction
	var backupTransactions []*model.Transaction

	if forMineBlockTransactionDtos != nil {
		for _, transactionDto := range forMineBlockTransactionDtos {
			defer func() {
				if e := recover(); e != nil {
					transactionHash := TransactionDtoTool.CalculateTransactionHash(transactionDto)
					LogUtil.Error("类型转换异常,将从挖矿交易数据库中删除该交易["+transactionHash+"]。", e)
					unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
				}
			}()
			transaction := blockchainDatabase.TransactionDto2Transaction(transactionDto)
			transactions = append(transactions, transaction)
		}
	}

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}
	for _, transaction := range backupTransactions {
		checkTransaction := blockchainDatabase.CheckTransaction(transaction)
		if checkTransaction {
			transactions = append(transactions, transaction)
		} else {
			transactionHash := TransactionTool.CalculateTransactionHash(transaction)
			LogUtil.Debug("交易不能被挖矿,将从挖矿交易数据库中删除该交易。交易哈希" + transactionHash)
			unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
		}
	}

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}
	//防止双花
	transactionOutputIdSet := make(map[string]bool)
	for _, transaction := range backupTransactions {
		inputs := transaction.Inputs
		if inputs != nil {
			canAdd := true
			for _, transactionInput := range inputs {
				unspentTransactionOutput := transactionInput.UnspentTransactionOutput
				transactionOutputId := TransactionTool.GetTransactionOutputId(unspentTransactionOutput)
				_, exists := transactionOutputIdSet[transactionOutputId]
				if exists {
					canAdd = false
					break
				} else {
					transactionOutputIdSet[transactionOutputId] = true
				}
			}
			if canAdd {
				transactions = append(transactions, transaction)
			} else {
				transactionHash := TransactionTool.CalculateTransactionHash(transaction)
				LogUtil.Debug("交易不能被挖矿,将从挖矿交易数据库中删除该交易。交易哈希" + transactionHash)
				unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
			}
		}
	}

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}
	//防止一个地址被用多次
	addressSet := make(map[string]bool)
	for _, transaction := range backupTransactions {
		outputs := transaction.Outputs
		if outputs != nil {
			canAdd := true
			for _, output := range outputs {
				address := output.Address
				_, exists := addressSet[address]
				if exists {
					canAdd = false
					break
				} else {
					addressSet[address] = true
				}
			}
			if canAdd {
				transactions = append(transactions, transaction)
			} else {
				transactionHash := TransactionTool.CalculateTransactionHash(transaction)
				LogUtil.Debug("交易不能被挖矿,将从挖矿交易数据库中删除该交易。交易哈希" + transactionHash)
				unconfirmedTransactionDatabase.DeleteByTransactionHash(transactionHash)
			}
		}
	}

	//按照费率(每字符的手续费)从大到小排序交易
	//TransactionTool.sortByTransactionFeeRateDescend(transactions);

	backupTransactions = []*model.Transaction{}
	backupTransactions = append(backupTransactions, transactions...)
	transactions = []*model.Transaction{}
	//到此时，剩余交易都是经过验证的了，且按照交易费率从大到小排列了。
	//尽可能多的获取交易
	size := uint64(0)
	for i := 0; i < len(backupTransactions); i++ {
		//序号从0开始，加一。
		//留给挖矿交易一个位置，减一。
		if uint64(i+1) > BlockSetting.BLOCK_MAX_TRANSACTION_COUNT-1 {
			break
		}
		transaction := backupTransactions[i]
		size += SizeTool.CalculateTransactionSize(transaction)
		if size > BlockSetting.BLOCK_MAX_CHARACTER_COUNT {
			break
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
