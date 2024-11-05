package _str

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
