package http

import (
	"encoding/json"
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/brushknight/proviant/internal/pkg/product_category"
	"github.com/brushknight/proviant/internal/pkg/service"
	"github.com/brushknight/proviant/internal/pkg/stock"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	productRepo         *product.Repository
	listRepo            *list.Repository
	categoryRepo        *category.Repository
	productCategoryRepo *product_category.Repository
	stockRepo           *stock.Repository
	relationService     *service.RelationService
	router              *mux.Router
	l                   i18n.Localizer
}

func (s *Server) Run(hostPort string) error {
	return http.ListenAndServe(hostPort, s.router)
}

func (s *Server) parseJSON(r *http.Request, model interface{}) error {
	return json.NewDecoder(r.Body).Decode(model)
}

func (s *Server) getLocale(r *http.Request) i18n.Locale {
	return i18n.LocaleFromString(r.Header.Get("User-Locale"))
}


func (s *Server) handleBadRequest(w http.ResponseWriter,locale i18n.Locale,  error string, params ...interface{}) {
	m := i18n.NewMessage(error, params)
	response := Response{
		Status: BadRequest,
		Error:  s.l.T(m, locale),
	}

	s.jsonResponse(w, response)
}

func (s *Server) handleError(w http.ResponseWriter, locale i18n.Locale, error errors.CustomError) {
	response := Response{
		Status: error.Code(),
		Error:  s.l.T(error.Message(), locale),
	}

	s.jsonResponse(w, response)
}

func (s *Server) jsonResponse(w http.ResponseWriter, response Response) {
	payload, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	_, err = w.Write(payload)
	if err != nil {
		log.Println(err)
	}
	return
}

func NewServer(productRepo *product.Repository,
	listRepo *list.Repository,
	categoryRepo *category.Repository,
	productCategoryRepo *product_category.Repository,
	stockRepo *stock.Repository,
	relationService *service.RelationService,
	l i18n.Localizer) *Server {

	server := &Server{
		productRepo:         productRepo,
		listRepo:            listRepo,
		categoryRepo:        categoryRepo,
		productCategoryRepo: productCategoryRepo,
		stockRepo:           stockRepo,
		relationService:     relationService,
		l:                   l,
	}

	router := mux.NewRouter()

	// product routes
	router.HandleFunc("/api/v1/product/{id}/", server.getProduct).Methods("GET")
	router.HandleFunc("/api/v1/product/", server.getProducts).Methods("GET")
	router.HandleFunc("/api/v1/product/", server.createProduct).Methods("POST")
	router.HandleFunc("/api/v1/product/{id}/", server.updateProduct).Methods("PUT")
	router.HandleFunc("/api/v1/product/{id}/", server.deleteProduct).Methods("DELETE")
	// category routes
	router.HandleFunc("/api/v1/category/{id}/", server.getCategory).Methods("GET")
	router.HandleFunc("/api/v1/category/", server.getCategories).Methods("GET")
	router.HandleFunc("/api/v1/category/", server.createCategory).Methods("POST")
	router.HandleFunc("/api/v1/category/{id}/", server.updateCategory).Methods("PUT")
	router.HandleFunc("/api/v1/category/{id}/", server.deleteCategory).Methods("DELETE")
	// list routes
	router.HandleFunc("/api/v1/list/{id}/", server.getList).Methods("GET")
	router.HandleFunc("/api/v1/list/", server.getLists).Methods("GET")
	router.HandleFunc("/api/v1/list/", server.createList).Methods("POST")
	router.HandleFunc("/api/v1/list/{id}/", server.updateList).Methods("PUT")
	router.HandleFunc("/api/v1/list/{id}/", server.deleteList).Methods("DELETE")
	// stock routers
	router.HandleFunc("/api/v1/product/{id}/stock/", server.getStock).Methods("GET")
	router.HandleFunc("/api/v1/product/{id}/add/", server.addStock).Methods("POST")
	router.HandleFunc("/api/v1/product/{id}/consume/", server.consumeStock).Methods("POST")
	router.HandleFunc("/api/v1/product/{product_id}/stock/{id}/", server.deleteStock).Methods("DELETE")
	router.HandleFunc("/api/v1/i18n/missing/", server.getMissingTranslations).Methods("GET")

	router.PathPrefix("/static").Handler(http.FileServer(http.Dir("./public/")))

	spa := spaHandler{staticPath: "public", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	server.router = router

	return server
}
