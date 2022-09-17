package service

import (
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
)

type serviceWord struct {
	repository *repositories.Repository
}

func (s *serviceWord) AddWords(words models.WordsAdd, userId string) []string {
	return s.repository.AddWords(words, userId)
}

func (s *serviceWord) GetWords(userId string) (*models.WordsGet, error) {
	return s.repository.GetWords(userId)
}

func (s *serviceWord) UpdateWord(words models.WordsUpdate, userId string) error {
	return s.repository.UpdateWord(words, userId)
}

func (s *serviceWord) DeleteWords(words models.WordsDelete, userId string) {
	s.repository.DeleteWords(words, userId)
}
