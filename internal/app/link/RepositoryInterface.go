package link

type Repository interface {
	AddLink(link string, userID uint32) (int, error)
	GetLink(urlID string) (string, error)
	GetAllLinks(userID uint32) map[string]*Link
}
