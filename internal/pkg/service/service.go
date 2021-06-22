package service

import (
	"gitlab.com/behind-the-fridge/product/internal/pkg/category"
	"gitlab.com/behind-the-fridge/product/internal/pkg/list"
	"gitlab.com/behind-the-fridge/product/internal/pkg/product"
	"gitlab.com/behind-the-fridge/product/internal/pkg/product_category"
	"gitlab.com/behind-the-fridge/product/internal/pkg/stock"
)

type RelationService struct {
	productRepository  *product.Repository
	listRepository     *list.Repository
	categoryRepository *category.Repository
	stockRepository    *stock.Repository
	productCategoryRepository    *product_category.Repository
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

	productCategories := ps.productCategoryRepository.GetByProductId(id)

	productDTO.Categories = []category.DTO{}

	for _, productCategory := range productCategories{
		productDTO.CategoryIds = append(productDTO.CategoryIds, productCategory.CategoryId)
	}

	categories := ps.categoryRepository.GetByIds(productDTO.CategoryIds)

	categoriesDTOs := []category.DTO{}

	for _, c := range categories{
		categoriesDTOs = append(categoriesDTOs, category.ModelToDTO(c))
	}

	productDTO.Categories = categoriesDTOs

	return productDTO, nil
}

func (ps *RelationService) GetAllProducts() []product.DTO {

	models := ps.productRepository.GetAll()

	dtos := []product.DTO{}

	for _, model := range models{
		dtos = append(dtos, product.ModelToDTO(model))
	}

	for idx := range dtos{

		productCategories := ps.productCategoryRepository.GetByProductId(dtos[idx].Id)

		for _, productCategory := range productCategories{
			dtos[idx].CategoryIds = append(dtos[idx].CategoryIds, productCategory.CategoryId)
		}

	}

	return dtos
}

func (ps *RelationService) CreateProduct(dto product.DTO) error {

	_, err := ps.listRepository.Get(dto.ListId)

	if err != nil {
		return err
	}

	if len(dto.CategoryIds) != 0{
		for _, categoryId := range dto.CategoryIds{
			_, err := ps.categoryRepository.Get(categoryId)

			if err != nil {
				return err
			}
		}
	}

	p := ps.productRepository.Create(dto)

	if len(dto.CategoryIds) != 0 {
		ps.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return nil
}

func (ps *RelationService) UpdateProduct(dto product.DTO) error {

	_, err := ps.listRepository.Get(dto.ListId)

	if err != nil {
		return err
	}

	if len(dto.CategoryIds) != 0{
		for _, categoryId := range dto.CategoryIds{
			_, err := ps.categoryRepository.Get(categoryId)

			if err != nil {
				return err
			}
		}
	}

	p, err := ps.productRepository.Update(dto)

	if err != nil {
		return err
	}

	if len(dto.CategoryIds) != 0 {
		ps.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return nil
}

func (ps *RelationService) AddStock(dto stock.DTO) error {

	_, err := ps.productRepository.Get(dto.ProductId)

	if err != nil {
		return err
	}

	ps.stockRepository.Add(dto)

	return nil
}

func (ps *RelationService) ConsumeStock(dto stock.ConsumeDTO) error {

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
		productRepository:  productRepository,
		listRepository:     listRepository,
		categoryRepository: categoryRepository,
		stockRepository: stockRepository,
		productCategoryRepository: productCategoryRepository,
	}
}
