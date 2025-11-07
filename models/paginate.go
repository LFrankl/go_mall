package models

import (
	"baby/settings"
	"gorm.io/gorm"
)

func Paginate(db *gorm.DB, p int) (*gorm.DB, int, int, int, int) {
	//当前页数小于或者等于0，当前页数变成第一页
	if p <= 0 {
		p = 1
	}

	//计算所有数据总数和总页数
	var count int64
	db.Count(&count)
	PageCount := int(count) / settings.PageSize
	//存在余数
	if int(count)%settings.PageSize > 0 {
		PageCount++
	}

	//当前页数超出总页数，则当前页数变成总页数
	if p > PageCount {
		p = PageCount
	}
	previous := 1
	if p >= 0 {
		previous = p - 1
	}
	next := p + 1
	if next > PageCount {
		next = PageCount
	}
	//计算偏移量，用于数据查询
	offset := (p - 1) * settings.PageSize
	res := db.Offset(offset).Limit(settings.PageSize)
	return res, previous, next, int(count), PageCount
}
