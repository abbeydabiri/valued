package functions

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"strconv"
	"unicode"
)

func SortIntsAsc(s []int) {
	sort.Sort(sort.IntSlice(s))
}

func SortIntsDesc(s []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
}

func SortAsc(s []string) {
	sort.Sort(stringSlice(s))
}

func SortDesc(s []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(s)))
}

func Strings(s []string) {
	sort.Sort(stringSlice(s))
}

func Less(s, t string) bool {
	return LessRunes([]rune(s), []rune(t))
}

func LessRunes(s, t []rune) bool {
	nprefix := commonPrefix(s, t)
	if len(s) == nprefix && len(t) == nprefix {
		return false
	}
	sEnd := leadDigits(s[nprefix:]) + nprefix
	tEnd := leadDigits(t[nprefix:]) + nprefix
	if sEnd > nprefix || tEnd > nprefix {
		start := trailDigits(s[:nprefix])
		if sEnd-start > 0 && tEnd-start > 0 {
			sn := atoi(s[start:sEnd])
			tn := atoi(t[start:tEnd])
			if sn != tn {
				return sn < tn
			}
		}
	}
	switch {
	case len(s) == nprefix:
		return true
	case len(t) == nprefix:
		return false
	default:
		return s[nprefix] < t[nprefix]
	}
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func atoi(r []rune) uint64 {
	if len(r) < 1 {
		panic(errors.New("atoi got an empty slice"))
	}
	const cutoff = uint64((1<<64-1)/10 + 1)
	const maxVal = 1<<64 - 1

	var n uint64
	for _, d := range r {
		v := uint64(d - '0')
		if n >= cutoff {
			return 1<<64 - 1
		}
		n *= 10
		n1 := n + v
		if n1 < n || n1 > maxVal {
			return 1<<64 - 1
		}
		n = n1
	}
	return n
}

func commonPrefix(s, t []rune) int {
	for i := range s {
		if i >= len(t) {
			return len(t)
		}
		if s[i] != t[i] {
			return i
		}
	}
	return len(s)
}

func trailDigits(r []rune) int {
	for i := len(r) - 1; i >= 0; i-- {
		if !isDigit(r[i]) {
			return i + 1
		}
	}
	return 0
}

func leadDigits(r []rune) int {
	for i := range r {
		if !isDigit(r[i]) {
			return i
		}
	}
	return len(r)
}

type stringSlice []string

func (ss stringSlice) Len() int {
	return len(ss)
}

func (ss stringSlice) Less(i, j int) bool {
	return Less(ss[i], ss[j])
}

func (ss stringSlice) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

//Search Results Sorting

func HighlightSearchResult(sSearchText string, sTexttosearch string) (string, int) {
	nFrequency := 0
	if strings.Contains(sTexttosearch, sSearchText) && len(sTexttosearch) > 2 {
		rSearchText := regexp.MustCompile(sSearchText)
		nFrequency = len(rSearchText.FindStringSubmatchIndex(sTexttosearch))
		sTexttosearch = rSearchText.ReplaceAllString(sTexttosearch, "<u><b><i>"+sSearchText+"</i></b></u>")
	}
	return sTexttosearch, nFrequency
}

func SortSearchResult(cTagformat string, sSearchtextSplit []string, mapValuesTemp map[string]interface{}) map[string]interface{} {

	var aKeys []string
	mapValues := make(map[string]interface{})
	mapValuesSort := make(map[string]interface{})
	for cControl, mapXdoc := range mapValuesTemp {
		cFrequency := 1
		cFrequencyTemp := 0
		mapXdoc := mapXdoc.(map[string]interface{})

		cCode := fmt.Sprintf("%s", mapXdoc["code"])
		cTitle := fmt.Sprintf("%s", mapXdoc["title"])
		cDescription := fmt.Sprintf("%s", mapXdoc["description"])

		for _, sSearchtext := range sSearchtextSplit {
			if len(sSearchtext) > 0 {

				cCode, cFrequencyTemp = HighlightSearchResult(sSearchtext, cCode)
				cFrequency += cFrequencyTemp

				cTitle, cFrequencyTemp = HighlightSearchResult(sSearchtext, cTitle)
				cFrequency += cFrequencyTemp

				cDescription, cFrequencyTemp = HighlightSearchResult(sSearchtext, cDescription)
				cFrequency += cFrequencyTemp
			}
		}

		mapXdoc["code"] = cCode
		mapXdoc["title"] = cTitle
		mapXdoc["description"] = cDescription

		if mapXdoc["lastdate"] != nil && mapXdoc["lastdate"].(string) != "" {
			cControl = fmt.Sprintf("%s.%s", GetUnixString(mapXdoc["lastdate"].(string)), cControl)
		}

		cTag := fmt.Sprintf("%v%s", cFrequency, cControl)
		mapValuesSort[cTag] = mapXdoc
		aKeys = append(aKeys, cTag)
	}

	cNumber := 1
	// sort.Sort(sort.Reverse(sort.StringSlice(aKeys)))
	SortDesc(aKeys)
	for _, cPriority := range aKeys {
		cTag := fmt.Sprintf("%s#%s", cNumber, cTagformat)
		mapXdoc := mapValuesSort[cPriority].(map[string]interface{})
		mapXdoc["lineno"] = fmt.Sprintf("%v", cNumber)
		mapValues[cTag] = mapXdoc
		cNumber++
	}

	return mapValues
}

func SortMap(mapUnsorted map[string]interface{}) []string {
	var aKeys []string
	for cFieldname, iFieldvalue := range mapUnsorted {
		fieldType := reflect.TypeOf(iFieldvalue)
		if fieldType != nil {
			switch fieldType.Kind() {
			case reflect.Map:
				iFieldvalue = SortMap(iFieldvalue.(map[string]interface{}))
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

// Mostly based on http://golang.org/pkg/sort/#example__sortKeys
// Since name is a string, don't need to implement fancy structs

type By func(name1, name2 string) bool

func (by By) Sort(names []string) {
	ps := &nameSorter{
		names: names,
		by:    by,
	}
	sort.Sort(ps)
}

type nameSorter struct {
	names []string
	by    func(name1, name2 string) bool
}

func getNumber(name string) (new_s int) {
	// Need a key to sort the names by. In this case, it's the number contained in the name
	// Example: M23abcdeg ---> 23
	// Of course, this tends to break hideously once complicated stuff is involved:
	// For example, 'ülda123dmwak142.e2dööööwq,' what do you sort by here?
	// Could be 123, could be 142, or 2, could be 1231422 (this takes the last case)

	s := make([]string, 0)
	for _, element := range name {
		if unicode.IsNumber(element) {
			s = append(s, string(element))
		}
	}

	new_s, _ = strconv.Atoi(strings.Join(s, ""))
	// if err != nil {
	// 	log.Fatal(err) // Just die
	// }
	return new_s
}

// Need some inbuilt methods
func (s *nameSorter) Swap(i, j int) {
	s.names[i], s.names[j] = s.names[j], s.names[i]
}
func (s *nameSorter) Len() int {
	return len(s.names)
}
func (s *nameSorter) Less(i, j int) bool {
	return s.by(s.names[i], s.names[j])
}
