package vc_grouper

// Skills info from master data field "skills" These match to the string in the files:
// MsgSkillName_en.strb
// MsgSkillDesc_en.strb - shown on the card
// MsgSkillFire_en.strb - used during battle
type Skill struct {
	// skill id
	_id,
	// level type for skill upgrade costs
	level_type,
	// skill type
	_type,
	// id for timing
	timing_id,
	// max procs
	max_count,
	// cond scene
	cond_scene_id,
	// cond side
	cond_side_id,
	// cond
	cond_id,
	// king series
	king_series_id,
	// king id
	king_id,
	// cond param
	cond_param,
	// default proc rate
	default_ratio,
	// max proc rate
	max_ratio,
	// date accessible
	public_start_datetime,
	public_end_datetime,
	// effect info
	effect_id,
	effect_param,
	effect_param_2,
	effect_param_3,
	effect_param_4,
	effect_param_5,
	effect_default_value,
	effect_max_value,
	// target info
	target_scope_id,
	target_logic_id,
	target_param,
	// animation info
	animation_id int

	//skill name from strings file
	name,
	// description from strings file
	description,
	// fire text from strings file
	fire string
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
