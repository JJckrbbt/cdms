package api

import (
	//"context"
	"log/slog"
	"net/http"
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

// CombinedChargebackStats holds all the new chargeback stats for the dashboard.
type CombinedChargebackStats struct {
	StatusSummary []db.GetChargebackStatusSummaryRow `json:"status_summary"`
	TimeWindows   map[string]TimeWindowStats         `json:"time_windows"`
}

// HandleGetChargebackStats handles GET /api/dashboard/chargeback-stats
func (h *DashboardHandler) HandleGetChargebackStats(c echo.Context) error {
	ctx := c.Request().Context()

	statusSummary, err := h.queries.GetChargebackStatusSummary(ctx)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get chargeback status summary", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve status summary")
	}

	timeWindows := make(map[string]TimeWindowStats)
	windows := map[string]int{"7d": 7, "14d": 14, "21d": 21, "28d": 28}
	now := time.Now()

	for key, days := range windows {
		// Correctly calculate non-overlapping 7-day windows
		endDate := now.AddDate(0, 0, -(days - 7))
		startDate := now.AddDate(0, 0, -days)

		// Use pgtype.Timestamptz for the created_at fields, and Date for the others
		pgStartTimestamp := pgtype.Timestamptz{Time: startDate, Valid: true}
		pgEndTimestamp := pgtype.Timestamptz{Time: endDate, Valid: true}
		pgStartDate := pgtype.Date{Time: startDate, Valid: true}
		pgEndDate := pgtype.Date{Time: endDate, Valid: true}

		// Use the correct generated struct and field names
		newStats, err := h.queries.GetNewChargebackStatsForWindow(ctx, db.GetNewChargebackStatsForWindowParams{CreatedAt: pgStartTimestamp, CreatedAt_2: pgEndTimestamp})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get new chargebacks stats")
		}

		avgDaysToPFS, err := h.queries.GetAverageDaysToPFSForWindow(ctx, db.GetAverageDaysToPFSForWindowParams{PassedToPfsDate: pgStartDate, PassedToPfsDate_2: pgEndDate})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get average days to PFS")
		}
		avgDaysToPFSFloat, _ := avgDaysToPFS.Float64Value()

		avgDaysForPFSComplete, err := h.queries.GetAverageDaysForPFSCompletionForWindow(ctx, db.GetAverageDaysForPFSCompletionForWindowParams{PfsCompletionDate: pgStartDate, PfsCompletionDate_2: pgEndDate})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get average days for PFS completion")
		}
		avgDaysForPFSCompleteFloat, _ := avgDaysForPFSComplete.Float64Value()

		pfsCounts, err := h.queries.GetPFSCountsForWindow(ctx, db.GetPFSCountsForWindowParams{PassedToPfsDate: pgStartDate, PassedToPfsDate_2: pgEndDate})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get PFS status counts")
		}

		newItemsValueStr, _ := newStats.NewChargebacksValue.Value()

		timeWindows[key] = TimeWindowStats{
			NewItemsCount:         newStats.NewChargebacksCount,
			NewItemsValue:         newItemsValueStr.(string),
			AvgDaysToPFS:          avgDaysToPFSFloat.Float64,
			AvgDaysForPFSComplete: avgDaysForPFSCompleteFloat.Float64,
			PassedToPFS:           pfsCounts.PassedToPfsCount,
			CompletedByPFS:        pfsCounts.CompletedByPfsCount,
		}
	}

	stats := CombinedChargebackStats{
		StatusSummary: statusSummary,
		TimeWindows:   timeWindows,
	}

	return c.JSON(http.StatusOK, stats)
}
