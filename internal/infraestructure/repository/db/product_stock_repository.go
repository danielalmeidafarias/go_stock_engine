package db

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"gorm.io/gorm"
)

type ErrorMapper interface {
	MapErrorToDomain(err error, context string) *domain.Error
}

type ProductStockRepository struct {
	db          *gorm.DB
	dbErrMapper ErrorMapper
}

func NewProductStockRepository(gorm *gorm.DB, errMapper ErrorMapper) *ProductStockRepository {
	return &ProductStockRepository{db: gorm, dbErrMapper: errMapper}
}

func (r *ProductStockRepository) Create(in *entities.ProductStock) (string, *domain.Error) {
	model := MapProductStockToModel(in)

	if err := r.db.Create(model).Error; err != nil {
		return "", r.dbErrMapper.MapErrorToDomain(err, "failed to create product")
	}

	return model.ID, nil
}

func (r *ProductStockRepository) Update(in *entities.ProductStock) *domain.Error {
	model := MapProductStockToModel(in)

	result := r.db.Save(model)
	if result.Error != nil {
		return r.dbErrMapper.MapErrorToDomain(result.Error, "failed to update product")
	}

	if result.RowsAffected == 0 {
		return domain.NewError("product not found", domain.ErrNotFound)
	}

	return nil
}

func (r *ProductStockRepository) GetAll(pagination *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	var models []ProductStockModel

	query := r.db.Model(&ProductStockModel{})

	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, r.dbErrMapper.MapErrorToDomain(err, "failed to list products")
	}

	result := make([]*entities.ProductStock, len(models))
	for i := range models {
		result[i] = models[i].ToDomain()
	}

	return result, nil
}

func (r *ProductStockRepository) GetOneByID(id string) (*entities.ProductStock, *domain.Error) {
	var model ProductStockModel

	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, r.dbErrMapper.MapErrorToDomain(err, "failed to get product")
	}

	return model.ToDomain(), nil
}

func (r *ProductStockRepository) GetByCategory(category entities.ProductCategory, pagination *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	var models []ProductStockModel

	query := r.db.Where("category = ?", string(category))

	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, r.dbErrMapper.MapErrorToDomain(err, "failed to get products by category")
	}

	result := make([]*entities.ProductStock, len(models))
	for i := range models {
		result[i] = models[i].ToDomain()
	}

	return result, nil
}

func (r *ProductStockRepository) DeleteProductStock(id string) *domain.Error {

	result := r.db.Delete(&ProductStockModel{}, "id = ?", id)
	if result.Error != nil {
		return r.dbErrMapper.MapErrorToDomain(result.Error, "failed to delete product")
	}

	if result.RowsAffected == 0 {
		return domain.NewError("product not found", domain.ErrNotFound)
	}

	return nil
}
