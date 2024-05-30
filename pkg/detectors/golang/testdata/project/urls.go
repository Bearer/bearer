package main

import (
	"fmt"
	"log"
	"os"
)

var m = map[string]string{
	"ignore.domain.com":  `raw.example.com`,
	"ignore2.domain.com": "string.example.com",
}

func main() {
	url := "http://find.example.com"
	log.Println(url)

	url2 := "http://" + "simple." + "concat.example.com"

	hostname := m["ignore.domain.com"]
	log.Println(hostname)

	// TEST: nothing detected for plain function call
	f()

	// TEST: Sprintf
	supplierId := 1234
	url3 := fmt.Sprintf(
		"https://%s.example.com/%s"+"/api/shop/supplier/%d/%%5B",
		os.Getenv("DOMAIN_PREFIX"),
		"path",
		supplierId
	)
}
