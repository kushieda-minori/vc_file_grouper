package wiki

//Credentials Credentials to use for bot-login. your main login credentials will not work.
type Credentials struct {
	Username   string // must be set before login
	Password   string // must be set before login. Use the "bot key" see https://valkyriecrusade.fandom.com/wiki/Special:BotPasswords
	CSRFToken  string // called after login: call GET api.php?action=query&meta=tokens
	LoginToken string // called before login: call GET api.php?action=query&meta=tokens&format=json&type=login
}

//MyCreds credentials that are used for Bot commands sent to the wiki
var MyCreds Credentials = Credentials{}
