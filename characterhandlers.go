package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"zetsuboushita.net/vc_file_grouper/vc"
)

func characterTableHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	filter := func(character *vc.CardCharacter) (match bool) {
		match = true
		if len(qs) < 1 {
			return
		}
		card := character.FirstEvoCard(VcData)
		if card == nil {
			return false
		}
		if isThor := qs.Get("isThor"); isThor != "" {
			match = match && card.ThorSkillId1 > 0
		}
		if name := qs.Get("name"); name != "" {
			match = match && strings.Contains(strings.ToLower(card.Name), strings.ToLower(name))
		}
		if skillname := qs.Get("skillname"); skillname != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(VcData); skill1 != nil {
				s1 = skill1.Name != "" && strings.Contains(strings.ToLower(skill1.Name), strings.ToLower(skillname))
				//os.Stdout.WriteString(skill1.Name + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(VcData); skill2 != nil {
				s2 = skill2.Name != "" && strings.Contains(strings.ToLower(skill2.Name), strings.ToLower(skillname))
				//os.Stdout.WriteString(skill2.Name + " " + strconv.FormatBool(s2) + "\n")
			}
			match = match && (s1 || s2)
		}
		if skilldesc := qs.Get("skilldesc"); skilldesc != "" && match {
			var s1, s2 bool
			if skill1 := card.Skill1(VcData); skill1 != nil {
				s1 = skill1.Fire != "" && strings.Contains(strings.ToLower(skill1.Fire), strings.ToLower(skilldesc))
				//os.Stdout.WriteString(skill1.Fire + " " + strconv.FormatBool(s1) + "\n")
			}
			if skill2 := card.Skill2(VcData); skill2 != nil {
				s2 = skill2.Fire != "" && strings.Contains(strings.ToLower(skill2.Fire), strings.ToLower(skilldesc))
				//os.Stdout.WriteString(skill2.Fire + " " + strconv.FormatBool(s2) + "\n")
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
	chars := make([]vc.CardCharacter, len(VcData.CardCharacters))
	copy(chars, VcData.CardCharacters)

	sort.Slice(chars, func(i, j int) bool {
		first := chars[i]
		second := chars[j]

		firstCards := vc.CardList(first.Cards(VcData))
		secondCards := vc.CardList(second.Cards(VcData))

		maxFirst := firstCards.Latest()
		maxSecond := secondCards.Latest()

		if maxFirst == nil && maxSecond == nil {
			return first.Id > second.Id
		}
		if maxFirst == nil {
			return false
		}
		if maxSecond == nil {
			return true
		}

		return maxFirst.Id > maxSecond.Id
	})

	for _, character := range chars {
		if !filter(&character) {
			continue
		}

		card := character.FirstEvoCard(VcData)

		cardName := "N/A"

		if card != nil {
			if card.Name == "" {
				cardName = card.Image()
			} else {
				cardName = card.Name
			}
		}

		cardNos := ""
		for _, card := range character.Cards(VcData) {
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
			character.Id,
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

func characterDetailHandler(w http.ResponseWriter, r *http.Request) {
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
	cardId, err := strconv.Atoi(pathParts[2])
	if err != nil || cardId < 1 || cardId > len(VcData.CardCharacters) {
		http.Error(w, "Invalid character id "+pathParts[2], http.StatusNotFound)
		return
	}

	character := vc.CardCharacterScan(cardId, VcData.CardCharacters)

	cards := character.Cards(VcData)

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
		skill1 := card.Skill1(VcData)
		if skill1 == nil {
			skill1 = &vc.Skill{}
		}
		// skill2 := card.Skill2(VcData)
		// skillS1 := card.SpecialSkill1(VcData)
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
			card.Id,
			card.CardNo,
			card.Name,
			card.EvolutionRank,
			card.LastEvolutionRank,
			card.EvolutionCardId,
			card.Rarity(),
			card.Element(),
			card.CardCharaId,
			card.DeckCost,
			card.DefaultOffense,
			card.DefaultDefense,
			card.DefaultFollower,
			card.MaxOffense,
			card.MaxDefense,
			card.MaxFollower,
			card.Skill1Name(VcData),
			card.SkillMin(VcData),
			card.SkillMax(VcData),
			card.SkillProcs(VcData),
			skill1.EffectDefaultValue,
			skill1.DefaultRatio,
			skill1.EffectMaxValue,
			skill1.MaxRatio,
			card.SkillTarget(VcData),
			card.SkillTargetLogic(VcData),
			card.Skill2Name(VcData),
			card.Skill3Name(VcData),
			card.ThorSkill1Name(VcData),
			card.SpecialSkill1Name(VcData),
			card.Description(VcData),
			card.Friendship(VcData),
			card.Login(VcData),
			card.Meet(VcData),
			card.BattleStart(VcData),
			card.BattleEnd(VcData),
			card.FriendshipMax(VcData),
			card.FriendshipEvent(VcData),
		)
	}

	io.WriteString(w, "</tbody></table></div></body></html>")
}
