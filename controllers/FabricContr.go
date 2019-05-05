package controllers

import (
	"encoding/json"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/traceability-system/fabric/fabricservice"
)

type FabricContr struct {
	Ctx iris.Context
}

//开发者可以在BeforeActivation方法中来处理请求定义
func (f *FabricContr) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("GET", "/info", "QueryInfo")
	//发布产品
	a.Handle("POST", "/postProduct", "PostProduct")
	//查询产品
	a.Handle("GET", "/queryProducts", "QueryProducts")
}

func (f *FabricContr) QueryProducts() interface{} {
	id := f.Ctx.URLParam("number")
	startKey = f.Ctx.URLParam("startKey")
	endKey := f.Ctx.URLParam("endKey")
	if(startKey && endKey && len(startKey) > 0 && len(endKey) > 0){
		return fabricservice.queryProductsRange(startkey,endKey)
	}else{
	   return fabricservice.queryProductNo(id)
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
	// return fabricservice.PostProduct(name,productor,addr,id,millPrice,price,desc,owner)
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
