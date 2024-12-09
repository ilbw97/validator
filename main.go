package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type TestData struct {
	Host           string `json:"host"`
	RequiredScheme bool   `json:"required_scheme"`
	Path           string `json:"path"`
}

var (
	mux       sync.Mutex
	testDatas = []TestData{}
)

func loadData() error {
	mux.Lock()
	defer mux.Unlock()

	// JSON 파일을 읽어오기
	fname := "./test.json"
	jsonData, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
		return err
	}

	testDatas = []TestData{}
	if err := json.Unmarshal(jsonData, &testDatas); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {

	err := loadData()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, testData := range testDatas {
		fmt.Println("-------------------------------------------------------------------------")
		if testData.Host != "" {
			if domainRes := DomainValidator(testData.Host); domainRes {
				fmt.Printf("DomainValidator said %s is valid domain \n\n", testData.Host)
			} else {
				fmt.Printf("DomainValidator said %s is invalid domain \n\n", testData.Host)
			}

			res := URIValidator(testData.Host, testData.RequiredScheme)
			if res {
				fmt.Printf("%v is valid url \n", testData.Host)
				isSub, err := IsSubdomain(testData.Host)
				if err != nil {
					fmt.Println(err)
				}
				if isSub {
					fmt.Printf("%v is sub domain \n", testData.Host)
				} else {
					fmt.Printf("%v is root domain \n", testData.Host)
				}
			} else {
				fmt.Printf("%v is invalid url \n", testData.Host)
			}
		}
		if testData.Path != "" {
			if isValid := PathValidator(testData.Path); isValid {
				fmt.Printf("%v is valid path \n", testData.Host)
			} else {
				fmt.Printf("%v is invalid path \n", testData.Host)
			}
		}
		fmt.Println("-------------------------------------------------------------------------")
	}
}
