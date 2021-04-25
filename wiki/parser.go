package wiki

import (
	"container/list"
	"errors"
	"regexp"
	"strings"
)

func getNextTempalteKey(runes *[]rune, rPos *int) (ret string, err error) {
	start := *rPos
	for ; (*rPos) < len(*runes); (*rPos)++ {
		r := (*runes)[*rPos]
		if r == '=' {
			ret = strings.TrimSpace(string((*runes)[start:(*rPos)]))
			(*rPos)++ // skip past the '='
			return
		}
		if r == '|' || (r == '}' && (*runes)[(*rPos)+1] == '}') {
			// position parameter instead of keyed.
			*rPos = start
			return "", nil
		}
	}
	return "", errors.New("Invalid page format. Unable to locate key separator")
}

var linebreakRegEx, _ = regexp.Compile(`(\s*[\r\n]\s*)+`)

func getNextTemplateValue(runes *[]rune, rPos *int) (ret string, err error) {
	bracketStack := list.List{}
	start := *rPos
	lenRunes := len(*runes)
	for ; (*rPos) < lenRunes; (*rPos)++ {
		r := (*runes)[*rPos]
		rLookAhead := (*runes)[(*rPos):intMin(lenRunes, (*rPos)+4)]
		if bracketStack.Len() == 0 && (r == '|' || (r == '}' && rLookAhead[1] == '}')) {
			ret = strings.TrimSpace(string((*runes)[start:(*rPos)]))
			ret = linebreakRegEx.ReplaceAllString(ret, " ")
			ret = strings.ReplaceAll(ret, " |", "|")
			(*rPos)-- // step back so the main parser sees the separator
			return
		}
		if r == '{' {
			bracketStack.PushFront(r)
		} else if r == '[' {
			bracketStack.PushFront(r)
		} else if r == '}' && bracketStack.Front() != nil && bracketStack.Front().Value == '{' {
			toRemove := bracketStack.Front()
			bracketStack.Remove(toRemove)
		} else if r == ']' && bracketStack.Front() != nil && bracketStack.Front().Value == '[' {
			toRemove := bracketStack.Front()
			bracketStack.Remove(toRemove)
		} else if runeSame(rLookAhead, []rune("<!--")) {
			// take HTML comments as part of the field value
			// <!-- -->
			(*rPos) += 4
			//skip to '-->'
			skipPast(runes, rPos, "-->")
		}
	}
	return "", errors.New("Invalid page format. Unable to locate key separator")
}

func skipPast(runes *[]rune, rPos *int, find string) {
	findRunes := []rune(find)
	lenFind := len(findRunes)
	lenRunes := len(*runes)
	for ; (*rPos) < lenRunes-lenFind; (*rPos)++ {
		r := (*runes)[*rPos:intMin(lenRunes, (*rPos)+lenFind)]
		if runeSame(r, findRunes) {
			(*rPos) += lenFind
			return
		}
	}
}

func runeSame(r1, r2 []rune) bool {
	l1 := len(r1)
	l2 := len(r2)
	if l1 != l2 {
		return false
	}
	for i, r := range r1 {
		if r != r2[i] {
			return false
		}
	}
	return true
}

func intMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}
