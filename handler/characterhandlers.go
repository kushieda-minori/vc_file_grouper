package handler

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// CharacterTableHandler show character data in a table format
func CharacterTableHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	filter := func(character *vc.CardCharacter) (match bool) {
		match = true
		if len(qs) < 1 {
			return
		}
		card := character.FirstEvoCard()
		if card == nil {
			return false
		}
		if isThor := qs.Get("isThor"); isThor != "" {
			match = match && card.ThorSkillID1 > 0
		}
		if name := qs.Get("name"); name != "" {
			match = match && strings.Contains(strings.ToLower(card.Name), strings.ToLower(name))
		}
		if skillname := qs.Get("skillname"); skillname != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(); skill1 != nil {
				s1 = skill1.Name != "" && strings.Contains(strings.ToLower(skill1.Name), strings.ToLower(skillname))
				//log.Printf(skill1.Name + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(); skill2 != nil {
				s2 = skill2.Name != "" && strings.Contains(strings.ToLower(skill2.Name), strings.ToLower(skillname))
				//log.Printf(skill2.Name + " " + strconv.FormatBool(s2) + "\n")
			}
			match = match && (s1 || s2)
		}
		if skilldesc := qs.Get("skilldesc"); skilldesc != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(); skill1 != nil {
				s1 = skill1.Fire != "" && strings.Contains(strings.ToLower(skill1.Fire), strings.ToLower(skilldesc))
				//log.Printf(skill1.Fire + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(); skill2 != nil {
				s2 = skill2.Fire != "" && strings.Contains(strings.ToLower(skill2.Fire), strings.ToLower(skilldesc))
				//log.Printf(skill2.Fire + " " + strconv.FormatBool(s2) + "\n")
			}
			match = match && (s1 || s2)
		}
		return
	}
	// File header
	fmt.Fprintf(w, `<html><head><title>All Characters</title>
<style>table, th, td {border: 1px solid black;};</style>
</head>
<body>
<form method="GET">
<label for="f_name">Name:</label><input id="f_name" name="name" value="%s" />
<label for="f_skillname">Skill Name:</label><input id="f_skillname" name="skillname" value="%s" />
<label for="f_skilldesc">Skill Description:</label><input id="f_skilldesc" name="skilldesc" value="%s" />
<label for="f_skillisthor">Has Thor Skill:</label><input id="f_skillisthor" name="isThor" type="checkbox" value="checked" %s />
<button type="submit">Submit</button>
</form>
<div>
<table><thead><tr>
<th>_id</th>
<th>card_nos</th>
<th>name</th>
<th>Description</th>
<th>Friendship</th>
<th>Login</th>
<th>Meet</th>
<th>Battle Start</th>
<th>Battle End</th>
<th>Friendship Max</th>
<th>Friendship Event</th>
</tr></thead>
<tbody>
`,
		qs.Get("name"),
		qs.Get("skillname"),
		qs.Get("skilldesc"),
		qs.Get("isThor"),
	)

	// sort the characters by most recent card
	// copy the character data so we don't modify the inline global
	chars := make([]vc.CardCharacter, len(vc.Data.CardCharacters))
	copy(chars, vc.Data.CardCharacters)

	sort.Slice(chars, func(i, j int) bool {
		first := chars[i]
		second := chars[j]

		firstCards := vc.CardList(first.Cards())
		secondCards := vc.CardList(second.Cards())

		maxFirst := firstCards.Latest()
		maxSecond := secondCards.Latest()

		if maxFirst == nil && maxSecond == nil {
			return first.ID > second.ID
		}
		if maxFirst == nil {
			return false
		}
		if maxSecond == nil {
			return true
		}

		return maxFirst.ID > maxSecond.ID
	})

	for _, character := range chars {
		if !filter(&character) {
			continue
		}

		card := character.FirstEvoCard()

		cardName := "N/A"

		if card != nil {
			if card.Name == "" {
				cardName = card.Image()
			} else {
				cardName = card.Name
			}
		}

		cardNos := ""
		for _, card := range character.Cards() {
			if len(cardNos) > 0 {
				cardNos = cardNos + ", "
			}
			cardNos = cardNos + strconv.Itoa(card.CardNo)
		}

		fmt.Fprintf(w, "<tr><td>%d</td>"+
			"<td>%s</td>"+
			"<td><a href=\"/characters/detail/%[1]d\">%[3]s</a></td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td></tr>\n",
			character.ID,
			cardNos,
			cardName,
			character.Description,
			character.Friendship,
			character.Login,
			character.Meet,
			character.BattleStart,
			character.BattleEnd,
			character.FriendshipMax,
			character.FriendshipEvent,
		)
	}

	io.WriteString(w, "</tbody></table></div></body></html>")
}

// CharacterDetailHandler show character details
func CharacterDetailHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	var pathLen int
	if path[len(path)-1] == '/' {
		pathLen = len(path) - 1
	} else {
		pathLen = len(path)
	}

	pathParts := strings.Split(path[1:pathLen], "/")
	// "characters/detail/id"
	if len(pathParts) < 3 {
		http.Error(w, "Invalid character id ", http.StatusNotFound)
		return
	}
	charID, err := strconv.Atoi(pathParts[2])
	if err != nil || charID < 1 || charID > len(vc.Data.CardCharacters) {
		http.Error(w, "Invalid character id "+pathParts[2], http.StatusNotFound)
		return
	}

	character := vc.CardCharacterScan(charID)

	cards := character.Cards()

	if len(cards) == 0 {
		http.Error(w, "Character has no cards "+pathParts[2], http.StatusNotFound)
		return
	}

	// sort by Evolution Rank
	sort.Slice(cards, func(i, j int) bool {
		c1 := cards[i]
		c2 := cards[j]
		if c1.EvolutionRank == c2.EvolutionRank {
			return c1.CardNo < c2.CardNo
		}
		if c1.EvolutionRank < c2.EvolutionRank {
			return true
		}
		return false
	})

	fmt.Fprintf(w, `<html><head><title>All Cards</title>
<style>table, th, td {border: 1px solid black;};</style>
</head>
<body>
<div>
<table><thead><tr>
<th>_id</th>
<th>card_no</th>
<th>name</th>
<th>evolution_rank</th>
<th>max_evolution_rank</th>
<th>Next Evo</th>
<th>Rarity</th>
<th>Element</th>
<th>Character ID</th>
<th>deck_cost</th>
<th>default_offense</th>
<th>default_defense</th>
<th>default_follower</th>
<th>max_offense</th>
<th>max_defense</th>
<th>max_follower</th>
<th>Skill 1 Name</th>
<th>Skill Min</th>
<th>Skill Max</th>
<th>Skill Procs</th>
<th>Min Effect</th>
<th>Min Rate</th>
<th>Max Effect</th>
<th>Max Rate</th>
<th>Target Scope</th>
<th>Target Logic</th>
<th>Skill 2</th>
<th>Skill 3</th>
<th>Thor Skill</th>
<th>Skill Special</th>
<th>Description</th>
<th>Friendship</th>
<th>Login</th>
<th>Meet</th>
<th>Battle Start</th>
<th>Battle End</th>
<th>Friendship Max</th>
<th>Friendship Event</th>
</tr></thead>
<tbody>
`)

	for _, card := range cards {
		skill1 := card.Skill1()
		if skill1 == nil {
			skill1 = &vc.Skill{}
		}
		// skill2 := card.Skill2()
		// skillS1 := card.SpecialSkill1()
		fmt.Fprintf(w, "<tr><td>%d</td>"+
			"<td><a href=\"/cards/detail/%[1]d\">%05[2]d</a></td>"+
			"<td><a href=\"/cards/detail/%[1]d\">%[3]s</a></td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%d</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td>"+
			"<td>%s</td></tr>\n",
			card.ID,
			card.CardNo,
			card.Name,
			card.EvolutionRank,
			card.LastEvolutionRank,
			card.EvolutionCardID,
			card.Rarity(),
			card.Element(),
			card.CardCharaID,
			card.DeckCost,
			card.DefaultOffense,
			card.DefaultDefense,
			card.DefaultFollower,
			card.MaxOffense,
			card.MaxDefense,
			card.MaxFollower,
			card.Skill1Name(),
			card.SkillMin(),
			card.SkillMax(),
			card.SkillProcs(),
			skill1.EffectDefaultValue,
			skill1.DefaultRatio,
			skill1.EffectMaxValue,
			skill1.MaxRatio,
			card.SkillTarget(),
			card.SkillTargetLogic(),
			card.Skill2Name(),
			card.Skill3Name(),
			card.ThorSkill1Name(),
			card.SpecialSkill1Name(),
			card.Description(),
			card.Friendship(),
			card.Login(),
			card.Meet(),
			card.BattleStart(),
			card.BattleEnd(),
			card.FriendshipMax(),
			card.FriendshipEvent(),
		)
	}

	io.WriteString(w, "</tbody></table></div></body></html>")
}
