package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Kin-dza-dzaa/wordApi/internal/apierror"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

var (
	ErrAuthorizationFailed = errors.New("unauthorized")
	ErrUnmarshalFailed     = errors.New("unmarshall failed")
)

func (handlers *Handlers) AddWordsHandler() apierror.HttpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		wordsModel := new(models.WordsToAdd)
		if err := json.NewDecoder(r.Body).Decode(wordsModel); err != nil {
			return apierror.NewResponse("error", ErrUnmarshalFailed.Error(), http.StatusBadRequest)
		}
		userId, ok := r.Context().Value(KeyForToken(KeyForToken("user_id"))).(string)
		if !ok {
			return apierror.NewResponse("error", ErrAuthorizationFailed.Error(), http.StatusUnauthorized)

		}
		badWords, err := handlers.service.AddWords(r.Context(), *wordsModel, userId)
		if err != nil {
			return err
		}
		if len(badWords) != 0 {
			sliceOfBadWords := make([]string, 0, len(badWords))
			for i := range badWords {
				sliceOfBadWords = append(sliceOfBadWords, i)
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "some words weren't added", "bad_words": sliceOfBadWords})
			w.WriteHeader(http.StatusOK)
			return nil
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were added"})
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

func (handlers *Handlers) GetWordsHandler() apierror.HttpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			return apierror.NewResponse("error", ErrAuthorizationFailed.Error(), http.StatusUnauthorized)
		}
		
		wordsSlice := make([]models.Word, 0, 100)
		words := models.Words{
			Words: &wordsSlice,
		}
		
		err := handlers.service.GetWords(r.Context(), words, userId)
		if err != nil {
			return err
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "response": words, "message": "words were sent"})
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

func (handlers *Handlers) UpdateWordHandler() apierror.HttpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		wordsModel := new(models.WordToUpdate)
		if err := json.NewDecoder(r.Body).Decode(wordsModel); err != nil {
			return apierror.NewResponse("error", ErrUnmarshalFailed.Error(), http.StatusBadRequest)
		}
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			return apierror.NewResponse("error", ErrAuthorizationFailed.Error(), http.StatusUnauthorized)
		}
		if err := handlers.service.UpdateWord(r.Context(), *wordsModel, userId); err != nil {
			return err
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": wordsModel.OldWord + " was updated to " + wordsModel.NewWord})
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

func (handlers *Handlers) DeleteWordHandler() apierror.HttpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		wordsModel := new(models.WordsToDelete)
		if err := json.NewDecoder(r.Body).Decode(wordsModel); err != nil {
			return apierror.NewResponse("error", ErrUnmarshalFailed.Error(), http.StatusBadRequest)
		}
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			return apierror.NewResponse("error", ErrAuthorizationFailed.Error(), http.StatusUnauthorized)
		}
		if err := handlers.service.DeleteWords(r.Context(), *wordsModel, userId); err != nil {
			return err
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were deleted"})
		w.WriteHeader(http.StatusOK)
		return nil
	}
}

func (handlers *Handlers) UpdateStateHandler() apierror.HttpErrHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		wordsModel := new(models.StatesToUpdate)
		if err := json.NewDecoder(r.Body).Decode(wordsModel); err != nil {
			return apierror.NewResponse("error", ErrUnmarshalFailed.Error(), http.StatusBadRequest)
		}
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			return apierror.NewResponse("error", ErrAuthorizationFailed.Error(), http.StatusUnauthorized)
		}
		if err := handlers.service.UpdateState(r.Context(), *wordsModel, userId); err != nil {
			return err
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were updated"})
		w.WriteHeader(http.StatusOK)
		return nil
	}
}
