package http

import (
	"encoding/json"
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
}

func (s *Server) Run(hostPort string) error {
	return http.ListenAndServe(hostPort, s.router)
}

func (s *Server) parseJSON(r *http.Request, model interface{}) error {
	return json.NewDecoder(r.Body).Decode(model)
}

func (s *Server) JSONResponse(w http.ResponseWriter, response Response) {
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
	relationService *service.RelationService) *Server {

	server := &Server{
		productRepo:         productRepo,
		listRepo:            listRepo,
		categoryRepo:        categoryRepo,
		productCategoryRepo: productCategoryRepo,
		stockRepo:           stockRepo,
		relationService:     relationService,
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

	router.PathPrefix("/static").Handler(http.FileServer(http.Dir("./public/")))

	spa := spaHandler{staticPath: "public", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	server.router = router

	return server
}
