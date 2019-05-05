package fabricservice
import(
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)
type FabricService struct{
	ConnectionProfile string
	OrgID             string
	OrdererID         string
	ChannelID         string
	ChaincodeID       string
	ChaincodeVersion  string
	ChannelConfig     string
	ChaincodeGoPath   string
	ChaincodePath     string
	OrgAdmin          string
	OrgName           string
	OrgPeer0          string
	UserName          string
	initialized       bool
	channelCli             *channel.Client
	resmgmtCli             *resmgmt.Client //资源管理客户端(相当于管理员admin)，用以创建或更新通道
	fabsdk               *fabsdk.FabricSDK
	eventCli           *event.Client
}