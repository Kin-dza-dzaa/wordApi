package external

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	config "github.com/Kin-dza-dzaa/wordApi/configs"
	"github.com/Kin-dza-dzaa/wordApi/internal/apierror"
	"github.com/Kin-dza-dzaa/wordApi/internal/models"
)

var (
	ErrUpdateFailed = errors.New("bad word")
)

func getExamples(translateData []interface{}, tranlationPtr *models.Translation) {
	arrayOfEx, ok := translateData[2].([]interface{})	
	if !ok {
		return
	}
	arrayOfEx, ok = arrayOfEx[0].([]interface{})	
	if !ok {
		return 
	}
	for _, value := range arrayOfEx {
		arrayWithEx, ok := value.([]interface{})
		if !ok {
			return 
		}
		example, ok := arrayWithEx[1].(string)
		if !ok {
			return 
		}
		tranlationPtr.Examples = append(tranlationPtr.Examples, example)
	}
}

func getDefinitions(translateData []interface{}, tranlationPtr *models.Translation) {
	arrayOfDef, ok := translateData[1].([]interface{})	
	if !ok {
		return
	}
	arrayOfDef, ok = arrayOfDef[0].([]interface{})	
	if !ok {
		return 
	}
	tranlationPtr.DefinitionsWithExamples = map[string][][]string{}
	for _, value := range arrayOfDef {
		arrayOfDefPart, ok := value.([]interface{})
		if !ok {
			return 
		}
		partOfSpeach, ok := arrayOfDefPart[0].(string)
		if !ok {
			continue
		}
		arrayOfDefPart, ok = arrayOfDefPart[1].([]interface{})
		if !ok {
			return
		}
		var defWithExpamples [][]string
		for _, v := range arrayOfDefPart {
			definitions, ok := v.([]interface{})
			if !ok {
				return 
			}
			def, ok := definitions[0].(string)
			if !ok {
				return
			}
			if len(definitions) == 1 {
				defWithExpamples = append(defWithExpamples, []string{def})
				continue
			} else {
				_, ok := definitions[1].(string)
				if !ok {
					defWithExpamples = append(defWithExpamples, []string{def})
					continue
				}
			}
			example, ok := definitions[1].(string)
			if !ok {
				return
			}
			defWithExpamples = append(defWithExpamples, []string{def, example})
			
		}
		tranlationPtr.DefinitionsWithExamples[partOfSpeach] = defWithExpamples
	}
} 

func getTranslations(translateData []interface{}, translationPtr *models.Translation) error {
	translationData, ok := translateData[5].([]interface{})
	if !ok {
		return translationPtr
	}
	translationData, ok = translationData[0].([]interface{})
	if !ok {
		return translationPtr
	}
	translationPtr.Translations = map[string][]string{}
	for _, value := range translationData {
		partOfSpeechData, ok := value.([]interface{})
		if !ok {
			return translationPtr
		}
		partOfSpeech, ok := partOfSpeechData[0].(string)
		if !ok {
			continue
		}
		wordTranslations, ok := partOfSpeechData[1].([]interface{})
		if !ok {
			return translationPtr
		}
		var translationSlice []string
		for _, v := range wordTranslations {
			wordTranslation, ok := v.([]interface{})
			if !ok {
				return translationPtr
			}
			translation, ok := wordTranslation[0].(string)
			if !ok {
				return translationPtr
			}
			translationSlice = append(translationSlice, translation)
		} 
		translationPtr.Translations[partOfSpeech] = append(translationPtr.Translations[partOfSpeech], translationSlice...)
	}
	return nil
}

func getJson(jsonResponse []interface{}, sourceLanguage string, destinationLanguage string, word string) (*models.Translation, error) {
	var translation *models.Translation = new(models.Translation)
	translation.SourceLanguage, translation.DestinationLanguage, translation.Word = sourceLanguage, destinationLanguage, word
	translateData, ok := jsonResponse[3].([]interface{})
	if !ok {
		return new(models.Translation), translation
	}
	err := getTranslations(translateData, translation)
	if err != nil {
		return new(models.Translation), err
	}
	getDefinitions(translateData, translation)
	getExamples(translateData, translation)
	return translation, nil
}

func unmarshalJsonTwice(rawJson []byte) ([]interface{}, error) {
	var words []interface{}
	rawJson = bytes.Split(rawJson, []byte("\n"))[3]
	err := json.Unmarshal(rawJson, &words)
	if err != nil {
		return nil, err
	}
	words, ok := words[0].([]interface{})
	if !ok {
		return nil, errors.New("assertion error")
	}
	stringRes, ok := words[2].(string)
	if !ok {
		return nil, errors.New("assertion error")
	}
	words = []interface{}{}
	err = json.Unmarshal([]byte(stringRes), &words)
	if err != nil {
		return nil, err
	}
	if len(words) != 4 {
		return nil, errors.New("word doesn't exist")
	}
	return words, nil
}

func postReq(word, sourceLanguage, destinationLanguage string, config *config.Config) ([]byte, error) {
	var slice [][][]interface{} = [][][]interface{}{{{"MkEWBc", fmt.Sprintf(`[["%s","%s","%s",true],[null]]`, word, sourceLanguage, destinationLanguage), nil, "generic"}}}
	data, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(config.ExternalUrl, "application/x-www-form-urlencoded", bytes.NewReader(append([]byte("f.req="), data...)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("error: status code - %d", response.StatusCode)
	}
	buff := make([]byte, 1024)
	var acceptedJson []byte
	for {
		n, err := response.Body.Read(buff)
		if err == io.EOF {
			break
		} else {
			if err != nil {
				return nil,  err
			}
		}
		acceptedJson = append(acceptedJson, buff[:n]...)
	}
	return acceptedJson, nil
}

func GetTranlations(word string, transData *models.Translation, sourceLanguage, destinationLanguage string, config *config.Config, channelBadWords chan string, wg *sync.WaitGroup) {
	// DON'T USE WITHOUT Sync.WaitGroup in gorutine 
	// also keep in mind that you've to work ONLY with one slice element per gorutine
	// it thread safe only if two requirements above are accomplished
	defer wg.Done()
	rawJson, err := postReq(word, sourceLanguage, destinationLanguage, config)
	if err != nil {
		channelBadWords <- word
		return 
	}
	unmarshalledJson, err := unmarshalJsonTwice(rawJson)
	if err != nil {
		channelBadWords <- word
		return 
	}
	translation, err := getJson(unmarshalledJson, sourceLanguage, destinationLanguage, word)
	if err != nil {
		channelBadWords <- word
		return 
	}
	*transData = *translation
}

func GetTranlationsUpdate(word string, transData *models.Translation, sourceLanguage, destinationLanguage string, config *config.Config) error {
	rawJson, err := postReq(word, sourceLanguage, destinationLanguage, config)
	if err != nil {
		return apierror.NewResponse("error", ErrUpdateFailed.Error(), http.StatusBadRequest)
	}
	unmarshalledJson, err := unmarshalJsonTwice(rawJson)
	if err != nil {
		return apierror.NewResponse("error", ErrUpdateFailed.Error(), http.StatusBadRequest)
	}
	translation, err := getJson(unmarshalledJson, sourceLanguage, destinationLanguage, word)
	if err != nil {
		return apierror.NewResponse("error", ErrUpdateFailed.Error(), http.StatusBadRequest)
	}
	*transData = *translation
	return nil
}
