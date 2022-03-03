package link

import (
	"errors"
	"strconv"
)

type MemoryRepository struct {
	lastUrlId int
	storage   map[string]string
}

func Init() *MemoryRepository {
	repo := MemoryRepository{
		lastUrlId: 0,
		storage:   make(map[string]string),
	}

	return &repo
}

func (repo *MemoryRepository) AddLink(link string) int {
	repo.lastUrlId++
	repo.storage[strconv.Itoa(repo.lastUrlId)] = link

	return repo.lastUrlId
}

func (repo *MemoryRepository) GetLink(urlId string) (string, error) {
	if val, ok := repo.storage[urlId]; ok {
		return val, nil
	} else {
		return "", errors.New("record not found")
	}
}
