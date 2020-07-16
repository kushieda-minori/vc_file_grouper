package vc

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Skill info from master data field "skills" These match to the string in the files:
// MsgSkillName_en.strb
// MsgSkillDesc_en.strb - shown on the card
// MsgSkillFire_en.strb - used during battle
type Skill struct {
	ID           int `json:"_id"`            // skill id
	LevelType    int `json:"level_type"`     // level type for skill upgrade costs
	Type         int `json:"_type"`          // skill type
	TimingID     int `json:"timing_id"`      // id for timing
	MaxCount     int `json:"max_count"`      // max procs
	CondSceneID  int `json:"cond_scene_id"`  // cond scene
	CondSideID   int `json:"cond_side_id"`   // cond side
	CondID       int `json:"cond_id"`        // cond
	KingSeriesID int `json:"king_series_id"` // king series
	KingID       int `json:"king_id"`        // king id
	CondParam    int `json:"cond_param"`     // cond param
	DefaultRatio int `json:"default_ratio"`  // default proc rate
	MaxRatio     int `json:"max_ratio"`      // max proc rate
	// date accessible
	PublicStartDatetime Timestamp `json:"public_start_datetime"`
	PublicEndDatetime   Timestamp `json:"public_end_datetime"`
	// effect info
	EffectID           int `json:"effect_id"`
	EffectParam        int `json:"effect_param"`
	EffectParam2       int `json:"effect_param_2"`
	EffectParam3       int `json:"effect_param_3"`
	EffectParam4       int `json:"effect_param_4"`
	EffectParam5       int `json:"effect_param_5"`
	EffectDefaultValue int `json:"effect_default_value"`
	EffectMaxValue     int `json:"effect_max_value"`
	// target info
	TargetScopeID int `json:"target_scope_id"`
	TargetLogicID int `json:"target_logic_id"`
	TargetParam   int `json:"target_param"`
	// animation info
	AnimationID             int             `json:"animation_id"`
	ThorHammerAnimationType json.RawMessage `json:"thorhammer_animation_type"`

	Name         string `json:"name"`        //skill name from strings file
	Description  string `json:"description"` // description from strings file
	Fire         string `json:"fire"`        // fire text from strings file
	_skillLevels []SkillLevel
}

// SkillLevel information
type SkillLevel struct {
	ID        int `json:"_id"`        // skill id
	LevelType int `json:"level_type"` // level type for skill upgrade costs
	Level     int `json:"level"`      // level
	Medal     int `json:"medal"`      // medals to upgrade to the next level
}

// CustomSkillLevel level information for custom skills
type CustomSkillLevel struct {
	ID int `json:""`
}

// SkillCostIncrementPattern increment pattern for custom skill costs vs card levels
type SkillCostIncrementPattern struct {
	ID                int `json:""`
	CardLevelInterval int `json:"card_level_interval"`
	Increment         int `json:"increment"`
	Max               int `json:"max"`
}

// Expires true if the skill has an expiration date
func (s *Skill) Expires() bool {
	if s == nil {
		return false
	}
	return s.PublicEndDatetime.After(time.Time{})
}

// Effect of the skill (this is a visual effect, not the ability)
func (s *Skill) Effect() string {
	if s == nil {
		return ""
	}
	if s.EffectID > 0 {
		if val, ok := Effect[s.EffectID]; ok {
			return val
		}
		return "New/Unknown"
	}
	return ""
}

// Activations returns the number of times a skill can activate.
// a negative number indicates infinite procs
func (s *Skill) Activations() int {
	if s == nil {
		return 0
	}
	if s.MaxCount > 0 {
		// battle start skills seem to have random Max Count values. Force it to 1
		// since they can only proc once anyway
		if strings.Contains(strings.ToLower(s.SkillMin()), "battle start") {
			return 1
		}
		return s.MaxCount
	}
	return s.MaxCount
}

// ActivationString returns the number of times a skill can activate.
func (s *Skill) ActivationString() string {
	if s == nil {
		return ""
	}
	activations := s.Activations()
	if activations > 0 {
		return strconv.Itoa(activations)
	} else if strings.Contains(s.SkillMin(), "【Autoskill】") {
		return "Always On"
	} else {
		return "Infinite"
	}
}

// SkillMin minimum skill ability
func (s *Skill) SkillMin() string {
	if s == nil {
		return ""
	}
	return formatSkill(s.Description, s.EffectDefaultValue, s.DefaultRatio)
}

// SkillMax maximum skill ability
func (s *Skill) SkillMax() string {
	if s == nil {
		return ""
	}
	return formatSkill(s.Description, s.EffectMaxValue, s.MaxRatio)
}

// FireMin minimum skill ability fire information
func (s *Skill) FireMin() string {
	if s == nil {
		return ""
	}
	return formatSkill(s.Fire, s.EffectDefaultValue, -1)
}

// FireMax maximum skill ability fire information
func (s *Skill) FireMax() string {
	if s == nil {
		return ""
	}
	return formatSkill(s.Fire, s.EffectMaxValue, -1)
}

// Levels level information for the skill
func (s *Skill) Levels() []SkillLevel {
	if s == nil {
		return nil
	}

	if s._skillLevels == nil {
		s._skillLevels = make([]SkillLevel, 0)
		for _, sl := range Data.SkillLevels {
			if sl.LevelType == s.LevelType {
				s._skillLevels = append(s._skillLevels, sl)
			}
		}
	}
	return s._skillLevels
}

// TargetScope scope of the target (enemy or allies)
func (s *Skill) TargetScope() string {
	if s == nil {
		return ""
	}
	if s.TargetScopeID > 0 {
		if val, ok := TargetScope[s.TargetScopeID]; ok {
			return val
		}
		return "Unknown - " + strconv.Itoa(s.TargetScopeID)
	}
	return ""
}

// TargetLogic what target to hit
func (s *Skill) TargetLogic() string {
	if s == nil {
		return ""
	}
	if s.TargetLogicID > 0 {
		if val, ok := TargetLogic[s.TargetLogicID]; ok {
			return val
		}
		return "Unknown - " + strconv.Itoa(s.TargetLogicID)
	}
	return ""
}

// SkillScan searches for a skill by ID
func SkillScan(id int) *Skill {
	if id > 0 {
		l := len(Data.Skills)
		i := sort.Search(l, func(i int) bool { return Data.Skills[i].ID >= id })
		if i >= 0 && i < l && Data.Skills[i].ID == id {
			return &(Data.Skills[i])
		}
	}
	return nil
}

// MaxSkillID for the skill in the list
func MaxSkillID(skills []Skill) (max int) {
	max = 0
	for _, val := range skills {
		if val.ID > max {
			max = val.ID
		}
	}
	return
}

func formatSkill(descr string, effect, ratio int) (s string) {
	eff := strconv.Itoa(effect)
	r := strconv.Itoa(ratio)

	s = strings.ReplaceAll(descr, "{1:x}", eff)
	s = strings.ReplaceAll(s, "{1:}", eff)
	s = strings.ReplaceAll(s, "{1}", eff)

	s = strings.ReplaceAll(s, "{2:}", r)
	s = strings.ReplaceAll(s, "{2:x}", r)
	s = strings.ReplaceAll(s, "{2}", r)
	return s
}

// TargetScope to hit whom
var TargetScope = map[int]string{
	-1: "",
	0:  "",
	1:  "Allies",
	2:  "Enemies",
}

// TargetLogic Target Scope detail
var TargetLogic = map[int]string{
	1:  "Target Field",
	2:  "Lowest HP",
	3:  "Single Random",
	4:  "unknown, maybe max HP",
	8:  "Random Target",
	9:  "Self",
	10: "Opposing Element",
	12: "All Dead",
	13: "Single Dead, Random",
	14: "Same Element",
	16: "Random Target (Skill)",
	17: "Dead and Alive",
	20: "Random Target (Salvo)",
}

// Effect of the skill
var Effect = map[int]string{
	1:  "Heal",
	2:  "Deal Damage",
	3:  "Deal Element Damage",
	4:  "Turn Skip",
	5:  "ATK Up",
	6:  "ATK Down",
	7:  "DEF Up",
	8:  "DEF Down",
	10: "Battle EXP",
	11: "Cancel Buffs / Weak Effects",
	12: "Fully Ressurect",
	13: "Deal Element Damage and Suck",
	14: "Deal Damage and Suck",
	15: "Receive Rewards",
	16: "Unleash Skill",
	17: "Gold Conversion",
	20: "AW Hunt Point+",
	22: "Counter Attack",
	23: "Skill Nullification",
	24: "ATK • DEF Up",
	26: "Resurect and Heal",
	27: "Turn Skip / Deal Element Damage",
	30: "Knock out single Enemy",
	31: "Awakened Burst",
	32: "Elemental Wave",
	35: "Elemental Salvo",
	36: "Random Skill",
	38: "Single Salvo",
}
