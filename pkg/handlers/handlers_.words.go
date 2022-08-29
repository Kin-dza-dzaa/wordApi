package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

func (h *Handlers) AddWordsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var wordsModel models.Words
		if err := json.NewDecoder(r.Body).Decode(&wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected array with words"})
			return 
		}
		userId, ok := r.Context().Value(key).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "authorization error"})
			return
		}
		badWords := h.service.AddWords(wordsModel, userId)
		w.WriteHeader(http.StatusOK)
		if len(badWords) != 0 {
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": fmt.Sprintf("some words weren't added: %v", badWords)})
			return 
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were added"})
	})
} 

func (h *Handlers) GetWordsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		userId, ok := r.Context().Value(key).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "authorization error"})
			return
		}
		words, err := h.service.GetWords(userId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "words":words})
		w.WriteHeader(http.StatusOK)
	})
}

func (h *Handlers) UpdateWordHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var wordsModel models.Words
		if err := json.NewDecoder(r.Body).Decode(&wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected array of 2 words: [old_word, new_word]"})
			return 
		}
		userId, ok := r.Context().Value(key).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "authorization error"})
			return
		}
		if err := h.service.UpdateWord(wordsModel, userId); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": wordsModel.Words[0] + " was updated to " + wordsModel.Words[1]})
		w.WriteHeader(http.StatusOK)
	})
} 

func (h *Handlers) DeleteWordHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var wordsModel models.Words
		if err := json.NewDecoder(r.Body).Decode(&wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected array of 1 word: [delete_word]"})
			return 
		}
		userId, ok := r.Context().Value(key).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "authorization error"})
			return
		}
		h.service.DeleteWords(wordsModel, userId)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were delted"})
		w.WriteHeader(http.StatusOK)
	})
} 
