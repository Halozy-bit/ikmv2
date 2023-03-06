package api

type CatalogGetProduct struct {
	Page     int
	Category string `query:"category"`
	LastID   string `json:"last_id"`
}
