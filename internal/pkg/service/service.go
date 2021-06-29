package service

import (
	"github.com/brushknight/proviant/internal/errors"
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

func (ps *RelationService) GetProduct(id int) (product.DTO, *errors.CustomError) {

	p, err := ps.productRepository.Get(id)

	if err != nil {
		return product.DTO{}, err
	}

	productDTO := product.ModelToDTO(p)

	l, err := ps.listRepository.Get(productDTO.ListId)

	if err != nil {
		return product.DTO{}, err
	}

	productDTO.List = list.ModelToDTO(l)

	productCategories := ps.productCategoryRepository.GetByProductId(id)

	productDTO.Categories = []category.DTO{}

	for _, productCategory := range productCategories {
		productDTO.CategoryIds = append(productDTO.CategoryIds, productCategory.CategoryId)
	}

	categories := ps.categoryRepository.GetByIds(productDTO.CategoryIds)

	categoriesDTOs := []category.DTO{}

	for _, c := range categories {
		categoriesDTOs = append(categoriesDTOs, category.ModelToDTO(c))
	}

	productDTO.Categories = categoriesDTOs

	return productDTO, nil
}

func (ps *RelationService) GetAllProducts(query *product.Query) []product.DTO {

	models := ps.productRepository.GetAll(query)

	dtos := []product.DTO{}

	for _, model := range models {
		dtos = append(dtos, product.ModelToDTO(model))
	}

	// TODO apply category filtering here
	filteredDTOs := []product.DTO{}

	for idx := range dtos {

		dtos[idx].CategoryIds = []int{}
		dtos[idx].Categories = []interface{}{}

		productCategories := ps.productCategoryRepository.GetByProductId(dtos[idx].Id)

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

func (ps *RelationService) CreateProduct(dto product.CreateDTO) (product.DTO, *errors.CustomError) {

	_, err := ps.listRepository.Get(dto.ListId)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		for _, categoryId := range dto.CategoryIds {
			_, err := ps.categoryRepository.Get(categoryId)

			if err != nil {
				return product.DTO{}, err
			}
		}
	}

	p := ps.productRepository.Create(dto)

	if len(dto.CategoryIds) != 0 {
		ps.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return ps.GetProduct(p.Id)
}

func (ps *RelationService) UpdateProduct(dto product.DTO) (product.DTO, *errors.CustomError) {

	_, err := ps.listRepository.Get(dto.ListId)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		for _, categoryId := range dto.CategoryIds {
			_, err := ps.categoryRepository.Get(categoryId)

			if err != nil {
				return product.DTO{}, err
			}
		}
	}

	p, err := ps.productRepository.Update(dto)

	if err != nil {
		return product.DTO{}, err
	}

	if len(dto.CategoryIds) != 0 {
		ps.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return ps.GetProduct(p.Id)
}

func (ps *RelationService) AddStock(dto stock.DTO) (stock.Stock, *errors.CustomError) {

	_, err := ps.productRepository.Get(dto.ProductId)

	if err != nil {
		return stock.Stock{}, err
	}

	model := ps.stockRepository.Add(dto)

	return model, nil
}

func (ps *RelationService) ConsumeStock(dto stock.ConsumeDTO) *errors.CustomError {

	_, err := ps.productRepository.Get(dto.ProductId)

	if err != nil {
		return err
	}

	ps.stockRepository.Consume(dto)

	return nil
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
