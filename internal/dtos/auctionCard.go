package dtos

import (
	m "github.com/vladas9/backend-practice/internal/models"

	"fmt"
	"os"
	"time"
)

type AuctionCard struct {
	Id        int        `json:"id"`
	ImgSrc    string     `json:"img_src"`
	Title     string     `json:"title"`
	NumOfBids int        `json:"num_of_bids"`
	MaxBid    m.Decimal  `json:"max_bid"`
	EndDate   time.Time  `json:"end_date"`
	Category  m.Category `json:"category_name"`
}

func MapAuctionCard(
	id int,
	auction m.AuctionModel,
	item m.ItemModel) AuctionCard {

	card := AuctionCard{
		Id:        id,
		ImgSrc:    fmt.Sprintf("http://%s:%s/api/img/%s", os.Getenv("HOST"), os.Getenv("PORT"), item.Images[0]),
		Title:     item.Name,
		NumOfBids: int(auction.BidCount),
		MaxBid:    auction.CurrentBid,
		EndDate:   auction.EndTime,
		Category:  item.Category,
	}

	return card
}