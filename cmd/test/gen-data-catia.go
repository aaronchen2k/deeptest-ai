package main

import (
	"encoding/json"
	"flag"
	llmUtils "github.com/deeptest-com/deeptest-next/internal/pkg/libs/llm"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	name := ""

	flagSet := flag.NewFlagSet("deeptest", flag.ContinueOnError)
	flagSet.StringVar(&name, "n", "custom_catia.json", "")
	flagSet.Parse(os.Args[1:])

	number := 10
	mp := map[string][]string{
		"materials":  {"CR240*", "[CR270LA]", "aluminium,不锈钢", "CR240*,[CR270LA]"},
		"excel_path": {"D:\\Documents\\catiaVB\\TB.xlsx", "/home/aaron/excel_file.xlsx", "~/excel_file.xlsx", "~/excel.xlsx", "./excel.xlsx"},
	}

	workDir, _ := os.Getwd()
	srcPath := filepath.Join(workDir, "xdoc", "data", name)
	distPath := filepath.Join(workDir, "xdoc", "data", strings.Replace(name, ".json", "_out.json", -1))

	content := _file.ReadFileBuf(srcPath)

	var arr []llmUtils.InstructionItem
	json.Unmarshal(content, &arr)
	if len(arr) == 0 {
		return
	}

	inputSample := arr[0].Input
	outputSample := arr[0].Output
	var ret []llmUtils.InstructionItem

	for true {
		for _, item := range arr {
			if item.Input == "" {
				item.Input = inputSample
			}
			//item.Input = ""

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
