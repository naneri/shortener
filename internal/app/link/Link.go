package link

type Link struct {
	ID     int    `json:"id"`
	UserId uint32 `json:"userId"`
	URL    string `json:"url"`
}
