package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/proviant-io/core/internal/config"
	"github.com/proviant-io/core/internal/di"
	"github.com/proviant-io/core/internal/errors"
	"github.com/proviant-io/core/internal/i18n"
	"github.com/proviant-io/core/internal/pkg/category"
	"github.com/proviant-io/core/internal/pkg/list"
	"github.com/proviant-io/core/internal/pkg/product"
	"github.com/proviant-io/core/internal/pkg/product_category"
	"github.com/proviant-io/core/internal/pkg/service"
	"github.com/proviant-io/core/internal/pkg/stock"
	"log"
	"net/http"
	"strconv"
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
	cfg                 config.Config
	di                  *di.DI
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

func (s *Server) handleBadRequest(w http.ResponseWriter, locale i18n.Locale, error string, params ...interface{}) {
	m := i18n.NewMessage(error, params...)
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

func (s *Server) accountId(r *http.Request) int {

	accountHeader := r.Header.Get("AccountId")

	if accountHeader == "" {
		return 0
	}

	accountId, err := strconv.Atoi(accountHeader)

	if err != nil {
		log.Println(err)
		return -1
	}

	return accountId
}

func (s *Server) userId(r *http.Request) int {

	accountHeader := r.Header.Get("UserId")

	if accountHeader == "" {
		return 0
	}

	accountId, err := strconv.Atoi(accountHeader)

	if err != nil {
		log.Println(err)
		return -1
	}

	return accountId
}

func NewServer(productRepo *product.Repository,
	listRepo *list.Repository,
	categoryRepo *category.Repository,
	productCategoryRepo *product_category.Repository,
	stockRepo *stock.Repository,
	relationService *service.RelationService,
	l i18n.Localizer,
	i *di.DI) *Server {

	server := &Server{
		productRepo:         productRepo,
		listRepo:            listRepo,
		categoryRepo:        categoryRepo,
		productCategoryRepo: productCategoryRepo,
		stockRepo:           stockRepo,
		relationService:     relationService,
		l:                   l,
		di:                  i,
	}

	router := mux.NewRouter()

	apiV1Router := router.PathPrefix("/api/v1").Subrouter()

	// product routes
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{id}/", server.getProduct)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/", server.getProducts)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/", server.createProduct)).Methods("POST")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{id}/", server.updateProduct)).Methods("PUT")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{id}/", server.deleteProduct)).Methods("DELETE")
	// category routes
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/category/{id}/", server.getCategory)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/category/", server.getCategories)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/category/", server.createCategory)).Methods("POST")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/category/{id}/", server.updateCategory)).Methods("PUT")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/category/{id}/", server.deleteCategory)).Methods("DELETE")
	// list routes
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/list/{id}/", server.getList)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/list/", server.getLists)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/list/", server.createList)).Methods("POST")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/list/{id}/", server.updateList)).Methods("PUT")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/list/{id}/", server.deleteList)).Methods("DELETE")
	// stock routers
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{id}/stock/", server.getStock)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{id}/add/", server.addStock)).Methods("POST")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{id}/consume/", server.consumeStock)).Methods("POST")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{product_id}/stock/{id}/", server.deleteStock)).Methods("DELETE")
	// shopping list
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/shopping_list/", server.getShoppingLists)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/shopping_list/{id}/", server.getShoppingList)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/shopping_list/{id}/", server.addShoppingListItem)).Methods("POST")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/shopping_list/{list_id}/{id}/", server.updateShoppingListItem)).Methods("PUT")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/shopping_list/{list_id}/{id}/", server.deleteShoppingListItem)).Methods("DELETE")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/shopping_list/{list_id}/{id}/check/", server.checkShoppingListItem)).Methods("PUT")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/shopping_list/{list_id}/{id}/uncheck/", server.uncheckShoppingListItem)).Methods("PUT")
	// stock consumption log
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/product/{id}/consumption_log/", server.getConsumptionLog)).Methods("GET")



	// chore
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/i18n/missing/", server.getMissingTranslations)).Methods("GET")
	apiV1Router.HandleFunc(server.di.Apm.WrapHandleFunc("/version/", server.getVersion)).Methods("GET")

	userContentRouter := router.PathPrefix("/uc/").Subrouter()
	userContentRouter.HandleFunc(server.di.Apm.WrapHandleFunc("/img/{fileName}", server.getImage)).Methods("GET")

	if i.Cfg.Mode == config.ModeWeb {
		router.PathPrefix("/static").Handler(http.FileServer(http.Dir("./public/")))

		if i.Cfg.UserContent.Mode == config.UserContentModeLocal {
			router.PathPrefix("/content/").Handler(http.StripPrefix("/content/", http.FileServer(http.Dir(i.Cfg.UserContent.Location))))
		}

		spa := spaHandler{staticPath: "public", indexPath: "index.html"}
		router.PathPrefix("/").Handler(spa)
	}

	server.router = router

	return server
}
