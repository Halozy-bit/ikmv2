package helper

func MaxPage(maxPerPage, totalItem int) int {
	maxPage := totalItem / maxPerPage
	if totalItem%maxPerPage > 0 {
		maxPage += 1
	}
	return maxPage
}
