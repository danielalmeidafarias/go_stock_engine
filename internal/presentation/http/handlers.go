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

type errorResponse struct {
	Error string `json:"error" example:"error message"`
}

type createResponse struct {
	ID string `json:"id" example:"uuid"`
}

// productStockResponse represents a product stock item.
type productStockResponse struct {
	ID                string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name              string  `json:"name" example:"Engine Oil Filter"`
	Category          string  `json:"category" example:"engine"`
	CurrentStock      int     `json:"current_stock" example:"150"`
	MinimumStock      int     `json:"minimum_stock" example:"50"`
	AverageDailySales int     `json:"average_daily_sales" example:"10"`
	LeadTimeDays      int     `json:"lead_time_days" example:"7"`
	UnitCost          float64 `json:"unit_cost" example:"25.50"`
	CriticalityLevel  int     `json:"criticality_level" example:"3"`
}

// restockPriorityResponse represents a product restock priority.
type restockPriorityResponse struct {
	ExpectedConsumption int                  `json:"expected_consumption" example:"70"`
	ProjectedStock      int                  `json:"projected_stock" example:"-20"`
	IsRepositionNeeded  bool                 `json:"is_reposition_needed" example:"true"`
	UrgencyScore        int                  `json:"urgency_score" example:"210"`
	ProductStock        productStockResponse `json:"product_stock"`
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

// Create godoc
// @Summary      Create a product stock
// @Description  Creates a new product stock entry
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        request  body      createProductStockRequest  true  "Product stock data"
// @Success      201      {object}  createResponse
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /stock [post]
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

// GetAll godoc
// @Summary      List all product stocks
// @Description  Returns a paginated list of all product stocks
// @Tags         stock
// @Produce      json
// @Param        page   query     int  false  "Page number"   default(1)
// @Param        limit  query     int  false  "Items per page" default(20)
// @Success      200    {array}   productStockResponse
// @Failure      500    {object}  errorResponse
// @Router       /stock [get]
func (h *ProductStockHandler) GetAll(c *gin.Context) {
	pagination := parsePagination(c)

	products, domainErr := h.getAllUC.Execute(pagination)
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetOne godoc
// @Summary      Get a product stock by ID
// @Description  Returns a single product stock by its ID
// @Tags         stock
// @Produce      json
// @Param        id   path      string  true  "Product stock ID"
// @Success      200  {object}  productStockResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /stock/{id} [get]
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

// Update godoc
// @Summary      Update a product stock
// @Description  Partially updates a product stock by its ID
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "Product stock ID"
// @Param        request  body      updateProductStockRequest  true  "Fields to update"
// @Success      204      "No Content"
// @Failure      400      {object}  errorResponse
// @Failure      404      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /stock/{id} [put]
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

// Delete godoc
// @Summary      Delete a product stock
// @Description  Deletes a product stock by its ID
// @Tags         stock
// @Produce      json
// @Param        id   path      string  true  "Product stock ID"
// @Success      204  "No Content"
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /stock/{id} [delete]
func (h *ProductStockHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	domainErr := h.deleteUC.Execute(id)
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetByCategory godoc
// @Summary      Get product stocks by category
// @Description  Returns a paginated list of product stocks filtered by category
// @Tags         stock
// @Produce      json
// @Param        category  path      string  true   "Product category"
// @Param        page      query     int     false  "Page number"    default(1)
// @Param        limit     query     int     false  "Items per page" default(20)
// @Success      200       {array}   productStockResponse
// @Failure      400       {object}  errorResponse
// @Failure      500       {object}  errorResponse
// @Router       /stock/category/{category} [get]
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

// GetRestockPriorities godoc
// @Summary      Get restock priorities
// @Description  Returns a paginated list of products that need restocking, sorted by urgency
// @Tags         restock
// @Produce      json
// @Param        page   query     int  false  "Page number"    default(1)
// @Param        limit  query     int  false  "Items per page" default(20)
// @Success      200    {array}   restockPriorityResponse
// @Failure      500    {object}  errorResponse
// @Router       /restock/priorities [get]
func (h *ProductStockHandler) GetRestockPriorities(c *gin.Context) {
	pagination := parsePagination(c)

	priorities, domainErr := h.getPriorityUC.Execute(pagination)
	if domainErr != nil {
		c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), gin.H{"error": domainErr.Message})
		return
	}

	c.JSON(http.StatusOK, priorities)
}
