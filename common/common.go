package common

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"unicode/utf8"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

type Configuration struct {
	AddressUsers       string `json:"address_users"`
	AddressThreads     string `json:"address_threads"`
	DbNameUsers        string `json:"db_name_users"`
	DbNameThreads      string `json:"db_name_threads"`
	ShowSQL            bool   `json:"show_sql"`
	LogToFile          bool   `json:"log_to_file"`
	LogFileNameUsers   string `json:"log_file_name_users"`
	LogFileNameThreads string `json:"log_file_name_threads"`
}

const (
	ConfigFileName = "../config.json"
	DbDriver       = "postgres"
	DbParameter    = "dbname=%s user=%s password=%s host=localhost port=5432 sslmode=disable"
)

const (
	LogInfoPrefix    = "[INFO]"
	LogWarningPrefix = "[WARNING]"
	LogErrorPrefix   = "[ERROR]"
)

func LoadConfig() (config *Configuration, err error) {
	file, err := os.Open(ConfigFileName)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	config = &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		return
	}
	return
}

// set maxConn<=0 if use default
func OpenDb(
	dbName string,
	showSQL bool,
	maxConn int,
) (dbEngine *xorm.Engine, err error) {
	dbEngine, err = xorm.NewEngine(
		DbDriver,
		fmt.Sprintf(
			DbParameter,
			dbName,
			os.Getenv("DBUSER"),
			os.Getenv("DBPASS"),
		),
	)
	if err != nil {
		return
	}
	dbEngine.ShowSQL(showSQL)
	if maxConn <= 0 {
		maxConn = runtime.NumCPU()
	}
	dbEngine.SetMaxOpenConns(maxConn)
	return
}

func OpenLogger(logToFile bool, logFileName string) (logger *log.Logger, err error) {
	if logToFile {
		var file *os.File
		file, err = os.OpenFile(
			fmt.Sprintf("%s.log", logFileName),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0666,
		)
		if err != nil {
			return
		}
		logger = log.New(
			file,
			LogInfoPrefix,
			log.Ldate|log.Ltime|log.Lshortfile,
		)
	} else {
		logger = log.Default()
		logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}
	return
}

func LogInfo(logger *log.Logger) *log.Logger {
	logger.SetPrefix(LogInfoPrefix)
	return logger
}

func LogWarning(logger *log.Logger) *log.Logger {
	logger.SetPrefix(LogWarningPrefix)
	return logger
}

func LogError(logger *log.Logger) *log.Logger {
	logger.SetPrefix(LogErrorPrefix)
	return logger
}

func NewUuIdString() string {
	raw := uuid.New()
	return raw.String()
}

func Encrypt(plainText string) (crypted string) {
	asBytes := sha256.Sum256([]byte(plainText))
	crypted = fmt.Sprintf("%x", asBytes)
	return
}

func IsEmpty(str ...string) bool {
	for _, s := range str {
		if utf8.RuneCountInString(s) == 0 {
			return true
		}
	}
	return false
}
