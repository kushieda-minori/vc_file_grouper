package main

import (
	"encoding/csv"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func skillTableHandler(w http.ResponseWriter, r *http.Request) {

}

func skillCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=vcData-skills-"+strconv.Itoa(VcData.Version)+"_"+VcData.Common.UnixTime.Format(time.RFC3339)+".csv")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.Write([]string{"Id",
		"Name",
		"Description",
		"Fire",
		"SkillMin",
		"SkillMax",
		"FireMin",
		"LevelType",
		"Type",
		"TimingId",
		"MaxCount",
		"CondSceneId",
		"CondSideId",
		"CondId",
		"KingSeriesId",
		"KingId",
		"CondParam",
		"DefaultRatio",
		"MaxRatio",
		"PublicStartDatetime",
		"PublicEndDatetime",
		"EffectId",
		"EffectParam",
		"EffectParam2",
		"EffectParam3",
		"EffectParam4",
		"EffectParam5",
		"EffectDefaultValue",
		"EffectMaxValue",
		"TargetScopeId",
		"TargetScope",
		"TargetLogicId",
		"TargetLogic",
		"TargetParam",
		"AnimationId",
		"ThorHammerAnimationType",
	})
	for _, s := range VcData.Skills {
		var startDate, endDate string
		if s.PublicStartDatetime.IsZero() {
			startDate = "-1"
			endDate = "-1"
		} else {
			startDate = s.PublicStartDatetime.Format(time.RFC3339)
			endDate = s.PublicEndDatetime.Format(time.RFC3339)
		}
		err := cw.Write([]string{strconv.Itoa(s.Id),
			s.Name,
			s.Description,
			s.Fire,
			s.SkillMin(),
			s.SkillMax(),
			s.FireMin(),
			strconv.Itoa(s.LevelType),
			strconv.Itoa(s.Type),
			strconv.Itoa(s.TimingId),
			strconv.Itoa(s.MaxCount),
			strconv.Itoa(s.CondSceneId),
			strconv.Itoa(s.CondSideId),
			strconv.Itoa(s.CondId),
			strconv.Itoa(s.KingSeriesId),
			strconv.Itoa(s.KingId),
			strconv.Itoa(s.CondParam),
			strconv.Itoa(s.DefaultRatio),
			strconv.Itoa(s.MaxRatio),
			startDate,
			endDate,
			strconv.Itoa(s.EffectId),
			strconv.Itoa(s.EffectParam),
			strconv.Itoa(s.EffectParam2),
			strconv.Itoa(s.EffectParam3),
			strconv.Itoa(s.EffectParam4),
			strconv.Itoa(s.EffectParam5),
			strconv.Itoa(s.EffectDefaultValue),
			strconv.Itoa(s.EffectMaxValue),
			strconv.Itoa(s.TargetScopeId),
			s.TargetScope(),
			strconv.Itoa(s.TargetLogicId),
			s.TargetLogic(),
			strconv.Itoa(s.TargetParam),
			strconv.Itoa(s.AnimationId),
			strings.Replace(string(s.ThorHammerAnimationType[:]), "\"", "", -1),
		})
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
		}
	}
	cw.Flush()
}
