package wiki

import (
	"strings"
	"testing"
)

func TestParseWikiPage(t *testing.T) {
	testPage := `
Some Header

{{Card|element=cool|rarity=R|unknown key = Some UnknownValue|position 1|
	|quote misc 1 = {{Quote|my quote | another param}}
 | quote misc 2 = [[Quote|my quote {{!}} another param]]
	|quote misc 3 = [http://some.link/ nice text]
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
	if cardInfo.QuoteMisc1 != "{{Quote|my quote | another param}}" {
		t.Errorf("Invalid value for QuoteMisc1 found: `%s`", cardInfo.QuoteMisc1)
	}
	if cardInfo.QuoteMisc2 != "[[Quote|my quote {{!}} another param]]" {
		t.Errorf("Invalid value for QuoteMisc2 found: `%s`", cardInfo.QuoteMisc2)
	}
	if cardInfo.QuoteMisc3 != "[http://some.link/ nice text]" {
		t.Errorf("Invalid value for QuoteMisc3 found: `%s`", cardInfo.QuoteMisc3)
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

	expectedUnknownFields := 3 //

	if len(cardInfo.unknownFields) != expectedUnknownFields {
		t.Errorf("Expected %d unknown field. but found: %d", expectedUnknownFields, len(cardInfo.unknownFields))
	}

	val, ok := cardInfo.unknownFields["unknown key"]
	if !ok {
		t.Error("Unknown Key `unknown key` was not tracked")
	} else if val != "Some UnknownValue" {
		t.Errorf("Unknown Key `unknown key` had an unexpected value: `%s`", val)
	}

	val, ok = cardInfo.unknownFields["1"]
	if !ok {
		t.Error("Unknown Key `1` was not tracked")
	} else if val != "position 1" {
		t.Errorf("Unknown Key `1` had an unexpected value: `%s`", val)
	}

	val, ok = cardInfo.unknownFields["2"]
	if !ok {
		t.Error("Unknown Key `2` was not tracked")
	} else if val != "" {
		t.Errorf("Unknown Key `2` had an unexpected value: `%s`", val)
	}

	// validate tostring output
	actual := cardInfo.String()
	expected := `Some Header

{{Card
|element = cool
|rarity = R
|quote misc 1 = {{Quote|my quote | another param}}
|quote misc 2 = [[Quote|my quote {{!}} another param]]
|quote misc 3 = [http://some.link/ nice text]
|availability = {{Tooltip|Saintly Oracle's {{cool}}Celestial Stone Exchange|December 7th ー January 12th 2020}}
<!-- these fields were unknown to the bot, but have not been removed -->
|1 = position 1
|unknown key = Some UnknownValue
}}

some footer
`
	if actual != expected {
		t.Errorf("Actual string value did not match expected. Actual: `%s`", actual)
	}
}
