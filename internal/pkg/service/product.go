package service

import (
	"gitlab.com/behind-the-fridge/product/internal/pkg/category"
	"gitlab.com/behind-the-fridge/product/internal/pkg/list"
	"gitlab.com/behind-the-fridge/product/internal/pkg/product"
)


type ProductService struct {
	productRepository *product.Repository
	listRepository *list.Repository
	categoryRepository *category.Repository
}

func (ps *ProductService) GetProduct(id int) (product.DTO, error){

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

func NewProductService(productRepository *product.Repository, listRepository *list.Repository, categoryRepository *category.Repository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
		listRepository: listRepository,
		categoryRepository: categoryRepository,
	}
}