package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jjckrbbt/cdms/backend/internal/db" // Your sqlc package
	"github.com/jjckrbbt/cdms/backend/internal/logger"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of the db.Querier interface.
type MockQuerier struct {
	mock.Mock
}

// UPDATED: The mock method now returns a slice of the CORRECT struct type
func (m *MockQuerier) ListActiveChargebacks(ctx context.Context, params db.ListActiveChargebacksParams) ([]db.ActiveChargebacksWithVendorInfo, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]db.ActiveChargebacksWithVendorInfo), args.Error(1)
}

// UPDATED: The mock method now returns the CORRECT struct type
func (m *MockQuerier) GetActiveChargebackByID(ctx context.Context, id int64) (db.ActiveChargebacksWithVendorInfo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(db.ActiveChargebacksWithVendorInfo), args.Error(1)
}

// (You would add mock implementations for other Querier methods as you test them)

// --- List Chargebacks Test ---
func TestHandleListChargebacks(t *testing.T) {
	e := echo.New()
	logger.InitLogger("development")
	appLogger := logger.L()

	t.Run("Happy Path - successfully retrieves chargebacks", func(t *testing.T) {
		// --- Arrange ---
		mockQ := new(MockQuerier)

		// UPDATED: The dummy slice now uses the CORRECT struct type
		expectedChargebacks := []db.ActiveChargebacksWithVendorInfo{{}, {}}

		mockQ.On("ListActiveChargebacks", mock.Anything, db.ListActiveChargebacksParams{Limit: 10, Offset: 0}).
			Return(expectedChargebacks, nil).
			Once()

		handler := NewChargebackHandler(mockQ, appLogger)
		req := httptest.NewRequest(http.MethodGet, "/api/chargebacks?limit=10&page=1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// --- Act ---
		err := handler.HandleList(c)

		// --- Assert ---
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"id":0`)
		mockQ.AssertExpectations(t)
	})

	t.Run("Database error", func(t *testing.T) {
		// --- Arrange ---
		mockQ := new(MockQuerier)
		dbError := errors.New("a sudden database error")

		// UPDATED: The empty slice must also use the CORRECT struct type
		mockQ.On("ListActiveChargebacks", mock.Anything, mock.AnythingOfType("db.ListActiveChargebacksParams")).
			Return([]db.ActiveChargebacksWithVendorInfo{}, dbError).
			Once()

		handler := NewChargebackHandler(mockQ, appLogger)
		req := httptest.NewRequest(http.MethodGet, "/api/chargebacks", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// --- Act ---
		httpErr := handler.HandleList(c).(*echo.HTTPError)

		// --- Assert ---
		assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
		mockQ.AssertExpectations(t)
	})
}

// --- Get Chargeback By ID Test ---
func TestHandleGetChargebackByID(t *testing.T) {
	e := echo.New()
	logger.InitLogger("development")
	appLogger := logger.L()

	t.Run("Happy Path - successfully retrieves a single chargeback", func(t *testing.T) {
		// --- Arrange ---
		mockQ := new(MockQuerier)
		testID := int64(123)

		// UPDATED: The dummy struct now uses the CORRECT struct type
		expectedChargeback := db.ActiveChargebacksWithVendorInfo{ID: testID, Fund: "F-100"}

		mockQ.On("GetActiveChargebackByID", mock.Anything, testID).
			Return(expectedChargeback, nil).
			Once()

		handler := NewChargebackHandler(mockQ, appLogger)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", testID))

		// --- Act ---
		err := handler.HandleGetByID(c)

		// --- Assert ---
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"ID":123`)
		mockQ.AssertExpectations(t)
	})

	t.Run("Not Found error", func(t *testing.T) {
		// --- Arrange ---
		mockQ := new(MockQuerier)
		testID := int64(999)

		// UPDATED: The empty struct must also use the CORRECT struct type
		mockQ.On("GetActiveChargebackByID", mock.Anything, testID).
			Return(db.ActiveChargebacksWithVendorInfo{}, pgx.ErrNoRows).
			Once()

		handler := NewChargebackHandler(mockQ, appLogger)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(fmt.Sprintf("%d", testID))

		// --- Act ---
		httpErr := handler.HandleGetByID(c).(*echo.HTTPError)

		// --- Assert ---
		assert.Equal(t, http.StatusNotFound, httpErr.Code)
		mockQ.AssertExpectations(t)
	})

	// ... (The "Invalid ID format" test case remains the same and is correct) ...
}
