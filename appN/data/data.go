package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

var OS = ""
var ROOT = ""
var AssetRequest chan string
var AssetResponse chan []byte

func Asset(name string) (assetByte []byte, assetError error) {
	assetByte = nil
	assetError = nil

	if strings.HasSuffix(name, "/") {
		assetError = fmt.Errorf("Directory Listing Forbidden!!!")
	} else {
		nameSlice := strings.Split(name, "/")
		switch nameSlice[0] {
		default:
			assetError = fmt.Errorf("Directory %s is forbidden", nameSlice[0])

		case "conf", "files", "html", "email":
			filename := "data/" + name

			switch OS {
			case "ios":
			case "android":
				AssetRequest <- `{"action":"asset","filename":"` + filename + `"}`
				assetByte = <-AssetResponse
				if len(assetByte) == 0 {
					assetError = fmt.Errorf("File %s not found/is empty", filename)
				}

			default:
				assetByte, assetError = ioutil.ReadFile(filename)
			}
		}
	}

	return
}

func AssetDir(name string) (assetString []string, assetError error) {
	var assetByte []byte
	assetError = nil
	assetString = nil

	switch name {
	default:
		assetError = fmt.Errorf("Directory %s is forbidden", name)

	case "email", "html":
		filedir := "data/" + name
		switch OS {
		case "ios":
		case "android":
			AssetRequest <- `{"action":"assetDir","filedir":"` + filedir + `"}`
			assetByte = <-AssetResponse
			if len(assetByte) == 0 {
				assetError = fmt.Errorf("File %s not found/is empty", filedir)
			} else {
				assetMap := make(map[string]string)
				err := json.Unmarshal(assetByte, &assetMap)
				if err == nil {
					counter := 0
					assetString = make([]string, len(assetMap))
					for _, fileName := range assetMap {
						assetString[counter] = fileName
						counter++
					}
				}
			}

		default:
			fileInfos, err := ioutil.ReadDir(filedir)
			assetError = err

			assetString = make([]string, len(fileInfos))
			for counter, file := range fileInfos {
				assetString[counter] = file.Name()
			}
		}
	}

	return
}
