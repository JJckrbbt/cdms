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

// --- Request Structs for PATCH operations ---

// UserUpdateChargebackRequest defines fields a regular user can modify.
type UserUpdateChargebackRequest struct {
	CurrentStatus          *string `json:"current_status"`
	ReasonCode             *string `json:"reason_code"`
	Action                 *string `json:"action"`
	IssueInResearchDate    *string `json:"issue_in_research_date"` // YYYY-MM-DD
	ALCToRebill            *string `json:"alc_to_rebill"`
	TASToRebill            *string `json:"tas_to_rebill"`
	LineOfAccountingRebill *string `json:"line_of_accounting_rebill"`
	SpecialInstruction     *string `json:"special_instruction"`
	PassedToPSF            *string `json:"passed_to_psf"` // YYYY-MM-DD
}

// PFSUpdateChargebackRequest defines fields a PFS user can modify.
type PFSUpdateChargebackRequest struct {
	CurrentStatus      *string `json:"current_status"`
	PassedToPSF        *string `json:"passed_to_psf"` // YYYY-MM-DD
	NewIPACDocumentRef *string `json:"new_ipac_document_ref"`
	PFSCompletionDate  *string `json:"pfs_completion_date"` // YYYY-MM-DD
}

// AdminUpdateChargebackRequest defines fields an Admin can modify.
type AdminUpdateChargebackRequest struct {
	CurrentStatus          *string `json:"current_status"`
	ReasonCode             *string `json:"reason_code"`
	Action                 *string `json:"action"`
	IssueInResearchDate    *string `json:"issue_in_research_date"`
	ALCToRebill            *string `json:"alc_to_rebill"`
	TASToRebill            *string `json:"tas_to_rebill"`
	LineOfAccountingRebill *string `json:"line_of_accounting_rebill"`
	SpecialInstruction     *string `json:"special_instruction"`
	PassedToPSF            *string `json:"passed_to_psf"`
	PFSCompletionDate      *string `json:"pfs_completion_date"`
}

// Create Chargeback Struct defines fields an Admin uses to manually create Chargeback
type CreateChargebackRequest struct {
	Fund              string          `json:"fund"`
	BusinessLine      string          `json:"business_line"`
	Region            int16           `json:"region"`
	Program           string          `json:"program"`
	ALNum             int16           `json:"al_num"`
	SourceNum         string          `json:"source_num"`
	ALC               string          `json:"alc"`
	CustomerTAS       string          `json:"customer_tas"`
	TaskSubtask       string          `json:"task_subtask"`
	CustomerName      string          `json:"customer_name"`
	OrgCode           string          `json:"org_code"`
	DocumentDate      string          `json:"document_date"` // Expect dates as "YYYY-MM-DD" strings
	AccompDate        string          `json:"accomp_date"`
	ChargebackAmount  decimal.Decimal `json:"chargeback_amount"`
	Statement         string          `json:"statement"`
	BDDocNum          string          `json:"bd_doc_num"`
	Vendor            string          `json:"vendor"`
	LocationSystem    *string         `json:"location_system"`
	AgreementNum      *string         `json:"agreement_num"`
	Title             *string         `json:"title"`
	ClassID           *string         `json:"class_id"`
	AssignedRebillDRN *string         `json:"assigned_rebill_drn"`
	ArticlesServices  *string         `json:"articles_services"`
	CurrentStatus     *string         `json:"current_status"`
	ReasonCode        *string         `json:"reason_code"`
	Action            *string         `json:"action"`
}

type PaginatedChargebacksResponse struct {
	TotalCount int64                         `json:"total_count"`
	Data       []db.ListActiveChargebacksRow `json:"data"`
}

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

// HandleCreate handles POST /api/chargebacks
func (h *ChargebackHandler) HandleCreate(c echo.Context) error {
	var req CreateChargebackRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body "+err.Error())
	}

	docDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid document-date format, expected YYY-MM-DD")
	}
	accompDate, err := time.Parse("2006-01-02", req.AccompDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid accomp_date format, expected YYYY-MM-DD")
	}

	// Params for DB Call
	params := db.CreateChargebackParams{
		Fund:              db.ChargebackFund(req.Fund),
		BusinessLine:      db.ChargebackBusinessLine(req.BusinessLine),
		Region:            req.Region,
		Program:           req.Program,
		AlNum:             req.ALNum,
		SourceNum:         req.SourceNum,
		Alc:               req.ALC,
		CustomerTas:       req.CustomerTAS,
		TaskSubtask:       req.TaskSubtask,
		CustomerName:      req.CustomerName,
		OrgCode:           req.OrgCode,
		DocumentDate:      pgtype.Date{Time: docDate, Valid: true},
		AccompDate:        pgtype.Date{Time: accompDate, Valid: true},
		ChargebackAmount:  pgtype.Numeric{Int: req.ChargebackAmount.BigInt(), Valid: true},
		Statement:         req.Statement,
		BdDocNum:          req.BDDocNum,
		Vendor:            req.Vendor,
		LocationSystem:    pgtype.Text{String: derefString(req.LocationSystem), Valid: req.LocationSystem != nil},
		AgreementNum:      pgtype.Text{String: derefString(req.AgreementNum), Valid: req.AgreementNum != nil},
		Title:             pgtype.Text{String: derefString(req.Title), Valid: req.Title != nil},
		ClassID:           pgtype.Text{String: derefString(req.ClassID), Valid: req.ClassID != nil},
		AssignedRebillDrn: pgtype.Text{String: derefString(req.AssignedRebillDRN), Valid: req.AssignedRebillDRN != nil},
		ArticlesServices:  pgtype.Text{String: derefString(req.ArticlesServices), Valid: req.ArticlesServices != nil},
		CurrentStatus:     db.CdmsStatus(derefStringWithDefault(req.CurrentStatus, "Open")),
		ReasonCode:        db.NullChargebackReasonCode{ChargebackReasonCode: db.ChargebackReasonCode(derefString(req.ReasonCode)), Valid: req.ReasonCode != nil},
		Action:            db.NullChargebackAction{ChargebackAction: db.ChargebackAction(derefString(req.Action)), Valid: req.Action != nil},
	}

	chargeback, err := h.queries.CreateChargeback(c.Request().Context(), params)
	if err != nil {
		h.logger.ErrorContext(c.Request().Context(), "Failed to create chargeback in database", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create chargeback")
	}
	return c.JSON(http.StatusCreated, chargeback)
}

// HandleList handles GET /api/chargebacks
func (h *ChargebackHandler) HandleGetChargebacks(c echo.Context) error {
	ctx := c.Request().Context()

	// Check for Business Key Lookup query parameters
	bdDocNum := c.QueryParam("bd_doc_num")
	alNumStr := c.QueryParam("al_num")

	if bdDocNum != "" && alNumStr != "" {
		// --- Logic for getting a single item by business key ---
		h.logger.InfoContext(ctx, "Performing lookup by business key", "bd_doc_num", bdDocNum, "al_num", alNumStr)

		alNum, err := strconv.ParseInt(alNumStr, 10, 16)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid al_num format, must be a number.")
		}

		params := db.GetActiveChargebackByBusinessKeyParams{
			BdDocNum: bdDocNum,
			AlNum:    int16(alNum),
		}

		chargeback, err := h.queries.GetActiveChargebackByBusinessKey(ctx, params)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return echo.NewHTTPError(http.StatusNotFound, "Chargeback not found for the given business key")
			}
			h.logger.ErrorContext(ctx, "Failed to get chargeback by business key", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve chargeback")
		}

		return c.JSON(http.StatusOK, chargeback)
	}

	// --- Fallback to Pagination Logic for getting a list ---
	h.logger.InfoContext(ctx, "Performing paginated list lookup for chargebacks")

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = 50
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	listParams := db.ListActiveChargebacksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	chargebacks, err := h.queries.ListActiveChargebacks(ctx, listParams)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to list active chargebacks", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve chargebacks")
	}
	var totalCount int64
	if len(chargebacks) > 0 {
		totalCount = chargebacks[0].TotalCount
	}
	response := PaginatedChargebacksResponse{
		TotalCount: totalCount,
		Data:       chargebacks,
	}

	return c.JSON(http.StatusOK, response)
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

// HandleUpdate now contains logic to differentiate based on user role.
func (h *ChargebackHandler) HandleUpdate(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	// This role would come from validated JWT middleware in a real application
	userRole := c.Request().Header.Get("X-User-Role")
	if userRole == "" {
		userRole = "user" // Default to least privileged
	}

	// Fetch existing record to merge changes onto
	existing, err := h.queries.GetChargebackForUpdate(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "Chargeback not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve chargeback for update")
	}

	var updatedChargeback db.Chargeback
	var updateErr error

	switch userRole {
	case "admin":
		var req AdminUpdateChargebackRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		params := buildAdminUpdateParams(&req, &existing)
		updatedChargeback, updateErr = h.queries.AdminUpdateChargeback(ctx, params)

	case "pfs":
		var req PFSUpdateChargebackRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		params := buildPFSUpdateParams(&req, &existing)
		updatedChargeback, updateErr = h.queries.PFSUpdateChargeback(ctx, params)

	case "user":
		fallthrough
	default:
		var req UserUpdateChargebackRequest
		if err := c.Bind(&req); err != nil {
			return err
		}
		params := buildUserUpdateParams(&req, &existing)
		updatedChargeback, updateErr = h.queries.UserUpdateChargeback(ctx, params)
	}

	if updateErr != nil {
		h.logger.ErrorContext(ctx, "Failed to update chargeback", "error", updateErr, "id", id, "role", userRole)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update chargeback")
	}

	return c.JSON(http.StatusOK, updatedChargeback)
}
