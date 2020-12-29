package api

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

//URL root of the wiki
var URL string = "https://valkyriecrusade.fandom.com"

var client *http.Client = &http.Client{
	Jar: getCookieJar(),
}

func getCookieJar() *cookiejar.Jar {
	j, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: nil,
	})
	return j
}

//CardNameToWiki Converts a card name to a wiki name. The result is safe to use in URL paths
func CardNameToWiki(name string) string {
	return url.QueryEscape(strings.ReplaceAll(name, " ", "_"))
}
