package service

import (
	"fmt"
	"github.com/proviant-io/core/internal/config"
	"github.com/proviant-io/core/internal/di"
	"github.com/proviant-io/core/internal/errors"
	"github.com/proviant-io/core/internal/i18n"
	"github.com/proviant-io/core/internal/pkg/category"
	"github.com/proviant-io/core/internal/pkg/list"
	"github.com/proviant-io/core/internal/pkg/product"
	"github.com/proviant-io/core/internal/pkg/product_category"
	"github.com/proviant-io/core/internal/pkg/shopping"
	"github.com/proviant-io/core/internal/pkg/stock"
	"github.com/proviant-io/core/internal/utils"
	"log"
	"path"
	"strings"
)

type RelationService struct {
	productRepository         *product.Repository
	listRepository            *list.Repository
	categoryRepository        *category.Repository
	stockRepository           *stock.Repository
	productCategoryRepository *product_category.Repository
	di                        *di.DI
	config                    config.Config
}

func (s *RelationService) GetProduct(id int, accountId int) (product.DTO, *errors.CustomError) {

	p, err := s.productRepository.Get(id, accountId)

	if err != nil {
		return product.DTO{}, err
	}

	productDTO := product.ModelToDTO(p)

	l, err := s.listRepository.Get(productDTO.ListId, accountId)

	if err != nil {
		return product.DTO{}, err
	}

	productDTO.List = list.ModelToDTO(l)

	productCategories := s.productCategoryRepository.GetByProductId(id, accountId)

	productDTO.Categories = []category.DTO{}

	for _, productCategory := range productCategories {
		productDTO.CategoryIds = append(productDTO.CategoryIds, productCategory.CategoryId)
	}

	categories := s.categoryRepository.GetByIds(productDTO.CategoryIds, accountId)

	categoriesDTOs := []category.DTO{}

	for _, c := range categories {
		categoriesDTOs = append(categoriesDTOs, category.ModelToDTO(c))
	}

	productDTO.Categories = categoriesDTOs

	return productDTO, nil
}

func (s *RelationService) GetAllProducts(query *product.Query, accountId int) []product.DTO {

	models := s.productRepository.GetAll(query, accountId)

	dtos := []product.DTO{}

	for _, model := range models {
		dtos = append(dtos, product.ModelToDTO(model))
	}

	filteredDTOs := []product.DTO{}

	for idx := range dtos {

		dtos[idx].CategoryIds = []int{}
		dtos[idx].Categories = []interface{}{}

		productCategories := s.productCategoryRepository.GetByProductId(dtos[idx].Id, accountId)

		for _, productCategory := range productCategories {
			dtos[idx].CategoryIds = append(dtos[idx].CategoryIds, productCategory.CategoryId)
		}

		if query != nil && query.Category != 0 {
			if utils.ContainsInt(dtos[idx].CategoryIds, query.Category) {
				filteredDTOs = append(filteredDTOs, dtos[idx])
			}

		} else {
			filteredDTOs = append(filteredDTOs, dtos[idx])
		}

	}

	return filteredDTOs
}

func (s *RelationService) CreateProduct(dto product.CreateDTO, accountId int) (product.DTO, *errors.CustomError) {

	_, err := s.listRepository.Get(dto.ListId, accountId)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		for _, categoryId := range dto.CategoryIds {
			_, err := s.categoryRepository.Get(categoryId, accountId)

			if err != nil {
				return product.DTO{}, err
			}
		}
	}

	if dto.ImageBase64 != "" {
		imgPath, pureErr := s.di.ImageSaver.SaveBase64(dto.ImageBase64)
		if pureErr != nil {
			return product.DTO{}, errors.NewInternalServer(i18n.NewMessage(pureErr.Error()))
		}

		// convert imgPath into server accessable one
		imgPath = strings.Replace(imgPath, s.config.UserContent.Location, "", 1)
		imgPath = path.Join("/uc/img", imgPath)
		dto.Image = imgPath
	}

	p := s.productRepository.Create(dto, accountId)

	if len(dto.CategoryIds) != 0 {
		s.productCategoryRepository.Link(p.Id, dto.CategoryIds, accountId)
	}

	return s.GetProduct(p.Id, accountId)
}

func (s *RelationService) UpdateProduct(dto product.UpdateDTO, accountId int) (product.DTO, *errors.CustomError) {

	_, err := s.listRepository.Get(dto.ListId, accountId)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		for _, categoryId := range dto.CategoryIds {
			_, err := s.categoryRepository.Get(categoryId, accountId)

			if err != nil {
				return product.DTO{}, err
			}
		}
	}

	oldModel, err := s.productRepository.Get(dto.Id, accountId)

	if err != nil {
		return product.DTO{}, err
	}

	// stock should not be change via model update
	dto.Stock = oldModel.Stock

	// sanitize from custom urls
	if oldModel.Image != dto.Image {
		dto.Image = ""
	}

	if dto.ImageBase64 != "" {
		imgPath, pureErr := s.di.ImageSaver.SaveBase64(dto.ImageBase64)
		if pureErr != nil {
			return product.DTO{}, errors.NewInternalServer(i18n.NewMessage(pureErr.Error()))
		}

		// convert imgPath into server accessable
		imgPath = strings.Replace(imgPath, s.config.UserContent.Location, "", 1)
		imgPath = path.Join("/uc/img", imgPath)
		dto.Image = imgPath

		fileToRemove := strings.Replace(oldModel.Image, "/uc/img", "", 1)

		pureErr = s.di.ImageSaver.DeleteFile(fileToRemove)
		if pureErr != nil {
			log.Printf("cannot delete product image file: %s, %v\n", fileToRemove, pureErr)
		}
	}

	p, err := s.productRepository.UpdateFromDTO(dto, accountId)

	if err != nil {
		return product.DTO{}, err
	}

	// NOTE: here could be performance bottle neck
	s.productCategoryRepository.DeleteByProductId(p.Id, accountId)

	if len(dto.CategoryIds) != 0 {
		s.productCategoryRepository.Link(p.Id, dto.CategoryIds, accountId)
	}

	return s.GetProduct(p.Id, accountId)
}

func (s *RelationService) AddStock(dto stock.DTO, accountId int) (stock.Stock, *errors.CustomError) {

	p, err := s.productRepository.Get(dto.ProductId, accountId)

	if err != nil {
		return stock.Stock{}, err
	}

	model := s.stockRepository.Add(dto, accountId)

	p.Stock += dto.Quantity

	_, err = s.productRepository.Save(p, accountId)

	return model, err
}

func (s *RelationService) ConsumeStock(dto stock.ConsumeDTO, accountId int) *errors.CustomError {

	p, err := s.productRepository.Get(dto.ProductId, accountId)

	if err != nil {
		return err
	}

	s.stockRepository.Consume(dto, accountId)

	if dto.Quantity >= p.Stock {
		p.Stock = 0
	} else {
		p.Stock -= dto.Quantity
	}

	_, err = s.productRepository.Save(p, accountId)

	return nil
}

func (s *RelationService) DeleteStock(id int, accountId int) *errors.CustomError {

	st, err := s.stockRepository.Get(id, accountId)

	if err != nil {
		return err
	}

	p, err := s.productRepository.Get(st.ProductId, accountId)

	if err != nil {
		return err
	}

	err = s.stockRepository.Delete(id, accountId)

	if err != nil {
		return err
	}

	p.Stock -= st.Quantity

	if p.Stock < 0 {
		p.Stock = 0
	}

	_, err = s.productRepository.Save(p, accountId)

	return err
}

func (s *RelationService) DeleteProduct(id int, accountId int) *errors.CustomError {

	oldModel, err := s.productRepository.Get(id, accountId)

	if err != nil {
		return err
	}

	fileToRemove := strings.Replace(oldModel.Image, "/content", "", 1)

	pureErr := s.di.ImageSaver.DeleteFile(fileToRemove)
	if pureErr != nil {
		fmt.Printf("cannot delete product image file: %s, %v", fileToRemove, pureErr)
	}

	s.stockRepository.DeleteByProductId(id, accountId)

	s.productCategoryRepository.DeleteByProductId(id, accountId)

	err = s.productRepository.Delete(id, accountId)

	return err
}

func (s *RelationService) DeleteCategory(id int, accountId int) *errors.CustomError {

	s.productCategoryRepository.DeleteByCategory(id, accountId)

	return s.categoryRepository.Delete(id, accountId)
}

func (s *RelationService) DeleteList(id int, accountId int) *errors.CustomError {

	q := &product.Query{
		List: id,
	}

	models := s.productRepository.GetAll(q, accountId)

	if len(models) > 0 {
		return errors.NewErrBadRequest(i18n.NewMessage("You can't remove list with products. Clean products first."))
	}

	return s.listRepository.Delete(id, accountId)
}

func (s *RelationService) GetShoppingList(id, accountId int) (shopping.ListFilledDTO, *errors.CustomError) {

	listModel, err := s.di.ShoppingList.Get(id, accountId)

	if err != nil {
		return shopping.ListFilledDTO{}, err
	}

	items := s.di.ShoppingListItem.GetAllByList(listModel.Id, accountId)

	itemsDTO := []shopping.ItemDTO{}

	for _, item := range items {
		itemsDTO = append(itemsDTO, shopping.ItemToDTO(item))
	}

	dto := shopping.ListFilledDTO{
		Id:    listModel.Id,
		Title: listModel.Title,
		Items: itemsDTO,
	}

	return dto, nil
}

func (s *RelationService) AddShoppingListItem(id int, dto shopping.ItemDTO, accountId int) (shopping.ItemDTO, *errors.CustomError) {

	listModel, err := s.di.ShoppingList.Get(id, accountId)

	if err != nil {
		return shopping.ItemDTO{}, err
	}

	dto.ListId = listModel.Id

	item := s.di.ShoppingListItem.Create(dto, accountId)

	return shopping.ItemToDTO(item), nil
}

func (s *RelationService) UpdateShoppingListItem(id int, dto shopping.ItemDTO, accountId int) (shopping.ItemDTO, *errors.CustomError) {

	listModel, err := s.di.ShoppingList.Get(id, accountId)

	if err != nil {
		return shopping.ItemDTO{}, err
	}

	dto.ListId = listModel.Id

	item, err := s.di.ShoppingListItem.Update(dto.Id, dto, accountId)

	if err != nil {
		return shopping.ItemDTO{}, err
	}

	return shopping.ItemToDTO(item), nil
}


func (s *RelationService) UpdateCheckedShoppingListItem( id int, checked bool, accountId int) (shopping.ItemDTO, *errors.CustomError) {

	var item shopping.Item
	var err *errors.CustomError

	if checked {
		item, err = s.di.ShoppingListItem.Check(id, accountId)
	}else {
		item, err = s.di.ShoppingListItem.Uncheck(id, accountId)
	}

	if err != nil {
		return shopping.ItemDTO{}, err
	}

	return shopping.ItemToDTO(item), nil
}


func NewRelationService(productRepository *product.Repository,
	listRepository *list.Repository,
	categoryRepository *category.Repository,
	stockRepository *stock.Repository,
	productCategoryRepository *product_category.Repository,
	i *di.DI,
	config config.Config,
) *RelationService {
	return &RelationService{
		productRepository:         productRepository,
		listRepository:            listRepository,
		categoryRepository:        categoryRepository,
		stockRepository:           stockRepository,
		productCategoryRepository: productCategoryRepository,
		di:                        i,
		config:                    config,
	}
}
