package main

import (
	"context"
	"fmt"
	"github.com/naneri/shortener/cmd/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	// устанавливаем соединение с сервером
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	// получаем переменную интерфейсного типа UsersClient,
	// через которую будем отправлять сообщения
	c := proto.NewShortenerServiceClient(conn)

	// функция, в которой будем отправлять сообщения
	TestShortener(c)
}

func TestShortener(conn proto.ShortenerServiceClient) {
	resp, err := conn.AddURL(context.Background(), &proto.AddLinkRequest{Link: "https://yandex.ru"})

	if err != nil {
		log.Fatal(err)
	}
	if resp.Error != "" {
		fmt.Println(resp.Error)
	}

}
