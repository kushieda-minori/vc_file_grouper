package vc_grouper

import (
	"strconv"
	"strings"
)

// Skills info from master data field "skills" These match to the string in the files:
// MsgSkillName_en.strb
// MsgSkillDesc_en.strb - shown on the card
// MsgSkillFire_en.strb - used during battle
type Skill struct {
	// skill id
	Id int `json:"_id"`
	// level type for skill upgrade costs
	LevelType int `json:"level_type"`
	// skill type
	Type int `json:"_type"`
	// id for timing
	TimingId int `json:"timing_id"`
	// max procs
	MaxCount int `json:"max_count"`
	// cond scene
	CondSceneId int `json:"cond_scene_id"`
	// cond side
	CondSideId int `json:"cond_side_id"`
	// cond
	CondId int `json:"cond_id"`
	// king series
	KingSeriesId int `json:"king_series_id"`
	// king id
	KingId int `json:"king_id"`
	// cond param
	CondParam int `json:"cond_param"`
	// default proc rate
	DefaultRatio int `json:"default_ratio"`
	// max proc rate
	MaxRatio int `json:"max_ratio"`
	// date accessible
	PublicStartDatetime int `json:"public_start_datetime"`
	PublicEndDatetime   int `json:"public_end_datetime"`
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
	AnimationId int `json:"animation_id"`

	//skill name from strings file
	Name string `json:"-"`
	// description from strings file
	Description string `json:"-"`
	// fire text from strings file
	Fire string `json:"-"`
}

func (s *Skill) GetSkillMin() string {
	min := strings.Replace(s.Description, "{1:x}", strconv.Itoa(s.EffectDefaultValue), -1)
	r := strconv.Itoa(s.DefaultRatio)
	min = strings.Replace(min, "{2:}", r, -1)
	min = strings.Replace(min, "{2:x}", r, -1)
	min = strings.Replace(min, "{2}", r, -1)
	return min
}

func (s *Skill) GetSkillMax() string {
	max := strings.Replace(s.Description, "{1:x}", strconv.Itoa(s.EffectMaxValue), -1)
	r := strconv.Itoa(s.MaxRatio)
	max = strings.Replace(max, "{2:}", r, -1)
	max = strings.Replace(max, "{2:x}", r, -1)
	max = strings.Replace(max, "{2}", r, -1)
	return max
}

func (s *Skill) GetTargetScope() string {
	if val, ok := TargetScope[s.TargetScopeId]; ok {
		return val
	} else {
		return ""
	}
}

func (s *Skill) GetTargetLogic() string {
	if val, ok := TargetLogic[s.TargetLogicId]; ok {
		return val
	} else {
		return ""
	}
}

var TargetScope = map[int]string{
	-1: "N/A",
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
	16: "Random Target Skill",
	17: "Dead and Alive",
}
