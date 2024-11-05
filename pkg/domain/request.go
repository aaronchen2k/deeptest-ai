package _domain

import (
	"regexp"
	"strings"
)

var (
	SortMap = map[string]string{
		"ascend":  "ASC",
		"descend": "DESC",
		"":        "ASC",
	}
)

type PageReq struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Field    string `json:"field"`
	Order    string `json:"order"`
}

func (r *PageReq) ConvertParams() {
	r.Field = SnakeCase(r.Field)
	r.Order = SortMap[r.Order]
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
