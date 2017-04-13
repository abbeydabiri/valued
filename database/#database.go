package database

import (
	"database/sql"
	// "github.com/cznic/ql"
	_ "github.com/lib/pq"
	"github.com/pmylund/go-cache"
	"golang.org/x/crypto/nacl/secretbox"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	// "regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	keySize   = 32
	nonceSize = 24
)

type Database struct {
	OS          string
	OSfilepath  string
	AndroidPipe chan string

	Init, PGSQL bool
	Cache       *cache.Cache
	Error       error
	Rows        *sql.Rows
	Sql         string
	Column      map[string]interface{}
	connection  *sql.DB
}

func (this *Database) SpaceRemove(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (this *Database) SpaceReplace(str string, pattern rune) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return pattern
		}
		return r
	}, str)
}

func (this *Database) SetSession(GOSESSID string, page string, mapPage map[string]interface{}, lExpires bool) {

	if this.GetSession(GOSESSID, page) == nil {
		this.Cache.Set(page+GOSESSID, make(map[string]interface{}), -1)
	}

	if lExpires {
		this.Cache.Replace(page+GOSESSID, mapPage, cache.DefaultExpiration)
	} else {
		this.Cache.Replace(page+GOSESSID, mapPage, -1)
	}
}

func (this *Database) GetSession(GOSESSID string, page string) map[string]interface{} {
	mapPageInterface, _ := this.Cache.Get(page + GOSESSID)

	if mapPageInterface == nil {
		return nil
	}

	mapPage := mapPageInterface.(map[string]interface{})
	return mapPage
}

func (this *Database) keyNounce() (key *[keySize]byte, nonce *[nonceSize]byte) {
	fullPath := filepath.Dir(os.Args[0])
	fullPath = this.SpaceRemove(fullPath)
	fullPath = strings.Replace(fullPath, "/", "", -1)
	fullPath = strings.Replace(fullPath, "\\", "", -1)

	fullPath = base64.StdEncoding.EncodeToString([]byte(fullPath))
	nPower := int(60 / len(fullPath))
	if len(fullPath) < 60 {
		nCount := 0
		for nPower > nCount {
			fullPath += fullPath
			nCount++
		}
		fullPath = fullPath[0:60]
	}

	key = new([keySize]byte)
	copy(key[:], []byte(fullPath[0:32])[:keySize])

	nonce = new([nonceSize]byte)
	copy(nonce[:], []byte(fullPath[0:32][0:24])[:nonceSize])

	return
}

func (this *Database) Encrypt(in []byte) (out []byte) {
	key, nonce := this.keyNounce()
	out = secretbox.Seal(out, in, nonce, key)
	return
}

func (this *Database) Decrypt(in []byte) (out []byte) {
	key, nonce := this.keyNounce()
	out, _ = secretbox.Open(out, in, nonce, key)
	return
}

func (this *Database) loadCache(appCache string) {
	this.Cache = cache.New(5*time.Minute, 30*time.Second)

	jsonBytesEncrypted, _ := ioutil.ReadFile(appCache)
	if len(jsonBytesEncrypted) != 0 {

		jsonBytes := this.Decrypt(jsonBytesEncrypted)
		dbCache := make(map[string]cache.Item, 500)

		if jsonBytes != nil {
			err := json.Unmarshal(jsonBytes, &dbCache)
			if err != nil {
				log.Println("Json Error in Cache File " + appCache)
				log.Println(err.Error())
			} else {
				this.Cache = cache.NewFrom(5*time.Minute, 30*time.Second, dbCache)
			}
		}
	}
}

func (this *Database) loadConfig() (dbConfig map[string]string) {
	dbConfig = make(map[string]string)
	databaseJson := this.OSfilepath + "data/conf/database"
	jsonBytes, _ := ioutil.ReadFile(databaseJson)
	if len(jsonBytes) != 0 {
		err := json.Unmarshal(jsonBytes, &dbConfig)
		if err != nil {
			log.Println("Json Error in Database Configuration File " + databaseJson)
			log.Println(err.Error())
			os.Exit(1)
		}
	} else {
		log.Println("Invalid Database Configuration File " + databaseJson)
		os.Exit(1)
	}
	return dbConfig
}

func (this *Database) Basename(s string) string {
	n := strings.LastIndexByte(s, '.')
	if n > 0 {
		return s[:n]
	}
	return s
}

func (this *Database) Disconnect() {

	cacheJsonBytes, err := json.MarshalIndent(this.Cache.Items(), "", "")
	if err != nil {
		println(err.Error())
	}

	cacheJsonBytesEncrypted := this.Encrypt(cacheJsonBytes)

	ioutil.WriteFile(this.OSfilepath+this.Basename(filepath.Base(os.Args[0]))+".tmp", cacheJsonBytesEncrypted, 777)

	log.Printf("Disconnecting Database...")
	this.connection.Close()
}

func (this *Database) Connect() ([]byte, bool) {
	//Before Connecting Check Images folder exists - create if not exists
	imagesFilepath := this.OSfilepath + "images/"
	_, err := os.Stat(imagesFilepath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(imagesFilepath, 0777)
		}
	}

	//Before Connecting load Cache file
	appname := this.Basename(filepath.Base(os.Args[0]))
	this.loadCache(this.OSfilepath + appname + ".tmp")

	//Check if PGSQL == true or false
	if this.PGSQL {
		dbConfig := this.loadConfig()
		postgresConn := "host=" + dbConfig["hostname"] + " port=" + dbConfig["port"] + " dbname=" + dbConfig["database"] + " user=" + dbConfig["username"] + " password=" + dbConfig["password"] + " sslmode=disable connect_timeout=1"
		this.connection, _ = sql.Open("postgres", postgresConn)

	} /*else {

		//this.Init = false
		appname_mdb := appname + ".mdb"
		dbBytes, _ := ioutil.ReadFile(this.OSfilepath + appname_mdb)
		if len(dbBytes) == 0 {
			this.Init = true
		}

		qlFileDB := this.OSfilepath + appname_mdb
		if this.OSfilepath == "" {
			realDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
			qlFileDB = filepath.Join(realDir, appname_mdb)
		}

		ql.RegisterDriver()
		this.connection, _ = sql.Open("ql", qlFileDB)
	}*/

	err = this.connection.Ping()
	log.Printf("Connecting Database..")

	if err == nil {
		return nil, true
	}

	if !this.PGSQL {
		log.Printf("Database File: " + this.OSfilepath + appname + ".mdb")
	}
	log.Printf("Database: Error: Connection Error!:" + err.Error())
	return []byte(err.Error()), false
}

func (this *Database) Query(cSql string) (mapRes map[string]interface{}, lSuccess bool) {
	mapRes = make(map[string]interface{})
	lSuccess = false

	if cSql == "" {
		log.Printf("Database: Error: Missing Query!:", cSql)
		return
	}

	err := this.connection.Ping()
	if err != nil {
		log.Printf("Database: Error: Connection Error!:" + err.Error())
		return
	}

	dbTx, errTx := this.connection.Begin()
	if errTx != nil {
		log.Printf("Database: Error: Transaction Error!:" + err.Error())
		return
	}

	switch {
	default:
		_, err := dbTx.Exec(cSql)
		if err != nil {
			log.Println("Query:" + cSql)
			log.Printf("Database: Error: Query Error!:", err.Error())
			dbTx.Rollback()
			return
		}
		dbTx.Commit()
		return

	case strings.HasPrefix(strings.ToLower(cSql), "select"):
		rows, err := dbTx.Query(cSql)

		if err != nil {
			log.Println("Query:" + cSql)
			log.Printf("Database: Error: Query Error!:", err.Error())
			dbTx.Rollback()
			return
		}

		defer rows.Close()
		if rows.Err() != nil {
			log.Printf("Database: Error: Row Error!:", rows.Err().Error())
			dbTx.Rollback()
			return
		}

		counter := 1
		for rows.Next() {

			mapRow := make(map[string]interface{})
			aColumns, _ := rows.Columns()

			aScanArgs := make([]interface{}, len(aColumns))
			aColValues := make([]interface{}, len(aColumns))

			for nColIndex := range aColValues {
				aScanArgs[nColIndex] = &aColValues[nColIndex]
			}

			rows.Scan(aScanArgs...)
			for nColIndex, aColData := range aColValues {
				if aColData != nil {
					switch aColData.(type) {
					case []byte:
						mapRow[aColumns[nColIndex]] = string(aColData.([]byte))
					case string:
						mapRow[aColumns[nColIndex]] = aColData.(string)
					case int:
						mapRow[aColumns[nColIndex]] = aColData.(int)
					case int64:
						mapRow[aColumns[nColIndex]] = aColData.(int64)
					case float64:
						mapRow[aColumns[nColIndex]] = aColData.(float64)
					}
				} else {
					mapRow[aColumns[nColIndex]] = ""
				}
			}
			mapRes[strconv.Itoa(counter)] = mapRow
			lSuccess = true
			counter++
		}

		dbTx.Commit()
		return
	}
}
