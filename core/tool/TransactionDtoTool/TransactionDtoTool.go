package TransactionDtoTool

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworldcoin-go/core/tool/ScriptDtoTool"
	"helloworldcoin-go/crypto/AccountUtil"
	"helloworldcoin-go/crypto/Sha256Util"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/util/ByteUtil"
)

func CalculateTransactionHash(transaction *dto.TransactionDto) string {
	bytesTransaction := BytesTransaction(transaction, false)
	sha256Digest := Sha256Util.DoubleDigest(bytesTransaction)
	return ByteUtil.BytesToHexString(sha256Digest)
}

func SignatureHashAll(transactionDto *dto.TransactionDto) string {
	bytesTransaction := BytesTransaction(transactionDto, true)
	sha256Digest := Sha256Util.DoubleDigest(bytesTransaction)
	return ByteUtil.BytesToHexString(sha256Digest)
}

func Signature(privateKey string, transactionDto *dto.TransactionDto) string {
	signatureHashAll := SignatureHashAll(transactionDto)
	bytesSignatureHashAll := ByteUtil.HexStringToBytes(signatureHashAll)
	signature := AccountUtil.Signature(privateKey, bytesSignatureHashAll)
	return signature
}

func VerifySignature(transaction *dto.TransactionDto, publicKey string, bytesSignature []byte) bool {
	message := SignatureHashAll(transaction)
	bytesMessage := ByteUtil.HexStringToBytes(message)
	return AccountUtil.VerifySignature(publicKey, bytesMessage, bytesSignature)
}

//region Serialization and Deserialization
/**
 * Serialization: Convert TransactionDto into byte array. Requires that the resulting byte array can Convert into the original transaction.
 */
func BytesTransaction(transaction *dto.TransactionDto, omitInputScript bool) []byte {
	var bytesUnspentTransactionOutputs [][]byte
	inputs := transaction.Inputs
	for _, input := range inputs {
		bytesTransactionHash := ByteUtil.HexStringToBytes(input.TransactionHash)
		bytesTransactionOutputIndex := ByteUtil.Uint64ToBytes(input.TransactionOutputIndex)
		bytesUnspentTransactionOutput := ByteUtil.Concatenate(ByteUtil.ConcatenateLength(bytesTransactionHash),
			ByteUtil.ConcatenateLength(bytesTransactionOutputIndex))
		if !omitInputScript {
			bytesInputScript := ScriptDtoTool.BytesInputScript(input.InputScript)
			bytesUnspentTransactionOutput = ByteUtil.Concatenate(bytesUnspentTransactionOutput, ByteUtil.ConcatenateLength(bytesInputScript))
		}
		bytesUnspentTransactionOutputs = append(bytesUnspentTransactionOutputs, ByteUtil.ConcatenateLength(bytesUnspentTransactionOutput))
	}

	var bytesTransactionOutputs [][]byte
	outputs := transaction.Outputs
	for _, output := range outputs {
		bytesOutputScript := ScriptDtoTool.BytesOutputScript(output.OutputScript)
		bytesValue := ByteUtil.Uint64ToBytes(output.Value)
		bytesTransactionOutput := ByteUtil.Concatenate(ByteUtil.ConcatenateLength(bytesOutputScript), ByteUtil.ConcatenateLength(bytesValue))
		bytesTransactionOutputs = append(bytesTransactionOutputs, ByteUtil.ConcatenateLength(bytesTransactionOutput))
	}

	data := ByteUtil.Concatenate(ByteUtil.FlatAndConcatenateLength(bytesUnspentTransactionOutputs),
		ByteUtil.FlatAndConcatenateLength(bytesTransactionOutputs))
	return data
}

//endregion
