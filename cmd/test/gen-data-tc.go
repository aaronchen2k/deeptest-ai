package main

import (
	"encoding/json"
	llmUtils "github.com/deeptest-com/deeptest-next/internal/pkg/libs/llm"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	"os"
	"path/filepath"
)

func main() {
	number := 1000
	mp := map[string][]string{
		"part_name":    {"p_a06", "000288", "p_a06", "001", "Y", "abc"},
		"project_name": {"G58", "PRJ_01", "A8", "abc", "X", "p_01"},
	}

	workDir, _ := os.Getwd()
	srcPath := filepath.Join(workDir, "xdoc", "data", "custom_tc.json")
	distPath := filepath.Join(workDir, "xdoc", "data", "custom_tc_out.json")

	content := _file.ReadFileBuf(srcPath)

	var arr []llmUtils.InstructionItem
	json.Unmarshal(content, &arr)
	if len(arr) == 0 {
		return
	}

	outputSample := arr[0].Output
	var ret []llmUtils.InstructionItem

	for true {
		for _, item := range arr {
			if item.Output == "" {
				item.Output = outputSample
			}

			newItem := llmUtils.ReplaceFields(item, mp)

			ret = append(ret, newItem)

			if len(ret) >= number {
				goto BREAK
			}
		}
	}

BREAK:

	bytes, _ := json.MarshalIndent(ret, "", "    ")
	_file.WriteBytes(distPath, bytes)
}
