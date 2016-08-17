package vc

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Skills info from master data field "skills" These match to the string in the files:
// MsgSkillName_en.strb
// MsgSkillDesc_en.strb - shown on the card
// MsgSkillFire_en.strb - used during battle
type Skill struct {
	Id           int `json:"_id"`            // skill id
	LevelType    int `json:"level_type"`     // level type for skill upgrade costs
	Type         int `json:"_type"`          // skill type
	TimingId     int `json:"timing_id"`      // id for timing
	MaxCount     int `json:"max_count"`      // max procs
	CondSceneId  int `json:"cond_scene_id"`  // cond scene
	CondSideId   int `json:"cond_side_id"`   // cond side
	CondId       int `json:"cond_id"`        // cond
	KingSeriesId int `json:"king_series_id"` // king series
	KingId       int `json:"king_id"`        // king id
	CondParam    int `json:"cond_param"`     // cond param
	DefaultRatio int `json:"default_ratio"`  // default proc rate
	MaxRatio     int `json:"max_ratio"`      // max proc rate
	// date accessible
	PublicStartDatetime Timestamp `json:"public_start_datetime"`
	PublicEndDatetime   Timestamp `json:"public_end_datetime"`
	// effect info
	EffectId           int `json:"effect_id"`
	EffectParam        int `json:"effect_param"`
	EffectParam2       int `json:"effect_param_2"`
	EffectParam3       int `json:"effect_param_3"`
	EffectParam4       int `json:"effect_param_4"`
	EffectParam5       int `json:"effect_param_5"`
	EffectDefaultValue int `json:"effect_default_value"`
	EffectMaxValue     int `json:"effect_max_value"`
	// target info
	TargetScopeId int `json:"target_scope_id"`
	TargetLogicId int `json:"target_logic_id"`
	TargetParam   int `json:"target_param"`
	// animation info
	AnimationId             int             `json:"animation_id"`
	ThorHammerAnimationType json.RawMessage `json:"thorhammer_animation_type"`

	//skill name from strings file
	Name string `json:"name"`
	// description from strings file
	Description string `json:"description"`
	// fire text from strings file
	Fire string `json:"fire"`
}

func (s *Skill) Effect() string {
	if val, ok := Effect[s.EffectId]; ok {
		return val
	} else {
		return "New/Unknown"
	}
}

func (s *Skill) SkillMin() string {
	return formatSkill(s.Description, s.EffectDefaultValue, s.DefaultRatio)
}

func (s *Skill) SkillMax() string {
	return formatSkill(s.Description, s.EffectMaxValue, s.MaxRatio)
}

func (s *Skill) FireMin() string {
	return formatSkill(s.Fire, s.EffectDefaultValue, -1)
}

func (s *Skill) FireMax() string {
	return formatSkill(s.Fire, s.EffectMaxValue, -1)
}

func (s *Skill) TargetScope() string {
	if val, ok := TargetScope[s.TargetScopeId]; ok {
		return val
	} else {
		return ""
	}
}

func (s *Skill) TargetLogic() string {
	if val, ok := TargetLogic[s.TargetLogicId]; ok {
		return val
	} else {
		return ""
	}
}

func SkillScan(id int, skills []Skill) *Skill {
	if id <= 0 {
		return nil
	}
	if id < len(skills) && id == skills[id-1].Id {
		return &(skills[id-1])
	}
	for k, v := range skills {
		if id == v.Id {
			return &(skills[k])
		}
	}
	return nil
}

func formatSkill(descr string, effect, ratio int) (s string) {
	eff := strconv.Itoa(effect)
	r := strconv.Itoa(ratio)

	s = strings.Replace(descr, "{1:x}", eff, -1)
	s = strings.Replace(s, "{1:}", eff, -1)
	s = strings.Replace(s, "{1}", eff, -1)

	s = strings.Replace(s, "{2:}", r, -1)
	s = strings.Replace(s, "{2:x}", r, -1)
	s = strings.Replace(s, "{2}", r, -1)
	return s
}

var TargetScope = map[int]string{
	-1: "",
	1:  "Allies",
	2:  "Enemies",
}

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
