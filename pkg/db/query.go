package db

// use this one on repository layer, if you need limit and offset
func ExtractPaginationValue(page, limit int) (offset int, rowPerPage int) {
	offset = 0
	rowPerPage = 10

	if limit > 0 {
		rowPerPage = limit
	}
	if page > 0 {
		offset = (page - 1) * rowPerPage
	}

	return offset, rowPerPage
}
