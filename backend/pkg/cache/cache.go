package cache

const (
	// string from primitive.ObjectId
	TopCatalog    string = "top_catalog_id"
	BottomCatalog string = "bottom_catalog_id"

	CatalogNumberOfProduct string = "total_product"
)

func init() {
	Store(TopCatalog, "")
	Store(BottomCatalog, "")
}

var list = make(map[string]interface{})

func Store(key string, value interface{}) {
	list[key] = value
}

func Get(key string) interface{} {
	return list[key]
}
