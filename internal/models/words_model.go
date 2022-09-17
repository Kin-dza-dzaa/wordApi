package models

import external "github.com/Kin-dza-dzaa/wordApi/internal/external_call"

type WordsAdd struct {
	Words 					[][]string				`json:"words"`
}

type Word struct {
	Word      				string 				 	
	State 	  				int 				 	
	CollectionName 			string 				 	
	TransData 				external.Translation 	
}

type WordsGet struct {
	Words 					[]Word					`json:"words"`
}

type WordsUpdate struct {
	Words 					[]string				`json:"words"`
}

type WordsDelete struct {
	Words 					[][]string				`json:"words"`
}