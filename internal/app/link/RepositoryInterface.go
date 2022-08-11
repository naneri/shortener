package link

type Repository interface {
	// AddLink adds a link to the storage
	AddLink(link string, userID uint32) (int, error)
	// GetLink gets a link from the storage
	GetLink(urlID string) (string, error)
	// GetAllLinks gets all links from the storage
	GetAllLinks() (map[string]*Link, error)
	// DeleteLinks deletes all links with given ids from the storage
	DeleteLinks(ids []string) error
}
