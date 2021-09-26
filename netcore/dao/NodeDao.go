package dao

/*
 @author king 409060350@qq.com
*/

import (
	"helloworld-blockchain-go/netcore/configuration"
	"helloworld-blockchain-go/netcore/po"
	"helloworld-blockchain-go/util/ByteUtil"
	"helloworld-blockchain-go/util/EncodeDecodeTool"
	"helloworld-blockchain-go/util/FileUtil"
	"helloworld-blockchain-go/util/KvDbUtil"
	"helloworld-blockchain-go/util/StringUtil"
)

type NodeDao struct {
	netCoreConfiguration *configuration.NetCoreConfiguration
}

func NewNodeDao(netCoreConfiguration *configuration.NetCoreConfiguration) *NodeDao {
	return &NodeDao{netCoreConfiguration}
}

const NODE_DATABASE_NAME = "NodeDatabase"

func (n *NodeDao) QueryNode(ip string) *po.NodePo {
	nodePos := n.QueryAllNodes()
	if nodePos != nil {
		for _, nodePo := range nodePos {
			if StringUtil.Equals(ip, nodePo.Ip) {
				return nodePo
			}
		}
	}
	return nil
}

func (n *NodeDao) AddNode(node *po.NodePo) {
	KvDbUtil.Put(n.getNodeDatabasePath(), n.getKeyByNodePo(node), EncodeDecodeTool.Encode(node))
}

func (n *NodeDao) UpdateNode(node *po.NodePo) {
	KvDbUtil.Put(n.getNodeDatabasePath(), n.getKeyByNodePo(node), EncodeDecodeTool.Encode(node))
}

func (n *NodeDao) DeleteNode(ip string) {
	KvDbUtil.Delete(n.getNodeDatabasePath(), n.getKeyByIp(ip))
}

func (n *NodeDao) QueryAllNodes() []*po.NodePo {
	var nodePos []*po.NodePo
	//获取所有
	bytesNodePos := KvDbUtil.Gets(n.getNodeDatabasePath(), 1, 100000000)
	if bytesNodePos != nil {
		for e := bytesNodePos.Front(); e != nil; e = e.Next() {
			nodePo := EncodeDecodeTool.Decode(e.Value.([]byte), po.NodePo{}).(*po.NodePo)
			nodePos = append(nodePos, nodePo)
		}
	}
	return nodePos
}
func (n *NodeDao) getNodeDatabasePath() string {
	return FileUtil.NewPath(n.netCoreConfiguration.NetCorePath, NODE_DATABASE_NAME)
}
func (n *NodeDao) getKeyByNodePo(node *po.NodePo) []byte {
	return n.getKeyByIp(node.Ip)
}
func (n *NodeDao) getKeyByIp(ip string) []byte {
	return ByteUtil.StringToUtf8Bytes(ip)
}
