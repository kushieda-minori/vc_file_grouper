package bot

import (
	"fmt"
	"log"
	"regexp"
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
		if c.EvoIsAwoken() {
			sPrefix = "Awoken"
		} else if c.EvoIsReborn() {
			sPrefix = "Reborn"
		}
		skills = append(skills, newSkill(sPrefix, s))
	}
	// skill 2
	s = c.Skill2()
	if s != nil && !s.Expires() {
		log.Printf("Found skill-2 on card %d:%s", c.ID, c.Name)
		sPrefix := "Second Skill"
		if c.EvoIsAwoken() {
			sPrefix = "Awoken " + sPrefix
		} else if c.EvoIsReborn() {
			sPrefix = "Reborn " + sPrefix
		}
		skills = append(skills, newSkill(sPrefix, s))
	}

	if a := c.LastEvo().AwakensTo(); a != nil {
		// awakening for 99% of awoken cards
		log.Printf("Found awakening on card %d:%s -> %d:%s", c.ID, c.Name, a.ID, a.Name)
		tmp := newSkills(a)
		skills = append(skills, tmp...)
		// the recursion will catch any rebirths
	} else if a, ok := c.GetEvolutions()["G"]; ok && a.ID != c.ID && !c.EvoIsReborn() {
		// this should find awakenings for amal cards (like towers)
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
	activations := s.ActivationString()
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

func cleanSkill(name string) string {
	//bot icons:
	//'<:Passion:534375884480577536>'
	//'<:Cool:534375884598018049>'
	//'<:Light:534375885541998602>'
	//'<:Dark:534375884279382016>'
	//'<:Special:534375884493291530>'
	name = strings.ReplaceAll(name, "{{Passion}}", "<:Passion:534375884480577536>")
	name = strings.ReplaceAll(name, "{{Cool}}", "<:Cool:534375884598018049>")
	name = strings.ReplaceAll(name, "{{Light}}", "<:Light:534375885541998602>")
	name = strings.ReplaceAll(name, "{{Dark}}", "<:Dark:534375884279382016>")
	name = strings.ReplaceAll(name, "{{Special}}", "<:Special:534375884493291530>")
	name = strings.ReplaceAll(strings.ReplaceAll(name, "{{", ""), "}}", "")
	re := regexp.MustCompile(`\s+[/\\]\s+Max \d+ time(s)?`)
	name = re.ReplaceAllString(name, "")
	return name
}
