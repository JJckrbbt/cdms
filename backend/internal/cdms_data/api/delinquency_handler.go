package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)


type UserUpdateDelinquencyRequest struct {
	CurrentStatus *string `json:"current_status"`
}

type PFSUpdateDelinquencyRequest struct {
	CurrentStatus *string `json:"current_status"`
}

type AdminUpdateDelinquencyRequest struct {
	CurrentStatus *string `json:"current_status"`
}

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

type PaginatedDelinquenciesReponse struct {
	TotalCount int64                           `json:"total_count"`
	Data       []db.ListActiveDelinquenciesRow `json:"data"`
}

type CreateDelinquencyRequest struct {
	BusinessLine                string          `json:"business_line"`
	BilledTotalAmount           decimal.Decimal `json:"billed_total_amount"`
	PrincipleAmount             decimal.Decimal `json:"principle_amount"`
	InterestAmount              decimal.Decimal `json:"interest_amount"`
	PenaltyAmount               decimal.Decimal `json:"penalty_amount"`
	AdministrationChargesAmount decimal.Decimal `json:"administration_charges_amount"`
	DebitOutstandingAmount      decimal.Decimal `json:"debit_outstanding_amount"`
	CreditTotalAmount           decimal.Decimal `json:"credit_total_amount"`
	CreditOutstandingAmount     decimal.Decimal `json:"credit_outstanding_amount"`
	DocumentDate                string          `json:"document_date"`
	AddressCode                 string          `json:"address_code"`
	Vendor                      string          `json:"vendor"`
	DebtAppealForbearance       bool            `json:"debt_appeal_forbearance"`
	Statement                   string          `json:"statement"`
	DocumentNumber              string          `json:"document_number"`
	VendorCode                  string          `json:"vendor_code"`
	CollectionDueDate           string          `json:"collection_due_date"`
	OpenDate                    string          `json:"open_date"`
	Title                       *string         `json:"title"`
	CurrentStatus               *string         `json:"current_status"`
}

func (h *DelinquencyHandler) HandleCreate(c echo.Context) error {
	var req CreateDelinquencyRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body: "+err.Error())
	}

	docDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid document_date format, expected YYYY-MM-DD")
	}
	collectionDueDate, err := time.Parse("2006-01-02", req.CollectionDueDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid collection_due_date format, expected YYYY-MM-DD")
	}
	openDate, err := time.Parse("2006-01-02", req.OpenDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid open_date format, expected YYYY-MM-DD")
	}

	params := db.CreateDelinquencyParams{
		BusinessLine:                db.ChargebackBusinessLine(req.BusinessLine),
		BilledTotalAmount:           pgtype.Numeric{Int: req.BilledTotalAmount.BigInt(), Valid: true},
		PrincipleAmount:             pgtype.Numeric{Int: req.PrincipleAmount.BigInt(), Valid: true},
		InterestAmount:              pgtype.Numeric{Int: req.InterestAmount.BigInt(), Valid: true},
		PenaltyAmount:               pgtype.Numeric{Int: req.PenaltyAmount.BigInt(), Valid: true},
		AdministrationChargesAmount: pgtype.Numeric{Int: req.AdministrationChargesAmount.BigInt(), Valid: true},
		DebitOutstandingAmount:      pgtype.Numeric{Int: req.DebitOutstandingAmount.BigInt(), Valid: true},
		CreditTotalAmount:           pgtype.Numeric{Int: req.CreditTotalAmount.BigInt(), Valid: true},
		CreditOutstandingAmount:     pgtype.Numeric{Int: req.CreditOutstandingAmount.BigInt(), Valid: true},
		DocumentDate:                pgtype.Date{Time: docDate, Valid: true},
		AddressCode:                 req.AddressCode,
		Vendor:                      req.Vendor,
		DebtAppealForbearance:       req.DebtAppealForbearance,
		Statement:                   req.Statement,
		DocumentNumber:              req.DocumentNumber,
		VendorCode:                  req.VendorCode,
		CollectionDueDate:           pgtype.Date{Time: collectionDueDate, Valid: true},
		OpenDate:                    pgtype.Date{Time: openDate, Valid: true},
		Title:                       pgtype.Text{String: derefString(req.Title), Valid: req.Title != nil},
	}

	delinquency, err := h.queries.CreateDelinquency(c.Request().Context(), params)
	if err != nil {
		h.logger.ErrorContext(c.Request().Context(), "Failed to create delinquency in database", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create delinquency")
	}

	return c.JSON(http.StatusCreated, delinquency)
}

func (h *DelinquencyHandler) HandleGetDelinquencies(c echo.Context) error {
	ctx := c.Request().Context()

	documentNumber := c.QueryParam("documentNumber")

	if documentNumber != "" {
		h.logger.InfoContext(ctx, "Performing lookup by business key")

		delinquency, err := h.queries.GetActiveDelinquencyByBusinessKey(ctx, documentNumber)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, "Delinquency not found for given business key")
			}
			h.logger.ErrorContext(ctx, "Failedto get delinuquency by business key", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve delinquency")
		}

		return c.JSON(http.StatusOK, delinquency)
	}

	h.logger.InfoContext(ctx, "Performing paginated list lookup for delinquencies")

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
	var totalCount int64
	if len(delinquencies) > 0 {
		totalCount = delinquencies[0].TotalCount
	}
	response := PaginatedDelinquenciesReponse{
		TotalCount: totalCount,
		Data:       delinquencies,
	}

	return c.JSON(http.StatusOK, response)
}

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

func (h *DelinquencyHandler) HandleUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = "user" // Default to least privileged
	}

	existing, err := h.queries.GetDelinquencyForUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Delinquency not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve delinquency for update")
	}

	var updatedDelinquency db.Nonipac
	var updateErr error

	switch userRole {
	case "admin":
		var req AdminUpdateDelinquencyRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		params := db.AdminUpdateDelinquencyParams{
			ID:            id,
			CurrentStatus: existing.CurrentStatus,
		}
		if req.CurrentStatus != nil {
			params.CurrentStatus = db.CdmsStatus(*req.CurrentStatus)
		}
		updatedDelinquency, updateErr = h.queries.AdminUpdateDelinquency(ctx, params)

	case "pfs":
		var req PFSUpdateDelinquencyRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		params := db.PFSUpdateDelinquencyParams{
			ID:            id,
			CurrentStatus: existing.CurrentStatus,
		}
		if req.CurrentStatus != nil {
			params.CurrentStatus = db.CdmsStatus(*req.CurrentStatus)
		}
		updatedDelinquency, updateErr = h.queries.PFSUpdateDelinquency(ctx, params)

	case "user":
		fallthrough
	default:
		var req UserUpdateDelinquencyRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		params := db.UserUpdateDelinquencyParams{
			ID:            id,
			CurrentStatus: existing.CurrentStatus,
		}
		if req.CurrentStatus != nil {
			params.CurrentStatus = db.CdmsStatus(*req.CurrentStatus)
		}
		updatedDelinquency, updateErr = h.queries.UserUpdateDelinquency(ctx, params)
	}

	if updateErr != nil {
		h.logger.ErrorContext(ctx, "Failed to update delinquency", "error", updateErr, "id", id, "role", userRole)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update delinquency")
	}

	return c.JSON(http.StatusOK, updatedDelinquency)
}

func (h *DelinquencyHandler) HandleDelinquencyStatus(c echo.Context) error {
	ctx := c.Request().Context()

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		h.logger.WarnContext(ctx, "Invalid Delinquency ID format provided", "id_param", idParam, "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	statusHistory, err := h.queries.GetStatusHistoryForDelinquencies(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.logger.InfoContext(ctx, "No status history found for Delinquency ID", "delinquency_id", id)
			return echo.NewHTTPError(http.StatusNotFound, "Status history not found for the given ID")
		}
		h.logger.ErrorContext(ctx, "Failed to get status history for delinquency", "error", err, "delinquency_id", id)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve status history")
	}

	return c.JSON(http.StatusOK, statusHistory)
}
