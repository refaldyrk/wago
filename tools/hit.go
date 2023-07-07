package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Hit(urls string) string {
	// Mengirim permintaan GET ke API
	urlss := HitEndpointStringURL(urls)
	response, err := http.Get(urlss)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	// Membaca body respons
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(body)
}
