package vc

import (
	"log"
	"time"
)

// Garden the "garden" field lists some details about the kindoms available to the players
type Garden struct { // garden
	ID           int `json:"_id"`
	BlockX       int `json:"block_x"`
	BlockY       int `json:"block_y"`
	UnlockBlockX int `json:"unlock_block_x"`
	UnlockBlockY int `json:"unlock_block_y"`
	BgID         int `json:"bg_id"`
	Debris       int `json:"debris"`
	CastleID     int `json:"castle_id"`
}

// GardenDebris "garden_debris" lists information about clearing debris from your kingdom
type GardenDebris struct { // garden_debris
	ID           int `json:"_id"`
	GardenID     int `json:"garden_id"`
	StructureID  int `json:"structure_id"`
	X            int `json:"x"`
	Y            int `json:"y"`
	LevelCap     int `json:"level_cap"`
	UnlockAreaID int `json:"unlock_area_id"`
	Time         int `json:"time"`
	Coin         int `json:"coin"`
	Iron         int `json:"iron"`
	Ether        int `json:"ether"`
	Gem          int `json:"elixir"`
	Cash         int `json:"cash"`
	Exp          int `json:"exp"`
}

// Structure "structures" gives information about availability of for buildinds.
// The names of the structions in this list match those in the MsgBuildingName_en.strb file
type Structure struct { // structures
	ID              int    `json:"_id"`
	StructureTypeID int    `json:"structure_type_id"`
	MaxLv           int    `json:"max_lv"`
	UnlockCastleID  int    `json:"unlock_castle_id"`
	UnlockCastleLv  int    `json:"unlock_castle_lv"`
	UnlockAreaID    int    `json:"unlock_area_id"`
	BaseNum         int    `json:"base_num"`
	SizeX           int    `json:"size_x"`
	SizeY           int    `json:"size_y"`
	Order           int    `json:"order"`
	EventID         int    `json:"event_id"`
	Visitable       int    `json:"visitable"`
	Step            int    `json:"step"`
	Passable        int    `json:"passable"`
	Connectable     int    `json:"connectable"`
	Enable          int    `json:"enable"`
	Stockable       int    `json:"stockable"`
	Flag            int    `json:"flag"`
	GardenFlag      int    `json:"garden_flag"`
	ShopGroupDecoID int    `json:"shop_group_deco_id"`
	Name            string `json:"-"` // MsgBuildingName_en.strb
	Description     string `json:"-"` // MsgBuildingDesc_en.strb
	_levels         []StructureLevel
	_numCosts       []StructureCost
	_castleBonus    []CastleLevel
	_debris         *GardenDebris
}

// TextureIDs gets the texture IDs for this structure
func (s *Structure) TextureIDs() []int {
	levels := s.Levels()
	ret := make([]int, 0)
	seen := make(map[int]struct{}, 0)
	for _, level := range levels {
		if _, ok := seen[level.TexID]; !ok {
			ret = append(ret, level.TexID)
			seen[level.TexID] = struct{}{}
		}
		if level.Bank != nil {
			b := level.Bank
			if _, ok := seen[b.LowTexID]; !ok && b.LowTexID > 0 {
				ret = append(ret, b.LowTexID)
				seen[b.LowTexID] = struct{}{}
			}
			if _, ok := seen[b.MidTexID]; !ok && b.MidTexID > 0 {
				ret = append(ret, b.MidTexID)
				seen[b.MidTexID] = struct{}{}
			}
			if _, ok := seen[b.HiTexID]; !ok && b.HiTexID > 0 {
				ret = append(ret, b.HiTexID)
				seen[b.HiTexID] = struct{}{}
			}
		}
		if level.Resource != nil {
			r := level.Resource
			if _, ok := seen[r.JobTexID]; !ok && r.JobTexID > 0 {
				ret = append(ret, r.JobTexID)
				seen[r.JobTexID] = struct{}{}
			}
			if _, ok := seen[r.FixTexID]; !ok && r.FixTexID > 0 {
				ret = append(ret, r.FixTexID)
				seen[r.FixTexID] = struct{}{}
			}
		}
	}
	log.Printf("Tex IDs: %v", ret)
	return ret
}

// "event_structures" lists any structures available in the current event

// StructureLevel structure_level lists the level for the available structures
type StructureLevel struct { // structure_level
	ID            int            `json:"_id"`
	StructureID   int            `json:"structure_id"`
	Level         int            `json:"level"`
	TexID         int            `json:"tex_id"` // texture image id?
	LevelCap      int            `json:"level_cap"`
	UnlockAreaID  int            `json:"unlock_area_id"`
	Time          int            `json:"time"`
	BeginnerTime  int            `json:"beginner_time"`
	Coin          int            `json:"coin"`
	Iron          int            `json:"iron"`
	Ether         int            `json:"ether"`
	Cash          int            `json:"cash"`
	Gem           int            `json:"elixir"`
	BeginnerCoin  int            `json:"beginner_coin"`
	BeginnerIron  int            `json:"beginner_iron"`
	BeginnerEther int            `json:"beginner_ether"`
	BeginnerCash  int            `json:"beginner_cash"`
	BeginnerGem   int            `json:"beginner_elixir"`
	Price         int            `json:"price"`
	Exp           int            `json:"exp"`
	ItemID1       int            `json:"item_id_1"`
	ItemNum1      int            `json:"item_num_1"`
	ItemID2       int            `json:"item_id_2"`
	ItemNum2      int            `json:"item_num_2"`
	ItemID3       int            `json:"item_id_3"`
	ItemNum3      int            `json:"item_num_3"`
	Resource      *ResourceLevel `json:"-"`
	Bank          *BankLevel     `json:"-"`
	SpecialEffect *SpecialEffect `json:"-"`
}

// StructureCost cost to build a structure
type StructureCost struct { // structure_num_cost
	ID            int `json:"_id"`
	Num           int `json:"num"`
	StructureID   int `json:"structure_id"`
	Coin          int `json:"coin"`
	Iron          int `json:"iron"`
	Ether         int `json:"ether"`
	Cash          int `json:"cash"`
	Gem           int `json:"elixir"`
	BeginnerCoin  int `json:"beginner_coin"`
	BeginnerIron  int `json:"beginner_iron"`
	BeginnerEther int `json:"beginner_ether"`
	BeginnerCash  int `json:"beginner_cash"`
	BeginnerGem   int `json:"beginner_elixir"`
}

// ResourceLevel amount of resources needed to level up a structure
type ResourceLevel struct { // resource
	ID           int `json:"_id"`
	StructureID  int `json:"structure_id"`
	Level        int `json:"level"`
	IntervalTime int `json:"interval_time"`
	JobTexID     int `json:"job_tex_id"`
	FixTexID     int `json:"fix_tex_id"`
	ResourceID   int `json:"resource_id"`
	Income       int `json:"income"`
}

// BankLevel the amount of resouces a bank can hold at a certain level
type BankLevel struct { //bank_level
	ID          int `json:"_id"`
	StructureID int `json:"structure_id"`
	Level       int `json:"level"`
	LowTexID    int `json:"low_tex_id"` // low level texture image ID
	MidTexID    int `json:"mid_tex_id"` // mid level texture image ID
	HiTexID     int `json:"hi_tex_id"`  // high level texture image ID
	ResourceID  int `json:"resource_id"`
	Value       int `json:"value"`
}

// CastleLevel Castle level information
type CastleLevel struct { //castle_level
	ID                int `json:"_id"`
	CastleStructureID int `json:"castle_structure_id"`
	Level             int `json:"level"`
	StructureID       int `json:"structure_id"`
	BaseAdd           int `json:"base_add"`
	Max               int `json:"max"`
}

// DecoWarehouse Deco Warehouse info
type DecoWarehouse struct { //deco_warehouse
	ID        int `json:"_id"`
	Level     int `json:"level"`
	StockSize int `json:"stock_size"`
}

// DecoResource Deco Warehouse info
type DecoResource struct { //deco_resource
	ID           int `json:"_id"`
	StructureID  int `json:"structure_id"`
	Level        int `json:"level"`
	IntervalTime int `json:"interval_time"`
	ResourceID   int `json:"resource_id"`
	Income       int `json:"income"`
	Coefficient  int `json:"coefficient"`
	Base         int `json:"base"`
	Max          int `json:"max"`
	CollectTime  int `json:"collect_time"`
}

//SpecialEffect special effect definitions of structures
type SpecialEffect struct { // special_effect
	ID              int `json:"_id"`
	StructureID     int `json:"structure_id"`
	Level           int `json:"level"`
	SpecialEffectID int `json:"special_effect_id"`
	Param1          int `json:"param1"`
	Param2          int `json:"param2"`
	Param3          int `json:"param3"`
	Param4          int `json:"param4"`
}

// Levels of a structure
func (s *Structure) Levels() []StructureLevel {
	if s._levels == nil {
		s._levels = make([]StructureLevel, 0)
		for _, l := range Data.StructureLevels {
			if l.StructureID == s.ID {
				l.cacheResource()      // cache off the resource level for later
				l.cacheBank()          // cache off the bank level for later
				l.cacheSpecialEffect() // cache off the bank level for later
				s._levels = append(s._levels, l)
			}
		}
	}
	return s._levels
}

// PurchaseCosts costs of purchasing the structure
func (s *Structure) PurchaseCosts() []StructureCost {
	if s._numCosts == nil {
		s._numCosts = make([]StructureCost, 0)
		for _, p := range Data.StructureNumCosts {
			if p.StructureID == s.ID {
				s._numCosts = append(s._numCosts, p)
			}
		}
	}
	return s._numCosts
}

// IsResource returns true if the building is a resource building
func (s *Structure) IsResource() bool {
	for _, l := range s.Levels() {
		if l.Resource != nil {
			return true
		}
	}
	return false
}

// IsBank returns true if the structure is a resource bank
func (s *Structure) IsBank() bool {
	for _, l := range s.Levels() {
		if l.Bank != nil {
			return true
		}
	}
	return false
}

// CastleBonuses bonuses obtained by leveling up the castle
func (s *Structure) CastleBonuses() []CastleLevel {
	if s._castleBonus == nil {
		s._castleBonus = make([]CastleLevel, 0)
		for _, cl := range Data.CastleLevels {
			if cl.StructureID == s.ID {
				s._castleBonus = append(s._castleBonus, cl)
			}
		}
	}
	return s._castleBonus
}

// MaxQty Number a player can own
func (s *Structure) MaxQty() int {
	clbs := s.CastleBonuses()
	if len(clbs) > 0 {
		maxCLB := 0
		for i, clb := range clbs {
			if clb.Level > clbs[maxCLB].Level {
				maxCLB = i
			}
		}
		return clbs[maxCLB].Max
	}
	return 1
}

func (l *StructureLevel) cacheResource() {
	if l.Resource == nil {
		for i, sr := range Data.ResourceLevels {
			if sr.StructureID == l.StructureID && sr.Level == l.Level {
				l.Resource = &(Data.ResourceLevels[i])
				break
			}
		}
	}
}
func (l *StructureLevel) cacheBank() {
	if l.Bank == nil {
		for i, br := range Data.BankLevels {
			if br.StructureID == l.StructureID && br.Level == l.Level {
				l.Bank = &(Data.BankLevels[i])
				break
			}
		}
	}
}
func (l *StructureLevel) cacheSpecialEffect() {
	if l.SpecialEffect == nil {
		for i, br := range Data.SpecialEffects {
			if br.StructureID == l.StructureID && br.Level == l.Level {
				l.SpecialEffect = &(Data.SpecialEffects[i])
				break
			}
		}
	}
}

// Rate rate of resouce gained
func (sr *ResourceLevel) Rate() int {
	return sr.Income / sr.IntervalTime
}

// FillTime time until the resource is full and won't produce anymore
func (sr *ResourceLevel) FillTime() time.Duration {
	return time.Duration(float64(sr.Income)/float64(sr.Rate())*60.0) * time.Second
}

// StructureScan searches for a structure by ID
func StructureScan(id int) *Structure {
	if id > 0 {
		if id < len(Data.Structures) && Data.Structures[id-1].ID == id {
			return &Data.Structures[id-1]
		}
		for k, val := range Data.Structures {
			if val.ID == id {
				return &(Data.Structures[k])
			}
		}
	}
	return nil
}

// StructureType types of structures
var StructureType = []string{
	"",                    // 0
	"Farm",                // 1
	"Iron Works",          // 2
	"Ether Furnace",       // 3
	"Gold Storehouse",     // 4
	"Iron Storehouse",     // 5
	"Ether Storehouse",    // 6
	"Workshop",            // 7
	"Castle",              // 8
	"Market",              // 9
	"Decorations",         // 10
	"Player Improvements", // 11
	"Alliance",            // 12
	"Deco Storage",        // 13
	"Amusement",           // 14
	"Ward",                // 15
	"Debris",              // 16
	"Kingdom Gate",        // 17
	"Resort Hotel",        // 18
	"Awakening Lab",       // 19
	"Gem Mine",            // 20
	"Gem Storage",         // 21
	"Treasure Hunt",       // 22
}
