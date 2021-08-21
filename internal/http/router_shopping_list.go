package http

import (
	"github.com/gorilla/mux"
	"github.com/proviant-io/core/internal/pkg/shopping"
	"github.com/proviant-io/core/internal/utils"
	"net/http"
	"strconv"
)

func (s *Server) getShoppingLists(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)

	models := s.di.ShoppingList.GetAll(accountId)

	// hack to create first list. Will be removed when multiple shopping lists will be released
	if len(models) == 0 {
		model := s.di.ShoppingList.Create(shopping.ListDTO{
			Title: "Shopping list",
		}, accountId)

		models = append(models, model)
	}

	var dtos []shopping.ListDTO

	for _, model := range models {
		dtos = append(dtos, shopping.ListToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.jsonResponse(w, response)
}

func (s *Server) getShoppingList(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	data, customErr := s.relationService.GetShoppingList(id, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   data,
	}

	s.jsonResponse(w, response)
}

func (s *Server) addShoppingListItem(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	listIdString := vars["id"]

	if listIdString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	listId, err := strconv.Atoi(listIdString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	dto := shopping.ItemDTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	dto.Title = utils.ClearString(dto.Title)

	if dto.Title == "" {
		s.handleBadRequest(w, locale, "title should not be empty")
		return
	}

	data, customErr := s.relationService.AddShoppingListItem(listId, dto, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeCreated,
		Data:   data,
	}

	s.jsonResponse(w, response)
}

func (s *Server) updateShoppingListItem(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	listIdString := vars["list_id"]

	if listIdString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	listId, err := strconv.Atoi(listIdString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	idString := vars["id"]

	if idString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	dto := shopping.ItemDTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	dto.Id = id

	dto.Title = utils.ClearString(dto.Title)

	if dto.Title == "" {
		s.handleBadRequest(w, locale, "title should not be empty")
		return
	}

	data, customErr := s.relationService.UpdateShoppingListItem(listId, dto, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   data,
	}

	s.jsonResponse(w, response)
}
func (s *Server) deleteShoppingListItem(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	customErr := s.di.ShoppingListItem.Delete(id, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
	}

	s.jsonResponse(w, response)
}

func (s *Server) checkShoppingListItem(w http.ResponseWriter, r *http.Request) {
	s.updateCheckedShoppingListItem(w, r, true)
}

func (s *Server) uncheckShoppingListItem(w http.ResponseWriter, r *http.Request) {
	s.updateCheckedShoppingListItem(w, r, true)
}

func (s *Server) updateCheckedShoppingListItem(w http.ResponseWriter, r *http.Request, checked bool) {
	accountId := s.accountId(r)
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	listIdString := vars["list_id"]

	if listIdString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	_, err := strconv.Atoi(listIdString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	idString := vars["id"]

	if idString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	data, customErr := s.relationService.UpdateCheckedShoppingListItem(id, checked, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeCreated,
		Data:   data,
	}

	s.jsonResponse(w, response)
}