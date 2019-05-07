package fabricservice

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

type FabricService struct {
	// ChaincodeObj
	// OrgObj
	OrgID      string
	OrgAdmin   string
	OrgPeers   []string //组织节点
	TargetPeer string
	OrgAchor   string //通信描点
	UserName   string

	ChaincodeID      string
	ChaincodeVersion string
	GoPath           string
	ChaincodePath    string

	ConnectionProfile string
	channelID         string
	OrdererID         string
	OrdererEndpoint   string //orderer节点地址
	ChannelConfig     string
	Initialized       bool
	channelCli        *channel.Client
	resmgmtCli        *resmgmt.Client //资源管理客户端(相当于管理员admin)，用以创建或更新通道
	fsdk              *fabsdk.FabricSDK
	eventCli          *event.Client
}
type OrgObj struct {
	OrgID      string
	OrgAdmin   string
	OrgPeers   []string //组织节点
	TargetPeer string
	OrgAchor   string //通信描点
	UserName   string
}

type ChaincodeObj struct {
	ChaincodeID      string
	ChaincodeVersion string
	GoPath           string
	ChaincodePath    string
}

func (fs *FabricService) Initialize() error {
	if fs.Initialized {
		return errors.New("fab sdk 已经初始化")
	}
	fmt.Println("fabric sdk 初始化开始....")
	sdk, err := fabsdk.New(config.FromFile(fs.ConnectionProfile))
	if err != nil {
		return errors.WithMessage(err, "加载sdk配置文件失败！")
	}
	fs.fsdk = sdk
	fmt.Println("fabric sdk 初始化完成....")

	rmCtx := fs.fsdk.Context(fabsdk.WithUser(fs.OrgAdmin), fabsdk.WithOrg(fs.OrgID))
	if rmCtx == nil {
		return errors.New("资源管理器上下文创建失败！")
	}

	rmCli, err := resmgmt.New(rmCtx)
	if err != nil {
		return errors.WithMessage(err, "资源管理器客户端创建失败！")
	}
	fs.resmgmtCli = rmCli
	existed, err := ChannelIsExist(rmCli, fs.channelID, fs.TargetPeer)
	if err != nil {
		return errors.WithMessage(err, "查询通道是否存在失败")
	}
	if existed == true {
		fmt.Println("fabric 通道已存在....")
	} else {
		_, err := CreateChannel(fs.fsdk, fs.channelID, fs.ChannelConfig, fs.OrgAdmin, fs.OrgID, fs.OrdererID, fs.OrdererEndpoint)
		if err != nil {
			return errors.WithMessage(err, "创建通道失败！")
		}
		_, err = JoinChannel(rmCli, fs.channelID, fs.OrdererEndpoint, fs.OrgPeers)
		if err != nil {
			return errors.WithMessage(err, "加入通道失败！")
		}
	}
	fs.Initialized = true
	return nil
}

//安装并实例化链码
func (fs *FabricService) InstallAndInstantiateCC() error {
	//创建发送给peers的chaincode package
	ccPkg, err := packager.NewCCPackage(fs.ChaincodePath, fs.GoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("chaincode package created")
	ccHasInstalled := false
	//查询已安装的链码
	ccInstalledRes, err := fs.resmgmtCli.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(fs.TargetPeer))
	if err != nil {
		return errors.WithMessage(err, "failed to Query Installed chaincode")
	}
	if ccInstalledRes != nil {
		for _, cc := range ccInstalledRes.Chaincodes {
			if strings.EqualFold(cc.Name, fs.ChaincodeID) {
				ccHasInstalled = true
			}
		}
	}
	fmt.Println("ccHasInstall", ccHasInstalled)
	if !ccHasInstalled {
		//安装链码(智能合约)到org peers
		installCCReq := resmgmt.InstallCCRequest{Name: fs.ChaincodeID, Path: fs.ChaincodePath, Version: fs.ChaincodeVersion, Package: ccPkg}
		_, err = fs.resmgmtCli.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return errors.WithMessage(err, "failed to install chaincode")
		}
		fmt.Println("chaincode install success")
	} else {
		fmt.Println("Chaincode already exist")
	}
	ccHasInstantiated := false
	//查询已经实例化的链码
	ccInstantiatedResp, err := fs.resmgmtCli.QueryInstantiatedChaincodes(fs.channelID, resmgmt.WithTargetEndpoints(fs.OrgPeers...))
	if ccInstantiatedResp.Chaincodes != nil && len(ccInstantiatedResp.Chaincodes) > 0 {
		for _, chaincodeIno := range ccInstantiatedResp.Chaincodes {
			fmt.Println("chaincodexx", chaincodeIno)
			if strings.EqualFold(chaincodeIno.Name, fs.ChaincodeID) {
				ccHasInstantiated = true
			}
		}
	}
	// could not get chConfig cache reference:read configuration for channel peers failed
	fmt.Println("ccHasInstantiated", ccHasInstantiated)
	// Set up chaincode policy
	// ccPolicy := cauthdsl.SignedByAnyMember([]string{"fbi.citizens.com"})
	if !ccHasInstantiated {
		//msp名称，非域名
		ccPolicy := cauthdsl.SignedByAnyMember([]string{fs.OrgID})
		req := resmgmt.InstantiateCCRequest{Name: fs.ChaincodeID, Path: fs.GoPath, Version: fs.ChaincodeVersion, Policy: ccPolicy}

		resp, err := fs.resmgmtCli.InstantiateCC(fs.channelID, req)

		if err != nil || resp.TransactionID == "" {
			return errors.WithMessage(err, "failed to instantiate the chaincode")
		}
		fmt.Println("Chaincode instantiated successed;tx:", resp.TransactionID)
	} else {
		fmt.Println("chaincode has instantiated")
	}
	//channel Context用于查询和执行事务交易
	chCtx := fs.fsdk.ChannelContext(fs.channelID, fabsdk.WithUser(fs.OrgAdmin))
	fs.channelCli, err = channel.New(chCtx)
	if err != nil {
		return errors.WithMessage(err, "failed to create new Channel client")
	}
	fmt.Println("channel client created")
	//访问通道事件
	fs.eventCli, err = event.New(chCtx)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")
	fmt.Println("Chaincode Installation and Instantiation successful!")
	return nil
}

func (fs *FabricService) closeSdk() {
	fs.fsdk.Close()
}
