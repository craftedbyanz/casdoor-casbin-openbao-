package handler

import (
	"net/http"
	"strconv"
	"time"

	"casdoor-casbin-openbao/internal/auth"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

type Order struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Total       float64   `json:"total"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
}

// Mock data
var orders = []Order{
	{
		ID:          "ord_001",
		UserID:      "admin",
		ProductName: "Laptop Pro",
		Quantity:    1,
		Price:       1299.99,
		Total:       1299.99,
		Status:      "delivered",
		CreatedAt:   time.Now().Add(-48 * time.Hour),
		CreatedBy:   "admin",
	},
	{
		ID:          "ord_002",
		UserID:      "hihi",
		ProductName: "Wireless Mouse",
		Quantity:    2,
		Price:       29.99,
		Total:       59.98,
		Status:      "processing",
		CreatedAt:   time.Now().Add(-6 * time.Hour),
		CreatedBy:   "hihi",
	},
	{
		ID:          "ord_003",
		UserID:      "testuser",
		ProductName: "USB Cable",
		Quantity:    3,
		Price:       9.99,
		Total:       29.97,
		Status:      "shipped",
		CreatedAt:   time.Now().Add(-12 * time.Hour),
		CreatedBy:   "testuser",
	},
}

// GetOrders returns all orders (admin only)
func (h *OrderHandler) GetOrders(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders":     orders,
		"total":      len(orders),
		"message":    "All orders retrieved",
		"accessed_by": user.Name,
	})
}

// GetMyOrders returns current user's orders only
func (h *OrderHandler) GetMyOrders(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var userOrders []Order
	for _, order := range orders {
		if order.UserID == user.Name {
			userOrders = append(userOrders, order)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders":  userOrders,
		"total":   len(userOrders),
		"message": "Your orders retrieved",
		"user":    user.Name,
	})
}

// GetOrder returns specific order by ID
func (h *OrderHandler) GetOrder(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	orderID := c.Param("id")
	
	for _, order := range orders {
		if order.ID == orderID {
			// Check ownership - user can only see their own orders unless admin
			if !user.IsAdmin && order.UserID != user.Name {
				return echo.NewHTTPError(http.StatusForbidden, "can only access your own orders")
			}
			
			return c.JSON(http.StatusOK, map[string]interface{}{
				"order":      order,
				"message":    "Order retrieved",
				"accessed_by": user.Name,
			})
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, "order not found")
}

// CreateOrder creates a new order
func (h *OrderHandler) CreateOrder(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	newOrder := Order{
		ID:          "ord_" + strconv.FormatInt(time.Now().Unix(), 10),
		UserID:      user.Name,
		ProductName: req.ProductName,
		Quantity:    req.Quantity,
		Price:       req.Price,
		Total:       req.Price * float64(req.Quantity),
		Status:      "pending",
		CreatedAt:   time.Now(),
		CreatedBy:   user.Name,
	}

	orders = append(orders, newOrder)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"order":      newOrder,
		"message":    "Order created successfully",
		"created_by": user.Name,
	})
}

// UpdateOrderStatus updates order status (admin only)
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	orderID := c.Param("id")
	
	var req struct {
		Status string `json:"status"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	for i, order := range orders {
		if order.ID == orderID {
			orders[i].Status = req.Status
			
			return c.JSON(http.StatusOK, map[string]interface{}{
				"order":      orders[i],
				"message":    "Order status updated",
				"updated_by": user.Name,
			})
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, "order not found")
}