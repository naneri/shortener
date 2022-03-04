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

var filePath string

func InitFileRepo(fileName string) *FileRepository {
	filePath = fileName
	repo := FileRepository{
		lastUrlId: 0,
		storage:   make(map[string]string),
	}

	readAllLinks(fileName, &repo)

	return &repo
}

func readAllLinks(fileName string, repo *FileRepository) {
	linkConsumer, err := NewConsumer(fileName)
	if err != nil {
		log.Fatal(err)
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

func (repo *FileRepository) AddLink(link string) int {
	repo.lastUrlId++
	repo.storage[strconv.Itoa(repo.lastUrlId)] = link

	if filePath != "" {
		linkProducer, err := NewProducer(filePath)
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

	return repo.lastUrlId
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

func NewProducer(fileName string) (*producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
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

func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
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

func main() {
	//fileName := "events.log"
	//defer os.Remove(fileName)
	//producer, err := NewProducer(fileName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer producer.Close()
	//consumer, err := NewConsumer(fileName)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer consumer.Close()

	//for _, event := range links {
	//	if err := producer.WriteLink(event); err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	readedEvent, err := consumer.ReadLink()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(readedEvent)
	//}
}
