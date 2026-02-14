package go_orm

import (
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain/entities"
	"github.com/danielalmeidafarias/go_stock_engine/internal/infraestructure/repository/gorm/db"
)

type ProductStockRepository struct {
	gorm db.GoORM
}

func NewProductStockRepository(gorm db.GoORM) *ProductStockRepository {
	return &ProductStockRepository{gorm: gorm}
}

func (r *ProductStockRepository) Create(in *entities.ProductStock) (string, *domain.Error) {
	model := MapProductStockToModel(in)
	db := r.gorm.GetORM()

	if err := db.Create(model).Error; err != nil {
		return "", r.gorm.MapGormError(err, "failed to create product")
	}

	return model.ID, nil
}

func (r *ProductStockRepository) Update(in *entities.ProductStock) *domain.Error {
	model := MapProductStockToModel(in)
	db := r.gorm.GetORM()

	result := db.Save(model)
	if result.Error != nil {
		return r.gorm.MapGormError(result.Error, "failed to update product")
	}

	if result.RowsAffected == 0 {
		return domain.NewError("product not found", domain.ErrNotFound)
	}

	return nil
}

func (r *ProductStockRepository) GetAll(pagination *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	var models []ProductStockModel
	db := r.gorm.GetORM()

	query := db.Model(&ProductStockModel{})

	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, r.gorm.MapGormError(err, "failed to list products")
	}

	result := make([]*entities.ProductStock, len(models))
	for i := range models {
		result[i] = models[i].ToDomain()
	}

	return result, nil
}

func (r *ProductStockRepository) GetOneByID(id string) (*entities.ProductStock, *domain.Error) {
	var model ProductStockModel
	db := r.gorm.GetORM()

	if err := db.First(&model, "id = ?", id).Error; err != nil {
		return nil, r.gorm.MapGormError(err, "failed to get product")
	}

	return model.ToDomain(), nil
}

func (r *ProductStockRepository) GetByCategory(category entities.ProductCategory, pagination *domain.Pagination) ([]*entities.ProductStock, *domain.Error) {
	var models []ProductStockModel
	db := r.gorm.GetORM()

	query := db.Where("category = ?", string(category))

	if pagination != nil {
		offset := (pagination.Page - 1) * pagination.Limit
		query = query.Offset(offset).Limit(pagination.Limit)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, r.gorm.MapGormError(err, "failed to get products by category")
	}

	result := make([]*entities.ProductStock, len(models))
	for i := range models {
		result[i] = models[i].ToDomain()
	}

	return result, nil
}

func (r *ProductStockRepository) DeleteProductStock(id string) *domain.Error {
	db := r.gorm.GetORM()

	result := db.Delete(&ProductStockModel{}, "id = ?", id)
	if result.Error != nil {
		return r.gorm.MapGormError(result.Error, "failed to delete product")
	}

	if result.RowsAffected == 0 {
		return domain.NewError("product not found", domain.ErrNotFound)
	}

	return nil
}
