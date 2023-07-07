package tools

import "fmt"

const (
	API_BASE = "https://api.vhtear.com/"
	API_KEY  = "e6355682Feb4d1Fe4f43Fe912aFea12c2d16a9c8VHTear"
)

func HitEndpointStringURL(param string) string {
	return fmt.Sprintf("%s%s&apikey=%s", API_BASE, param, API_KEY)
}
