package fabricservice
import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/traceability-system/models/product"
	"time"
)

//发布产品

func (fs *FabricService) IssueProduct(name, number,addr, desc, millPrice, price, productor ,owner string) (string, error) {
	if len(number) <= 0 || len(name) <= 0 || len(millPrice) <= 0 || len(price) <= 0 || len(addr) <= 0 || len(owner) <= 0 || len(productor) <= 0 {
		return "", fmt.Errorf("args input error!")
	}
	product := product.Product{
		ObjectType: "product",
		Name:       name,
		Number:     number,
		MillPrice:  millPrice,
		Price:      price,
		Addr:      addr,
		Desc:	desc,
		Owner:      owner,
		Productor:  productor,
	}

	productBytes, err := json.Marshal(product)
	if err != nil || productBytes == nil {
		return "", fmt.Errorf("product marshal failed!")
	}

	eventID := "eventOfIssueProduct!"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in IssueProduct")
	reg, noti, err := fs.eventCli.RegisterChaincodeEvent(fs.ChaincodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("event register error!")
	}
	defer fs.eventCli.Unregister(reg)
	fmt.Println("productBytes",productBytes)
	var params [][]byte 
	// params = append(params,[]byte("IssueProduct"))
	params = append(params,[]byte(number))
	params = append(params,productBytes)
	resp, err := fs.channelCli.Execute(channel.Request{
		ChaincodeID:  fs.ChaincodeID,
		Fcn:          "issueProduct",
		Args:         params,
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to invoke issueProduct:%v", err)
	}

	//等待事件回调
	select {
	case ccEvent := <-noti:
		fmt.Printf("Received CC event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Printf("did NOT receive CC event for eventId(%s)\n", eventID)
	}
	return string(resp.TransactionID), nil
}

//产品转移
func (fs *FabricService) TransferProduct(nOwner string, number string, price string) (string, error) {
	if len(nOwner) <= 0 || len(number) <= 0 {
		return "", fmt.Errorf("args input error!")
	}
	eventID := "eventOfTransfer!"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in tranfer")
	reg, noti, err := fs.eventCli.RegisterChaincodeEvent(fs.ChaincodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("event register error!")
	}
	defer fs.eventCli.Unregister(reg)

	resp, err := fs.channelCli.Execute(channel.Request{
		ChaincodeID:  fs.ChaincodeID,
		Fcn:          "TransferProduct",
		Args:         [][]byte{[]byte(nOwner), []byte(number), []byte(price)},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to invoke TransferProduct:%v", err)
	}

	//等待事件回调
	select {
	case ccEvent := <-noti:
		fmt.Printf("Received CC event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Printf("did NOT receive CC event for eventId(%s)", eventID)
	}
	return string(resp.TransactionID), nil
}

//修改价格
func (fs *FabricService) AlterProductPrice(owner string, number string, price string) (string, error) {
	if len(owner) <= 0 || len(number) <= 0 || len(price) <= 0 {
		return "", fmt.Errorf("args input error!")
	}
	// var args []string
	// args = append(args,"AlterProductPrice")
	// args = append(args,owner)
	// args = append(args,number)
	// args = append(args,price)
	eventID := "eventOfAlterPrice"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in alter price")
	reg, noti, err := fs.eventCli.RegisterChaincodeEvent(fs.ChaincodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("event register error!")
	}
	defer fs.eventCli.Unregister(reg)

	resp, err := fs.channelCli.Execute(channel.Request{
		ChaincodeID:  fs.ChaincodeID,
		Fcn:          "AlterProductPrice",
		Args:         [][]byte{[]byte(owner), []byte(number), []byte(price)},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to invoke AlterProductPrice:%v", err)
	}

	//等待事件回调
	select {
	case ccEvent := <-noti:
		fmt.Printf("Received CC event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Printf("did NOT receive CC event for eventId(%s)", eventID)
	}
	return string(resp.TransactionID), nil
}

//查询指定范围的产品
func (fs *FabricService) QueryProductsRange(startKey string, endKey string) ([]byte, error) {
	if len(startKey) <= 0 || len(endKey) <= 0 {
		return nil, fmt.Errorf("参数有误！")
	}
	resp, err := fs.channelCli.Query(channel.Request{
		ChaincodeID: fs.ChaincodeID,
		Fcn:         "QueryProductRange",
		Args:        [][]byte{[]byte(startKey), []byte(endKey)},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to range query:%v", err)
	}
	fmt.Println("range query:", resp.Payload)
	//products := []Product{}

	return resp.Payload, nil
}

func (fs *FabricService) QueryProductNo(productNo string) ([]byte, error) {
	if len(productNo) <= 0 {
		return nil, fmt.Errorf("参数有误！")
	}
	resp, err := fs.channelCli.Query(channel.Request{
		ChaincodeID: fs.ChaincodeID,
		Fcn:         "QueryProductNo",
		Args:        [][]byte{[]byte(productNo)},
	})
	if err != nil {
		return nil,fmt.Errorf("failed to query:%v", err)
	}

	fmt.Println("productNo query resp:", resp.Payload)
	return resp.Payload, nil
}
