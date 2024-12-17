package _consts

import "os"

const (
	FilePthSep = string(os.PathSeparator)
)

var (
	SortMap = map[string]string{
		"ascend":  "ASC",
		"descend": "DESC",
		"":        "ASC",
	}
)
