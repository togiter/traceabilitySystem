package fabricservice

func PostProduct() interface{} {
	return "PostProduct"
}

func QueryProductNo(id string) interface{} {
	return "QueryProductNo"
}

func QueryProductsRange(startKey, endKey string) interface{} {
	return "QueryProductsRange"
}
