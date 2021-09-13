package StringsUtil

import (
	"helloworld-blockchain-go/util/StringUtil"
	"strings"
)

/*
 @author king 409060350@qq.com
*/

func HasDuplicateElement(datas *[]string) bool {
	visited := make(map[string]bool, 0)
	for i := 0; i < len(*datas); i++ {
		if visited[(*datas)[i]] == true {
			return true
		} else {
			visited[(*datas)[i]] = true
		}
	}
	return false
}

func Contains(values *[]string, value string) bool {
	if values != nil && len(*values) != 0 {
		for _, v := range *values {
			if v == value {
				return true
			}
		}
	}
	return false
}

func Split(values string, valueSeparator string) []string {
	if StringUtil.IsNullOrEmpty(values) {
		return []string{}
	}
	return strings.Split(values, valueSeparator)
}
