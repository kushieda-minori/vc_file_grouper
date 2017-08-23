package vc

import (
	"sort"
)

//"mst_guildbattle_schedule"
type GuildBattle struct {
	Id                 int               `json:"_id"`
	GuildBattleType    int               `json:"guild_battle_type"`
	GuildBingoId       int               `json:"guild_bingo_id"`
	EnableDayOfWeek    int               `json:"enable_day_of_week"`
	StartDatetime      Timestamp         `json:"start_datetime"`
	EndDatetime        Timestamp         `json:"end_datetime"`
	SkipDay            int               `json:"skip_day"`
	RoundScheduleGroup int               `json:"round_schedule_group"`
	MatchGuildCount    int               `json:"match_guild_count"`
	AreaAttribute1     int               `json:"area_attribute_1"`
	AreaAttribute2     int               `json:"area_attribute_2"`
	AreaAttribute3     int               `json:"area_attribute_3"`
	AreaAttribute4     int               `json:"area_attribute_4"`
	BannerId           int               `json:"banner_id"`
	individualRewards  []RankRewardSheet `json:"-"`
	rankRewards        []RankRewardSheet `json:"-"`
}

// "mst_guildbingo"
type GuildBingoBattle struct {
	Id                            int         `json:"_id"`
	ExchangeItemId                int         `json:"exchange_item_id"`
	ExchangeItemRemovalDate       Timestamp   `json:"exchange_item_removal_date"`
	RankingRewardDistributionDate Timestamp   `json:"ranking_reward_distribution_date"`
	ModelName                     string      `json:"model_name"`
	SheetGroupId                  int         `json:"sheet_group_id"`
	ExchangeRewardGroupId         int         `json:"exchange_reward_group_id"`      // links to GuildBingoExchangeReward.GroupId
	RoundRankingRewardGroupId     int         `json:"round_ranking_reward_group_id"` // links to GuildBingoRoundRankingReward.GroupId
	RoundJoinRewardGroupId        int         `json:"round_join_reward_group_id"`
	CellRewardGroupId             int         `json:"cell_reward_group_id"`
	LineRewardGroupId             int         `json:"line_reward_group_id"`
	DrawBingoBallGroupId          int         `json:"draw_bingo_ball_group_id"`
	WinnerPoint                   int         `json:"winner_point"`
	DefenseDeckWinNum             int         `json:"defense_deck_win_num"`
	DefenseDeckRewardNum          int         `json:"defense_deck_reward_num"`
	RoundJoinRewardNum            int         `json:"round_join_reward_num"`
	MvpRewardNum                  int         `json:"mvp_reward_num"`
	RoundKillKingRewardNum        int         `json:"round_kill_king_reward_num"`
	KingCellRewardGroupId         int         `json:"king_cell_reward_group_id"`
	KingSeriesId                  int         `json:"king_series_id"`
	CampaignId                    int         `json:"campaign_id"`
	archwitches                   []Archwitch `json:"-"`
}

// ABB item exhange
// "mst_guildbingo_exchange_reward"
type GuildBingoExchangeReward struct {
	Id         int `json:"_id"`
	GroupId    int `json:"group_id"`    // reward group
	RequireNum int `json:"require_num"` // cost
	RewardType int `json:"reward_type"` // 1 is card, 2 is item
	RewardId   int `json:"reward_id"`
	Num        int `json:"num"` // number obtained from exchange
	IsPickup   int `json:"is_pickup"`
}

//AUB win rewards
// "mst_guildbattle_win_reward"
type GuildAUBWinReward struct {
	Id         int `json:"_id"`
	SheetId    int `json:"sheet_id"`
	Win        int `json:"win"`
	ItemId     int `json:"item_id"`
	FragmentId int `json:"fragment_id"`
	CardId     int `json:"card_id"`
	Num        int `json:"num"`
}

// "mst_guildbattle_point_reward"
type GuildBattleRewardRef struct {
	Id                          int `json:"_id"`
	EventId                     int `json:"event_id"` // links to GuildBattle.id
	SheetId                     int `json:"sheet_id"` // links to VCFile.GuildBattleRankingRewards[].SheetId
	RankingSheetId              int `json:"ranking_sheet_id"`
	IndividualRankingSheetId    int `json:"individual_ranking_sheet_id"` // links to GuildBattleIndividualPoint.SheetId
	MidIndividualRankingSheetId int `json:"mid_individual_ranking_sheet_id"`
	GuildWinSheetId             int `json:"guild_win_sheet_id"`
	MidBonusDistributionDate    int `json:"mid_bonus_distribution_date"`
}

// "mst_guildbingo_round_ranking_reward"
// round ranking rewards
type GuildBingoRoundRankingReward struct {
	Id        int `json:"_id"`
	GroupId   int `json:"group_id"`
	Rank      int `json:"rank"`
	RewardNum int `json:"reward_num"`
}

func (g *GuildBattle) BingoBattle(v *VcFile) *GuildBingoBattle {
	if g.GuildBingoId > 0 {
		l := len(v.GuildBingoBattles)
		i := sort.Search(l, func(i int) bool { return v.GuildBingoBattles[i].Id >= g.GuildBingoId })
		if i >= 0 && i < l && v.GuildBingoBattles[i].Id == g.GuildBingoId {
			return &(v.GuildBingoBattles[i])
		}
	}
	return nil
}

func (gb *GuildBingoBattle) Archwitches(v *VcFile) []Archwitch {
	if gb.archwitches == nil {
		gb.archwitches = make([]Archwitch, 0)
		for _, a := range v.Archwitches {
			if gb.KingSeriesId == a.KingSeriesId {
				gb.archwitches = append(gb.archwitches, a)
			}
		}
	}
	return gb.archwitches
}

func (g *GuildBingoBattle) ExchangeRewards(v *VcFile) []GuildBingoExchangeReward {
	set := make([]GuildBingoExchangeReward, 0)
	if g.ExchangeRewardGroupId > 0 {
		for _, val := range v.GuildBingoExchangeRewards {
			if val.GroupId == g.ExchangeRewardGroupId {
				set = append(set, val)
			}
		}
	}
	sort.Sort(GuildBingoExchangeRewardByTypeAndId(set))
	return set
}

func (gb *GuildBattle) IndividualRewards(v *VcFile) []RankRewardSheet {
	if gb.individualRewards == nil {
		gb.individualRewards = make([]RankRewardSheet, 0)
		rewards := gb.rewards(v)
		for _, ipr := range v.GuildBattleIndividualPoints {
			if rewards.SheetId == ipr.SheetId {
				gb.individualRewards = append(gb.individualRewards, ipr)
			}
		}
	}
	return gb.individualRewards
}

func (gb *GuildBattle) RankRewards(v *VcFile) []RankRewardSheet {
	if gb.rankRewards == nil {
		gb.rankRewards = make([]RankRewardSheet, 0)
		rewards := gb.rewards(v)
		for _, rr := range v.GuildBattleRankingRewards {
			if rewards.IndividualRankingSheetId == rr.SheetId {
				gb.rankRewards = append(gb.rankRewards, rr)
			}
		}
	}
	return gb.rankRewards
}

func (gb *GuildBattle) rewards(v *VcFile) *GuildBattleRewardRef {
	for _, rref := range v.GuildBattleRewardRefs {
		if gb.Id == rref.EventId {
			return &rref
		}
	}
	return nil
}

func GuildBattleScan(id int, battles []GuildBattle) *GuildBattle {
	if id <= 0 {
		return nil
	}
	l := len(battles)
	i := sort.Search(l, func(i int) bool { return battles[i].Id >= id })
	if i >= 0 && i < l && battles[i].Id == id {
		return &(battles[i])
	}
	return nil
}

type GuildBingoExchangeRewardByTypeAndId []GuildBingoExchangeReward

func (d GuildBingoExchangeRewardByTypeAndId) Len() int {
	return len(d)
}

func (d GuildBingoExchangeRewardByTypeAndId) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d GuildBingoExchangeRewardByTypeAndId) Less(i, j int) bool {
	if d[i].RewardType < d[j].RewardType {
		return true
	}
	if d[i].RewardType == d[j].RewardType {
		if d[i].RewardId > d[j].RewardId {
			return true
		}
	}
	return false
}
