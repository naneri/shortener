package link

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type FileRepository struct {
	fileStorage *os.File
	lastURLID   int
	storage     map[string]string
}

type Link struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func InitFileRepo(file *os.File) (*FileRepository, error) {
	repo := FileRepository{
		lastURLID:   0,
		storage:     make(map[string]string),
		fileStorage: file,
	}

	if file != nil {
		err := readAllLinks(file, &repo)
		if err != nil {
			return nil, err
		}
	}

	return &repo, nil
}

func readAllLinks(file *os.File, repo *FileRepository) error {
	linkConsumer := NewConsumer(file)

	for {
		readedLink, err := linkConsumer.ReadLink()
		if err != nil {
			if err == io.EOF {
				fmt.Println("finished processing the file")
				break
			} else {
				return err
			}
		}

		repo.storage[strconv.Itoa(readedLink.ID)] = readedLink.URL
	}

	return nil
}

func (repo *FileRepository) AddLink(link string) (int, error) {
	repo.lastURLID++
	repo.storage[strconv.Itoa(repo.lastURLID)] = link

	if repo.fileStorage != nil {
		linkProducer, err := NewProducer(repo.fileStorage)
		if err != nil {
			log.Fatal(err)
		}

		newLink := Link{
			ID:  repo.lastURLID,
			URL: link,
		}

		if err := linkProducer.WriteLink(&newLink); err != nil {
			log.Fatal(err)
		}
	}

	return repo.lastURLID, nil
}

func (repo *FileRepository) GetLink(urlID string) (string, error) {
	if val, ok := repo.storage[urlID]; ok {
		return val, nil
	} else {
		return "", errors.New("record not found")
	}
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(file *os.File) (*producer, error) {
	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteLink(link *Link) error {
	return p.encoder.Encode(&link)
}

func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(file *os.File) *consumer {
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}
}

func (c *consumer) ReadLink() (*Link, error) {
	link := &Link{}
	if err := c.decoder.Decode(&link); err != nil {
		return nil, err
	}
	return link, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}
