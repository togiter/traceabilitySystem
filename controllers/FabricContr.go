package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/traceability-system/fabric/fabricservice"
)

type FabricContr struct {
	Ctx    iris.Context
	Fabric fabricservice.FabricService
}

func NewFabric() fabricservice.FabricService {
	fab := fabricservice.FabricService{
		// fabricservice.ChaincodeObj{
		ChaincodeID:      "productchaincode",
		ChaincodeVersion: "v0",
		GoPath:           os.Getenv("GOPATH"),
		ChaincodePath:    "github.com/traceability-system/fabric/chaincode/",
		// },
		// fabricservice.OrgObj{
		OrgID:      "Org1MSP",
		OrgAdmin:   "Admin",
		OrgPeers:   []string{"peer0.org1.example.com", "peer1.org1.example.com"}, //组织节点
		TargetPeer: "peer0.org1.example.com",
		OrgAchor:   "peer0.org1.example.com", //通信描点
		UserName:   "UserName",
		// },
		ChannelConfig: "github.com/traceability-system/fabric/configs/artifacts/traceability-system.tx",
		ConnectionProfile:os.Getenv("GOPATH") + "/src/github.com/traceability-system/fabric/configs/trace-sys.yaml",
	}
	err := fab.Initialize()
	if err != nil {
		fmt.Println(err)
	}
	return fab
}

//开发者可以在BeforeActivation方法中来处理请求定义
func (f *FabricContr) BeforeActivation(a mvc.BeforeActivation) {
	if(f.Fabric.Initialized == false){
		f.Fabric = NewFabric()
	}
	a.Handle("GET", "/info", "QueryInfo")
	//发布产品
	a.Handle("POST", "/postProduct", "PostProduct")
	//查询产品
	a.Handle("GET", "/queryProducts", "QueryProducts")
}

func (f *FabricContr) QueryProducts() interface{} {
	id := f.Ctx.URLParam("number")
	startKey := f.Ctx.URLParam("startKey")
	endKey := f.Ctx.URLParam("endKey")
	if len(startKey) > 0 && len(endKey) > 0 {
		result, _ := f.Fabric.QueryProductsRange(startKey, endKey)
		return result
	} else {
		result, _ := f.Fabric.QueryProductNo(id)
		return result
	}
	// return id
}

func (f *FabricContr) PostProduct() interface{} {
	name := f.Ctx.FormValue("name")
	productor := f.Ctx.FormValue("productor") //厂家
	addr := f.Ctx.FormValue("addr")           //产地
	id := f.Ctx.FormValue("number")
	desc := f.Ctx.FormValue("desc")
	millPrice := f.Ctx.FormValue("millPrice") //出厂价格
	price := f.Ctx.FormValue("price")
	owner := f.Ctx.FormValue("owner")
	aMap := map[string]string{
		"name":      name,
		"productor": productor,
		"addr":      addr,
		"id":        id,
		"millPrice": millPrice,
		"price":     price,
		"owner":     owner,
		"desc":      desc,
	}

	result, err := json.Marshal(aMap)
	if err != nil {
		return err
	}
	// return fabricservice.IssueProduct(name,productor,addr,id,millPrice,price,desc,owner)
	return result
}

func (f *FabricContr) QueryInfo() interface{} {
	return map[string]string{
		"name": "pikaqiu",
		"type": "animate",
		"age":  "1999",
	}
}

func (f *FabricContr) PostFabric() interface{} {
	name := f.Ctx.FormValue("name")
	return map[string]string{"name": name}
}
