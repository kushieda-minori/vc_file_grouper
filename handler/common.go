package handler

import (
	"strings"

	"zetsuboushita.net/vc_file_grouper/vc"
)

func removeDuplicates(a []string) []string {
	result := []string{}
	seen := map[string]string{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

func cleanCustomSkillRecipe(name string) string {
	ret := ""
	lower := strings.ToLower(vc.CleanCustomSkillNoImage(name))
	if strings.Contains(lower, "all enemies") {
		ret += "aoe"
	}
	if strings.Contains(lower, "stop") {
		ret += "+ts "
	} else if ret != "" {
		ret += " "
	}
	if strings.Contains(lower, "fixed") {
		ret += "fixed "
	}
	if strings.Contains(lower, "proportional") {
		ret += "proportional "
	}
	if strings.Contains(lower, "awoken burst") {
		ret += "awoken burst "
	} else if strings.Contains(lower, "recover") {
		ret += "recover "
	} else if strings.Contains(lower, "own atk up") {
		ret += "own atk up "
	}
	if strings.Contains(lower, "passion") {
		ret += "passion "
	}
	if strings.Contains(lower, "cool") {
		ret += "cool "
	}
	if strings.Contains(lower, "light") {
		ret += "light "
	}
	if strings.Contains(lower, "dark") {
		ret += "dark "
	}
	if strings.Contains(lower, "1") {
		ret += "1"
	}
	if strings.Contains(lower, "2") {
		ret += "2"
	}
	if strings.Contains(lower, "3") {
		ret += "3"
	}
	return ret
}

func cleanTicketName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.Replace(ret, "ticket", "", -1)
	ret = strings.Replace(ret, "summon", "", -1)
	ret = strings.Replace(ret, "guaranteed", "", -1)
	ret = strings.Replace(ret, "★★★", "3star", -1)
	ret = strings.Replace(ret, "★★", "2star", -1)
	ret = strings.Replace(ret, "★", "1star", -1)
	ret = strings.TrimSpace(ret)
	return ret
}

func cleanItemName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.Replace(ret, "valkyrie", "", -1)
	ret = strings.Replace(ret, " ", "", -1)
	ret = strings.TrimSpace(ret)
	return ret
}

func cleanArcanaName(name string) string {
	ret := strings.ToLower(name)
	ret = strings.Replace(ret, "arcana's", "", -1)
	ret = strings.Replace(ret, "arcana", "", -1)
	ret = strings.Replace(ret, "%", "", -1)
	ret = strings.Replace(ret, "+", "", -1)
	ret = strings.Replace(ret, " ", "", -1)
	ret = strings.Replace(ret, "forced", "", 1)
	ret = strings.Replace(ret, "strongdef", "def", 1)
	ret = strings.TrimSpace(ret)
	return ret
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isChecked(values []string, e string) string {
	if contains(values, e) {
		return "checked"
	}
	return ""
}
