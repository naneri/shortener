package link

type Repository interface {
	AddLink(link string) (int, error)
	GetLink(urlId string) (string, error)
}
