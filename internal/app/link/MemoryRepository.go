package link

import (
	"errors"
	"strconv"
)

type MemoryRepository struct {
	lastURLID int
	storage   map[string]string
}

func InitMemoryRepo() *MemoryRepository {
	repo := MemoryRepository{
		lastURLID: 0,
		storage:   make(map[string]string),
	}

	return &repo
}

func (repo *MemoryRepository) AddLink(link string, userID uint32) (int, error) {
	repo.lastURLID++
	repo.storage[strconv.Itoa(repo.lastURLID)] = link

	return repo.lastURLID, nil
}

func (repo *MemoryRepository) GetLink(urlID string) (string, error) {
	if val, ok := repo.storage[urlID]; ok {
		return val, nil
	} else {
		return "", errors.New("record not found")
	}
}

func (repo *MemoryRepository) DeleteLinks(ids []string) error {

	return nil
}

// GetAllLinks - this is just a fake method to comply with Interface
func (repo *MemoryRepository) GetAllLinks() (map[string]*Link, error) {
	links := make(map[string]*Link)

	return links, nil
}
