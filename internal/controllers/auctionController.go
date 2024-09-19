package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"github.com/vladas9/backend-practice/internal/dtos"
	s "github.com/vladas9/backend-practice/internal/services"
	"github.com/vladas9/backend-practice/internal/utils"
)

func (c *Controller) GetAuctions(w http.ResponseWriter, r *http.Request) error {
	fail := func(err error) error {
		return fmt.Errorf("GetAuctions controller: %w", err)
	}
	if err := r.ParseForm(); err != nil {
		return fail(err)
	}

	offsetStr := r.FormValue("offset")
	leangthStr := r.FormValue("limit")
	minPriceStr := r.FormValue("min_price")
	maxPriceStr := r.FormValue("max_price")
	categoryStr := r.FormValue("category")
	conditionStr := r.FormValue("lotcondition")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return fail(err)
	}

	leangth, err := strconv.Atoi(leangthStr)
	if err != nil {
		return fail(err)
	}

	maxPrice, err := decimal.NewFromString(maxPriceStr)
	if err != nil {
		return fail(err)
	}
	minPrice, err := decimal.NewFromString(minPriceStr)
	if err != nil {
		return fail(err)
	}

	params := s.AuctionParams{
		Offset:    offset,
		Len:       leangth,
		MaxPrice:  maxPrice,
		MinPrice:  minPrice,
		Category:  categoryStr,
		Condition: conditionStr,
	}
	problems := params.Validate()
	if problems != nil {
		return &ApiError{Status: 400, ErrorMsg: problems}
	}

	resp, err := c.service.GetAuctions(params)
	if err != nil {
		return fmt.Errorf("AuctionController: %w", err)
	}
	var cards []dtos.AuctionCard
	utils.Logger.Info("AuctionController: getAuctions:", resp)
	for i, respItem := range resp {
		cards = append(cards, dtos.MapAuctionRespToCard(i+1, respItem))
	}
	return WriteJSON(w, http.StatusOK, cards)
}

func (c *Controller) GetAuction(w http.ResponseWriter, r *Response) error {

	//auctions, err := c.service.GetAuctionById() // use internal/dtos/auctionFull.go and the auction service
	return nil
}
