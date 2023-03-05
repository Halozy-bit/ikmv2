package api

// Count Total Product Next Page
// calculate the total product shown on the next page
func CountTtlProductNxtPage(page, total_product int) int {
	nextPage := (page - 1) * MaxProductPerPage
	nextPage = total_product - nextPage
	if nextPage > 16 {
		nextPage = MaxProductPerPage
	}

	return nextPage
}
