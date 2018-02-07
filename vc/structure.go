package vc

import (
	"time"
)

// the "garden" field lists some details about the kindoms available to the players
type Garden struct { // garden
	Id           int `json:"_id"`
	BlockX       int `json:"block_x"`
	BlockY       int `json:"block_y"`
	UnlockBlockX int `json:"unlock_block_x"`
	UnlockBlockY int `json:"unlock_block_y"`
	BgId         int `json:"bg_id"`
	Debris       int `json:"debris"`
	CastleId     int `json:"castle_id"`
}

//"garden_debris" lists information about clearing debris from your kingdom
type GardenDebris struct { // garden_debris
	Id           int `json:"_id"`
	GardenId     int `json:"garden_id"`
	StructureId  int `json:"structure_id"`
	X            int `json:"x"`
	Y            int `json:"y"`
	LevelCap     int `json:"level_cap"`
	UnlockAreaId int `json:"unlock_area_id"`
	Time         int `json:"time"`
	Coin         int `json:"coin"`
	Iron         int `json:"iron"`
	Ether        int `json:"ether"`
	Cash         int `json:"cash"`
	Exp          int `json:"exp"`
}

// "structures" gives information about availability of for buildinds.
// The names of the structions in this list match those in the MsgBuildingName_en.strb file
type Structure struct { // structures
	Id              int    `json:"_id"`
	StructureTypeId int    `json:"structure_type_id"`
	MaxLv           int    `json:"max_lv"`
	UnlockCastleId  int    `json:"unlock_castle_id"`
	UnlockCastleLv  int    `json:"unlock_castle_lv"`
	UnlockAreaId    int    `json:"unlock_area_id"`
	BaseNum         int    `json:"base_num"`
	SizeX           int    `json:"size_x"`
	SizeY           int    `json:"size_y"`
	Order           int    `json:"order"`
	EventId         int    `json:"event_id"`
	Visitable       int    `json:"visitable"`
	Step            int    `json:"step"`
	Passable        int    `json:"passable"`
	Connectable     int    `json:"connectable"`
	Enable          int    `json:"enable"`
	Stockable       int    `json:"stockable"`
	Flag            int    `json:"flag"`
	GardenFlag      int    `json:"garden_flag"`
	ShopGroupDecoId int    `json:"shop_group_deco_id"`
	Name            string `json:"-"` // MsgBuildingName_en.strb
	Description     string `json:"-"` // MsgBuildingDesc_en.strb
	_levels         []StructureLevel
	_numCosts       []StructureCost
	_castleBonus    []CastleLevel
}

// "event_structures" lists any structures available in the current event

// structure_level lists the level for the available structures
type StructureLevel struct { // structure_level
	Id           int            `json:"_id"`
	StructureId  int            `json:"structure_id"`
	Level        int            `json:"level"`
	TexId        int            `json:"tex_id"`
	LevelCap     int            `json:"level_cap"`
	UnlockAreaId int            `json:"unlock_area_id"`
	Time         int            `json:"time"`
	BeginnerTime int            `json:"beginner_time"`
	Coin         int            `json:"coin"`
	Iron         int            `json:"iron"`
	Ether        int            `json:"ether"`
	Cash         int            `json:"cash"`
	Elixir       int            `json:"elixir"`
	Price        int            `json:"price"`
	Exp          int            `json:"exp"`
	Resource     *ResourceLevel `json:"-"`
	Bank         *BankLevel     `json:"-"`
}

type StructureCost struct { // structure_num_cost
	Id          int `json:"_id"`
	Num         int `json:"num"`
	StructureId int `json:"structure_id"`
	Coin        int `json:"coin"`
	Iron        int `json:"iron"`
	Ether       int `json:"ether"`
	Cash        int `json:"cash"`
	Elixir      int `json:"elixir"`
}

type ResourceLevel struct { // resource
	Id           int `json:"_id"`
	StructureId  int `json:"structure_id"`
	Level        int `json:"level"`
	IntervalTime int `json:"interval_time"`
	JobTexId     int `json:"job_tex_id"`
	FixTexId     int `json:"fix_tex_id"`
	ResourceId   int `json:"resource_id"`
	Income       int `json:"income"`
}

type BankLevel struct { //bank_level
	Id          int `json:"_id"`
	StructureId int `json:"structure_id"`
	Level       int `json:"level"`
	LowTexId    int `json:"low_tex_id"`
	MidTexId    int `json:"mid_tex_id"`
	HiTexId     int `json:"hi_tex_id"`
	ResourceId  int `json:"resource_id"`
	Value       int `json:"value"`
}

type CastleLevel struct { //castle_level
	Id                int `json:"_id"`
	CastleStructureId int `json:"castle_structure_id"`
	Level             int `json:"level"`
	StructureId       int `json:"structure_id"`
	BaseAdd           int `json:"base_add"`
	Max               int `json:"max"`
}

func (s *Structure) Levels(v *VcFile) []StructureLevel {
	if s._levels == nil {
		s._levels = make([]StructureLevel, 0)
		for _, l := range v.StructureLevels {
			if l.StructureId == s.Id {
				l.cacheResource(v) // cache off the resource level for later
				l.cacheBank(v)     // cache off the bank level for later
				s._levels = append(s._levels, l)
			}
		}
	}
	return s._levels
}

func (s *Structure) PurchaseCosts(v *VcFile) []StructureCost {
	if s._numCosts == nil {
		s._numCosts = make([]StructureCost, 0)
		for _, p := range v.StructureNumCosts {
			if p.StructureId == s.Id {
				s._numCosts = append(s._numCosts, p)
			}
		}
	}
	return s._numCosts
}

func (s *Structure) IsResource() bool {
	return s.StructureTypeId == 1 ||
		s.StructureTypeId == 2 ||
		s.StructureTypeId == 3 ||
		s.StructureTypeId == 20
}

func (s *Structure) IsBank() bool {
	return s.StructureTypeId == 4 ||
		s.StructureTypeId == 5 ||
		s.StructureTypeId == 6 ||
		s.StructureTypeId == 21
}

func (s *Structure) CastleBonuses(v *VcFile) []CastleLevel {
	if s._castleBonus == nil {
		s._castleBonus = make([]CastleLevel, 0)
		for _, cl := range v.CastleLevels {
			if cl.StructureId == s.Id {
				s._castleBonus = append(s._castleBonus, cl)
			}
		}
	}
	return s._castleBonus
}

func (s *Structure) MaxQty(v *VcFile) int {
	maxCLB := 0
	clbs := s.CastleBonuses(v)
	for i, clb := range clbs {
		if clb.Level > clbs[maxCLB].Level {
			maxCLB = i
		}
	}
	return clbs[maxCLB].Max
}

func (l *StructureLevel) cacheResource(v *VcFile) {
	if l.Resource == nil {
		for i, sr := range v.ResourceLevels {
			if sr.StructureId == l.StructureId && sr.Level == l.Level {
				l.Resource = &v.ResourceLevels[i]
				break
			}
		}
	}
}
func (l *StructureLevel) cacheBank(v *VcFile) {
	if l.Bank == nil {
		for i, br := range v.BankLevels {
			if br.StructureId == l.StructureId && br.Level == l.Level {
				l.Bank = &v.BankLevels[i]
				break
			}
		}
	}
}

func (sr *ResourceLevel) Rate() int {
	return sr.Income / sr.IntervalTime
}

func (sr *ResourceLevel) FillTime() time.Duration {
	return time.Duration(float64(sr.Income)/float64(sr.Rate())*60.0) * time.Second
}

func StructureScan(id int, v *VcFile) *Structure {
	if id > 0 {
		if id < len(v.Structures) && v.Structures[id-1].Id == id {
			return &v.Structures[id-1]
		}
		for k, val := range v.Structures {
			if val.Id == id {
				return &v.Structures[k]
			}
		}
	}
	return nil
}
