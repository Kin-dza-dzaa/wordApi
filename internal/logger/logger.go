package logger

import 
(
	"github.com/rs/zerolog"
	"os"
)

func GetWriter() (*os.File, error) {
	writer, err := os.OpenFile("./internal/logger/loggs.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

func Getlogger(w *os.File) *zerolog.Logger{
	mylogger := zerolog.New(w).With().Caller().Timestamp().Logger()
	return &mylogger
}
