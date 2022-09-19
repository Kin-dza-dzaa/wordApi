package models

import external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"

type WordToAdd struct {
	Word      				string 				 		`json:"word"`
	CollectionName 			string 				 		`json:"collection_name"`
}

type WordsAdd struct {
	Words 					[]WordToAdd					`json:"words"`
}

type WordToUpdate struct {
	OldWord      			string 				 		`json:"old_word"`
	NewWord      			string 				 		`json:"new_word"`
	CollectionName 			string 				 		`json:"collection_name"`
}

type WordsDelete struct {
	Words 					[][]string					`json:"words"`
}

type Word struct {
	Word      				string 				 		`json:"word"`
	State 	  				int 				 		`json:"state"`
	CollectionName 			string 				 		`json:"collection_name"`
	TransData 				external.Translation 		`json:"trans_data"`
}

type WordsGet struct {
	Words 					[]Word					 	`json:"words"`
}
