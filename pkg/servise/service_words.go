package service

import (
	external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
	"github.com/Kin-dza-dzaa/wordApi/pkg/repositories"
)

type serviceWord struct {
	repository *repositories.Repository
}

func (s *serviceWord) AddWords(words models.Words, userId string) []string {
	return s.repository.AddWords(words, userId)
}

func (s *serviceWord) GetWords(userId string) (*[]external.Translation, error) {
	return s.repository.GetWords(userId)
}

func (s *serviceWord) UpdateWord(words models.Words, userId string) error {
	return s.repository.UpdateWord(words, userId)
}

func (s *serviceWord) DeleteWords(words models.Words, userId string) {
	s.repository.DeleteWords(words, userId)
}
