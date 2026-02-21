package common

type Meta struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type MetaWithCount struct {
	Meta
	Total int `json:"total,omitempty"`
}

func NewMeta(paramPage, paramLimit int) Meta {
	var page, limit int
	if paramPage != 0 && paramLimit != 0 {
		page = paramPage
		limit = paramLimit
	}

	return Meta{
		Page:  page,
		Limit: limit,
	}
}

func NewMetaWithCount(page, limit, total int) MetaWithCount {
	meta := NewMeta(page, limit)

	return MetaWithCount{
		Meta:  meta,
		Total: total,
	}
}
