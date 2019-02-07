package nobu

import (
	"fmt"
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
	skills := make([]Skill, 0)

	// skill 1
	s := c.Skill1()
	if s != nil {
		activations := getActivations(s)
		sPrefix := "Skill"
		if len(c.Rarity()) == 3 {
			if c.Rarity()[0] == 'G' {
				sPrefix = "Awoken"
			} else if c.Rarity()[0] == 'X' {
				sPrefix = "Reborn"
			}
		}
		skills = append(skills, Skill{
			Name: fmt.Sprintf("%s: %s", sPrefix, s.Name), // Skill, Second Skill, Awoken, Awoken Second Skill
			Value: fmt.Sprintf("Activations: %s\nMin Level Effect: %s\nMax Level Effect: %s",
				activations,
				cleanSkill(s.SkillMin()),
				cleanSkill(s.SkillMax()),
			),
		})
	}
	// skill 2
	s = c.Skill2()
	if s != nil && !s.Expires() {
		activations := getActivations(s)
		sPrefix := "Second Skill"
		if len(c.Rarity()) == 3 {
			if c.Rarity()[0] == 'G' {
				sPrefix = "Awoken " + sPrefix
			} else if c.Rarity()[0] == 'X' {
				sPrefix = "Reborn " + sPrefix
			}
		}
		skills = append(skills, Skill{
			Name: fmt.Sprintf("%s: %s", sPrefix, s.Name), // Skill, Second Skill, Awoken, Awoken Second Skill
			Value: fmt.Sprintf("Activations: %s\nMin Level Effect: %s\nMax Level Effect: %s",
				activations,
				cleanSkill(s.SkillMin()),
				cleanSkill(s.SkillMax()),
			),
		})
	}

	if a := c.LastEvo().AwakensTo(); a != nil {
		tmp := newSkills(a)
		skills = append(skills, tmp...)
		// the recursion will catch any rebirths
	}
	// will only pick up rebirths if we are looking at the awoken card.
	if r := c.RebirthsTo(); r != nil {
		tmp := newSkills(r)
		skills = append(skills, tmp...)
	}

	return skills
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
	name = strings.Replace(strings.Replace(name, "{{", "", -1), "}}", "", -1)
	re := regexp.MustCompile(`\s+[/\\]\s+Max \d+ time(s)?`)
	name = re.ReplaceAllString(name, "")
	return name
}
