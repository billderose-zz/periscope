package main

import "sort"
import "net/http"
import "fmt"
import "io/ioutil"
import "regexp"
import "strings"

const ResourceUrl = "https://s3-us-west-2.amazonaws.com/periscope-public/alice-in-wonderland.txt"

var PunctuationRegex = regexp.MustCompile("[[:punct:]]")
var WhitespaceRegex = regexp.MustCompile("[[:space]]")
var DashRegex = regexp.MustCompile("-")
var QuoteRegex = regexp.MustCompile(`['"]`)

type bagOfWords map[string]int

type sortedBagOfWords struct {
	m bagOfWords
	s []string
}

func (sm *sortedBagOfWords) Len() int {
	return len(sm.m)
}

func (sm *sortedBagOfWords) Less(i, j int) bool {
	countOne := sm.m[sm.s[i]]
	countTwo := sm.m[sm.s[j]]
	if (countOne == countTwo) {
		return strings.Compare(sm.s[i], sm.s[j]) < 0
	}
	return countTwo < countOne
}

func (sm *sortedBagOfWords) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}



func (m bagOfWords) sort() sortedBagOfWords {
	sm := sortedBagOfWords{m: m, s: make([]string, len(m))}
	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(&sm)
	return sm
}	

func sanitize(body []byte) string {
	unquoted := QuoteRegex.ReplaceAll(body, []byte(""))
	undashed := DashRegex.ReplaceAll(unquoted, []byte(" "))
	unpunct := PunctuationRegex.ReplaceAll(undashed, []byte(" "))
	unspaced := WhitespaceRegex.ReplaceAll(unpunct, []byte(" "))

	return strings.ToLower(string(unspaced))
}

func tokenize(body []byte) []string {
	return strings.Split(sanitize(body), " ")
}

func bagOfWordsify(tokens []string) bagOfWords {
	bagOfWords := make(bagOfWords)
	for _, t := range tokens {
		bagOfWords[t]++
	}
	return bagOfWords 
}

func main() {
    unsortedMap := make(bagOfWords)
	unsortedMap["a"] = 1
	resp, err := http.Get(ResourceUrl)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Unable to download text: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to download text: %s", err)
	}
	bagOfWords := bagOfWordsify(tokenize(body))
	for word, count := range bagOfWords {
		fmt.Printf("%s: %s\n", word, strings.Repeat("*", count))
	}
	sortedBagOfWords := bagOfWords.sort()
	for _, word := range sortedBagOfWords.s {
		fmt.Printf("%s: %s\n", word, strings.Repeat("*", sortedBagOfWords.m[word]))
	}
}