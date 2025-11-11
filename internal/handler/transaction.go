package handler

import (
	"net/http"
	"strconv"
	"time"

	"casdoor-casbin-openbao/internal/auth"
	"github.com/labstack/echo/v4"
)

type TransactionHandler struct{}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{}
}

type Transaction struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   string    `json:"created_by"`
}

// Mock data - In real app, this would be from database
var transactions = []Transaction{
	{
		ID:          "txn_001",
		UserID:      "admin",
		Amount:      1000.50,
		Type:        "deposit",
		Status:      "completed",
		Description: "Initial deposit",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		CreatedBy:   "admin",
	},
	{
		ID:          "txn_002", 
		UserID:      "hihi",
		Amount:      -250.00,
		Type:        "withdrawal",
		Status:      "pending",
		Description: "ATM withdrawal",
		CreatedAt:   time.Now().Add(-2 * time.Hour),
		CreatedBy:   "hihi",
	},
	{
		ID:          "txn_003",
		UserID:      "testuser",
		Amount:      500.75,
		Type:        "transfer",
		Status:      "completed", 
		Description: "Transfer from savings",
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		CreatedBy:   "testuser",
	},
}

// GetTransactions returns all transactions (admin only)
func (h *TransactionHandler) GetTransactions(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": transactions,
		"total":        len(transactions),
		"message":      "All transactions retrieved",
		"accessed_by":  user.Name,
	})
}

// GetMyTransactions returns current user's transactions only
func (h *TransactionHandler) GetMyTransactions(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var userTransactions []Transaction
	for _, txn := range transactions {
		if txn.UserID == user.Name {
			userTransactions = append(userTransactions, txn)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactions": userTransactions,
		"total":        len(userTransactions),
		"message":      "Your transactions retrieved",
		"user":         user.Name,
	})
}

// GetTransaction returns specific transaction by ID
func (h *TransactionHandler) GetTransaction(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	txnID := c.Param("id")
	
	for _, txn := range transactions {
		if txn.ID == txnID {
			// Check ownership - user can only see their own transactions unless admin
			if !user.IsAdmin && txn.UserID != user.Name {
				return echo.NewHTTPError(http.StatusForbidden, "can only access your own transactions")
			}
			
			return c.JSON(http.StatusOK, map[string]interface{}{
				"transaction": txn,
				"message":     "Transaction retrieved",
				"accessed_by": user.Name,
			})
		}
	}

	return echo.NewHTTPError(http.StatusNotFound, "transaction not found")
}

// CreateTransaction creates a new transaction
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	user, ok := auth.GetUserFromContext(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req struct {
		Amount      float64 `json:"amount"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	newTxn := Transaction{
		ID:          "txn_" + strconv.FormatInt(time.Now().Unix(), 10),
		UserID:      user.Name,
		Amount:      req.Amount,
		Type:        req.Type,
		Status:      "pending",
		Description: req.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   user.Name,
	}

	transactions = append(transactions, newTxn)

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"transaction": newTxn,
		"message":     "Transaction created successfully",
		"created_by":  user.Name,
	})
}