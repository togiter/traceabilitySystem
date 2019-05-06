/*
*1.安装链码
*2.实例化链码
*3.创建并加入通道
*4.更新配置
 */
package fabricservice

import (
	"fmt"

	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

const (
	adminUser       = "Admin"
	ordererOrg      = "OrderOrg"
	ordererEndpoint = ""
)

var orgExpectedPeers = map[string]int{
	"Org1": 2,
	"Org2": 2,
}

/*
*获取组织msp
*orgName 组织名称
*
 */
func orgMSPID(sdk *fabsdk.FabricSDK, orgName string) (string, error) {
	//获取sdk的配置
	configBackkend, err := sdk.Config()
	if err != nil {
		return "", errors.WithMessage(err, "获取sdk配置失败！")
	}
	//获取节点配置
	endpointConfig, err := fab.ConfigFromBackend(configBackkend)
	if err != nil {
		return "", errors.WithMessage(err, "获取端点配置失败")
	}

	mspID, ok := comm.MSPID(endpointConfig, orgName)
	if !ok {
		return "", errors.New("查找MSP ID 失败！")
	}

	return mspID, nil
}

/*
*获取多组织上下文
*
*
*
 */
// func orgContexts(sdk *fabsdk.FabricSDK, user fabsdk.ContextOption, orgNames []string) ([]*OrgContext, error) {
// 	orgCtxs := make([]*OrgContext, len(orgNames))
// 	for i, orgName := range orgNames {
// 		cliCtx := sdk.Context(user, fabsdk.WithOrg(orgName))
// 		resMgmt, err := resmgmt.New(cliCtx)
// 		if err != nil {
// 			return nil, errors.WithMessage(err, "创建资源管理器失败！！！")
// 		}
// 		//获取期望节点数
// 		expectedPeers, ok := orgExpectedPeers[i]
// 		if !ok {
// 			return nil, errors.WithMessage(err, "unknown org name")
// 		}
// 		peers, err := DiscoverLocalPeers(cliCtx, expectedPeers)
// 		if err != nil {
// 			return nil, errors.WithMessage(err, "本地节点未指定")
// 		}
// 		orgCtx := OrgContext{
// 			OrgID:       orgName,
// 			CtxProvider: cliCtx,
// 			ResMgmt:     resMgmt,
// 			Peers:       peers,
// 		}
// 		orgContexts[i] = &orgCtx
// 	}
// 	return orgCtxs, nil
// }

/*
*指定组织成员背书签名策略
*@param orgName 指定组织名称
*返回背书策略
 */
func prepareOrgPolicy(sdk *fabsdk.FabricSDK, orgName string) (string, error) {
	mspID, err := orgMSPID(sdk, orgName)
	if err != nil {
		return "", errors.WithMessage(err, "MSP ID could not be determined")
	}

	return fmt.Sprintf("AND('%s.member')", mspID), nil
}

/*
*生成链码（安装及初始化）
*
*
 */
func doChaincode(ccID, ccPath, ccVersion, goPath, chId, orgMSP string, orgResMgmt *resmgmt.Client) (bool, error) {
	ccPkg, err := packager.NewCCPackage(ccPath, goPath)
	if err != nil {
		return false, errors.WithMessage(err, "打包链码失败！")
	}
	//安装链码到指定peers
	installCCReq := resmgmt.InstallCCRequest{
		Name:    ccID,
		Path:    ccPath,
		Version: ccVersion,
		Package: ccPkg,
	}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return false, errors.WithMessage(err, "安装链码失败！")
	}

	//设置链码策略
	ccPolicy := cauthdsl.SignedByAnyMember([]string{orgMSP})
	//实例化链码
	resp, err := orgResMgmt.InstantiateCC(
		chId,
		resmgmt.InstantiateCCRequest{Name: ccID, Path: ccPath, Version: ccVersion, Args: nil, Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
	fmt.Println("实例化结果：", resp)
	if err != nil {
		return false, errors.WithMessage(err, "实例化链码失败！")
	}
	return true, nil
}

/*
*检查指定链码在指定通道上是否已经实例化
*@param channelID 通道名称
*@param ccName 链码名称
*@param ccVersion 链码版本
 */

func isCCInstantiated(resMgmt *resmgmt.Client, channelID, ccName, ccVersion string) (bool, error) {
	resp, err := resMgmt.QueryInstantiatedChaincodes(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return false, errors.WithMessage(err, "查询实例化链码失败！")
	}
	for _, chaincode := range resp.Chaincodes {
		fmt.Println("查询实例化链码", chaincode.Name, chaincode.Version)
		if chaincode.Name == ccName && chaincode.Version == ccVersion {
			return true, nil
		}
	}
	return false, nil
}

/*
*是否已经加入通道
*
*
 */

/**
*创建通道，orgId的admin在ordererEndpoint创建ch
*@param ch 通道名称 chPath 通道配置路径
*@param admin 创建者
*@param orgId 所在组织
*@param ordererEndpoint 排序节点地址
 */
func CreateChannel(sdk *fabsdk.FabricSDK, ch, chPath, admin, orgName, ordererOrg, ordererEndpoint string) (bool, error) {
	cliCtx := sdk.Context(fabsdk.WithUser(admin), fabsdk.WithOrg(ordererOrg))
	rmCli, err := resmgmt.New(cliCtx)
	if err != nil {
		return false, errors.WithMessage(err, "创建通道-资源管理器失败！")
	}
	//获取管理员身份
	mspCli, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		return false, errors.WithMessage(err, "获取msp客户端失败！！")
	}
	adminId, err := mspCli.GetSigningIdentity(admin)
	if err != nil {
		return false, errors.WithMessage(err, "获取管理员身份Id失败！！")
	}

	//创建通道请求
	chReq := resmgmt.SaveChannelRequest{
		ChannelID:         ch,
		ChannelConfigPath: chPath,
		SigningIdentities: []msp.SigningIdentity{adminId}}

	if _, err = rmCli.SaveChannel(chReq,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(ordererEndpoint)); err != nil {
		return false, errors.WithMessage(err, "创建通道失败！")
	}
	return true, nil
}

/**
*加入通道，
*@param ch 通道名称
*@param peers 要加入通道的节点
*@param admin管理员
*@param orgId 所在组织
*@param ordererEndpoint 排序节点地址
 */
func JoinChannel(sdk *fabsdk.FabricSDK, ch, admin, orgId, ordererEndpoint string, peers []string) (bool, error) {
	//准备客户端上下文
	cliCtx := sdk.Context(fabsdk.WithUser(admin), fabsdk.WithOrg(orgId))
	//创建资源管理器(join channels,install,instantiate,upgrade chaincodes)
	rmgmtCli, err := resmgmt.New(cliCtx)
	if err != nil {
		return false, errors.WithMessage(err, "join channel，创建资源管理器失败！")
	}
	if err := rmgmtCli.JoinChannel(
		ch,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithTargetEndpoints(peers...),
		resmgmt.WithOrdererEndpoint(ordererEndpoint)); err != nil {
		return false, errors.WithMessage(err, "加入通道失败！")
	}
	return true, nil
}
