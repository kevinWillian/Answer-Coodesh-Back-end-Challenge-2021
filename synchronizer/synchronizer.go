package synchronizer

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/kevinWillian/Answer-Coodesh-Back-end-Challenge-2021/models"
)

var Totalextdb = 0
var Totalsync = 0
var BaseUrl = "https://api.spaceflightnewsapi.net/v3"

func SyncWithSpaceflightnewsapi(db *gorm.DB) {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	httpresp, err := client.Get(BaseUrl + "/articles/count")

	if err != nil {
		log.Println("Erro ao obter contagem do space flight news api")
		return
	}

	json.NewDecoder(httpresp.Body).Decode(&Totalextdb)

	var dbTotalsync models.Config

	db.Where("key = ?", "Totalsync").First(&dbTotalsync)

	if dbTotalsync.Key != "" {
		Totalsync, err = strconv.Atoi(dbTotalsync.Vaule)
		if err != nil {
			return
		}
	} else {
		Totalsync = 0
	}

	if Totalsync == 0 {

		for i := (Totalsync / 500); Totalsync < Totalextdb; i++ {

			var arts []models.Article

			start := Totalsync
			limit := 500 - (Totalsync % 500)

			httpresp, err = client.Get(BaseUrl + "/articles?_limit=" + strconv.Itoa(limit) + "&_start=" + strconv.Itoa(start) + "&_sort=id")

			if err != nil {
				log.Println(err.Error())
				return
			}

			json.NewDecoder(httpresp.Body).Decode(&arts)

			err = db.Create(&arts).Error

			if err != nil {
				log.Fatalln(err)
			}

			Totalsync += len(arts)

			configTotalsync := models.Config{
				Key:   "Totalsync",
				Vaule: strconv.Itoa(Totalsync),
			}

			db.Save(&configTotalsync)

			progress := (Totalsync * 100) / Totalextdb

			log.Println("Sincronizacao de banco em " + strconv.Itoa(progress) + "%")

		}

	}

}

func UpdateSyncWithSpaceflightnewsapi(db *gorm.DB) {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	httpresp, err := client.Get(BaseUrl + "/articles/count")

	if err != nil {
		return
	}

	json.NewDecoder(httpresp.Body).Decode(&Totalextdb)

	var dbTotalsync models.Config

	db.Where("key = ?", "Totalsync").First(&dbTotalsync)

	if dbTotalsync.Key != "" {
		Totalsync, err = strconv.Atoi(dbTotalsync.Vaule)
		if err != nil {
			return
		}
	} else {
		Totalsync = 0
	}

	var addd = 0

	synced := false

	for start := 0; !synced; start++ {

		var arts []models.Article
		var art models.Article

		useurl := BaseUrl + "/articles?_limit=1" + "&_start=" + strconv.Itoa(start) + "&_sort=id%3Adesc"
		httpresp, err = client.Get(useurl)

		if err != nil {
			log.Println(err.Error())
			return
		}

		json.NewDecoder(httpresp.Body).Decode(&arts)

		art = arts[0]

		var dbart models.Article

		db.Where("synced_id = ?", art.SyncedID).First(&dbart)

		if dbart.ID == 0 {
			var dbigart models.IgnoredArticle

			db.Where("synced_id = ?", art.SyncedID).First(&dbigart)

			if dbigart.ID == 0 {
				err = db.Create(&art).Error
				if err != nil {
					log.Println(err.Error())
					return
				}
				addd++
				Totalsync++

				configTotalsync := models.Config{
					Key:   "Totalsync",
					Vaule: strconv.Itoa(Totalsync),
				}

				db.Save(&configTotalsync)
			} else {
				synced = true
			}
		} else {
			synced = true
		}

	}

}
