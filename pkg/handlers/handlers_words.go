package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

func (handlers *Handlers) AddWordsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wordsModel models.WordsToAdd
		if err := handlers.customUnmarshall(r, &wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": `expected ojbect: {"words": [{"word": string, "state": int, "collection_name": string}, ...]}`})
			return
		}
		userId, ok := r.Context().Value(KeyForToken(KeyForToken("user_id"))).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "invalid token"})
			return
		}
		badWords := handlers.service.AddWords(&wordsModel, userId)
		w.WriteHeader(http.StatusOK)
		if len(badWords) != 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message":"some words weren't added", "bad_words": badWords})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were added"})
	})
}

func (handlers *Handlers) GetWordsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "invalid token"})
			return
		}
		words, err := handlers.service.GetWords(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "response": words, "message": "words were sent"})
		w.WriteHeader(http.StatusOK)
	})
}

func (handlers *Handlers) UpdateWordHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wordsModel models.WordToUpdate
		if err := handlers.customUnmarshall(r, &wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": `expected object: {"old_word": string, "new_word": string, "collection_name": string}`})
			return
		}
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "invalid token"})
			return
		}
		if err := handlers.service.UpdateWord(&wordsModel, userId); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": wordsModel.OldWord + " was updated to " + wordsModel.NewWord})
		w.WriteHeader(http.StatusOK)
	})
}

func (handlers *Handlers) DeleteWordHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wordsModel models.WordsToDelete
		if err := handlers.customUnmarshall(r, &wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": `expected object: {"words": [{"collection_name": string, "word": string}...]}`})
			return
		}
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "invalid token"})
			return
		}
		handlers.service.DeleteWords(&wordsModel, userId)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were delted"})
		w.WriteHeader(http.StatusOK)
	})
}

func (handlers *Handlers) UpdateStateHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wordsModel models.StatesToUpdate
		if err := handlers.customUnmarshall(r, &wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": `expected object: {"collection_name": string, "words": [{"word": string, "new_state": number}]}`})
			return
		}
		userId, ok := r.Context().Value(KeyForToken("user_id")).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "invalid token"})
			return
		}
		handlers.service.UpdateState(&wordsModel, userId)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were updated"})
		w.WriteHeader(http.StatusOK)
	})
}

func (hanlders *Handlers) customUnmarshall(r *http.Request, target interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return err
	}
	switch v := target.(type) {
		case *models.WordsToAdd:
			if v.Words == nil {
				return errors.New("")
			}
		case *models.WordsToDelete:
			if v.Words == nil {
				return errors.New("")
			}
		case *models.StatesToUpdate:
			if v.Words == nil {
				return errors.New("")
			}		
		case *models.WordToUpdate:
			break
		default:
			return errors.New("")
	}
	return nil
}