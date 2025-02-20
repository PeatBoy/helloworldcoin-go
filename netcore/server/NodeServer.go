package server

/*
 @author x.king xdotking@gmail.com
*/

import (
	"fmt"
	"helloworldcoin-go/core"
	"helloworldcoin-go/core/tool/Model2DtoTool"
	"helloworldcoin-go/netcore-dto/dto"
	"helloworldcoin-go/netcore-dto/dto/API"
	"helloworldcoin-go/netcore/configuration"
	"helloworldcoin-go/netcore/model"
	"helloworldcoin-go/netcore/service"
	"helloworldcoin-go/setting/BlockSetting"
	"helloworldcoin-go/util/JsonUtil"
	"helloworldcoin-go/util/LogUtil"
	"helloworldcoin-go/util/SystemUtil"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type NodeServer struct {
	blockchainCore       *core.BlockchainCore
	nodeService          *service.NodeService
	netCoreConfiguration *configuration.NetCoreConfiguration
}

func NewNodeServer(netCoreConfiguration *configuration.NetCoreConfiguration, blockchainCore *core.BlockchainCore, nodeService *service.NodeService) *NodeServer {
	return &NodeServer{blockchainCore, nodeService, netCoreConfiguration}
}

func (b *NodeServer) Start() {
	http.HandleFunc(API.PING, b.ping)
	http.HandleFunc(API.GET_NODES, b.getNodes)
	http.HandleFunc(API.GET_BLOCK, b.getBlock)
	http.HandleFunc(API.POST_BLOCK, b.postBlock)
	http.HandleFunc(API.POST_BLOCKCHAIN_HEIGHT, b.postBlockchainHeight)
	http.HandleFunc(API.GET_BLOCKCHAIN_HEIGHT, b.getBlockchainHeight)
	http.HandleFunc(API.POST_TRANSACTION, b.postTransaction)
	http.HandleFunc(API.GET_UNCONFIRMED_TRANSACTIONS, b.getUnconfirmedTransactions)

	e := http.ListenAndServe(":8888", nil)
	if e != nil {
		SystemUtil.ErrorExit("blockchain node server can not start.", e)
	}
}

func (b *NodeServer) getBlock(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetBlockRequest{}).(*dto.GetBlockRequest)
	blockByBlockHeight := b.blockchainCore.QueryBlockByBlockHeight(request.BlockHeight)
	block := Model2DtoTool.Block2BlockDto(blockByBlockHeight)
	var response dto.GetBlockResponse
	response.Block = block
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *NodeServer) postBlock(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PostBlockRequest{}).(*dto.PostBlockRequest)
	b.blockchainCore.AddBlockDto(request.Block)
	var response dto.PostBlockResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *NodeServer) postBlockchainHeight(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PostBlockchainHeightRequest{}).(*dto.PostBlockchainHeightRequest)
	requestIp := getIP(req)

	var node model.Node
	node.Ip = requestIp
	node.BlockchainHeight = request.BlockchainHeight
	b.nodeService.UpdateNode(&node)
	var response dto.PostBlockchainHeightResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *NodeServer) getBlockchainHeight(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetBlockchainHeightRequest{}).(*dto.GetBlockchainHeightRequest)
	fmt.Println(request)
	blockchainHeight := b.blockchainCore.QueryBlockchainHeight()
	var response dto.GetBlockchainHeightResponse
	response.BlockchainHeight = blockchainHeight
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *NodeServer) postTransaction(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PostTransactionRequest{}).(*dto.PostTransactionRequest)
	b.blockchainCore.PostTransaction(request.Transaction)
	var response dto.PostTransactionResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}
func (b *NodeServer) getUnconfirmedTransactions(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetUnconfirmedTransactionsRequest{}).(*dto.GetUnconfirmedTransactionsRequest)
	fmt.Println(request)
	unconfirmedTransactionDatabase := b.blockchainCore.GetUnconfirmedTransactionDatabase()
	transactions := unconfirmedTransactionDatabase.SelectTransactions(1, BlockSetting.BLOCK_MAX_TRANSACTION_COUNT)
	var response dto.GetUnconfirmedTransactionsResponse
	response.Transactions = transactions
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}

func (b *NodeServer) ping(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.PingRequest{}).(*dto.PingRequest)
	requestIp := getIP(req)
	fmt.Println(request)

	//将ping的来路作为区块链节点
	if b.netCoreConfiguration.IsAutoSearchNode() {
		var node model.Node
		node.Ip = requestIp
		node.BlockchainHeight = 0
		b.nodeService.AddNode(&node)
		LogUtil.Debug("发现节点[" + requestIp + "]在Ping本地节点，已将发现的节点放入了节点数据库。")
	}
	var response dto.PingResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}

func (b *NodeServer) getNodes(w http.ResponseWriter, req *http.Request) {
	requestBody, _ := ioutil.ReadAll(req.Body)
	request := JsonUtil.ToObject(string(requestBody), dto.GetNodesRequest{}).(*dto.GetNodesRequest)
	fmt.Println(request)
	allNodes := b.nodeService.QueryAllNodes()
	var nodes []dto.NodeDto
	if allNodes != nil {
		for _, node := range allNodes {
			var n dto.NodeDto
			n.Ip = node.Ip
			nodes = append(nodes, n)
		}
	}
	var response dto.GetNodesResponse
	w.Header().Set("content-type", "text/json")
	io.WriteString(w, JsonUtil.ToString(response))
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "localhost"
	}
	if net.ParseIP(ip) != nil {
		return ip
	}
	return "localhost"
}
