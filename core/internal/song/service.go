package song

import (
	"effect-mobile/pkg/res"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

const timeParseLayout = "2006-01-02"

type Service struct {
	songRepo *Repository
	logger   *slog.Logger
}

func NewService(songRepo *Repository, logger *slog.Logger) *Service {
	return &Service{
		songRepo: songRepo,
		logger:   logger,
	}
}

func (s *Service) Create(dto *SongCreateDto) (string, res.Error) {

	releaseDate, rawErr := time.Parse(timeParseLayout, dto.ReleaseDate)
	if rawErr != nil {
		s.logger.Error(rawErr.Error())
		return "", res.NewError(rawErr, 422)
	}

	song := NewSong(dto.Song, dto.Group, dto.Link, releaseDate)
	return s.songRepo.Create(&song)
}

func (s *Service) Delete(id string) res.Error {

	_, err := s.songRepo.CheckExistsBySongOrId("", id)

	if err != nil {
		return err
	}

	return s.songRepo.Delete(id)
}

func (s *Service) Update(id string, dto *SongUpdateDto) res.Error {

	_, err := s.songRepo.CheckExistsBySongOrId("", id)

	if err != nil {
		return err
	}
	releaseDate, rawErr := time.Parse(timeParseLayout, dto.ReleaseDate)
	if rawErr != nil {
		s.logger.Error(rawErr.Error())
		return res.NewError(rawErr, 422)
	}

	dto.Id = id
	dto.ReleaseDate = releaseDate.Format(time.RFC3339)

	return s.songRepo.Update(dto)
}

func (s *Service) Info(group, song string) (*SongModel, res.Error) {
	query := s.getOneQuery("AND \"group\" = $1 AND song = $2")
	return s.songRepo.FindOne(query, group, song)
}

func (s *Service) FindOne(id string) (*SongModel, res.Error) {
	query := s.getOneQuery("AND s.id = $1")
	return s.songRepo.FindOne(query, id)
}

func (s *Service) Listing(pagParams *res.PaginationParams) (*res.PaginationResponse[SongShowModel], res.Error) {

	var songs []SongModel
	var total int
	var err res.Error

	wait := sync.WaitGroup{}
	wait.Add(2)

	listingQuery, totalQuery, args := pagParams.BuildSqlFromParams("songs", false)

	go func() {
		defer wait.Done()
		songs, err = s.songRepo.FindMany(listingQuery, args...)
	}()

	go func() {
		defer wait.Done()
		total, err = s.songRepo.FindCount(totalQuery, args...)
	}()

	wait.Wait()

	if err != nil {
		return nil, err
	}

	songsShow := MapModelsToShow(songs)

	return res.NewPagResp(pagParams, songsShow, total), nil

}

func (s *Service) getOneQuery(addWhere string) string {

	return fmt.Sprintf(`
		SELECT s.*, string_agg(c.couplet, '\n ' ORDER BY c.couplet_num ASC) as text FROM songs s
		LEFT JOIN couplets c ON c.song_id = s.id
		WHERE deleted_at IS NULL %s 
		GROUP BY s.id 
		limit 1
	`, addWhere)
}
