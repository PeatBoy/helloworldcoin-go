package core

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworld-blockchain-go/core/model"
	"helloworld-blockchain-go/core/model/script/BooleanCode"
	"helloworld-blockchain-go/core/model/script/OperationCode"
	"helloworld-blockchain-go/core/tool/TransactionTool"
	"helloworld-blockchain-go/crypto/AccountUtil"
	"helloworld-blockchain-go/util/ByteUtil"
	"helloworld-blockchain-go/util/StringStack"
	"helloworld-blockchain-go/util/StringUtil"
)

type VirtualMachine struct {
}

func NewVirtualMachine() *VirtualMachine {
	var virtualMachine VirtualMachine
	return &virtualMachine
}

func (this *VirtualMachine) Execute(transactionEnvironment *model.Transaction, script *model.Script) *model.Result {

	stack := StringStack.NewStringStack()
	for i := 0; i < len(*script); i++ {
		operationCode := (*script)[i]
		bytesOperationCode := ByteUtil.HexStringToBytes(operationCode)
		if ByteUtil.Equals(OperationCode.OP_DUP.Code, bytesOperationCode) {
			if stack.Size() < 1 {
				panic("指令运行异常")
			}
			stack.Push(*stack.Peek())
		} else if ByteUtil.Equals(OperationCode.OP_HASH160.Code, bytesOperationCode) {
			if stack.Size() < 1 {
				panic("指令运行异常")
			}
			publicKey := stack.Pop()
			publicKeyHash := AccountUtil.PublicKeyHashFromPublicKey(*publicKey)
			stack.Push(publicKeyHash)
		} else if ByteUtil.Equals(OperationCode.OP_EQUALVERIFY.Code, bytesOperationCode) {
			if stack.Size() < 2 {
				panic("指令运行异常")
			}
			if !StringUtil.Equals(*stack.Pop(), *stack.Pop()) {
				panic("脚本执行失败")
			}
		} else if ByteUtil.Equals(OperationCode.OP_CHECKSIG.Code, bytesOperationCode) {
			if stack.Size() < 2 {
				panic("指令运行异常")
			}
			publicKey := stack.Pop()
			signature := stack.Pop()
			bytesSignature := ByteUtil.HexStringToBytes(*signature)
			verifySignatureSuccess := TransactionTool.VerifySignature(transactionEnvironment, *publicKey, bytesSignature)
			if !verifySignatureSuccess {
				panic("脚本执行失败")
			}
			stack.Push(ByteUtil.BytesToHexString(BooleanCode.TRUE.Code))
		} else if ByteUtil.Equals(OperationCode.OP_PUSHDATA.Code, bytesOperationCode) {
			if len(*script) < i+2 {
				panic("指令运行异常")
			}
			i++
			stack.Push((*script)[i])
		} else {
			panic("不能识别的操作码")
		}
	}
	return stack
}
