package couplet

import (
	"database/sql"
	"effect-mobile/internal/song"
	"effect-mobile/pkg/res"
	"fmt"
	"log/slog"
	"sync"
)

type Service struct {
	coupletRepo *Repository
	songRepo    *song.Repository
	logger      *slog.Logger
}

func NewService(coupletRepo *Repository, songRepo *song.Repository, logger *slog.Logger) *Service {
	return &Service{
		coupletRepo: coupletRepo,
		songRepo:    songRepo,
		logger:      logger,
	}
}

func (s *Service) Listing(songId string, pagParams *res.PaginationParams) (*res.PaginationResponse[CoupletModel], res.Error) {

	var couplets []CoupletModel
	var total int
	var err res.Error

	wait := sync.WaitGroup{}
	wait.Add(2)

	// TODO
	listingQuery, totalQuery, args := pagParams.BuildSqlFromParams("couplets", true)

	go func() {
		defer wait.Done()
		couplets, err = s.coupletRepo.FindMany(listingQuery, args...)
	}()

	go func() {
		defer wait.Done()
		total, err = s.coupletRepo.FindCount(totalQuery, args...)
	}()

	wait.Wait()

	if err != nil {
		return nil, err
	}

	return res.NewPagResp(pagParams, couplets, total), nil

}

func (s *Service) Create(trx *sql.Tx, songId string, dto CoupletCreateDto, withContinue bool) res.Error {

	hasSong, _ := s.songRepo.FindOne("SELECT id FROM songs where id = $1 limit 1", songId)

	if hasSong.Id == "" {
		return res.NewErrorWithMessage("Song not found!", 404)
	}

	coupletsCount := 0
	if withContinue {
		countWhere := "SELECT COUNT(*) as total FROM couplets WHERE song_id = $1"
		total, err := s.coupletRepo.FindCount(countWhere, songId)
		if err != nil {
			return err
		}
		coupletsCount = total
	}

	newCouplets := dto.MapToCreate(coupletsCount, songId)

	return s.coupletRepo.Create(trx, newCouplets)
}

func (s *Service) RewriteCouplets(songId string, dto CoupletCreateDto) res.Error {
	trx, rawErr := s.coupletRepo.GetTrx()

	if rawErr != nil {
		s.logger.Error(rawErr.Error())
		return res.NewError(rawErr, 500)
	}

	err := s.Delete(trx, songId, []string{})

	if err != nil {
		s.logger.Info(fmt.Sprintf("Result rollback: %v", trx.Rollback()))
		s.logger.Error(err.Error())
		return err
	}

	err = s.Create(trx, songId, dto, false)

	s.logger.Info(fmt.Sprintf("Result commit: %v", trx.Commit()))

	return err

}

func (s *Service) Delete(trx *sql.Tx, songId string, coupletsIds []string) res.Error {
	return s.coupletRepo.Delete(trx, songId, coupletsIds)
}
