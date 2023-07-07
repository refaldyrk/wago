package command

import "strings"

func GetCommand(str string) string {
	if !strings.HasPrefix(str, "/") {
		return "standard"
	}
	str = strings.TrimPrefix(str, "/")
	return str
}
