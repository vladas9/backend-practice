package services

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	dto "github.com/vladas9/backend-practice/internal/dtos"
	m "github.com/vladas9/backend-practice/internal/models"
	r "github.com/vladas9/backend-practice/internal/repository"
	u "github.com/vladas9/backend-practice/internal/utils"
)

type AuctionService interface {
	NewAuction(dto *dto.AuctionFull) error
	GetAuctionCards(params AuctionCardParams) ([]dto.AuctionCard, error)
	GetFullAuctionById(id uuid.UUID) (*dto.AuctionFull, error)
	GetAuctionTable(params AuctionTableParams) ([]dto.AuctionTable, error)
}

// type auctionService struct{ *Service }
//
//	func NewAuctionService(s *Service) AuctionService {
//		return &auctionService{s}
//	}

type AuctionCardParams struct {
	Category           string
	Condition          string
	Offset, Len        int
	MinPrice, MaxPrice m.Decimal
}

func (a AuctionCardParams) Validate() Problems {
	problems := Problems{}

	if a.Len <= 0 {
		problems["len"] = "must be greater than 0"
	}
	if a.Offset < 0 {
		problems["offset"] = "cannot be negative"
	}
	if !a.MaxPrice.IsZero() && a.MaxPrice.Compare(a.MinPrice) == -1 {
		problems["max_price"] = "max price cannot be less than min price"
	}
	if ok := m.IsCategory(a.Category); !ok {
		problems["category"] = fmt.Sprintf("%s does not exist", a.Category)
	}
	if ok := m.IsCondition(a.Condition); !ok {
		problems["condition"] = fmt.Sprintf("%s does not exist", a.Condition)
	}

	if len(problems) > 0 {
		return problems
	}
	return nil
}

type AuctionTableParams struct {
	UserId        uuid.UUID
	Limit, Offset int
}

func (p AuctionTableParams) Validate() Problems {
	problems := Problems{}

	if p.Limit <= 0 {
		problems["limit"] = "must be greater than 0"
	}
	if p.Offset < 0 {
		problems["offset"] = "cannot be negative"
	}

	if len(problems) > 0 {
		return problems
	}
	return nil
}

func (s *Service) NewAuction(dto *dto.AuctionFull) error { return nil }

//func (s *auctionService) NewAuction(dto *dto.AuctionFull) error {
//	err := s.store.WithTx(func(stx *r.StoreTx) error {
//		if itemId, err := stx.ItemRepo().Insert(item); err != nil {
//			return err
//		} else {
//			auct.ItemId = itemId
//		}
//		if err := stx.AuctionRepo().Insert(auct); err != nil {
//			return err
//		}
//		return nil
//	})
//	return fail(err)
//}

func (s *Service) GetAuctionTable(params AuctionTableParams) ([]dto.AuctionTable, error) {
	var auctions []*m.AuctionDetails
	var err error
	err = s.store.WithTx(func(stx *r.StoreTx) error {
		auctions, err = getUserAuctions(stx, params, withItem, withMaxBidder)
		return err
	})
	if err != nil {
		return nil, fail(err)
	}
	var table []dto.AuctionTable
	for _, a := range auctions {
		table = append(table, *dto.MapAuctionTable(a))
	}
	return table, nil
}

func getUserAuctions(stx *r.StoreTx, params AuctionTableParams, opts ...auctOpt) (auctions []*m.AuctionDetails, err error) {
	var auctModels []*m.AuctionModel
	auctModels, err = stx.AuctionRepo().GetAllByUserId(params.UserId, params.Limit, params.Offset)
	var auct *m.AuctionDetails
	for _, auctModel := range auctModels {
		auct, err = getAuctionDetails(auctModel, stx, opts...)
		if err != nil {
			return nil, fail(err)
		}
		auctions = append(auctions, auct)
	}
	return auctions, nil
}

func (s *Service) GetAuctionCards(params AuctionCardParams) ([]dto.AuctionCard, error) {
	var err error
	var auctDetails []*m.AuctionDetails
	err = s.store.WithTx(func(stx *r.StoreTx) error {
		if auctDetails, err = getAuctions(stx, params, withItem); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, fail(err)
	}

	var cards []dto.AuctionCard
	for _, a := range auctDetails {
		cards = append(cards, *dto.MapAuctionCard(a))
	}
	return cards, nil
}

func getAuctions(stx *r.StoreTx, params AuctionCardParams, opts ...auctOpt) (auctions []*m.AuctionDetails, err error) {
	auctionRepo := stx.AuctionRepo()
	var auctModels []*m.AuctionModel
	if params.MaxPrice.IsZero() {
		auctModels, err = auctionRepo.GetAll(params.Offset, params.Len)
	} else {
		auctModels, err = auctionRepo.GetAllFiltered(
			params.Offset, params.Len,
			params.MinPrice, params.MaxPrice)
		u.Logger.Info("getAuctionsWith in:", auctions)
	}
	if err != nil {
		return nil, fail(err)
	}
	for _, auctModel := range auctModels {
		auct, err := getAuctionDetails(auctModel, stx, opts...)
		if err != nil {
			return nil, fail(err)
		} else if auct.ItemHas(params.Condition, params.Category) {
			auctions = append(auctions, auct)
		}
	}
	return auctions, nil
}

func (s *Service) GetFullAuctionById(id uuid.UUID) (*dto.AuctionFull, error) {
	var err error
	var auct *m.AuctionDetails
	err = s.store.WithTx(func(stx *r.StoreTx) error {
		auctModel, err := stx.AuctionRepo().GetById(id)
		if err != nil {
			return err
		}
		auct, err = getAuctionDetails(auctModel, stx, withItem, withBids)
		return err
	})
	if err != nil {
		return nil, fail(err)
	}
	return dto.MapAuctionRespToFull(auct), nil
}

type auctOpt func(stx *r.StoreTx, auct *m.AuctionDetails) (*m.AuctionDetails, error)

func getAuctionDetails(auct *m.AuctionModel, stx *r.StoreTx, opts ...auctOpt) (*m.AuctionDetails, error) {
	var err error
	details := m.NewAuctionDetails(auct)
	for _, opt := range opts {
		details, err = opt(stx, details)
		if err != nil {
			return nil, err
		}
	}
	return details, nil
}

func withBids(stx *r.StoreTx, auct *m.AuctionDetails) (*m.AuctionDetails, error) {
	var err error
	auct.BidList, err = stx.BidRepo().GetAllFor(auct.Auction)
	if err != nil {
		return nil, fail(err)
	}
	return auct, nil
}

func withItem(stx *r.StoreTx, auct *m.AuctionDetails) (*m.AuctionDetails, error) {
	var err error
	auct.Item, err = stx.ItemRepo().GetById(auct.Auction.ItemId)
	if err != nil {
		return nil, fail(err)
	}
	return auct, nil
}

func withMaxBidder(stx *r.StoreTx, auct *m.AuctionDetails) (*m.AuctionDetails, error) {
	var err error
	auct.MaxBidder, err = stx.UserRepo().GetById(auct.Auction.MaxBidderId)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fail(err)
		}
		auct.MaxBidder = &m.UserModel{}
	}
	return auct, nil
}
