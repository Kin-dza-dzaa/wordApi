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
		var wordsModel models.WordsAdd
		if err := json.NewDecoder(r.Body).Decode(&wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected ojbect like this: {words: [{word: string, state: int, collection_name: string}, ...]}"})
			return 
		}
		// for _, v := range wordsModel.Words {
		// 	if len(v) != 3 {
		// 		w.WriteHeader(http.StatusBadRequest)
		// 		json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected 2 dimensional array: {words: [[word, state, collection_name], ...]}"})
		// 		return
		// 	}
		// }
		userId, ok := r.Context().Value(KEY).(string)
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
		userId, ok := r.Context().Value(KEY).(string)
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
		var wordsModel models.WordToUpdate
		if err := json.NewDecoder(r.Body).Decode(&wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected array of 4 words: [old_word, new_word, new_state: string, new_collection: string]"})
			return 
		}
		// if len(wordsModel.Words) != 4 {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected array of 4 words: [old_word, new_word, new_state: string, new_collection: string]"})
		// 	return 
		// }
		userId, ok := r.Context().Value(KEY).(string)
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
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": wordsModel.OldWord + " was updated to " + wordsModel.NewWord})
		w.WriteHeader(http.StatusOK)
	})
} 

func (h *Handlers) DeleteWordHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		var wordsModel models.WordsDelete
		if err := json.NewDecoder(r.Body).Decode(&wordsModel); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected 2 demensional array of words: [[delete_word: string, collection_name: string]...]"})
			return 
		}
		userId, ok := r.Context().Value(KEY).(string)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "authorization error"})
			return
		}
		for _, v := range wordsModel.Words {
			if len(v) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{"result": "error", "message": "expected array of words: [[delete_word: string, collection_name: string]...]"})
				return
			}
		}
		h.service.DeleteWords(wordsModel, userId)
		json.NewEncoder(w).Encode(map[string]interface{}{"result": "ok", "message": "words were delted"})
		w.WriteHeader(http.StatusOK)
	})
} 
