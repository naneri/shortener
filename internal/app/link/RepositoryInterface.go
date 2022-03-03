package link

type Repository interface {
	AddLink(link string) int
	GetLink(urlId string) (string, error)
}
