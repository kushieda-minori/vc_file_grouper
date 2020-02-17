package vc

import (
	"sort"
)

// GuildBattle "mst_guildbattle_schedule"
type GuildBattle struct {
	ID                 int       `json:"_id"`
	GuildBattleType    int       `json:"guild_battle_type"`
	GuildBingoID       int       `json:"guild_bingo_id"`
	EnableDayOfWeek    int       `json:"enable_day_of_week"`
	StartDatetime      Timestamp `json:"start_datetime"`
	EndDatetime        Timestamp `json:"end_datetime"`
	SkipDay            int       `json:"skip_day"`
	RoundScheduleGroup int       `json:"round_schedule_group"`
	MatchGuildCount    int       `json:"match_guild_count"`
	AreaAttribute1     int       `json:"area_attribute_1"`
	AreaAttribute2     int       `json:"area_attribute_2"`
	AreaAttribute3     int       `json:"area_attribute_3"`
	AreaAttribute4     int       `json:"area_attribute_4"`
	BannerID           int       `json:"banner_id"`
	individualRewards  []RankRewardSheet
	rankRewards        []RankRewardSheet
}

// GuildBingoBattle "mst_guildbingo"
type GuildBingoBattle struct {
	ID                            int       `json:"_id"`
	ExchangeItemID                int       `json:"exchange_item_id"`
	ExchangeItemRemovalDate       Timestamp `json:"exchange_item_removal_date"`
	RankingRewardDistributionDate Timestamp `json:"ranking_reward_distribution_date"`
	ModelName                     string    `json:"model_name"`
	SheetGroupID                  int       `json:"sheet_group_id"`
	ExchangeRewardGroupID         int       `json:"exchange_reward_group_id"`      // links to GuildBingoExchangeReward.GroupID
	RoundRankingRewardGroupID     int       `json:"round_ranking_reward_group_id"` // links to GuildBingoRoundRankingReward.GroupID
	RoundJoinRewardGroupID        int       `json:"round_join_reward_group_id"`
	CellRewardGroupID             int       `json:"cell_reward_group_id"`
	LineRewardGroupID             int       `json:"line_reward_group_id"`
	DrawBingoBallGroupID          int       `json:"draw_bingo_ball_group_id"`
	WinnerPoint                   int       `json:"winner_point"`
	DefenseDeckWinNum             int       `json:"defense_deck_win_num"`
	DefenseDeckRewardNum          int       `json:"defense_deck_reward_num"`
	RoundJoinRewardNum            int       `json:"round_join_reward_num"`
	MvpRewardNum                  int       `json:"mvp_reward_num"`
	RoundKillKingRewardNum        int       `json:"round_kill_king_reward_num"`
	KingCellRewardGroupID         int       `json:"king_cell_reward_group_id"`
	KingSeriesID                  int       `json:"king_series_id"`
	CampaignID                    int       `json:"campaign_id"`
	LoserPoint                    int       `json:"loser_point"`
	archwitches                   []Archwitch
	_campaigns                    []GuildBingoPointCampaign
}

// GuildBingoExchangeReward ABB item exhange
// "mst_guildbingo_exchange_reward"
type GuildBingoExchangeReward struct {
	ID            int `json:"_id"`
	GroupID       int `json:"group_id"`    // reward group
	RequireNum    int `json:"require_num"` // cost
	RewardType    int `json:"reward_type"` // 1 is card, 2 is item
	RewardID      int `json:"reward_id"`
	ExchangeLimit int `json:"exchange_limit"`
	Num           int `json:"num"` // number obtained from exchange
	IsPickup      int `json:"is_pickup"`
}

// GuildAUBWinReward AUB win rewards
// "mst_guildbattle_win_reward"
type GuildAUBWinReward struct {
	ID         int `json:"_id"`
	SheetID    int `json:"sheet_id"`
	Win        int `json:"win"`
	ItemID     int `json:"item_id"`
	FragmentID int `json:"fragment_id"`
	CardID     int `json:"card_id"`
	Num        int `json:"num"`
}

// GuildBattleRewardRef "mst_guildbattle_point_reward"
type GuildBattleRewardRef struct {
	ID                          int `json:"_id"`
	EventID                     int `json:"event_id"` // links to GuildBattle.id
	SheetID                     int `json:"sheet_id"` // links to VCFile.GuildBattleRankingRewards[].SheetID
	RankingSheetID              int `json:"ranking_sheet_id"`
	IndividualRankingSheetID    int `json:"individual_ranking_sheet_id"` // links to GuildBattleIndividualPoint.SheetID
	MidIndividualRankingSheetID int `json:"mid_individual_ranking_sheet_id"`
	GuildWinSheetID             int `json:"guild_win_sheet_id"`
	MidBonusDistributionDate    int `json:"mid_bonus_distribution_date"`
}

// GuildBingoRoundRankingReward "mst_guildbingo_round_ranking_reward"
// round ranking rewards
type GuildBingoRoundRankingReward struct {
	ID        int `json:"_id"`
	GroupID   int `json:"group_id"`
	Rank      int `json:"rank"`
	RewardNum int `json:"reward_num"`
}

// GuildBingoPointCampaign campain start/end for ABB double point days
type GuildBingoPointCampaign struct {
	ID            int       `json:"_id"`
	CampaignID    int       `json:"campaign_id"`
	StartDatetime Timestamp `json:"start_datetime"`
	EndDatetime   Timestamp `json:"end_datetime"`
	Multiple      int       `json:"multiple"`
}

// BingoBattle Bingo battle information
func (g *GuildBattle) BingoBattle() *GuildBingoBattle {
	if g == nil {
		return nil
	}
	if g.GuildBingoID > 0 {
		l := len(Data.GuildBingoBattles)
		i := sort.Search(l, func(i int) bool { return Data.GuildBingoBattles[i].ID >= g.GuildBingoID })
		if i >= 0 && i < l && Data.GuildBingoBattles[i].ID == g.GuildBingoID {
			return &(Data.GuildBingoBattles[i])
		}
	}
	return nil
}

// ExchangeRewards rewards for item exchanges for this battle
func (g *GuildBingoBattle) ExchangeRewards() []GuildBingoExchangeReward {
	if g == nil {
		return []GuildBingoExchangeReward{}
	}
	set := make([]GuildBingoExchangeReward, 0)
	if g.ExchangeRewardGroupID > 0 {
		for _, val := range Data.GuildBingoExchangeRewards {
			if val.GroupID == g.ExchangeRewardGroupID {
				set = append(set, val)
			}
		}
	}
	sort.Sort(GuildBingoExchangeRewardByTypeAndID(set))
	return set
}

// Campaigns extra point campaigns for this ABB
func (g *GuildBingoBattle) Campaigns() []GuildBingoPointCampaign {
	if g == nil {
		return []GuildBingoPointCampaign{}
	}
	if g._campaigns == nil {
		g._campaigns = make([]GuildBingoPointCampaign, 0)
		for _, a := range Data.GuildBingoPointCampaigns {
			if g.CampaignID == a.CampaignID {
				g._campaigns = append(g._campaigns, a)
			}
		}
	}
	return g._campaigns
}

// IndividualRewards individual rewards for this event
func (g *GuildBattle) IndividualRewards() []RankRewardSheet {
	if g == nil {
		return []RankRewardSheet{}
	}
	if g.individualRewards == nil {
		g.individualRewards = make([]RankRewardSheet, 0)
		rewards := g.rewards()
		for _, ipr := range Data.GuildBattleIndividualPoints {
			if rewards.SheetID == ipr.SheetID {
				g.individualRewards = append(g.individualRewards, ipr)
			}
		}
	}
	return g.individualRewards
}

// RankRewards for this event
func (g *GuildBattle) RankRewards() []RankRewardSheet {
	if g == nil {
		return []RankRewardSheet{}
	}
	if g.rankRewards == nil {
		g.rankRewards = make([]RankRewardSheet, 0)
		rewards := g.rewards()
		for _, rr := range Data.GuildBattleRankingRewards {
			if rewards.IndividualRankingSheetID == rr.SheetID {
				g.rankRewards = append(g.rankRewards, rr)
			}
		}
	}
	return g.rankRewards
}

func (g *GuildBattle) rewards() *GuildBattleRewardRef {
	if g == nil {
		return nil
	}
	for _, rref := range Data.GuildBattleRewardRefs {
		if g.ID == rref.EventID {
			return &rref
		}
	}
	return nil
}

// GuildBattleScan searches for a guild battle by ID
func GuildBattleScan(id int) *GuildBattle {
	if id <= 0 {
		return nil
	}
	l := len(Data.GuildBattles)
	i := sort.Search(l, func(i int) bool { return Data.GuildBattles[i].ID >= id })
	if i >= 0 && i < l && Data.GuildBattles[i].ID == id {
		return &(Data.GuildBattles[i])
	}
	return nil
}

// GuildBingoExchangeRewardByTypeAndID sort interface
type GuildBingoExchangeRewardByTypeAndID []GuildBingoExchangeReward

func (d GuildBingoExchangeRewardByTypeAndID) Len() int {
	return len(d)
}

func (d GuildBingoExchangeRewardByTypeAndID) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d GuildBingoExchangeRewardByTypeAndID) Less(i, j int) bool {
	if d[i].RewardType < d[j].RewardType {
		return true
	}
	if d[i].RewardType == d[j].RewardType {
		if d[i].RewardID > d[j].RewardID {
			return true
		}
	}
	return false
}
