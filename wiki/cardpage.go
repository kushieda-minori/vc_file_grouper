package wiki

import (
	"encoding/json"
	"errors"
	"strings"
)

//CardPage represents a wiki page that is for Card information
type CardPage struct {
	CardInfo   CardFlat
	PageHeader string
	PageFooter string
}

func (c *CardPage) String() (ret string) {
	if c == nil {
		return ""
	}

	if c.PageHeader != "" {
		ret += c.PageHeader + "\n\n"
	}

	ret += c.CardInfo.String()

	if c.PageFooter != "" {
		ret += "\n" + c.PageFooter + "\n"
	}

	return
}

//ParseCardPage Parses a wiki page into a card. returns `nil` if there is no Card template definition in the page.
func ParseCardPage(pageText string) (ret CardPage, err error) {
	pageText = strings.TrimSpace(pageText)
	// lowercase all the text for comparison reasons.
	pageLower := strings.ToLower(pageText)
	cardIdx := strings.Index(pageLower, "{{card")
	if pageText == "" || cardIdx < 0 {
		cardIdx = strings.Index(pageLower, "{{template:card")
		if cardIdx < 0 {
			err = errors.New("Unable to find card template on page: " + pageText)
			return
		}
	}
	if cardIdx > 0 {
		ret.PageHeader = strings.TrimSpace(pageText[:cardIdx])
	}

	// convert the page Card template to a map
	pageContentMap, cardEndIdx, err := parseCard(pageText[cardIdx:])
	if err != nil {
		return
	}

	// convert the parsed map to JSON
	datajson, err := json.Marshal(pageContentMap)
	if err != nil {
		return
	}
	// unmarshall the JSON into our Card object
	err = json.Unmarshal(datajson, &(ret.CardInfo))
	if err != nil {
		return
	}

	for k, v := range pageContentMap {
		if !cardFieldIsKnown(k) {
			if ret.CardInfo.unknownFields == nil {
				ret.CardInfo.unknownFields = make(map[string]string)
			}
			ret.CardInfo.unknownFields[k] = v
		}
	}

	if cardEndIdx+cardIdx < len(pageText) {
		ret.PageFooter = strings.TrimSpace(pageText[cardEndIdx+cardIdx:])
	}

	return
}
