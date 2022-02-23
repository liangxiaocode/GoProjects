package cmd

import (
	"PassageOne/internal/word"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ModeLower = iota + 1 // 1
	ModeUpper
	ModeUnderLineToUpperCamelCase
	ModeUnderLineToLowerCamelCase
	ModeCamelCaseToUnderLine
)

var desc = strings.Join([]string{
	"转换命令如下:",
	"1.单词全部转化为大写",
	"2.单词全部转化为小写",
	"3.下划线转化为小写驼峰",
	"4.下划线转化为大写驼峰",
	"5.驼峰单词转化为下划线",
}, "\n")

var wordCmd = &cobra.Command{
	Use:   "word",
	Short: "单词格式转换",
	Long:  desc,
	Run: func(cmd *cobra.Command, args []string) {
		var content string
		switch mode {
		case ModeLower:
			content = word.ToLower(str)
		case ModeUpper:
			content = word.ToUpper(str)
		case ModeUnderLineToLowerCamelCase:
			content = word.UnderLineToLowerCamelCase(str)
		case ModeUnderLineToUpperCamelCase:
			content = word.UnderLineToUpperCamelCase(str)
		case ModeCamelCaseToUnderLine:
			content = word.CamelCaseToUnderLine(str)
		default:
			log.Fatalf("不支持该模式转换,请执行help word查看帮助提示")
		}
		log.Printf("输出转换后结果:%s\n", content)
	},
}

var str string
var mode int8

// init函数在编译的时候就会调用
func init() {
	wordCmd.Flags().StringVarP(&str, "str", "s", "", "请输入单词内容")
	wordCmd.Flags().Int8VarP(&mode, "mode", "m", 0, "请输入单词内容")
}
