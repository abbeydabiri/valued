package functions

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func LoadLogFile() {
	fileName := fmt.Sprintf("/%s_%d_%d_%d.log", filepath.Base(os.Args[0]), time.Now().Year(), time.Now().Month(), time.Now().Day())
	filePath := fmt.Sprintf("log/%d/%d/%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	writeFile(fileName, filePath, []byte(``))

	logfile, err := os.OpenFile(filePath+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln("Failed to open log file", ":", err)
	}
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(logfile)
}

func LineCounter(ioReader io.Reader) (int, error) {
	byteBuffer := make([]byte, 8196)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := ioReader.Read(byteBuffer)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(byteBuffer[:c], lineSep)

		if err == io.EOF {
			break
		}
	}

	return count, nil
}

func SaveFile(fileName, filePath string, fileBytes []byte) string {

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filePath, 0777)
		} else {
			return ""
		}
	}

	filePathName := fmt.Sprintf("%s/%s", filePath, fileName)
	if filePath == "" {
		filePathName = fileName
	}

	if len(fileBytes) > 0 {
		file, err := os.Create(filePathName)
		defer file.Close()
		if err != nil {
			log.Println("Failed Create Error", ":", err)
			return ""
		}
		_, err = file.Write(fileBytes)

		if err != nil {
			log.Println("File Write Error: ", err)
			return ""
		}
	}

	return filePathName
}

func SaveImage(fileName, OSFilePath string, fileBytes []byte) string {
	fileName = fmt.Sprintf("/%s", fileName)
	filePath := fmt.Sprintf("images/%d/%d/%d", time.Now().Year(), time.Now().Month(), time.Now().Day())

	if writeFile(fileName, OSFilePath+filePath, fileBytes) {
		return filePath + fileName
	}

	return ""
}

func writeFile(fileName, filePath string, fileBytes []byte) bool {

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(filePath, 0777)
		} else {
			return false
		}
	}

	if len(fileBytes) > 0 {
		file, err := os.Create(filePath + fileName)
		defer file.Close()
		if err != nil {
			log.Println("Failed Create Error", ":", err)
			return false
		}
		_, err = file.Write(fileBytes)

		if err != nil {
			log.Println("File Write Error: ", err)
			return false
		}
	}
	return true
}

func UploadFile(aUploadFields []string, aUploadTypes map[string]string, httpReq *http.Request) (mapFiles map[string]interface{}, sMessage string) {

	mbSize := uint(2)
	mbUint := uint(1000000)
	mbInt64 := int64(1000000)
	mbShifter := uint(20)
	maxSize := mbSize << mbShifter

	mapFiles = make(map[string]interface{})
	for _, sFieldname := range aUploadFields {
		fileReader, fileHeader, err := httpReq.FormFile(sFieldname)
		if err == nil {
			fileBytes, _ := ioutil.ReadAll(fileReader)
			fileType := http.DetectContentType(fileBytes)

			// if httpReq.ContentLength > int64(maxSize) {
			lAbort := false
			if int64(len(fileBytes)) > int64(maxSize) {
				lAbort = true
				sMessage += fmt.Sprintf("Selected file=> %s (%vmb) <br> is larger than allowed limit (%vmb) <br>", fileHeader.Filename, httpReq.ContentLength/mbInt64, maxSize/mbUint)
			}

			aUploadTypesError := make(map[string]bool)
			for _, sFiletype := range aUploadTypes {
				if !strings.Contains(fileType, sFiletype) {
					aUploadTypesError[sFiletype] = true
				}
			}

			if len(aUploadTypesError) == len(aUploadTypes) {
				lAbort = true
				sMessage += fmt.Sprintf("File => %s (%s) is not allowed <br>", fileHeader.Filename, fileType)
			}

			if !lAbort {
				if len(fileBytes) > 0 {
					tempFile := make(map[string]interface{})
					tempFile["filetype"] = http.DetectContentType(fileBytes)
					tempFile["filename"] = fileHeader.Filename
					tempFile["filesize"] = len(fileBytes)
					tempFile["filebytes"] = fileBytes
					mapFiles[sFieldname] = tempFile
				} else {
					sMessage += fmt.Sprintf("File => %s is empty!! <br>", fileHeader.Filename)
				}
			}
		}
	}
	return
}
