package fabricservice

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type FabricService struct {
	ChaincodeObj
	OrgObj
	ConnectionProfile string
	OrdererID         string
	OrdererEndpoint   string //orderer节点地址
	ChannelConfig     string
	initialized       bool
	channelCli        *channel.Client
	resmgmtCli        *resmgmt.Client //资源管理客户端(相当于管理员admin)，用以创建或更新通道
	fabsdk            *fabsdk.FabricSDK
	eventCli          *event.Client
}
type OrgObj struct {
	OrgID    string
	OrgAdmin string
	OrgPeers []string //组织节点
	OrgAchor string   //通信描点
	UserName string
}

type ChaincodeObj struct {
	ChaincodeID      string
	ChaincodeVersion string
	GoPath           string
	ChaincodePath    string
}
