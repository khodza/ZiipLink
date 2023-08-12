package services

import (
	"fmt"
	"zipinit/internal/config"
	"zipinit/internal/storage"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"
)

type RepositoryInterface interface {
	SaveAndGenerateRandomStrings(count int) error
	SaveUrl(url string, alias string) (string, error)
	GetUrl(alias string) (string, error)
	GetAnyUnusedAlias() (string, error)
	MakeAliasUsed(alias string) error
	SaveNewAlias(alias string) error
}

type Service struct {
	db  RepositoryInterface
	log *slog.Logger
	cfg *config.Config
}

func NewService(db RepositoryInterface, log *slog.Logger, cfg *config.Config) *Service {
	return &Service{
		db:  db,
		log: log,
		cfg: cfg,
	}
}

func (s *Service) SaveUrl(url string, providedAlias string) (string, error) {
	const op = "service.SaveUrl"
	// If no alias is provided, get the first unused alias and update it as used
	if providedAlias == "" {
		alias, err := s.db.GetAnyUnusedAlias()
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		//make alias as used
		err = s.db.MakeAliasUsed(alias)
		if err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
		providedAlias = alias
	} else {
		//if provided alias is not empty, save it to the db
		if err := s.db.SaveNewAlias(providedAlias); err != nil {
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}
	// Save the URL with the provided or generated alias
	if _, err := s.db.SaveUrl(url, providedAlias); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", storage.ErrUrlExists
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	shortLink := fmt.Sprintf("%s/%s", s.cfg.Domain, providedAlias)

	return shortLink, nil
}

func (s *Service) GetUrl(alias string) (string, error) {
	const op = "service.GetUrl"
	url, err := s.db.GetUrl(alias)
	if err != nil {
		if err == storage.ErrURLNotFound {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Service) SaveAndGenerateRandomStrings(count int) error {
	const op = "service.SaveAndGenerateRandomStrings"
	if err := s.db.SaveAndGenerateRandomStrings(count); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}
