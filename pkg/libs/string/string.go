package _str

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"regexp"
	"strconv"
	"strings"
)

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss [][]string) zap.Field {
	return zap.Array(key, StringsArray(ss))
}

type StringsArray [][]string

func (ss StringsArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		for ii := range ss[i] {
			arr.AppendString(ss[i][ii])
		}
	}
	return nil
}

func Join(strs ...string) string {
	var builder strings.Builder
	if len(strs) == 0 {
		return ""
	}
	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}

func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func FindInArr(val string, arr []string) bool {
	for _, i := range arr {
		if val == i {
			return true
		}
	}

	return false
}
