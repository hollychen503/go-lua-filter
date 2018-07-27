package main

import (
	"fmt"
	"log"
	"time"

	"github.com/buger/jsonparser"
)

//func pure(data []byte) bool {

func pure(data string) bool { // 如果用string copy 呢？
	// There is `GetInt` and `GetBoolean` helpers if you exactly know key data type
	name, err := jsonparser.GetString([]byte(data), "person", "name", "first")
	if err != nil {
		panic(err)
	}
	if name != "holly" {
		return false
	}

	return true
}

func main() {

	data := []byte(`{
		"person": {
		  "name": {
			"first": "holly",
			"last": "Bugaev",
			"fullName": "Leonid Bugaev"
		  },
		  "github": {
			"handle": "buger",
			"followers": 109
		  },
		  "avatars": [
			{ "url": "https://avatars1.githubusercontent.com/u/14009?v=3&s=460", "type": "thumbnail" }
		  ]
		},
		"company": {
		  "name": "Acme"
		}
	  }`)

	imax := 1000
	jmax := 10000
	bgn := time.Now().UnixNano()

	for i := 0; i < imax; i++ {

		for j := 0; j < jmax; j++ {
			err := pure(string(data))
			if !err {
				log.Panic("not pure")
			}
		}
		log.Println(i)
	}
	end := time.Now().UnixNano()
	diff := end - bgn

	fmt.Println("total req:", imax*jmax)
	fmt.Println("used ", diff, "nano seconds, ", (float64(diff) / (float64)(time.Second)), "seconds")
	fmt.Println(float64(imax*jmax)/(float64(diff)/(float64)(time.Second)), " req per second")

}
