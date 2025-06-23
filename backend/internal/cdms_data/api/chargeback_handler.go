package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/labstack/echo/v4"
)

type ChargebackHandler struct {
	queries db.Querier
	logger  *slog.Logger
}

func NewChargebackHandler(q db.Querier, logger *slog.Logger) *ChargebackHandler {
	return &ChargebackHandler{
		queries: q,
		logger:  logger.With("component", "chargeback_handler"),
	}
}

// HandleList handles GET /api/chargebacks
func (h *ChargebackHandler) HandleList(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	params := db.ListActiveChargebacksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	chargebacks, err := h.queries.ListActiveChargebacks(c.Request().Context(), params)
	if err != nil {
		h.logger.ErrorContext(c.Request().Context(), "Failed to list active chargebacks", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve chargebacks")
	}

	return c.JSON(http.StatusOK, chargebacks)
}

// HandleGetByID handles GET /api/chargebacks/{id}
func (h *ChargebackHandler) HandleGetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format, ID must be a number.")
	}

	chargeback, err := h.queries.GetActiveChargebackByID(c.Request().Context(), id)
	if err != nil {
		// Check if the error is a "not found" error to return a 404
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Chargeback not found")
		}
		h.logger.ErrorContext(c.Request().Context(), "Failed to get chargeback by ID", "error", err, "id", idParam)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve chargeback")
	}

	return c.JSON(http.StatusOK, chargeback)
}
