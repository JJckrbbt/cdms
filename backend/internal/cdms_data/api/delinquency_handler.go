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

type DelinquencyHandler struct {
	queries db.Querier
	logger  *slog.Logger
}

func NewDelinquencyHandler(q db.Querier, logger *slog.Logger) *DelinquencyHandler {
	return &DelinquencyHandler{
		queries: q,
		logger:  logger.With("component", "delinquency_handler"),
	}
}

func (h *DelinquencyHandler) HandleList(c echo.Context) error {
	// ... (pagination logic is the same)
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	params := db.ListActiveDelinquenciesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	delinquencies, err := h.queries.ListActiveDelinquencies(c.Request().Context(), params)
	if err != nil {
		h.logger.ErrorContext(c.Request().Context(), "Failed to list active delinquencies", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve delinquencies")
	}

	return c.JSON(http.StatusOK, delinquencies)
}

// NEW: Add the handler for getting a single delinquency by ID.
func (h *DelinquencyHandler) HandleGetByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format, ID must be a number.")
	}

	delinquency, err := h.queries.GetActiveDelinquencyByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Delinquency not found")
		}
		h.logger.ErrorContext(c.Request().Context(), "Failed to get delinquency by ID", "error", err, "id", idParam)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve delinquency")
	}

	return c.JSON(http.StatusOK, delinquency)
}
