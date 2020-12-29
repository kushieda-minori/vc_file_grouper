package api

import "net/url"

//Logout Removes authentication tokens
func Logout() {
	client.Get(URL + "api.php?action=logout&token=" + url.QueryEscape(MyCreds.LoginToken))
	MyCreds.CSRFToken = ""
	MyCreds.LoginToken = ""
	client.Jar = getCookieJar()
}
