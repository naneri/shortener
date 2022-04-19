package link

type Repository interface {
	AddLink(link string) (int, error)
	GetLink(urlID string) (string, error)
}
