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

func (repo *MemoryRepository) AddLink(link string) (int, error) {
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
