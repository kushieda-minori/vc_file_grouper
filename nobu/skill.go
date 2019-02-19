package nobu

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// Skill skill as known by nobu bot
type Skill struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// NewSkills generates a skill array for a given card
func newSkills(c *vc.Card) []Skill {
	log.Printf("looking for skills on card %d:%s", c.ID, c.Name)
	skills := make([]Skill, 0)

	// skill 1
	s := c.Skill1()
	if s != nil {
		log.Printf("Found skill-1 on card %d:%s", c.ID, c.Name)
		sPrefix := "Skill"
		if len(c.Rarity()) == 3 {
			if c.Rarity()[0] == 'G' {
				sPrefix = "Awoken"
			} else if c.Rarity()[0] == 'X' {
				sPrefix = "Reborn"
			}
		}
		skills = append(skills, newSkill(sPrefix, s))
	}
	// skill 2
	s = c.Skill2()
	if s != nil && !s.Expires() {
		log.Printf("Found skill-2 on card %d:%s", c.ID, c.Name)
		sPrefix := "Second Skill"
		if len(c.Rarity()) == 3 {
			if c.Rarity()[0] == 'G' {
				sPrefix = "Awoken " + sPrefix
			} else if c.Rarity()[0] == 'X' {
				sPrefix = "Reborn " + sPrefix
			}
		}
		skills = append(skills, newSkill(sPrefix, s))
	}

	if a := c.LastEvo().AwakensTo(); a != nil {
		log.Printf("Found awakening on card %d:%s -> %d:%s", c.ID, c.Name, a.ID, a.Name)
		tmp := newSkills(a)
		skills = append(skills, tmp...)
		// the recursion will catch any rebirths
	}
	// will only pick up rebirths if we are looking at the awoken card.
	if r := c.RebirthsTo(); r != nil {
		log.Printf("Found rebirth on card %d:%s -> %d:%s", c.ID, c.Name, r.ID, r.Name)
		tmp := newSkills(r)
		skills = append(skills, tmp...)
	}

	log.Printf("Found %d skills total on card %d:%s", len(skills), c.ID, c.Name)
	return skills
}

func newSkill(sPrefix string, s *vc.Skill) Skill {
	activations := getActivations(s)
	sMin := cleanSkill(s.SkillMin())
	sMax := cleanSkill(s.SkillMax())
	skill := ""
	if sMin == sMax {
		skill = fmt.Sprintf("Activations: %s\nEffect: %s", activations, sMin)
	} else {
		skill = fmt.Sprintf("Activations: %s\nMin Level Effect: %s\nMax Level Effect: %s",
			activations,
			sMin,
			sMax,
		)
	}
	return Skill{
		Name:  fmt.Sprintf("%s: %s", sPrefix, s.Name), // Skill, Second Skill, Awoken, Awoken Second Skill
		Value: skill,
	}
}

func getActivations(s *vc.Skill) string {
	if s.MaxCount > 0 {
		return strconv.Itoa(s.MaxCount)
	} else if strings.Contains(s.SkillMin(), "【Autoskill】") {
		return "Always On"
	} else {
		return "Infinite"
	}
}

func cleanSkill(name string) string {
	//bot icons:
	//'<:Passion:534375884480577536>'
	//'<:Cool:534375884598018049>'
	//'<:Light:534375885541998602>'
	//'<:Dark:534375884279382016>'
	//'<:Special:534375884493291530>'
	name = strings.Replace(name, "{{Passion}}", "<:Passion:534375884480577536>", -1)
	name = strings.Replace(name, "{{Cool}}", "<:Cool:534375884598018049>", -1)
	name = strings.Replace(name, "{{Light}}", "<:Light:534375885541998602>", -1)
	name = strings.Replace(name, "{{Dark}}", "<:Dark:534375884279382016>", -1)
	name = strings.Replace(name, "{{Special}}", "<:Special:534375884493291530>", -1)
	name = strings.Replace(strings.Replace(name, "{{", "", -1), "}}", "", -1)
	re := regexp.MustCompile(`\s+[/\\]\s+Max \d+ time(s)?`)
	name = re.ReplaceAllString(name, "")
	return name
}
