package hearken

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	letterThenDigit = regexp.MustCompile(`(\pL)(\pN)`)
	digitThenLetter = regexp.MustCompile(`(\pN)(\pL)`)
	notWordOrSpace  = regexp.MustCompile(`[^\pL\pN\s]`)
	runsOfSpace     = regexp.MustCompile(`\s+`)
)

// Normalize puts spoken text and the names in a home into one comparable
// form: lower case, punctuation removed, single spaces, and number words
// written as digits ("Test Twenty Two!" becomes "test 22"). Words a speech
// recogniser often hears instead of a number count as that number: "to"
// and "too" are 2, "for" is 4.
func Normalize(value string) string {
	text := strings.ToLower(value)
	text = letterThenDigit.ReplaceAllString(text, "${1} ${2}")
	text = digitThenLetter.ReplaceAllString(text, "${1} ${2}")
	text = notWordOrSpace.ReplaceAllString(text, " ")
	text = strings.TrimSpace(runsOfSpace.ReplaceAllString(text, " "))

	var words []string
	for _, word := range strings.Split(text, " ") {
		number, isNumber := numberWords[word]
		if !isNumber {
			words = append(words, word)
			continue
		}
		last := len(words) - 1
		previous, err := -1, error(nil)
		if last >= 0 {
			previous, err = strconv.Atoi(words[last])
		}
		switch {
		case last < 0 || err != nil:
			words = append(words, strconv.Itoa(number))
		case previous >= 20 && previous%10 == 0 && number > 0 && number < 10:
			words[last] = strconv.Itoa(previous + number)
		default:
			words[last] += strconv.Itoa(number)
		}
	}
	return strings.Join(words, " ")
}

var numberWords = map[string]int{
	"zero": 0, "oh": 0,
	"one": 1, "won": 1,
	"two": 2, "to": 2, "too": 2,
	"three": 3,
	"four":  4, "for": 4, "fore": 4,
	"five": 5, "six": 6, "seven": 7,
	"eight": 8, "ate": 8,
	"nine": 9, "ten": 10, "eleven": 11, "twelve": 12, "thirteen": 13,
	"fourteen": 14, "fifteen": 15, "sixteen": 16, "seventeen": 17,
	"eighteen": 18, "nineteen": 19,
	"twenty": 20, "thirty": 30, "forty": 40, "fifty": 50,
	"sixty": 60, "seventy": 70, "eighty": 80, "ninety": 90,
}
