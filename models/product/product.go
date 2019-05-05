package product
type Product struct {
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
	Number     string `json:"number"`    //产品编号
	MillPrice  string `json:"millPrice"` //出厂价格，不可改变
	Price      string `json:"price"`
	Addr      string `json:"addr"`
	Desc	 string `json:"desc"`
	Owner      string `json:"owner"`     //产品拥有者
	Productor  string `json:"productor"` //厂家
}