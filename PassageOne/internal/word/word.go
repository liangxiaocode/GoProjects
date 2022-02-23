package word

import (
	"strings"
	"unicode"
)

// 封装一层：小->大
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// 封装一层：大->小
func ToLower(s string) string {
	return strings.ToLower(s)
}

func UnderLineToUpperCamelCase(s string) string {
	// 用空格替换里面的下划线,-1表示替换所用"_"
	s = strings.Replace(s, "_", " ", -1)
	// 将首字母改成大写
	s = strings.Title(s)
	// 用空字符串替换空格
	return strings.Replace(s, " ", "", -1)
}

func UnderLineToLowerCamelCase(s string) string {
	s = UnderLineToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// 驼峰单词转下划线
func CamelCaseToUnderLine(s string) string {
	// 用rune处理保证汉字的输入
	var output []rune
	for i, r := range s {
		if i == 0 {
			// 首字母需要改成小写
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}
