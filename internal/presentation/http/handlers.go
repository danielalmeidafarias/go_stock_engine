package http

import (
	"net/http"
	"strconv"

	usecases "github.com/danielalmeidafarias/go_stock_engine/internal/application"
	"github.com/danielalmeidafarias/go_stock_engine/internal/domain"
	"github.com/gin-gonic/gin"
)

type ProductStockHandler struct {
	createUC        *usecases.CreateProductStockUseCase
	getAllUC        *usecases.GetAllProductStockUseCase
	getOneUC        *usecases.GetOneProductStockUseCase
	updateUC        *usecases.UpdateProductStockUseCase
	deleteUC        *usecases.DeleteProductStockUseCase
	getByCategoryUC *usecases.GetByCategoryProductStockUseCase
	getPriorityUC   *usecases.GetProductPriorityUseCase
}

func NewProductStockHandler(
	createUC *usecases.CreateProductStockUseCase,
	getAllUC *usecases.GetAllProductStockUseCase,
	getOneUC *usecases.GetOneProductStockUseCase,
	updateUC *usecases.UpdateProductStockUseCase,
	deleteUC *usecases.DeleteProductStockUseCase,
	getByCategoryUC *usecases.GetByCategoryProductStockUseCase,
	getPriorityUC *usecases.GetProductPriorityUseCase,
) *ProductStockHandler {
	return &ProductStockHandler{
		createUC:        createUC,
		getAllUC:        getAllUC,
		getOneUC:        getOneUC,
		updateUC:        updateUC,
		deleteUC:        deleteUC,
		getByCategoryUC: getByCategoryUC,
		getPriorityUC:   getPriorityUC,
	}
}

func mapErrorToHTTPStatus(errCode domain.ErrorCode) int {
	switch errCode {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrConflict:
		return http.StatusConflict
	case domain.ErrBadRequest:
		return http.StatusBadRequest
	case domain.ErrInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func parsePagination(c *gin.Context) domain.Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	return domain.Pagination{
		Page:  page,
		Limit: limit,
	}
}

type createProductStockRequest struct {
	Name              string  `json:"name" binding:"required"`
	Category          string  `json:"category" binding:"required"`
	CurrentStock      int     `json:"current_stock"`
	MinimumStock      int     `json:"minimum_stock"`
	AverageDailySales int     `json:"average_daily_sales"`
	LeadTimeDays      int     `json:"lead_time_days"`
	UnitCost          float64 `json:"unit_cost" binding:"required"`
	CriticalityLevel  int     `json:"criticality_level" binding:"required"`
}

func (h *ProductStockHandler) Create(c *gin.Context) {
	var req createProductStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, domainErr := h.createUC.Execute(usecases.CreateProductStockDTO{
		Name:              req.Name,
		Category:          req.Category,
		CurrentStock:      req.CurrentStock,
		MinimumStock:      req.MinimumStock,
		AverageDailySales: req.AverageDailySales,
		LeadTimeDays:      req.LeadTimeDays,
		UnitCost:          req.UnitCost,
		CriticalityLevel:  req.CriticalityLevel,
	})
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *ProductStockHandler) GetAll(c *gin.Context) {
	pagination := parsePagination(c)

	products, domainErr := h.getAllUC.Execute(pagination)
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductStockHandler) GetOne(c *gin.Context) {
	id := c.Param("id")

	product, domainErr := h.getOneUC.Execute(id)
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusOK, product)
}

type updateProductStockRequest struct {
	CurrentStock      *int     `json:"current_stock"`
	MinimumStock      *int     `json:"minimum_stock"`
	AverageDailySales *int     `json:"average_daily_sales"`
	LeadTimeDays      *int     `json:"lead_time_days"`
	UnitCost          *float64 `json:"unit_cost"`
	CriticalityLevel  *int     `json:"criticality_level"`
}

func (h *ProductStockHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req updateProductStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domainErr := h.updateUC.Execute(usecases.UpdateProductStockDTO{
		ID:                id,
		CurrentStock:      req.CurrentStock,
		MinimumStock:      req.MinimumStock,
		AverageDailySales: req.AverageDailySales,
		LeadTimeDays:      req.LeadTimeDays,
		UnitCost:          req.UnitCost,
		CriticalityLevel:  req.CriticalityLevel,
	})
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ProductStockHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	domainErr := h.deleteUC.Execute(id)
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *ProductStockHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	pagination := parsePagination(c)

	products, domainErr := h.getByCategoryUC.Execute(usecases.GetByCategoryDTO{
		Category:   category,
		Pagination: pagination,
	})
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductStockHandler) GetRestockPriorities(c *gin.Context) {
	pagination := parsePagination(c)

	priorities, domainErr := h.getPriorityUC.Execute(pagination)
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusOK, priorities)
}
