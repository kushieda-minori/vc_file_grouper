package api

import (
	"errors"
	"vc_file_grouper/wiki"
)

//EditCardPage Edits a card page
func EditCardPage(cp *wiki.CardPage) (err error) {
	// verify basic page information is available
	if cp == nil {
		err = errors.New("Page is nil")
		return
	}
	if cp.PageName == "" {
		err = errors.New("Page name can not be blank")
		return
	}

	// verify we are logged into the wiki API
	if MyCreds.LoginToken == "" {
		err = Login()
		if err != nil {
			return
		}
	}

	return errors.New("not implemeted")
}
