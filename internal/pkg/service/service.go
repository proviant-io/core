package service

import (
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/brushknight/proviant/internal/pkg/product_category"
	"github.com/brushknight/proviant/internal/pkg/stock"
	"github.com/brushknight/proviant/internal/utils"
)

type RelationService struct {
	productRepository         *product.Repository
	listRepository            *list.Repository
	categoryRepository        *category.Repository
	stockRepository           *stock.Repository
	productCategoryRepository *product_category.Repository
}

func (s *RelationService) GetProduct(id int) (product.DTO, *errors.CustomError) {

	p, err := s.productRepository.Get(id)

	if err != nil {
		return product.DTO{}, err
	}

	productDTO := product.ModelToDTO(p)

	l, err := s.listRepository.Get(productDTO.ListId)

	if err != nil {
		return product.DTO{}, err
	}

	productDTO.List = list.ModelToDTO(l)

	productCategories := s.productCategoryRepository.GetByProductId(id)

	productDTO.Categories = []category.DTO{}

	for _, productCategory := range productCategories {
		productDTO.CategoryIds = append(productDTO.CategoryIds, productCategory.CategoryId)
	}

	categories := s.categoryRepository.GetByIds(productDTO.CategoryIds)

	categoriesDTOs := []category.DTO{}

	for _, c := range categories {
		categoriesDTOs = append(categoriesDTOs, category.ModelToDTO(c))
	}

	productDTO.Categories = categoriesDTOs

	return productDTO, nil
}

func (s *RelationService) GetAllProducts(query *product.Query) []product.DTO {

	models := s.productRepository.GetAll(query)

	dtos := []product.DTO{}

	for _, model := range models {
		dtos = append(dtos, product.ModelToDTO(model))
	}

	// TODO apply category filtering here
	filteredDTOs := []product.DTO{}

	for idx := range dtos {

		dtos[idx].CategoryIds = []int{}
		dtos[idx].Categories = []interface{}{}

		productCategories := s.productCategoryRepository.GetByProductId(dtos[idx].Id)

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

func (s *RelationService) CreateProduct(dto product.CreateDTO) (product.DTO, *errors.CustomError) {

	_, err := s.listRepository.Get(dto.ListId)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		for _, categoryId := range dto.CategoryIds {
			_, err := s.categoryRepository.Get(categoryId)

			if err != nil {
				return product.DTO{}, err
			}
		}
	}

	p := s.productRepository.Create(dto)

	if len(dto.CategoryIds) != 0 {
		s.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return s.GetProduct(p.Id)
}

func (s *RelationService) UpdateProduct(dto product.DTO) (product.DTO, *errors.CustomError) {

	_, err := s.listRepository.Get(dto.ListId)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		for _, categoryId := range dto.CategoryIds {
			_, err := s.categoryRepository.Get(categoryId)

			if err != nil {
				return product.DTO{}, err
			}
		}
	}

	p, err := s.productRepository.Update(dto)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		s.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return s.GetProduct(p.Id)
}

func (s *RelationService) AddStock(dto stock.DTO) (stock.Stock, *errors.CustomError) {

	_, err := s.productRepository.Get(dto.ProductId)

	if err != nil {
		return stock.Stock{}, err
	}

	model := s.stockRepository.Add(dto)

	return model, nil
}

func (s *RelationService) ConsumeStock(dto stock.ConsumeDTO) *errors.CustomError {

	_, err := s.productRepository.Get(dto.ProductId)

	if err != nil {
		return err
	}

	s.stockRepository.Consume(dto)

	return nil
}

func (s * RelationService) DeleteProduct(id int) *errors.CustomError{

	s.stockRepository.DeleteByProductId(id)

	s.productCategoryRepository.DeleteByProductId(id)

	err := s.productRepository.Delete(id)

	return err
}

func (s *RelationService) DeleteCategory(id int) *errors.CustomError{

	s.productCategoryRepository.DeleteByCategory(id)

	return s.categoryRepository.Delete(id)
}

func (s *RelationService) DeleteList(id int) *errors.CustomError{

	q := &product.Query{
		List: id,
	}

	models := s.productRepository.GetAll(q)

	if len(models) > 0 {
		return errors.NewErrBadRequest(i18n.NewMessage("You can't remove list with products. Clean products first."))
	}

	return s.listRepository.Delete(id)
}

func NewRelationService(productRepository *product.Repository,
	listRepository *list.Repository,
	categoryRepository *category.Repository,
	stockRepository *stock.Repository,
	productCategoryRepository *product_category.Repository) *RelationService {
	return &RelationService{
		productRepository:         productRepository,
		listRepository:            listRepository,
		categoryRepository:        categoryRepository,
		stockRepository:           stockRepository,
		productCategoryRepository: productCategoryRepository,
	}
}
