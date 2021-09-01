package JsonUtil

import (
	"encoding/json"
	"helloworld-blockchain-go/dto"
	"helloworld-blockchain-go/netcore/po"
)

func ToString(object interface{}) string {
	jsonString, _ := json.Marshal(object)
	return string(jsonString)
}

func ToObject(jsonString string, object interface{}) interface{} {
	_0001, ok := object.(dto.GetBlockRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0001)
		return &_0001
	}
	_1001, ok := object.(dto.GetBlockResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1001)
		return &_1001
	}
	_0002, ok := object.(dto.PostBlockRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0002)
		return &_0002
	}
	_1002, ok := object.(dto.PostBlockResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1002)
		return &_1002
	}
	_0003, ok := object.(dto.GetUnconfirmedTransactionsRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0003)
		return &_0003
	}
	_1003, ok := object.(dto.GetUnconfirmedTransactionsResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1003)
		return &_1003
	}
	_0004, ok := object.(dto.GetBlockchainHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0004)
		return &_0004
	}
	_1004, ok := object.(dto.GetBlockchainHeightResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1004)
		return &_1004
	}
	_0005, ok := object.(dto.PostTransactionRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0005)
		return &_0005
	}
	_1005, ok := object.(dto.PostTransactionResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1005)
		return &_1005
	}
	_0006, ok := object.(dto.PostBlockchainHeightRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0006)
		return &_0006
	}
	_1006, ok := object.(dto.PostBlockchainHeightResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1006)
		return &_1006
	}
	_0007, ok := object.(dto.PingRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0007)
		return &_0007
	}
	_1007, ok := object.(dto.PingResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1007)
		return &_1007
	}
	_0008, ok := object.(dto.GetNodesRequest)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0008)
		return &_0008
	}
	_1008, ok := object.(dto.GetNodesResponse)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1008)
		return &_1008
	}
	_0009, ok := object.(dto.BlockDto)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0009)
		return &_0009
	}
	_1009, ok := object.(dto.BlockDto)
	if ok {
		json.Unmarshal([]byte(jsonString), &_1009)
		return &_1009
	}

	_0010, ok := object.(po.NodePo)
	if ok {
		json.Unmarshal([]byte(jsonString), &_0010)
		return &_0010
	}
	panic("JsonUtil.ToObject can not recognize object type")
}
