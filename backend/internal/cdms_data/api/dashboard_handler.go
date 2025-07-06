// dashboard_handler.go
package api

import (
	"log/slog"
	"net/http"
	"strconv" // Make sure this is imported
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	queries db.Querier
	logger  *slog.Logger
}

func NewDashboardHandler(q db.Querier, logger *slog.Logger) *DashboardHandler {
	return &DashboardHandler{
		queries: q,
		logger:  logger.With("component", "dashboard_handler"),
	}
}

// TimeWindowStats defines stats for a specific period (e.g., 7 days).
type TimeWindowStats struct {
	NewItemsCount         int64   `json:"new_items_count"`
	NewItemsValue         string  `json:"new_items_value"`
	AvgDaysToPFS          float64 `json:"avg_days_to_pfs"`
	AvgDaysForPFSComplete float64 `json:"avg_days_for_pfs_complete"`
	PassedToPFS           int64   `json:"passed_to_pfs"`
	CompletedByPFS        int64   `json:"completed_by_pfs"`
}

// DashboardStats now holds all the combined dashboard statistics for both chargebacks and delinquencies.
type DashboardStats struct {
	ChargebackStatusSummary []db.GetChargebackStatusSummaryRow            `json:"chargeback_status_summary"`
	ChargebackTimeWindows   map[string]TimeWindowStats                    `json:"chargeback_time_windows"`
	NonipacStatusSummary    []db.GetNonipacStatusSummaryRow               `json:"nonipac_status_summary"`
	NonipacAgingSchedule    []db.GetNonipacAgingScheduleByBusinessLineRow `json:"nonipac_aging_schedule"`
}

// HandleGetDashboardStats handles GET /api/dashboard/stats
// This replaces HandleGetChargebackStats and now fetches all dashboard data.
func (h *DashboardHandler) HandleGetDashboardStats(c echo.Context) error {
	ctx := c.Request().Context()

	// --- Fetch Chargeback Data ---
	chargebackStatusSummary, err := h.queries.GetChargebackStatusSummary(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get chargeback status summary for dashboard", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve chargeback status summary")
	}

	chargebackTimeWindows := make(map[string]TimeWindowStats)
	windows := map[string]int{"7d": 7, "14d": 14, "21d": 21, "28d": 28}
	now := time.Now()

	for key, days := range windows {
		endDate := now.AddDate(0, 0, -(days - 7))
		startDate := now.AddDate(0, 0, -days)

		pgStartTimestamp := pgtype.Timestamptz{Time: startDate, Valid: true}
		pgEndTimestamp := pgtype.Timestamptz{Time: endDate, Valid: true}
		pgStartDate := pgtype.Date{Time: startDate, Valid: true}
		pgEndDate := pgtype.Date{Time: endDate, Valid: true}

		newStats, err := h.queries.GetNewChargebackStatsForWindow(ctx, db.GetNewChargebackStatsForWindowParams{CreatedAt: pgStartTimestamp, CreatedAt_2: pgEndTimestamp})
		if err != nil {
			h.logger.ErrorContext(ctx, "Failed to get new chargebacks stats for dashboard", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve new chargebacks stats")
		}

		// Fixed: GetAverageDaysToPFSForWindow now returns a string
		avgDaysToPFSStr, err := h.queries.GetAverageDaysToPFSForWindow(ctx, db.GetAverageDaysToPFSForWindowParams{PassedToPfsDate: pgStartDate, PassedToPfsDate_2: pgEndDate})
		if err != nil {
			h.logger.ErrorContext(ctx, "Failed to get average days to PFS for dashboard", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve average days to PFS")
		}
		avgDaysToPFSFloat, parseErr := strconv.ParseFloat(avgDaysToPFSStr, 64) // Parse string to float64
		if parseErr != nil {
			h.logger.ErrorContext(ctx, "Failed to parse avgDaysToPFS string to float", "value", avgDaysToPFSStr, "error", parseErr)
			avgDaysToPFSFloat = 0.0 // Default to 0 on parse error
		}

		// Fixed: GetAverageDaysForPFSCompletionForWindow now returns a string
		avgDaysForPFSCompleteStr, err := h.queries.GetAverageDaysForPFSCompletionForWindow(ctx, db.GetAverageDaysForPFSCompletionForWindowParams{PfsCompletionDate: pgStartDate, PfsCompletionDate_2: pgEndDate})
		if err != nil {
			h.logger.ErrorContext(ctx, "Failed to get average days for PFS completion for dashboard", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve average days for PFS completion")
		}
		avgDaysForPFSCompleteFloat, parseErr := strconv.ParseFloat(avgDaysForPFSCompleteStr, 64) // Parse string to float64
		if parseErr != nil {
			h.logger.ErrorContext(ctx, "Failed to parse avgDaysForPFSComplete string to float", "value", avgDaysForPFSCompleteStr, "error", parseErr)
			avgDaysForPFSCompleteFloat = 0.0 // Default to 0 on parse error
		}

		pfsCounts, err := h.queries.GetPFSCountsForWindow(ctx, db.GetPFSCountsForWindowParams{PassedToPfsDate: pgStartDate, PassedToPfsDate_2: pgEndDate})
		if err != nil {
			h.logger.ErrorContext(ctx, "Failed to get PFS status counts for dashboard", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve PFS status counts")
		}

		newItemsValueStr := newStats.NewChargebacksValue // Already a string from SQL cast

		chargebackTimeWindows[key] = TimeWindowStats{
			NewItemsCount:         newStats.NewChargebacksCount,
			NewItemsValue:         newItemsValueStr,
			AvgDaysToPFS:          avgDaysToPFSFloat,
			AvgDaysForPFSComplete: avgDaysForPFSCompleteFloat,
			PassedToPFS:           pfsCounts.PassedToPfsCount,
			CompletedByPFS:        pfsCounts.CompletedByPfsCount,
		}
	}

	// --- Fetch Non-IPAC (Delinquency) Data ---
	nonipacStatusSummary, err := h.queries.GetNonipacStatusSummary(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get non-ipac status summary for dashboard", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve non-ipac status summary")
	}

	nonipacAgingSchedule, err := h.queries.GetNonipacAgingScheduleByBusinessLine(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get non-ipac aging schedule for dashboard", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve non-ipac aging schedule")
	}

	// --- Combine and Return Data ---
	stats := DashboardStats{
		ChargebackStatusSummary: chargebackStatusSummary,
		ChargebackTimeWindows:   chargebackTimeWindows,
		NonipacStatusSummary:    nonipacStatusSummary,
		NonipacAgingSchedule:    nonipacAgingSchedule,
	}

	return c.JSON(http.StatusOK, stats)
}
