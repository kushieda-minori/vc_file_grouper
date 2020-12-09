package wiki

import (
	"strings"
	"testing"
)

func TestParseWikiPage(t *testing.T) {
	testPage := `
Some Header

{{Card|element=cool|rarity=R|unknown key = Some UnknownValue
	|quote misc 1 = {{Quote|my quote | sanother param}}
|availability={{Tooltip|Saintly Oracle's {{cool}}Celestial Stone Exchange|December 7th ー January 12th 2020}}}}

some footer
	`
	cardInfo, err := ParseWikiPage(testPage)
	if err != nil {
		t.Errorf("Parse returned an error: %s", err.Error())
		return
	}
	if cardInfo == nil {
		t.Error("Expected Card Info to have a value, but was nil")
		return
	}
	if cardInfo.Element != "cool" {
		t.Errorf("Invalid value for Element found: `%s`", cardInfo.Element)
	}
	if cardInfo.Rarity != "R" {
		t.Errorf("Invalid value for Rarity found: `%s`", cardInfo.Rarity)
	}
	if cardInfo.QuoteMisc1 != "{{Quote|my quote | sanother param}}" {
		t.Errorf("Invalid value for QuoteMisc1 found: `%s`", cardInfo.QuoteMisc1)
	}
	if cardInfo.Availability != "{{Tooltip|Saintly Oracle's {{cool}}Celestial Stone Exchange|December 7th ー January 12th 2020}}" {
		t.Errorf("Invalid value for Availability found: `%s`", cardInfo.Availability)
	}
	if strings.TrimSpace(cardInfo.pageHeader) != "Some Header" {
		t.Errorf("Invalid value for pageHeader found: `%s`", cardInfo.pageHeader)
	}
	if strings.TrimSpace(cardInfo.pageFooter) != "some footer" {
		t.Errorf("Invalid value for pageFooter found: `%s`", cardInfo.pageFooter)
	}
	val, ok := cardInfo.unknownFields["unknown key"]
	if !ok {
		t.Error("Unknown Key `unknown key` was not tracked")
	} else if val != "Some UnknownValue" {
		t.Errorf("Unknown Key `unknown key` had an unexpected value: `%s`", val)
	}
}
