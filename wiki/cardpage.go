package wiki

import (
	"encoding/json"
	"errors"
	"strings"
)

//CardPage represents a wiki page that is for Card information
type CardPage struct {
	PageName   string
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

//Parse Parses a wiki page into a card. returns `nil` if there is no Card template definition in the page.
func (c *CardPage) Parse(pageText string) (err error) {
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
		c.PageHeader = strings.TrimSpace(pageText[:cardIdx])
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
	err = json.Unmarshal(datajson, &(c.CardInfo))
	if err != nil {
		return
	}

	for k, v := range pageContentMap {
		if !cardFieldIsKnown(k) {
			if c.CardInfo.unknownFields == nil {
				c.CardInfo.unknownFields = make(map[string]string)
			}
			c.CardInfo.unknownFields[k] = v
		}
	}

	if cardEndIdx+cardIdx < len(pageText) {
		c.PageFooter = strings.TrimSpace(pageText[cardEndIdx+cardIdx:])
	}

	return
}
