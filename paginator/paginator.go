package paginator

import (
	"errors"
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidOrder = errors.New("invalid order value")
)

// =====================
// Params
// =====================

type Config struct {
	Page     int
	PageSize int

	Sort  string
	Order string // asc / desc

	// 排序白名单（防 SQL 注入）
	SortWhitelist map[string]struct{}
}

// =====================
// Result
// =====================

type Result[T any] struct {
	List       []T   `json:"list"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPage  int   `json:"total_page"`
	TotalCount int64 `json:"total_count"`
}

// =====================
// Scope（核心扩展点）
// =====================

type Scope func(*gorm.DB) *gorm.DB

// =====================
// Paginate（核心方法）
// =====================

func Paginate[T any](
	db *gorm.DB,
	params Config,
	out *[]T,
	scopes ...Scope,
) (*Result[T], error) {

	params = normalize(params)

	db = applyScopes(db, scopes...)

	// count
	var total int64
	if err := db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, err
	}

	// sort safety check
	if !isSafeSort(params) {
		params.Sort = "id"
	}

	offset := (params.Page - 1) * params.PageSize

	query := db.
		Preload(clause.Associations).
		Order(params.Sort + " " + params.Order).
		Limit(params.PageSize).
		Offset(offset)

	if err := query.Find(out).Error; err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(params.PageSize)))
	if totalPage == 0 {
		totalPage = 1
	}

	return &Result[T]{
		List:       *out,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalCount: total,
		TotalPage:  totalPage,
	}, nil
}

// =====================
// Scope chain
// =====================

func applyScopes(db *gorm.DB, scopes ...Scope) *gorm.DB {
	for _, s := range scopes {
		db = s(db)
	}
	return db
}

// =====================
// normalize params
// =====================

func normalize(p Config) Config {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 200 {
		p.PageSize = 200
	}

	p.Order = strings.ToLower(p.Order)
	if p.Order != "asc" && p.Order != "desc" {
		p.Order = "desc"
	}

	if p.Sort == "" {
		p.Sort = "id"
	}

	return p
}

// =====================
// security check
// =====================

func isSafeSort(p Config) bool {
	if len(p.SortWhitelist) == 0 {
		return true
	}
	_, ok := p.SortWhitelist[p.Sort]
	return ok
}
