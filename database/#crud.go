package database

import (
	//
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Crud struct{}

func (this Crud) ReadFile(filename string) []byte {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	return file
}

func (this Crud) GetSystemTime() string {
	return time.Now().Format("02/01/2006 15:04:05 MST")
	// return time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")
}

func (this Crud) RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func (this Crud) Workflows() []string {
	return []string{"draft", "active", "inactive"}
}

func (this Crud) Fields() map[string]interface{} {
	fields := make(map[string]interface{})

	fields["code"] = ""
	fields["title"] = ""
	fields["description"] = ""

	fields["control"] = ""
	fields["workflow"] = ""

	fields["createdby"] = ""
	fields["updatedby"] = ""

	fields["createdate"] = ""
	fields["updatedate"] = ""

	return fields
}

func (this Crud) InitializeXdoc(cTablename string, aFields map[string]interface{}, curdb Database) (map[string]interface{}, bool) {

	// cSQL := "DROP TABLE %s cascade;"
	cSQL := "DROP TABLE %s ;"
	cSQL = fmt.Sprintf(cSQL, cTablename)
	curdb.Query(cSQL)

	sqlCreateIndex := ""
	sqlCreateTable := "CREATE TABLE " + cTablename + " ("

	sqlCreateIndexFormat := "\n CREATE INDEX %s_%s_%s on %s (%s);"

	for cFieldname, iFieldvalue := range aFields {
		fieldType := reflect.TypeOf(iFieldvalue)

		switch fieldType.Kind() {
		default:
			sqlCreateTable += cFieldname + " bytea, "

		case reflect.Bool:
			sqlCreateTable += cFieldname + " bool, "

		case reflect.String:
			sqlCreateTable += cFieldname + " text, "

		case reflect.Int, reflect.Int64:
			sqlCreateTable += cFieldname + " int, "

		case reflect.Float32, reflect.Float64:
			sqlCreateTable += cFieldname + " float8, "
		}

		sqlCreateIndex += fmt.Sprintf(sqlCreateIndexFormat, cTablename, cFieldname, this.RandomString(3), cTablename, cFieldname)
	}
	sqlCreateTable = strings.TrimSuffix(sqlCreateTable, ", ") + ");"

	// println(sqlCreateTable + sqlCreateIndex)
	return curdb.Query(sqlCreateTable + sqlCreateIndex)
}

func (this Crud) WriteXdoc(cUsername string, cTablename string, aTemplate map[string]interface{}, aFields map[string]interface{}, curdb Database) string {

	if aFields["control"] == nil || aFields["control"].(string) == "" {
		aFields["createdby"] = cUsername
		aFields["createdate"] = this.GetSystemTime()

		aFields["updatedby"] = aFields["createdby"]
		aFields["updatedate"] = aFields["createdate"]

		aFields["control"] = this.NextControl(cTablename, curdb)
		if aFields["workflow"] == nil || aFields["workflow"].(string) == "" {
			aFields["workflow"] = "draft"
		}

		if aFields["code"] == nil {
			aFields["code"] = ""
		}

	} else {
		aFields["updatedby"] = cUsername
		aFields["updatedate"] = this.GetSystemTime()
	}

	curdb.Query(this.updateXdoc(cTablename, aTemplate, aFields))
	return aFields["control"].(string)
}

func (this Crud) updateXdoc(cTablename string, aTemplate map[string]interface{}, aFields map[string]interface{}) string {
	sqlUpdateFields := ""
	for cFieldname, iFieldvalue := range aFields {
		if iFieldvalue != nil {
			if aTemplate[cFieldname] != nil {
				fieldType := reflect.TypeOf(iFieldvalue)
				switch fieldType.Kind() {
				default:
					sqlUpdateFields += fmt.Sprintf(`%s = '%v', `, cFieldname, iFieldvalue)

				case reflect.String:
					sqlUpdateFields += fmt.Sprintf(`%s = '%s', `, cFieldname, iFieldvalue)

				case reflect.Int, reflect.Int64, reflect.Float32, reflect.Float64:
					sqlUpdateFields += fmt.Sprintf(`%s = %v, `, cFieldname, iFieldvalue)
				}
			} else {
				log.Println(`Table->'` + cTablename + `' field->'` + cFieldname + `' does not exist!`)
			}
		}
	}

	if sqlUpdateFields == "" {
		return ""
	}

	sqlUpdateFields = strings.TrimSuffix(sqlUpdateFields, ", ")
	sqlUpdateTable := `update  %s set %s where control = '%s' `
	return fmt.Sprintf(sqlUpdateTable, cTablename, sqlUpdateFields, aFields["control"].(string))
}

func (this Crud) Delete(cTablename string, nControl int64, curdb Database) string {
	sqlDelete := `delete from %s where control = '%s' `

	return fmt.Sprintf(sqlDelete, cTablename, nControl)
}

func (this Crud) NextControl(cTablename string, curdb Database) string {

	sControl := "1.00000000000"
	cSql := fmt.Sprintf("SELECT control FROM %s ORDER BY control DESC LIMIT 1", cTablename)
	mapRes, _ := curdb.Query(cSql)

	if mapRes["1"] != nil {
		sControl = mapRes["1"].(map[string]interface{})["control"].(string)
	}

	sControl = this.incrementControl(sControl)

	cSql = fmt.Sprintf(`insert into %s (control) values ('%s')`, cTablename, sControl)
	curdb.Query(cSql)
	return sControl
}

func (this Crud) incrementControl(sControl string) string {

	nPointerpos := strings.Index(sControl, ".")
	sControl = strings.Replace(sControl, ".", "", 1)
	nControlInt, _ := strconv.ParseInt(sControl, 0, 64)

	nControlInt++
	sControl = strconv.FormatInt(nControlInt, 10)
	for len(sControl) < 12 {
		sControl += "0"
	}

	sControl = sControl[0:nPointerpos] + "." + sControl[nPointerpos:]
	return sControl
}
