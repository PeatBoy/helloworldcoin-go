package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworld-blockchain-go/core/tool/TransactionDtoTool"
	"helloworld-blockchain-go/netcore-dto/dto"
	"helloworld-blockchain-go/util/ByteUtil"
	"helloworld-blockchain-go/util/EncodeDecodeTool"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/JsonUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
	"helloworld-blockchain-go/util/LogUtil"
)

const UNCONFIRMED_TRANSACTION_DATABASE_NAME = "UnconfirmedTransactionDatabase"

type UnconfirmedTransactionDatabase struct {
	coreConfiguration *CoreConfiguration
}

func NewUnconfirmedTransactionDatabase(coreConfiguration *CoreConfiguration) *UnconfirmedTransactionDatabase {
	var unconfirmedTransactionDatabase UnconfirmedTransactionDatabase
	unconfirmedTransactionDatabase.coreConfiguration = coreConfiguration
	return &unconfirmedTransactionDatabase
}

func (u *UnconfirmedTransactionDatabase) InsertTransaction(transaction *dto.TransactionDto) bool {
	defer func() {
		//it's mean default return value is false
		if e := recover(); e != nil {
			LogUtil.Error("交易["+JsonUtil.ToString(transaction)+"]放入交易池异常。", e)
		}
	}()
	transactionHash := TransactionDtoTool.CalculateTransactionHash(transaction)
	KvDbUtil.Put(u.getUnconfirmedTransactionDatabasePath(), u.getKey(transactionHash), EncodeDecodeTool.Encode(transaction))
	return true
}

func (u *UnconfirmedTransactionDatabase) SelectTransactions(from uint64, size uint64) []*dto.TransactionDto {
	var transactionDtos []*dto.TransactionDto
	bytesTransactionDtos := KvDbUtil.Gets(u.getUnconfirmedTransactionDatabasePath(), from, size)
	if bytesTransactionDtos != nil {
		for e := bytesTransactionDtos.Front(); e != nil; e = e.Next() {
			transactionDto := EncodeDecodeTool.Decode(e.Value.([]byte), dto.TransactionDto{}).(*dto.TransactionDto)
			transactionDtos = append(transactionDtos, transactionDto)
		}
	}
	return transactionDtos
}

func (u *UnconfirmedTransactionDatabase) DeleteByTransactionHash(transactionHash string) {
	KvDbUtil.Delete(u.getUnconfirmedTransactionDatabasePath(), u.getKey(transactionHash))
}

func (u *UnconfirmedTransactionDatabase) SelectTransactionByTransactionHash(transactionHash string) *dto.TransactionDto {
	byteTransactionDto := KvDbUtil.Get(u.getUnconfirmedTransactionDatabasePath(), u.getKey(transactionHash))
	if byteTransactionDto == nil {
		return nil
	}
	return EncodeDecodeTool.Decode(byteTransactionDto, dto.TransactionDto{}).(*dto.TransactionDto)
}

func (u *UnconfirmedTransactionDatabase) getUnconfirmedTransactionDatabasePath() string {
	return FileUtil.NewPath(u.coreConfiguration.corePath, UNCONFIRMED_TRANSACTION_DATABASE_NAME)
}

func (u *UnconfirmedTransactionDatabase) getKey(transactionHash string) []byte {
	return ByteUtil.HexStringToBytes(transactionHash)
}
