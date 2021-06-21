package service

import (
	"gitlab.com/behind-the-fridge/product/internal/pkg/category"
	"gitlab.com/behind-the-fridge/product/internal/pkg/list"
	"gitlab.com/behind-the-fridge/product/internal/pkg/product"
	"gitlab.com/behind-the-fridge/product/internal/pkg/stock"
)

type RelationService struct {
	productRepository  *product.Repository
	listRepository     *list.Repository
	categoryRepository *category.Repository
	stockRepository    *stock.Repository
}

func (ps *RelationService) GetProduct(id int) (product.DTO, error) {

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

	return productDTO, nil
}

func (ps *RelationService) CreateProduct(dto product.DTO) error {

	_, err := ps.listRepository.Get(dto.ListId)

	if err != nil {
		return err
	}

	ps.productRepository.Create(dto)

	return nil
}

func (ps *RelationService) AddStock(dto stock.DTO) error {

	_, err := ps.productRepository.Get(dto.Id)

	if err != nil {
		return err
	}

	ps.stockRepository.Add(dto)

	return nil
}

func (ps *RelationService) ConsumeStock(dto stock.DTO) error {

	_, err := ps.productRepository.Get(dto.Id)

	if err != nil {
		return err
	}

	ps.stockRepository.Consume(dto)

	return nil
}

func NewRelationService(productRepository *product.Repository, listRepository *list.Repository, categoryRepository *category.Repository, stockRepository *stock.Repository) *RelationService {
	return &RelationService{
		productRepository:  productRepository,
		listRepository:     listRepository,
		categoryRepository: categoryRepository,
		stockRepository: stockRepository,
	}
}
