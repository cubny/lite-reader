package api_test

const (
	tcOK            = "ok"
	tcServiceErr    = "service returns error"
	tcInvalidFeedID = "invalid feed id"

	feedsPath       = "/feeds"
	feed1Path       = "/feeds/1"
	feedFetchPath   = "/feeds/1/fetch"
	item1Path       = "/items/1"
	signupPath      = "/signup"
	setupPath       = "/setup"
	setupStatusPath = "/setup/status"
	settingsPath    = "/settings"

	feedTitle = "title"
	feedDesc  = "description"
	feedLink  = "link"
	feedDir   = "dir"

	respInvalidFeedID    = `{"error":{"code":422,"details":"Invalid params - invalid feed id"}}`
	respInternalUnread   = `{"error":{"code":500, "details":"Internal error - cannot get unread items"}}`
	respInvalidReqBody   = `{"error":{"code":400,"details":"Bad Request - invalid request body"}}`
	respAllowSignupTrue  = `{"allow_signup":true}`
	respAllowSignupFalse = `{"allow_signup":false}`
)
