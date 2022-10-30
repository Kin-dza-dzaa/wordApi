package repositories

import (
	"context"
	"testing"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestIfWordInDb(t *testing.T) {
	testCases := []struct {
		name       string
		word       string
		result     bool
		beforeTest func(MockPool *pgxpoolmock.MockPgxIface, word string, result bool)
	}{
		{
			name:   "word_in_db",
			word:   "test",
			result: true,
			beforeTest: func(MockPool *pgxpoolmock.MockPgxIface, word string, result bool) {
				MockPool.EXPECT().QueryRow(gomock.Any(), queryIfDbHasWord, word).Return(pgxpoolmock.NewRow(result)).Times(1)
			},
		},
		{
			name:   "word_not_in_db",
			word:   "test",
			result: false,
			beforeTest: func(MockPool *pgxpoolmock.MockPgxIface, word string, result bool) {
				MockPool.EXPECT().QueryRow(gomock.Any(), queryIfDbHasWord, word).Return(pgxpoolmock.NewRow(result)).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repo := NewRepositoryWord(MockPool)
			tc.beforeTest(MockPool, tc.word, tc.result)
			var result bool
			err := repo.IfWordInDb(context.TODO(), tc.word, &result)
			assert.Nil(t, err)
			assert.Equal(t, tc.result, result)
		})
	}
}

func TestIfUserHasWord(t *testing.T) {
	testCases := []struct {
		name           string
		word           string
		userId         string
		collectionName string
		result         bool
		beforeTest     func(MockPool *pgxpoolmock.MockPgxIface, userId string, word string, collectionName string, result bool)
	}{
		{
			name:   "user_has_word",
			word:   "test",
			userId: "testId",
			collectionName: "test",
			result: true,
			beforeTest: func(MockPool *pgxpoolmock.MockPgxIface, userId string, word string, collectionName string, result bool) {
				MockPool.EXPECT().QueryRow(gomock.Any(), queryIfUserHasWord, userId, word, collectionName).Return(pgxpoolmock.NewRow(result)).Times(1)
			},
		},
		{
			name:   "user_not_have_word",
			word:   "test",
			userId: "testId",
			collectionName: "test",
			result: false,
			beforeTest: func(MockPool *pgxpoolmock.MockPgxIface, userId string, word string, collectionName string, result bool) {
				MockPool.EXPECT().QueryRow(gomock.Any(), queryIfUserHasWord, userId, word, collectionName).Return(pgxpoolmock.NewRow(result)).Times(1)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			MockPool := pgxpoolmock.NewMockPgxIface(ctrl)
			repo := NewRepositoryWord(MockPool)
			tc.beforeTest(MockPool, tc.userId, tc.word, tc.collectionName, tc.result)
			var result bool
			err := repo.IfUserHasWord(context.TODO(), tc.word, tc.collectionName, &result, tc.userId)
			assert.Nil(t, err)
			assert.Equal(t, tc.result, result)
		})
	}
}

