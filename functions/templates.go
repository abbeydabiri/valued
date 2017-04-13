package functions

import (
	"valued/data"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	cDir = "html/"
)

type Templates struct {
}

func (this Templates) Alert(cType string, cMessage string) string {
	jsonBytesModal := []byte(`{
				"alert" : {	
						"type" : "` + cType + `",
						"message" : "` + cMessage + `"
					}
			}`)
	mapModalAlert := make(map[string]interface{})
	json.Unmarshal(jsonBytesModal, &mapModalAlert)
	return `"alert" :` + strconv.Quote(string(this.Generate(mapModalAlert, nil)))
}

func (this Templates) ModalAlert(cTitle string, cType string, cMessage string) string {
	jsonBytesModal := []byte(`{
				"alert-modal" : {	
						"title" : "` + cTitle + `",
						"type" : "` + cType + `",
						"message" : "` + cMessage + `"
					}
			}`)
	mapModalAlert := make(map[string]interface{})
	json.Unmarshal(jsonBytesModal, &mapModalAlert)
	return `"modal" :` + strconv.Quote(string(this.Generate(mapModalAlert, nil)))
}

func (this Templates) Generate(jsonMap map[string]interface{}, jsonMapReplace map[string]interface{}) []byte {
	htmlBytes := []byte(``)
	switch {
	case jsonMap["template"] != nil:
		htmlBytes = this.GetTemplate(jsonMap["template"].(string))

	default:
		aSortedMap := this.SortMap(jsonMap)
		for _, cComponentname := range aSortedMap {
			aComponent := jsonMap[cComponentname]

			if jsonMapReplace != nil {
				if jsonMapReplace[cComponentname] != nil {
					aComponent = jsonMapReplace[cComponentname]
				}
			}

			htmlBytes = append(htmlBytes, this.buildComponent(cComponentname, aComponent, jsonMapReplace)...)
		}
		emptyTags, _ := regexp.Compile(`\[@.\S*?@\]`)
		htmlBytes = []byte(emptyTags.ReplaceAllString(string(htmlBytes), ""))
	}
	return htmlBytes
}

func (this Templates) stripHashtag(cString string) string {
	nHashpos := strings.Index(cString, "#")
	if nHashpos != -1 {
		cString = cString[nHashpos+1:]
	}
	return cString
}

func (this Templates) buildComponent(cComponent string, aComponent interface{}, jsonMapReplace interface{}) []byte {
	filename := cDir + this.stripHashtag(cComponent)

	componentBytes := []byte(``)
	componentType := reflect.TypeOf(aComponent)
	switch componentType.Kind() {
	case reflect.Map:

		aComponentFields := make(map[string]interface{})
		aSortedMap := this.SortMap(aComponent.(map[string]interface{}))

		// if strings.Contains(cComponent, "#app-reward-list") {
		// 	fmt.Printf("filename: %s == %s \n", filename, cComponent)
		// }

		for _, cFieldname := range aSortedMap {

			iFieldvalue := aComponent.(map[string]interface{})[cFieldname]

			if jsonMapReplace != nil {
				if jsonMapReplace.(map[string]interface{})[cFieldname] != nil {
					iFieldvalue = jsonMapReplace.(map[string]interface{})[cFieldname]
				}
			}

			if iFieldvalue == nil {
				log.Printf("Field: %s is nil\n", cFieldname)
				continue
			}
			// log.Println(cFieldname)
			fieldType := reflect.TypeOf(iFieldvalue)
			switch fieldType.Kind() {
			case reflect.Map:
				iFieldvalue = string(this.buildComponent(cFieldname, iFieldvalue, jsonMapReplace))

				/*case reflect.String:
				if jsonMapReplace != nil {
					if cFieldname == "value" {
						if jsonMapReplace.(map[string]interface{})[iFieldvalue.(string)] != nil {
							iFieldvalue = jsonMapReplace.(map[string]interface{})[iFieldvalue.(string)]
						} else {
							iFieldvalue = "value"
						}
					}
				}*/
			}

			cFieldnameStripped := this.stripHashtag(cFieldname)
			if aComponentFields[cFieldnameStripped] != nil {
				aComponentFields[cFieldnameStripped] = aComponentFields[cFieldnameStripped].(string) + iFieldvalue.(string)
			} else {

				aComponentFields[cFieldnameStripped] = fmt.Sprintf("%v", iFieldvalue)
			}
		}

		componentBytes, _ = data.Asset(filename)

		// componentBytes = this.ReadFile(filename)
		if componentBytes != nil {
			componentString := string(componentBytes)

			for cFieldname, iFieldvalue := range aComponentFields {

				tagsToReplace, _ := regexp.Compile(`\[@` + cFieldname + `@\]`)
				if cFieldname != iFieldvalue.(string) {
					componentString = tagsToReplace.ReplaceAllString(componentString, iFieldvalue.(string))
				} else {
					componentString = tagsToReplace.ReplaceAllString(componentString, "")
				}
			}
			return []byte(componentString)
		}
	}

	return nil
}

func (this Templates) parseTemplate(cComponent string, aTemplate map[string]interface{}) map[string]interface{} {

	aComponentFields := aTemplate[cComponent]
	componentType := reflect.TypeOf(aComponentFields)
	switch componentType.Kind() {
	case reflect.Map:
		for cComponentname := range aComponentFields.(map[string]interface{}) {
			if aTemplate[cComponentname] != nil {
				aComponentFields.(map[string]interface{})[cComponentname] = this.parseTemplate(cComponentname, aTemplate)
			}
		}
	}

	return aComponentFields.(map[string]interface{})
}

func (this Templates) GetTemplate(template string) []byte {

	tplFiles, err := data.AssetDir("html")
	if err != nil {
		log.Printf("Template: Error: GetTemplate Json Error!:" + err.Error())
		return nil
	}

	aPage := make(map[string]interface{})
	aTemplate := make(map[string]interface{})

	for _, fileName := range tplFiles {
		aFields := make(map[string]interface{})

		fileContent, _ := data.Asset(cDir + fileName)
		fileName = strings.Replace(fileName, ".tpl", "", -1)

		fieldTags, _ := regexp.Compile(`\[|@|\]`)
		fileFields, _ := regexp.Compile(`\[@.\S*?@\]`)
		for _, match := range fileFields.FindAllString(string(fileContent), -1) {
			result := fieldTags.ReplaceAllString(match, "")
			aFields[result] = result
		}
		aTemplate[fileName] = aFields
	}

	if aTemplate[template] != nil {
		aComponent := make(map[string]interface{})
		for cComponent := range aTemplate[template].(map[string]interface{}) {
			if aTemplate[cComponent] != nil {
				aComponent[cComponent] = this.parseTemplate(cComponent, aTemplate)
			} else {
				aComponent[cComponent] = cComponent
			}
		}
		aPage[template] = aComponent
	}

	jsonBytes, err := json.MarshalIndent(aPage, "", "    ")
	if err != nil {
		log.Printf("Template: Error: GetTemplate Json Error!:" + err.Error())
		return nil
	}
	return jsonBytes
}

func (this Templates) ReadFile(filename string) []byte {
	// println("Reading File -> " + filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return file
}

func (this Templates) SortMap(mapUnsorted map[string]interface{}) []string {
	var aKeys []string
	for cFieldname, iFieldvalue := range mapUnsorted {
		fieldType := reflect.TypeOf(iFieldvalue)
		if fieldType != nil {
			switch fieldType.Kind() {
			case reflect.Map:
				iFieldvalue = this.SortMap(iFieldvalue.(map[string]interface{}))
			}
		}
		aKeys = append(aKeys, cFieldname)
	}

	//// sort.Strings(aKeys)
	// Strings(aKeys)

	name_number := func(name1, name2 string) bool {
		return getNumber(name1) < getNumber(name2)
	}

	By(name_number).Sort(aKeys)
	return aKeys
}
