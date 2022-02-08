package netcore

/*
 @author x.king xdotking@gmail.com
*/

import (
	"helloworld-blockchain-go/netcore-client/client"
	"helloworld-blockchain-go/netcore-dto/dto"
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/service"
	"helloworld-blockchain-go/util/LogUtil"
	"helloworld-blockchain-go/util/ThreadUtil"
)

type BlockchainHeightSearcher struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
	nodeService          *service.NodeService
}

func NewBlockchainHeightSearcher(netCoreConfiguration *configuration.NetCoreConfiguration, nodeService *service.NodeService) *BlockchainHeightSearcher {
	var blockchainHeightSearcher BlockchainHeightSearcher
	blockchainHeightSearcher.netCoreConfiguration = netCoreConfiguration
	blockchainHeightSearcher.nodeService = nodeService
	return &blockchainHeightSearcher
}

func (b *BlockchainHeightSearcher) start() {
	defer func() {
		if e := recover(); e != nil {
			LogUtil.Error("在区块链网络中搜索节点的高度异常", e)
		}
	}()
	for {
		b.searchBlockchainHeight()
		ThreadUtil.MillisecondSleep(b.netCoreConfiguration.GetBlockchainHeightSearchTimeInterval())
	}
}

func (b *BlockchainHeightSearcher) searchBlockchainHeight() {
	nodes := b.nodeService.QueryAllNodes()
	if nodes == nil || len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		nodeClient := client.NewNodeClient(node.Ip)
		var getBlockchainHeightRequest dto.GetBlockchainHeightRequest
		getBlockchainHeightResponse := nodeClient.GetBlockchainHeight(getBlockchainHeightRequest)
		if getBlockchainHeightResponse != nil {
			node.BlockchainHeight = getBlockchainHeightResponse.BlockchainHeight
			b.nodeService.UpdateNode(node)
		}
	}
}
