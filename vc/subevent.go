package vc

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
