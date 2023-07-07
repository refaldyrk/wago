package api

import (
	"fmt"
	"wago/tools"

	"github.com/tidwall/gjson"
)

func Wikipedia(query string) string {
	body := tools.Hit(fmt.Sprintf("wikipedia?query=%s", query))
	value := gjson.Get(body, "result.Info")
	return value.String()
}
