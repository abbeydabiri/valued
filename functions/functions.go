package functions

import (
	"fmt"

	"math"
	"math/rand"

	"html"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	//httppostjson
	"bytes"
	"io/ioutil"
	"net/http"

	"regexp"
)

func CamelCase(word string) string {
	return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
}

func ThousandSeperator(num float64) string {
	numString := fmt.Sprintf("%v", num)
	re := regexp.MustCompile("(\\d+)(\\d{3})")
	for {
		formatted := re.ReplaceAllString(numString, "$1,$2")
		if formatted == numString {
			return formatted
		}
		numString = formatted
	}
}

func Round(input float64) float64 {
	if input < 0 {
		return math.Ceil(input - 0.5)
	}
	return math.Floor(input + 0.5)
}

func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}

func RoundDown(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Floor(digit)
	newVal = round / pow
	return
}

func HttpPostJSON(jsonURL string, jsonStr []byte) []byte {
	if len(jsonStr) == 0 {
		return nil
	}

	httpReq, _ := http.NewRequest("POST", jsonURL, bytes.NewBuffer(jsonStr))
	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Content-Length", strconv.Itoa(len(jsonStr)))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		println("HttpPostJSON error: " + err.Error())
		return nil
	}

	resBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return resBody
}

func TrimEscape(value string) string {
	return strings.TrimSpace(html.EscapeString(value))
}

func ReverseString(value string) string {
	// Convert string to rune slice.
	// ... This method works on the level of runes, not bytes.
	data := []rune(value)
	result := []rune{}

	// Add runes in reverse order.
	for i := len(data) - 1; i >= 0; i-- {
		result = append(result, data[i])
	}

	// Return new string.
	return string(result)
}

func SpaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func SpaceReplace(str string, pattern rune) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return pattern
		}
		return r
	}, str)
}

func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func Int64EncodeToString(nControl int64) string {

	var ALPHABET = "23456789BCDFGHJKLMNPQRSTVWXYZ"
	var BASE = int64(29)

	var sEncoded string
	for nControl >= BASE {
		nDiv := nControl / BASE
		nMod := nControl - (BASE * nDiv)

		sEncoded = string(ALPHABET[nMod]) + sEncoded
		nControl = nDiv
	}

	if nControl > 0 {
		sEncoded = string(ALPHABET[nControl]) + sEncoded
	}

	return sEncoded
}

func StringDecodeToInt64(sEncoded string) int64 {

	var ALPHABET = "23456789BCDFGHJKLMNPQRSTVWXYZ"

	var nMulti int
	var nDecoded int64

	nAlphaLen := len(ALPHABET)
	nMulti = 1
	for len(sEncoded) > 0 {

		sDigit := string(sEncoded[len(sEncoded)-1])
		nStrPos := strings.Index(ALPHABET, sDigit)
		nDecoded += int64(nMulti * nStrPos)
		nMulti = nMulti * nAlphaLen
		sEncoded = sEncoded[0 : len(sEncoded)-1]
	}

	return nDecoded
}

// Seconds-based time units
const (
	Minute   = 60
	Hour     = 60 * Minute
	Day      = 24 * Hour
	Week     = 7 * Day
	Month    = 30 * Day
	Year     = 12 * Month
	LongTime = 37 * Year
)

var magnitudes = []struct {
	d      int64
	format string
	divby  int64
}{
	{1, "now", 1},
	{2, "1 second %s", 1},
	{Minute, "%d seconds %s", 1},
	{2 * Minute, "1 minute %s", 1},
	{Hour, "%d minutes %s", Minute},
	{2 * Hour, "1 hour %s", 1},
	{Day, "%d hours %s", Hour},
	{2 * Day, "1 day %s", 1},
	{Week, "%d days %s", Day},
	{2 * Week, "1 week %s", 1},
	{Month, "%d weeks %s", Week},
	{2 * Month, "1 month %s", 1},
	{Year, "%d months %s", Month},
	{18 * Month, "1 year %s", 1},
	{2 * Year, "2 years %s", 1},
	{LongTime, "%d years %s", Year},
	{math.MaxInt64, "a long while %s", 1},
}

// GetDuration formats a time into a relative string.
//
// It takes two times and two labels.  In addition to the generic time
// delta string (e.g. 5 minutes), the labels are used applied so that
// the label corresponding to the smaller time is applied.
//
// GetDuration(timeInPast, timeInFuture, "earlier", "later") -> "3 weeks earlier"
// func GetDuration(a, b time.Time, cLabelpast, cLabelfuture string) string {

func GetTime() time.Time {
	return time.Now()
}

func GetTimeFormat(curTime time.Time) string {
	return curTime.Format("02/01/2006 15:04:05 MST")
}

func GetSystemTime() string {
	return time.Now().Format("02/01/2006 15:04:05 MST")
}

func GetSystemDate() string {
	return time.Now().Format("02/01/2006")
}

func GetUnixString(cTime string) string {
	cFormat := "02/01/2006 15:04:05 MST"
	tTime, err := time.Parse(cFormat, cTime)
	if err != nil {
		return ""
	}

	return strconv.FormatInt(tTime.Unix(), 10)
}

func GetDifferenceInYears(cTimeCur string, cTimePast string) (years int) {

	years = int(0)
	cFormat := "02/01/2006"

	if cTimeCur == "" {
		cTimeCur = time.Now().Format(cFormat)
	}
	tTimeCur, errCur := time.Parse(cFormat, cTimeCur)
	if errCur != nil {
		return
	}

	if cTimePast == "" {
		cTimePast = time.Now().Format(cFormat)
	}
	tTimePast, errPast := time.Parse(cFormat, cTimePast)
	if errPast != nil {
		return
	}

	if errCur == nil && errPast == nil {

		if tTimePast.After(tTimeCur) {
			tTimeCur, tTimePast = tTimePast, tTimeCur
		}

		yearCur, _, _ := tTimeCur.Date()
		yearPast, _, _ := tTimePast.Date()

		years = int(yearCur - yearPast)
	}

	return
}

func GetDifferenceInMonths(cTimeCur string, cTimePast string) (months int) {

	months = int(0)
	cFormat := "02/01/2006"

	if cTimeCur == "" {
		cTimeCur = time.Now().Format(cFormat)
	}
	tTimeCur, errCur := time.Parse(cFormat, cTimeCur)
	if errCur != nil {
		return
	}

	if cTimePast == "" {
		cTimePast = time.Now().Format(cFormat)
	}
	tTimePast, errPast := time.Parse(cFormat, cTimePast)
	if errPast != nil {
		return
	}

	if errCur == nil && errPast == nil {

		if tTimePast.After(tTimeCur) {
			tTimeCur, tTimePast = tTimePast, tTimeCur
		}

		_, monthCur, _ := tTimeCur.Date()
		_, monthPast, _ := tTimePast.Date()

		months = int(monthCur - monthPast)
	}

	return
}

func GetDifferenceInSeconds(cTimeCur string, cTimePast string) (seconds int) {

	seconds = int(0)
	cFormat := "02/01/2006"

	if cTimeCur == "" {
		cTimeCur = time.Now().Format(cFormat)
	}
	tTimeCur, errCur := time.Parse(cFormat, cTimeCur)
	if errCur != nil {
		return
	}

	if cTimePast == "" {
		cTimePast = time.Now().Format(cFormat)
	}
	tTimePast, errPast := time.Parse(cFormat, cTimePast)
	if errPast != nil {
		return
	}

	if errCur == nil && errPast == nil {

		secondsCur := tTimeCur.Unix()
		secondsPast := tTimePast.Unix()
		seconds = int(secondsCur - secondsPast)
	}

	return
}

func GetDuration(cPast, cPresent string) string {
	cLabel := "ago"
	cFormat := "02/01/2006 15:04:05 MST"

	tPast, err := time.Parse(cFormat, cPast)
	if err != nil {
		return ""
	}

	tPresent, err := time.Parse(cFormat, cPresent)
	if err != nil {
		return ""
	}

	diff := tPresent.Unix() - tPast.Unix()

	after := tPast.After(tPresent)
	if after {
		cLabel = "from now"
		diff = tPast.Unix() - tPresent.Unix()
	}

	n := sort.Search(len(magnitudes), func(i int) bool {
		return magnitudes[i].d > diff
	})

	mag := magnitudes[n]
	args := []interface{}{}
	escaped := false
	for _, ch := range mag.format {
		if escaped {
			switch ch {
			case '%':
			case 's':
				args = append(args, cLabel)
			case 'd':
				args = append(args, diff/mag.divby)
			}
			escaped = false
		} else {
			escaped = ch == '%'
		}
	}
	return fmt.Sprintf(mag.format, args...)
}
