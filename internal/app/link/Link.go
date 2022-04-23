package link

type Link struct {
	ID     int    `json:"id"`
	UserID uint32 `json:"userId"`
	URL    string `json:"url"`
}

//
//type DBLink struct {
//	ID 		uint
//	UserID 	uint32
//	URL 	string
//}
