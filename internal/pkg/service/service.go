package service

import (
	"fmt"
	"github.com/brushknight/proviant/internal/config"
	"github.com/brushknight/proviant/internal/errors"
	"github.com/brushknight/proviant/internal/i18n"
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/brushknight/proviant/internal/pkg/image"
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/brushknight/proviant/internal/pkg/product_category"
	"github.com/brushknight/proviant/internal/pkg/stock"
	"github.com/brushknight/proviant/internal/utils"
	"path"
	"strings"
)

type RelationService struct {
	productRepository         *product.Repository
	listRepository            *list.Repository
	categoryRepository        *category.Repository
	stockRepository           *stock.Repository
	productCategoryRepository *product_category.Repository
	imageSaver                image.Saver
	config                    config.Config
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

	if dto.ImageBase64 != "" {
		imgPath, pureErr := s.imageSaver.SaveBase64(dto.ImageBase64)
		if pureErr != nil {
			return product.DTO{}, errors.NewInternalServer(i18n.NewMessage(pureErr.Error()))
		}

		// convert imgPath into server accessable one
		imgPath = strings.Replace(imgPath, s.config.UserContent.Location, "", 1)
		imgPath = path.Join("/content", imgPath)
		dto.Image = imgPath
	}

	p := s.productRepository.Create(dto)

	if len(dto.CategoryIds) != 0 {
		s.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return s.GetProduct(p.Id)
}

func (s *RelationService) UpdateProduct(dto product.UpdateDTO) (product.DTO, *errors.CustomError) {

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

	oldModel, err := s.productRepository.Get(dto.Id)

	if err != nil {
		return product.DTO{}, err
	}

	// todo - remove old images

	if dto.ImageBase64 != "" {
		imgPath, pureErr := s.imageSaver.SaveBase64(dto.ImageBase64)
		if pureErr != nil {
			return product.DTO{}, errors.NewInternalServer(i18n.NewMessage(pureErr.Error()))
		}

		// convert imgPath into server accessable
		imgPath = strings.Replace(imgPath, s.config.UserContent.Location, "", 1)
		imgPath = path.Join("/content", imgPath)
		dto.Image = imgPath

		fileToRemove := strings.Replace(oldModel.Image, "/content", "", 1)

		pureErr = s.imageSaver.DeleteFile(fileToRemove)
		if pureErr != nil {
			fmt.Printf("cannot delete product image file: %s, %v", fileToRemove, pureErr)
		}
	}

	p, err := s.productRepository.UpdateFromDTO(dto)

	if err != nil {
		return product.DTO{}, err
	}

	// NOTE: here could be performance bottle neck
	s.productCategoryRepository.DeleteByProductId(p.Id)

	if len(dto.CategoryIds) != 0 {
		s.productCategoryRepository.Link(p.Id, dto.CategoryIds)
	}

	return s.GetProduct(p.Id)
}

func (s *RelationService) AddStock(dto stock.DTO) (stock.Stock, *errors.CustomError) {

	p, err := s.productRepository.Get(dto.ProductId)

	if err != nil {
		return stock.Stock{}, err
	}

	model := s.stockRepository.Add(dto)

	p.Stock += dto.Quantity

	_, err = s.productRepository.Save(p)

	return model, err
}

func (s *RelationService) ConsumeStock(dto stock.ConsumeDTO) *errors.CustomError {

	p, err := s.productRepository.Get(dto.ProductId)

	if err != nil {
		return err
	}

	s.stockRepository.Consume(dto)

	if dto.Quantity >= p.Stock {
		p.Stock = 0
	} else {
		p.Stock -= dto.Quantity
	}

	_, err = s.productRepository.Save(p)

	return nil
}

func (s *RelationService) DeleteStock(id int) *errors.CustomError {

	st, err := s.stockRepository.Get(id)

	if err != nil {
		return err
	}

	p, err := s.productRepository.Get(st.ProductId)

	if err != nil {
		return err
	}

	err = s.stockRepository.Delete(id)

	if err != nil {
		return err
	}

	p.Stock -= st.Quantity

	if p.Stock < 0 {
		p.Stock = 0
	}

	_, err = s.productRepository.Save(p)

	return err
}

func (s *RelationService) DeleteProduct(id int) *errors.CustomError {

	oldModel, err := s.productRepository.Get(id)

	if err != nil {
		return err
	}

	fileToRemove := strings.Replace(oldModel.Image, "/content", "", 1)

	pureErr := s.imageSaver.DeleteFile(fileToRemove)
	if pureErr != nil {
		fmt.Printf("cannot delete product image file: %s, %v", fileToRemove, pureErr)
	}

	s.stockRepository.DeleteByProductId(id)

	s.productCategoryRepository.DeleteByProductId(id)

	err = s.productRepository.Delete(id)

	return err
}

func (s *RelationService) DeleteCategory(id int) *errors.CustomError {

	s.productCategoryRepository.DeleteByCategory(id)

	return s.categoryRepository.Delete(id)
}

func (s *RelationService) DeleteList(id int) *errors.CustomError {

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
	productCategoryRepository *product_category.Repository,
	imageSaver image.Saver,
	config config.Config,
) *RelationService {
	return &RelationService{
		productRepository:         productRepository,
		listRepository:            listRepository,
		categoryRepository:        categoryRepository,
		stockRepository:           stockRepository,
		productCategoryRepository: productCategoryRepository,
		imageSaver:                imageSaver,
		config:                    config,
	}
}
