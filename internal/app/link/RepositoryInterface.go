package link

type Repository interface {
	AddLink(link string, userId uint32) (int, error)
	GetLink(urlID string) (string, error)
}
