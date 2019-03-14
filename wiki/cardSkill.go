package wiki

import (
	"fmt"
	"html"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

//CardSkill skill information
type CardSkill struct {
	EvoID        string
	Name         string
	IDMod        string
	Activations  int
	MinEffect    string
	MaxEffect    string
	RandomSkills []string
	Expiration   *time.Time
}

type tmpSkillHolder struct {
	Skill         *vc.Skill
	SkillNum      int
	SkillFirstEvo string
}
type tmpSkillsSeen map[*vc.Skill]tmpSkillHolder

func getSkills(c *vc.Card, evoKey string, tSkillsSeen *tmpSkillsSeen) []CardSkill {
	ret := make([]CardSkill, 0, 4)
	// ignore mid-evo skills as they are always the same as the first evo.
	if c == nil || c.EvoIsMidOf4() {
		return ret
	}

	addSkill := func(s *vc.Skill, ls *vc.Skill, num int, mod string) {
		if s == nil && ls == nil {
			return
		}
		if _, seen := (*tSkillsSeen)[s]; seen {
			return
		}
		(*tSkillsSeen)[s] = tmpSkillHolder{
			Skill:         s,
			SkillNum:      num,
			SkillFirstEvo: evoKey,
		}
		// need to find if this is an evo-maxed skill
		min := s.SkillMin()
		max := s.SkillMax()
		// thor skills use the "Fire" text
		if mod == "t" {
			min = s.FireMax() + " / 100% chance"
			max = s.FireMax() + " / 100% chance"
		}
		if ls != nil {
			// thor skills use the "Fire" text
			if mod == "t" {
				max = ls.FireMax() + " / 100% chance"
			} else {
				max = ls.SkillMax()
			}
			(*tSkillsSeen)[ls] = tmpSkillHolder{
				Skill:         ls,
				SkillNum:      num,
				SkillFirstEvo: evoKey,
			}
		}
		tSkill := CardSkill{
			EvoID:        evoKey,
			IDMod:        mod,
			Name:         s.Name,
			Activations:  s.MaxCount,
			MinEffect:    min,
			MaxEffect:    max,
			RandomSkills: make([]string, 0, 5),
		}
		if s.EffectID == 36 {
			// Random Skill
			for _, v := range []int{s.EffectParam, s.EffectParam2, s.EffectParam3, s.EffectParam4, s.EffectParam5} {
				rs := vc.SkillScan(v)
				if rs != nil {
					tSkill.RandomSkills = append(tSkill.RandomSkills, rs.FireMin())
				}
			}
		}
		if s.Expires() {
			tSkill.Expiration = &s.PublicEndDatetime.Time
		}
		ret = append(ret, tSkill)
	}

	var lastEvo *vc.Card
	if c.EvoIsFirst() {
		lastEvo = c.LastEvo()
	}
	addSkill(c.Skill1(), lastEvo.Skill1(), 1, "")
	addSkill(c.Skill2(), lastEvo.Skill2(), 2, "2")
	addSkill(c.Skill3(), lastEvo.Skill3(), 3, "3")
	addSkill(c.ThorSkill1(), lastEvo.ThorSkill1(), 4, "t")

	return ret
}

// actual output here
func (s CardSkill) String() string {
	evoMod := s.EvoID + s.IDMod
	if s.EvoID == "0" || s.EvoID == "H" {
		evoMod = s.IDMod
	}
	evoMod += " "

	lvl10 := ""
	if s.MinEffect != s.MaxEffect && s.MaxEffect != "" {
		lvl10 = fmt.Sprintf("\n|skill %slv10 = %s",
			evoMod,
			html.EscapeString(strings.Replace(s.MaxEffect, "\n", "<br />", -1)),
		)
	}

	ret := fmt.Sprintf(`|skill %[1]s= %[2]s
|skill %[1]slv1 = %[3]s%[4]s
|procs %[1]s= %[5]d
`,
		evoMod,
		html.EscapeString(s.Name),
		html.EscapeString(strings.Replace(s.MinEffect, "\n", "<br />", -1)),
		lvl10,
		s.Activations,
	)

	for k, v := range s.RandomSkills {
		ret += fmt.Sprintf("|random %s%d = %s \n",
			evoMod,
			k+1,
			html.EscapeString(strings.Replace(v, "\n", "<br />", -1)),
		)
	}

	if s.Expiration != nil {
		ret += fmt.Sprintf("|skill %send = %v\n", evoMod, s.Expiration)
	}
	return ret
}
