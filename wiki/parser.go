package wiki

import (
	"container/list"
	"errors"
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

func getNextTemplateValue(runes *[]rune, rPos *int) (ret string, err error) {
	bracketStack := list.List{}
	start := *rPos
	for ; (*rPos) < len(*runes); (*rPos)++ {
		r := (*runes)[*rPos]
		if bracketStack.Len() == 0 && (r == '|' || (r == '}' && (*runes)[(*rPos)+1] == '}')) {
			ret = strings.TrimSpace(string((*runes)[start:(*rPos)]))
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
		}
	}
	return "", errors.New("Invalid page format. Unable to locate key separator")
}
