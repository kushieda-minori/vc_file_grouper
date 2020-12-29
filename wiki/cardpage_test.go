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
	|quote misc 3 = [http://some.link/ nice text]<!-- comment {[|]} || } { -->
|availability={{Tooltip|Saintly Oracle's {{cool}}Celestial Stone Exchange|December 7th ー January 12th 2020}}}}



some footer


	`
	cardPage := CardPage{
		PageName: "Test Page",
	}
	err := cardPage.Parse(testPage)
	if err != nil {
		t.Errorf("Parse returned an error: %s", err.Error())
		return
	}
	if cardPage.PageName != "Test Page" {
		t.Errorf("Invalid value for PageName found: `%s`", cardPage.PageName)
	}
	if cardPage.CardInfo.Element != "cool" {
		t.Errorf("Invalid value for Element found: `%s`", cardPage.CardInfo.Element)
	}
	if cardPage.CardInfo.Rarity != "R" {
		t.Errorf("Invalid value for Rarity found: `%s`", cardPage.CardInfo.Rarity)
	}
	if cardPage.CardInfo.QuoteMisc1 != "{{Quote|my quote| another param}}" {
		t.Errorf("Invalid value for QuoteMisc1 found: `%s`", cardPage.CardInfo.QuoteMisc1)
	}
	if cardPage.CardInfo.QuoteMisc2 != "[[Quote|my quote {{!}} another param]]" {
		t.Errorf("Invalid value for QuoteMisc2 found: `%s`", cardPage.CardInfo.QuoteMisc2)
	}
	if cardPage.CardInfo.QuoteMisc3 != "[http://some.link/ nice text]<!-- comment {[|]}|| } { -->" {
		t.Errorf("Invalid value for QuoteMisc3 found: `%s`", cardPage.CardInfo.QuoteMisc3)
	}
	if cardPage.CardInfo.Availability != "{{Tooltip|Saintly Oracle's {{cool}}Celestial Stone Exchange|December 7th ー January 12th 2020}}" {
		t.Errorf("Invalid value for Availability found: `%s`", cardPage.CardInfo.Availability)
	}
	if strings.TrimSpace(cardPage.PageHeader) != "Some Header" {
		t.Errorf("Invalid value for pageHeader found: `%s`", cardPage.PageHeader)
	}
	if strings.TrimSpace(cardPage.PageFooter) != "some footer" {
		t.Errorf("Invalid value for pageFooter found: `%s`", cardPage.PageFooter)
	}

	expectedUnknownFields := 3 //

	if len(cardPage.CardInfo.unknownFields) != expectedUnknownFields {
		t.Errorf("Expected %d unknown field. but found: %d", expectedUnknownFields, len(cardPage.CardInfo.unknownFields))
	}

	val, ok := cardPage.CardInfo.unknownFields["unknown key"]
	if !ok {
		t.Error("Unknown Key `unknown key` was not tracked")
	} else if val != "Some UnknownValue" {
		t.Errorf("Unknown Key `unknown key` had an unexpected value: `%s`", val)
	}

	val, ok = cardPage.CardInfo.unknownFields["1"]
	if !ok {
		t.Error("Unknown Key `1` was not tracked")
	} else if val != "position 1" {
		t.Errorf("Unknown Key `1` had an unexpected value: `%s`", val)
	}

	val, ok = cardPage.CardInfo.unknownFields["2"]
	if !ok {
		t.Error("Unknown Key `2` was not tracked")
	} else if val != "" {
		t.Errorf("Unknown Key `2` had an unexpected value: `%s`", val)
	}

	// validate tostring output
	actual := cardPage.String()
	expected := `Some Header

{{Card
|element = cool
|rarity = R
|quote misc 1 = {{Quote|my quote| another param}}
|quote misc 2 = [[Quote|my quote {{!}} another param]]
|quote misc 3 = [http://some.link/ nice text]<!-- comment {[|]}|| } { -->
|availability = {{Tooltip|Saintly Oracle's {{cool}}Celestial Stone Exchange|December 7th ー January 12th 2020}}

|1 = position 1
|unknown key = Some UnknownValue
}}

some footer
`
	if actual != expected {
		t.Errorf("Actual string value did not match expected. Actual: `%s`", actual)
	}
}
