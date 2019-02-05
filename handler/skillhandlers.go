package handler

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"zetsuboushita.net/vc_file_grouper/vc"
)

// SkillTableHandler does nothing
func SkillTableHandler(w http.ResponseWriter, r *http.Request) {

}

// SkillCsvHandler outputs skills as a csv file
func SkillCsvHandler(w http.ResponseWriter, r *http.Request) {
	// File header
	w.Header().Set("Content-Disposition", "attachment; filename=vcData-skills-"+strconv.Itoa(vc.Data.Version)+"_"+vc.Data.Common.UnixTime.Format(time.RFC3339)+".csv")
	w.Header().Set("Content-Type", "text/csv")
	cw := csv.NewWriter(w)
	cw.UseCRLF = true
	cw.Write([]string{"ID",
		"Name",
		"Description",
		"Fire",
		"SkillMin",
		"SkillMax",
		"FireMin",
		"LevelType",
		"Type",
		"TimingID",
		"MaxCount",
		"CondSceneID",
		"CondSideID",
		"CondID",
		"KingSeriesID",
		"KingID",
		"CondParam",
		"DefaultRatio",
		"MaxRatio",
		"PublicStartDatetime",
		"PublicEndDatetime",
		"EffectID",
		"EffectParam",
		"EffectParam2",
		"EffectParam3",
		"EffectParam4",
		"EffectParam5",
		"EffectDefaultValue",
		"EffectMaxValue",
		"TargetScopeID",
		"TargetScope",
		"TargetLogicID",
		"TargetLogic",
		"TargetParam",
		"AnimationID",
		"ThorHammerAnimationType",
	})
	for _, s := range vc.Data.Skills {
		var startDate, endDate string
		if s.PublicStartDatetime.IsZero() {
			startDate = "-1"
			endDate = "-1"
		} else {
			startDate = s.PublicStartDatetime.Format(time.RFC3339)
			endDate = s.PublicEndDatetime.Format(time.RFC3339)
		}
		err := cw.Write([]string{strconv.Itoa(s.ID),
			s.Name,
			s.Description,
			s.Fire,
			s.SkillMin(),
			s.SkillMax(),
			s.FireMin(),
			strconv.Itoa(s.LevelType),
			strconv.Itoa(s.Type),
			strconv.Itoa(s.TimingID),
			strconv.Itoa(s.MaxCount),
			strconv.Itoa(s.CondSceneID),
			strconv.Itoa(s.CondSideID),
			strconv.Itoa(s.CondID),
			strconv.Itoa(s.KingSeriesID),
			strconv.Itoa(s.KingID),
			strconv.Itoa(s.CondParam),
			strconv.Itoa(s.DefaultRatio),
			strconv.Itoa(s.MaxRatio),
			startDate,
			endDate,
			strconv.Itoa(s.EffectID),
			strconv.Itoa(s.EffectParam),
			strconv.Itoa(s.EffectParam2),
			strconv.Itoa(s.EffectParam3),
			strconv.Itoa(s.EffectParam4),
			strconv.Itoa(s.EffectParam5),
			strconv.Itoa(s.EffectDefaultValue),
			strconv.Itoa(s.EffectMaxValue),
			strconv.Itoa(s.TargetScopeID),
			s.TargetScope(),
			strconv.Itoa(s.TargetLogicID),
			s.TargetLogic(),
			strconv.Itoa(s.TargetParam),
			strconv.Itoa(s.AnimationID),
			strings.Replace(string(s.ThorHammerAnimationType[:]), "\"", "", -1),
		})
		if err != nil {
			log.Printf(err.Error() + "\n")
		}
	}
	cw.Flush()
}
