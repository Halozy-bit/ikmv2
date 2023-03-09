package api

import "github.com/ikmv2/backend/pkg/helper"

// Count Total Product Next Page
// calculate the total product shown on the next page
func CountTtlProductNxtPage(page, total_product int) int {
	nextPage := (page - 1) * helper.MaxProductPerPage
	nextPage = total_product - nextPage
	if nextPage > 16 {
		nextPage = helper.MaxProductPerPage
	}

	return nextPage
}
