package service

import (
	"gitlab.com/behind-the-fridge/product/internal/pkg/category"
	"gitlab.com/behind-the-fridge/product/internal/pkg/list"
	"gitlab.com/behind-the-fridge/product/internal/pkg/product"
)


type RelationService struct {
	productRepository *product.Repository
	listRepository *list.Repository
	categoryRepository *category.Repository
}

func (ps *RelationService) GetProduct(id int) (product.DTO, error){

	p, err :=  ps.productRepository.Get(id)

	if err != nil{
		return product.DTO{}, err
	}

	productDTO := product.ModelToDTO(p)

	l, err := ps.listRepository.Get(productDTO.ListId)

	if err != nil{
		return product.DTO{}, err
	}

	productDTO.List = list.ModelToDTO(l)

	return productDTO, nil
}

func NewRelationService(productRepository *product.Repository, listRepository *list.Repository, categoryRepository *category.Repository) *RelationService {
	return &RelationService{
		productRepository: productRepository,
		listRepository: listRepository,
		categoryRepository: categoryRepository,
	}
}