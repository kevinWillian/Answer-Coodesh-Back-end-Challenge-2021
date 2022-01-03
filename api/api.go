package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kevinWillian/Answer-Coodesh-Back-end-Challenge-2021/models"
	"gorm.io/gorm"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/rs/cors"
)

func Paginate(alimit int, apage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := apage

		if page == 0 {
			page = 1
		}

		pageSize := alimit

		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
}

func StartRoutes(db *gorm.DB) {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	corsOpts := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsOpts.Handler)
	r.Use(middleware.Timeout(60 * time.Second))

	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("Back-end Challenge 2021 🏅 - Space Flight News"))
	})

	r.HandleFunc("/articles/sync/{id}", func(rw http.ResponseWriter, r *http.Request) {

		var resp []byte
		var err error

		id := chi.URLParam(r, "id")

		var idhost models.Article
		db.Where("synced_id = ?", id).First(&idhost)

		if idhost.ID == 0 {
			rw.Write([]byte("Registro não encontrado"))
			return
		}

		switch r.Method {
		case "PUT":
			var putArt models.Article

			err = json.NewDecoder(r.Body).Decode(&putArt)

			if err != nil {
				return
			}

			putArt.ID = idhost.ID

			if putArt.SyncedID == 0 {
				_id, err := strconv.Atoi(id)
				if err != nil {
					return
				}
				putArt.SyncedID = uint32(_id)
			}

			err = db.Save(&putArt).Error

			if err != nil {
				rw.Write([]byte("Falhar ao atualizar"))
				return
			}
		case "DELETE":
			igart := models.IgnoredArticle{
				SyncedID:    idhost.SyncedID,
				Featured:    idhost.Featured,
				Title:       idhost.Title,
				Url:         idhost.Url,
				ImageUrl:    idhost.ImageUrl,
				NewsSite:    idhost.NewsSite,
				Summary:     idhost.Summary,
				PublishedAt: idhost.PublishedAt,
			}
			err = db.Create(&igart).Error
			if err != nil {
				rw.Write([]byte("Falhar ao deletar"))
				return
			}
			err = db.Unscoped().Delete(&idhost).Error
			if err != nil {
				rw.Write([]byte("Falhar ao deletar"))
				return
			}
		default:
			resp, err = json.Marshal(idhost)

			if err != nil {
				return
			}
			rw.Write(resp)
		}
	})

	r.HandleFunc("/articles/{id}", func(rw http.ResponseWriter, r *http.Request) {

		var resp []byte
		var err error

		id := chi.URLParam(r, "id")

		var idhost models.Article
		db.Where("id = ?", id).First(&idhost)

		if idhost.ID == 0 {
			rw.Write([]byte("Registro não encontrado"))
			return
		}

		switch r.Method {
		case "PUT":

			var putArt models.Article

			err = json.NewDecoder(r.Body).Decode(&putArt)

			if err != nil {
				return
			}

			if putArt.ID == 0 {
				_id, err := strconv.Atoi(id)
				if err != nil {
					return
				}
				putArt.ID = uint(_id)
			}

			err = db.Save(&putArt).Error

			if err != nil {
				rw.Write([]byte("Falhar ao atualizar"))
				return
			}

		case "DELETE":
			if idhost.SyncedID != 0 {
				igart := models.IgnoredArticle{
					SyncedID:    idhost.SyncedID,
					Featured:    idhost.Featured,
					Title:       idhost.Title,
					Url:         idhost.Url,
					ImageUrl:    idhost.ImageUrl,
					NewsSite:    idhost.NewsSite,
					Summary:     idhost.Summary,
					PublishedAt: idhost.PublishedAt,
				}
				err = db.Create(&igart).Error
				if err != nil {
					rw.Write([]byte("Falhar ao deletar"))
					return
				}
			}
			err = db.Unscoped().Delete(&idhost).Error
			if err != nil {
				rw.Write([]byte("Falhar ao deletar"))
				return
			}
		default:
			resp, err = json.Marshal(idhost)

			if err != nil {
				return
			}
			rw.Write(resp)
		}
	})

	r.HandleFunc("/articles", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			rw.WriteHeader(http.StatusOK)
			http.Redirect(rw, r, "http://localhost:3000/articlesl=10s=1", 301)
		case "POST":
			rw.WriteHeader(http.StatusOK)

			var postArt models.Article

			err := json.NewDecoder(r.Body).Decode(&postArt)

			if err != nil {
				return
			}

			err = db.Create(&postArt).Error

			if err != nil {
				return
			}

			rw.Write([]byte(postArt.Title + " Adicionado!"))

		default:
			http.Redirect(rw, r, "http://localhost:3000/articlesl=10s=1", 301)
		}
	})

	r.HandleFunc("/articlesl={limit}s={start}", func(rw http.ResponseWriter, r *http.Request) {
		pglimit := chi.URLParam(r, "limit")
		pgstart := chi.URLParam(r, "start")

		intpglimit, err := strconv.Atoi(pglimit)

		if err != nil {
			log.Println("o limite da página nãoo é um número")
		}

		intpgstart, err := strconv.Atoi(pgstart)

		if err != nil {
			log.Println("o número da página não é um número")
		}

		var as []models.Article

		db.Scopes(Paginate(intpglimit, intpgstart)).Find(&as)

		resp, err := json.Marshal(as)

		if err != nil {
			return
		}

		rw.Write(resp)
	})

	log.Fatal(http.ListenAndServe(":3000", r))

}
