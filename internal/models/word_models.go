package models

import (
	"time"
)

type Translation struct {
	SourceLanguage          string                `json:"source_language"`
	DestinationLanguage     string                `json:"destination_language"`
	Word                    string                `json:"word"`
	Translations            map[string][]string   `json:"transltions"`
	DefinitionsWithExamples map[string][][]string `json:"definitions_with_examples,omitempty"`
	Examples                []string              `json:"examples,omitempty"`
}

func (tr *Translation) Error() string {
	return "response doesn't have any translation data"
}

type WordToAdd struct {
	Word                string       `json:"word"`
	CollectionName      string       `json:"collection_name"`
	TimeOfLastRepeating time.Time    `json:"-"`
	TransData           *Translation `json:"-"`
}

type WordsToAdd struct {
	Words []WordToAdd `json:"words"`
}

type Word struct {
	Word                string      `json:"word"`
	State               int         `json:"state"`
	CollectionName      string      `json:"collection_name"`
	TimeOfLastRepeating time.Time   `json:"time_of_last_repeating"`
	TransData           Translation `json:"trans_data"`
}

type Words struct {
	Words *[]Word `json:"words"`
}

type WordToUpdate struct {
	OldWord             string       `json:"old_word"`
	NewWord             string       `json:"new_word"`
	CollectionName      string       `json:"collection_name"`
	TimeOfLastRepeating time.Time    `json:"-"`
	TransData           *Translation `json:"-"`
}

type WordsToDelete struct {
	Words []WordToDelete `json:"words"`
}

type WordToDelete struct {
	Word           string `json:"word"`
	CollectionName string `json:"collection_name"`
}

type stateToUpdate struct {
	Word                string    `json:"word"`
	NewState            int       `json:"new_state"`
	TimeOfLastRepeating time.Time `json:"-"`
}

type StatesToUpdate struct {
	CollectionName string          `json:"collection_name"`
	Words          []stateToUpdate `json:"words"`
}
