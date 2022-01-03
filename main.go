package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/kevinWillian/Answer-Coodesh-Back-end-Challenge-2021/api"
	"github.com/kevinWillian/Answer-Coodesh-Back-end-Challenge-2021/models"
	"github.com/kevinWillian/Answer-Coodesh-Back-end-Challenge-2021/synchronizer"

	"github.com/robfig/cron/v3"
)

type ExtConfig struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     string
	Sslmode  string
	TimeZone string
}

func Connected() (ok bool) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	_, err := client.Get("https://www.google.com.br/")

	return err == nil
}

var dsn string
var db *gorm.DB
var dberr error

func ConnectDB(cfg ExtConfig) {
	dsn = "host=" + cfg.Host + " user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.Dbname + " port=" + cfg.Port + " sslmode=" + cfg.Sslmode + " TimeZone=" + cfg.TimeZone
	db, dberr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func migreteDB() error {
	return db.AutoMigrate(&models.Launche{}, &models.Event{}, &models.Article{}, &models.Config{}, &models.IgnoredArticle{})
}

func main() {

	if !Connected() {
		log.Fatalln("Verifique sua conexão com a internet, necessária para o funcionamento deste software!")
		return
	}

	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	dbconfig := ExtConfig{}
	err := decoder.Decode(&dbconfig)

	if err != nil {
		log.Println("erro na sintax do 'conf.json'")
		log.Fatalln(err.Error())
		return
	}

	ConnectDB(dbconfig)

	if dberr != nil {
		log.Println("Falha na inicializacao, o seguinte erro occoreu: " + dberr.Error())
		return
	}

	dberr = migreteDB()

	if dberr != nil {
		log.Println("Falha na criacao do banco de dados, o seguinte erro occoreu: " + dberr.Error())
		return
	}

	synchronizer.SyncWithSpaceflightnewsapi(db)

	c := cron.New()

	c.AddFunc("* 9 * * *", func() {
		synchronizer.UpdateSyncWithSpaceflightnewsapi(db)
	})

	c.Start()

	api.StartRoutes(db)
}
