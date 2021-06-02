package vc

import (
	"path/filepath"
	"strings"
)

//SubEvent fields on all new sub-event types
type SubEvent struct {
	ID                             int       `json:"_id"`
	ScenarioID                     int       `json:"scenario_id"`
	RankingRewardGroupID           int       `json:"ranking_reward_group_id"`
	ArrivalRewardGroupID           int       `json:"arrival_point_reward_group_id"`
	URLSchemeID                    int       `json:"url_scheme_id"`
	PublicStartDatetime            Timestamp `json:"public_start_datetime"`
	PublicEndDatetime              Timestamp `json:"public_end_datetime"`
	RankingStart                   Timestamp `json:"ranking_start_datetime"`
	RankingEnd                     Timestamp `json:"ranking_end_datetime"`
	EnemySymbolID                  int       `json:"enemy_symbol_id"`
	EventGatchaID                  int       `json:"eventgacha_id"`
	ElementalAlignmentBonuxGroupID int       `json:"elemental_alignment_bonus_group_id"`
	SymbolBonusGroupID             int       `json:"symbol_bonus_group_id"`
}

//GetURL Gets an event's URL
func (se *SubEvent) GetURL() string {
	if se == nil {
		return ""
	}
	url := URLSchemeScan(se.URLSchemeID)
	if url == nil {
		return ""
	}
	return url.Android
}
func (se *SubEvent) GetScenarioHtml(eventTitle, eventType string) (ret string, err error) {
	if se.ScenarioID < 0 {
		return
	}
	// string file location vcRoot/scenario/MsgScenarioString_<lang>.strb
	var lines []string
	lines, err = ReadStringFileFilter(filepath.Join(FilePath, "scenario", eventType, "MsgScenarioString_"+LangPack+".strb"), false)
	if err != nil {
		return
	}
	ret += "<html>\n<head>\n<title>" + eventTitle + " Story</title>\n</head>\n<body>\n<h1>" + eventTitle + "</h1>\n"

	scenario := 1
	llines := len(lines)
	for l, line := range lines {
		if scenario == se.ScenarioID || (se.ScenarioID == 0 && scenario == 1) {
			if line == "" {
				ret += "\n"
			}
			if strings.HasPrefix(line, "Chapter") {
				ret += "\n"
				ls := filterStoryLine(line)
				for _, l := range ls {
					if strings.HasPrefix(l, "Chapter") {
						ret += "<h2>" + l + "</h2>\n"
					} else {
						ret += l + "\n"
					}
				}
			} else {
				ls := filterStoryLine(line)
				if len(ls) > 0 {
					ret += "<dl>\n<dt>[[SPEAKER]]</dt>\n"
					ret += "<dd>"
					for _, l := range ls {
						ret += l + "<br/>\n"
					}
					ret += "</dd>\n<dl>\n"
				}
			}
		}
		if strings.Contains(line, "To be continued……") && l != llines-1 {
			scenario++
		}
	}
	ret += "\n</body>\n</html>"
	return
}

func filterStoryLine(line string) []string {
	line = strings.ReplaceAll(line, "\n", " ")
	line = strings.ReplaceAll(line, "  ", " ")
	line = strings.ReplaceAll(line, "<i><break>", "\n")
	line = strings.TrimSpace(line)
	lines := strings.Split(line, "\n")

	ll := len(lines)
	ret := make([]string, 0, ll)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			ret = append(ret, l)
		}
	}
	return ret
}
