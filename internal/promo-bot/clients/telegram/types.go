package telegram

type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	ID      int              `json:"update_id"`
	Message *IncomingMessage `json:"message"`
}

type IncomingMessage struct {
	Text string `json:"text"`
	User User   `json:"from"`
	Chat Chat   `json:"chat"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}
