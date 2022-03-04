package link

import (
	"encoding/json"
	"fmt"
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

func InitFileRepo(fileName string) *FileRepository {
	repo := FileRepository{
		lastUrlId: 0,
		storage:   make(map[string]string),
	}

	defer os.Remove(fileName)
	producer, err := NewProducer(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	consumer, err := NewConsumer(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()

	return &repo
}

func readAllLinks(consumer consumer) {
	for {
		readedEvent, err := consumer.ReadEvent()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(readedEvent)
	}
}

func (repo FileRepository) AddLink(link string) int {
	repo.lastUrlId++
	repo.storage[strconv.Itoa(repo.lastUrlId)] = link

	return repo.lastUrlId
	panic("implement me")
}

func (repo FileRepository) GetLink(urlId string) (string, error) {
	panic("implement me")
}

///*type Event struct {
//	ID       uint    `json:"id"`
//	CarModel string  `json:"car_model"`
//	Price    float64 `json:"price"`
//*/}

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

func (c *consumer) ReadEvent() (*Link, error) {
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
	//	readedEvent, err := consumer.ReadEvent()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(readedEvent)
	//}
}
