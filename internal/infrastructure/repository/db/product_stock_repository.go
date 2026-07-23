package db

import (
	"context"

	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"gorm.io/gorm"
)

type ProductStockRepository struct {
	db *gorm.DB
}

func NewProductStockRepository(gorm *gorm.DB) *ProductStockRepository {
	return &ProductStockRepository{db: gorm}
}

func (r *ProductStockRepository) Create(ctx context.Context, in *domain.ProductStock) (string, *domain.Error) {
	model := MapProductStockToModel(in)

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return "", r.MapErrorToDomain(err, "failed to create product")
	}

	return model.ID, nil
}

func (r *ProductStockRepository) Update(ctx context.Context, in *domain.ProductStock) *domain.Error {
	model := MapProductStockToModel(in)

	result := r.db.WithContext(ctx).Save(model)
	if result.Error != nil {
		return r.MapErrorToDomain(result.Error, "failed to update product")
	}

	if result.RowsAffected == 0 {
		return domain.NewError("product not found", domain.ErrNotFound)
	}

	return nil
}

func (r *ProductStockRepository) GetAll(ctx context.Context, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
	var models []ProductStockModel

	query := r.db.WithContext(ctx).Model(&ProductStockModel{})

	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, r.MapErrorToDomain(err, "failed to list products")
	}

	result := make([]*domain.ProductStock, len(models))
	for i := range models {
		result[i] = models[i].ToDomain()
	}

	return result, nil
}

func (r *ProductStockRepository) GetOneByID(ctx context.Context, id string) (*domain.ProductStock, *domain.Error) {
	var model ProductStockModel

	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return nil, r.MapErrorToDomain(err, "failed to get product")
	}

	return model.ToDomain(), nil
}

func (r *ProductStockRepository) GetByCategory(ctx context.Context, category string, pagination *domain.Pagination) ([]*domain.ProductStock, *domain.Error) {
	var models []ProductStockModel

	query := r.db.WithContext(ctx).Where("category = ?", category)

	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, r.MapErrorToDomain(err, "failed to get products by category")
	}

	result := make([]*domain.ProductStock, len(models))
	for i := range models {
		result[i] = models[i].ToDomain()
	}

	return result, nil
}

func (r *ProductStockRepository) DeleteProductStock(ctx context.Context, id string) *domain.Error {

	result := r.db.WithContext(ctx).Delete(&ProductStockModel{}, "id = ?", id)
	if result.Error != nil {
		return r.MapErrorToDomain(result.Error, "failed to delete product")
	}

	if result.RowsAffected == 0 {
		return domain.NewError("product not found", domain.ErrNotFound)
	}

	return nil
}
