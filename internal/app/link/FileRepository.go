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
	lastUrlId int
	storage   map[string]string
}

type Link struct {
	ID  int    `json:"id"`
	Url string `json:"url"`
}

var fileStorage *os.File

func InitFileRepo(file *os.File) *FileRepository {
	fileStorage = file
	repo := FileRepository{
		lastUrlId: 0,
		storage:   make(map[string]string),
	}

	if file != nil {
		readAllLinks(file, &repo)
	}

	return &repo
}

func readAllLinks(file *os.File, repo *FileRepository) {
	linkConsumer, fileOpenErr := NewConsumer(file)

	if fileOpenErr != nil {
		log.Fatal(fileOpenErr)
	}

	defer linkConsumer.Close()

	for {
		readedLink, err := linkConsumer.ReadLink()
		if err != nil {
			if err == io.EOF {
				fmt.Println("finished processing the file")
				break
			} else {
				log.Fatal(err)
			}
		}

		repo.storage[strconv.Itoa(readedLink.ID)] = readedLink.Url
	}
}

func (repo *FileRepository) AddLink(link string) (int, error) {
	repo.lastUrlId++
	repo.storage[strconv.Itoa(repo.lastUrlId)] = link

	if fileStorage != nil {
		linkProducer, err := NewProducer(fileStorage)
		if err != nil {
			log.Fatal(err)
		}
		defer linkProducer.Close()

		newLink := Link{
			ID:  repo.lastUrlId,
			Url: link,
		}

		if err := linkProducer.WriteLink(&newLink); err != nil {
			log.Fatal(err)
		}
	}

	return repo.lastUrlId, nil
}

func (repo *FileRepository) GetLink(urlId string) (string, error) {
	if val, ok := repo.storage[urlId]; ok {
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

func NewConsumer(file *os.File) (*consumer, error) {
	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
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
