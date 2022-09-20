package main

import (
	"context"
	"fmt"
	"github.com/naneri/shortener/cmd/grpc/config"
	"github.com/naneri/shortener/cmd/grpc/proto"
	"github.com/naneri/shortener/internal/app/link"
	"log"
	"strconv"
)

type ShortenerServer struct {
	proto.UnimplementedShortenerServiceServer
	LinkRepository link.Repository // <--
	Config         config.Config
}

func (s *ShortenerServer) AddUrl(ctx context.Context, in *proto.AddLinkRequest) (*proto.AddLinkResponse, error) {
	var response proto.AddLinkResponse
	log.Println("request received")

	_, err := s.LinkRepository.AddLink(in.Link, 0)
	if err != nil {
		response.Error = fmt.Sprintf("error storing the link: " + err.Error())
	}

	return &response, nil
}

func (s *ShortenerServer) GetUrl(ctx context.Context, in *proto.GetLinkRequest) (*proto.GetLinkResponse, error) {
	var response proto.GetLinkResponse

	fullLink, err := s.LinkRepository.GetLink(in.UrlId)
	if err != nil {
		response.Error = fmt.Sprintf("error storing the link: " + err.Error())
	} else {
		response.FullUrl = fullLink
	}

	return &response, nil
}

func (s *ShortenerServer) ListUserUrls(ctx context.Context, in *proto.GetUserUrlsRequest) (*proto.GetUserUrlsResponse, error) {
	var response proto.GetUserUrlsResponse

	u64, err := strconv.ParseUint(in.UserId, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	parsedUserID := uint32(u64)

	links, dbErr := s.LinkRepository.GetAllLinks()

	if dbErr != nil {
		log.Println("Error getting data from the storage: " + dbErr.Error())
	}

	for _, userLink := range links {
		if userLink.UserID == parsedUserID {
			protoUserLink := proto.StoredLink{
				Id:      strconv.Itoa(userLink.ID),
				FullUrl: userLink.URL,
			}
			response.Links = append(response.Links, &protoUserLink)
		}
	}

	return &response, nil
}

func (s ShortenerServer) DeleteUserUrls(ctx context.Context, in *proto.DeleteUserUrlsRequest) (*proto.DeleteUserUrlsResponse, error) {
	var resp proto.DeleteUserUrlsResponse

	err := s.LinkRepository.DeleteLinks(in.UrlIds)
	if err != nil {
		resp.Error = fmt.Sprintln("error deleting URLS: ", err.Error())
	}

	return &resp, nil
}

func (s ShortenerServer) ShortenBatch(ctx context.Context, in *proto.ShortenBatchRequest) (*proto.ShortenBatchResponse, error) {
	var resp proto.ShortenBatchResponse

	for _, batchLink := range in.BatchLinks {
		lastURLID, addErr := s.LinkRepository.AddLink(batchLink.OriginalUrl, 0)

		if addErr != nil {
			resp.Error = fmt.Sprintln("error storing a link: ", addErr.Error())
			return &resp, nil
		}

		resp.StoredBatchLinks = append(resp.StoredBatchLinks, &proto.StoredBatchLink{
			CorrelationId: batchLink.CorrelationId,
			ShortUrl:      generateShortLink(lastURLID, s.Config.BaseURL),
		})
	}

	return &resp, nil
}

func generateShortLink(lastURLID int, baseURL string) string {
	return fmt.Sprintf("%s/%d", baseURL, lastURLID)
}
