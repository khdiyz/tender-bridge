package helper

// returns limit and offset
func GetPagination(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return limit, (page - 1) * limit
}
