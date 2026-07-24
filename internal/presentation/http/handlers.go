package http

import (
	"context"
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

func respondError(c *gin.Context, domainErr *domain.Error) {
	if respondTimeout(c) {
		return
	}

	c.JSON(mapErrorToHTTPStatus(domainErr.ErrCode), errorResponseDTO{Error: domainErr.Message})
}

func respondTimeout(c *gin.Context) bool {
	if c.Request.Context().Err() != context.DeadlineExceeded {
		return false
	}

	c.JSON(http.StatusGatewayTimeout, errorResponseDTO{Error: "request timed out"})
	return true
}

func parsePagination(c *gin.Context) domain.Pagination {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "0"))

	return domain.Pagination{
		Page:  page,
		Limit: limit,
	}
}

// Create godoc
// @Summary      Create a product stock
// @Description  Creates a new product stock entry
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        request  body      createProductStockRequestDTO  true  "Product stock data"
// @Success      201      {object}  createResponseDTO
// @Failure      400      {object}  errorResponseDTO
// @Failure      500      {object}  errorResponseDTO
// @Failure      504      {object}  errorResponseDTO
// @Router       /stock [post]
func (h *ProductStockHandler) Create(c *gin.Context) {
	var req createProductStockRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseDTO{Error: err.Error()})
		return
	}

	id, domainErr := h.createUC.Execute(c.Request.Context(), usecases.CreateProductStockDTO{
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
		respondError(c, domainErr)
		return
	}
	if respondTimeout(c) {
		return
	}

	c.JSON(http.StatusCreated, createResponseDTO{ID: id})
}

// GetAll godoc
// @Summary      List all product stocks
// @Description  Returns a paginated list of all product stocks
// @Tags         stock
// @Produce      json
// @Param        page   query     int  false  "Page number"   default(1)
// @Param        limit  query     int  false  "Items per page" default(20)
// @Success      200    {array}   productStockResponseDTO
// @Failure      500    {object}  errorResponseDTO
// @Failure      504    {object}  errorResponseDTO
// @Router       /stock [get]
func (h *ProductStockHandler) GetAll(c *gin.Context) {
	pagination := parsePagination(c)

	products, domainErr := h.getAllUC.Execute(c.Request.Context(), pagination)
	if domainErr != nil {
		respondError(c, domainErr)
		return
	}
	if respondTimeout(c) {
		return
	}

	c.JSON(http.StatusOK, toProductStockResponseDTOs(products))
}

// GetOne godoc
// @Summary      Get a product stock by ID
// @Description  Returns a single product stock by its ID
// @Tags         stock
// @Produce      json
// @Param        id   path      string  true  "Product stock ID"
// @Success      200  {object}  productStockResponseDTO
// @Failure      400  {object}  errorResponseDTO
// @Failure      404  {object}  errorResponseDTO
// @Failure      500  {object}  errorResponseDTO
// @Failure      504  {object}  errorResponseDTO
// @Router       /stock/{id} [get]
func (h *ProductStockHandler) GetOne(c *gin.Context) {
	id := c.Param("id")

	product, domainErr := h.getOneUC.Execute(c.Request.Context(), id)
	if domainErr != nil {
		respondError(c, domainErr)
		return
	}
	if respondTimeout(c) {
		return
	}

	c.JSON(http.StatusOK, toProductStockResponseDTO(product))
}

// Update godoc
// @Summary      Update a product stock
// @Description  Partially updates a product stock by its ID
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "Product stock ID"
// @Param        request  body      updateProductStockRequestDTO  true  "Fields to update"
// @Success      204      "No Content"
// @Failure      400      {object}  errorResponseDTO
// @Failure      404      {object}  errorResponseDTO
// @Failure      500      {object}  errorResponseDTO
// @Failure      504      {object}  errorResponseDTO
// @Router       /stock/{id} [put]
func (h *ProductStockHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req updateProductStockRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseDTO{Error: err.Error()})
		return
	}

	domainErr := h.updateUC.Execute(c.Request.Context(), usecases.UpdateProductStockDTO{
		ID:                id,
		CurrentStock:      req.CurrentStock,
		MinimumStock:      req.MinimumStock,
		AverageDailySales: req.AverageDailySales,
		LeadTimeDays:      req.LeadTimeDays,
		UnitCost:          req.UnitCost,
		CriticalityLevel:  req.CriticalityLevel,
	})
	if domainErr != nil {
		respondError(c, domainErr)
		return
	}
	if respondTimeout(c) {
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
// @Failure      400  {object}  errorResponseDTO
// @Failure      404  {object}  errorResponseDTO
// @Failure      500  {object}  errorResponseDTO
// @Failure      504  {object}  errorResponseDTO
// @Router       /stock/{id} [delete]
func (h *ProductStockHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	domainErr := h.deleteUC.Execute(c.Request.Context(), id)
	if domainErr != nil {
		respondError(c, domainErr)
		return
	}
	if respondTimeout(c) {
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
// @Success      200       {array}   productStockResponseDTO
// @Failure      400       {object}  errorResponseDTO
// @Failure      500       {object}  errorResponseDTO
// @Failure      504       {object}  errorResponseDTO
// @Router       /stock/category/{category} [get]
func (h *ProductStockHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	pagination := parsePagination(c)

	products, domainErr := h.getByCategoryUC.Execute(c.Request.Context(), usecases.GetByCategoryDTO{
		Category:   category,
		Pagination: pagination,
	})
	if domainErr != nil {
		respondError(c, domainErr)
		return
	}
	if respondTimeout(c) {
		return
	}

	c.JSON(http.StatusOK, toProductStockResponseDTOs(products))
}

// GetRestockPriorities godoc
// @Summary      Get restock priorities
// @Description  Returns a paginated list of products that need restocking, sorted by urgency
// @Tags         restock
// @Produce      json
// @Param        page   query     int  false  "Page number"    default(1)
// @Param        limit  query     int  false  "Items per page" default(20)
// @Success      200    {object}  restockPrioritiesResponseDTO
// @Failure      500    {object}  errorResponseDTO
// @Failure      504    {object}  errorResponseDTO
// @Router       /restock/priorities [get]
func (h *ProductStockHandler) GetRestockPriorities(c *gin.Context) {
	pagination := parsePagination(c)

	priorities, domainErr := h.getPriorityUC.Execute(c.Request.Context(), pagination)
	if domainErr != nil {
		respondError(c, domainErr)
		return
	}
	if respondTimeout(c) {
		return
	}

	c.JSON(http.StatusOK, toRestockPrioritiesResponseDTO(priorities))
}
