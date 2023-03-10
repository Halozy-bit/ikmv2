package helper

func MaxPage(maxPerPage, totalItem int) int {
	maxPage := totalItem / maxPerPage
	if totalItem%maxPerPage > 0 {
		maxPage += 1
	}
	return maxPage
}

// Count Total Product Next Page
// calculate the total product shown on the next page
// WARN
// dont use 0 value
func CountTtlProductNxtPage(page, total_product int) int {
	nextPage := (page - 1) * MaxProductPerPage
	nextPage = total_product - nextPage
	if nextPage > MaxProductPerPage {
		nextPage = MaxProductPerPage
	}

	return nextPage
}
