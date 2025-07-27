package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

// --- Structs for API responses and request bodies ---

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Org       string `json:"org"`
	IsActive  bool   `json:"is_active"`
}

type UpdateUserRolesRequest struct {
	RoleIDs []int32 `json:"role_ids"`
}

type UpdateUserBusinessLinesRequest struct {
	BusinessLines []string `json:"business_lines"`
}

// UserResponse provides a consistent structure for user data sent to the frontend.
type UserResponse struct {
	ID        int64              `json:"id"`
	Email     string             `json:"email"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Org       db.UserOrg         `json:"org"`
	IsActive  bool               `json:"is_active"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	Roles     pq.StringArray     `json:"roles"`
}

// --- UserHandler and Constructor ---

type UserHandler struct {
	queries db.Querier
	logger  *slog.Logger
}

func NewUserHandler(q db.Querier, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		queries: q,
		logger:  logger.With("component", "user_handler"),
	}
}

// --- Handlers ---

// HandleGetMe handles /api/users/me, returns user details/permissions.
func (h *UserHandler) HandleGetMe(c echo.Context) error {
	userFromContext, ok := c.Get("user").(db.CdmsUser)
	if !ok {
		h.logger.ErrorContext(c.Request().Context(), "Could not retrieve user from context")
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not retrieve user from context")
	}

	fullUser, err := h.queries.GetUserWithAuthorizationContext(c.Request().Context(), userFromContext.ID)
	if err != nil {
		h.logger.ErrorContext(c.Request().Context(), "Failed to get user authorization context", "error", err, "user_id", userFromContext.ID)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve user permissions")
	}

	return c.JSON(http.StatusOK, fullUser)
}

// HandleListUsers is the handler for GET /api/admin/users.
func (h *UserHandler) HandleListUsers(c echo.Context) error {
	ctx := c.Request().Context()
	currentUser, ok := c.Get("user_context").(db.GetUserWithAuthorizationContextRow)
	if !ok {
		h.logger.ErrorContext(ctx, "User context not available in HandleListUsers")
		return echo.NewHTTPError(http.StatusInternalServerError, "User context not available")
	}

	canViewGlobal := hasPermission(currentUser.Permissions, "roles:assign_global")
	canViewScoped := hasPermission(currentUser.Permissions, "users:view_scoped")
	if !canViewGlobal && !canViewScoped {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to view users")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 50
	}
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	responseUsers := []UserResponse{}
	var err error

	if canViewGlobal {
		h.logger.InfoContext(ctx, "Fetching all users for global admin", "admin_id", currentUser.ID)
		users, queryErr := h.queries.ListAllUsers(ctx, db.ListAllUsersParams{
			Limit:  int32(limit),
			Offset: int32(offset),
		})
		err = queryErr
		for _, u := range users {
			responseUsers = append(responseUsers, UserResponse{
				ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName,
				Org: u.Org, IsActive: u.IsActive, CreatedAt: u.CreatedAt, Roles: u.Roles,
			})
		}
	} else {
		businessLinesTyped := make([]db.ChargebackBusinessLine, len(currentUser.BusinessLines))
		for i, s := range currentUser.BusinessLines {
			businessLinesTyped[i] = db.ChargebackBusinessLine(s)
		}

		h.logger.InfoContext(ctx, "Fetching scoped user list for business line admin", "admin_id", currentUser.ID, "scope", businessLinesTyped)
		users, queryErr := h.queries.ListUsersByBusinessLines(ctx, db.ListUsersByBusinessLinesParams{
			Limit:         int32(limit),
			Offset:        int32(offset),
			BusinessLines: businessLinesTyped,
		})
		err = queryErr
		for _, u := range users {
			responseUsers = append(responseUsers, UserResponse{
				ID: u.ID, Email: u.Email, FirstName: u.FirstName, LastName: u.LastName,
				Org: u.Org, IsActive: u.IsActive, CreatedAt: u.CreatedAt, Roles: u.Roles,
			})
		}
	}

	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to list users", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}

	return c.JSON(http.StatusOK, responseUsers)
}

// HandleUpdateUser updates a user's profile information.
func (h *UserHandler) HandleUpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID format")
	}

	currentUser, _ := c.Get("user_context").(db.GetUserWithAuthorizationContextRow)
	if !hasPermission(currentUser.Permissions, "users:edit") {
		return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to edit users")
	}

	if !hasPermission(currentUser.Permissions, "roles:assign_global") {
		targetUser, err := h.queries.GetUserWithAuthorizationContext(ctx, targetUserID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Target user not found")
		}
		if !isSubset(targetUser.BusinessLines, currentUser.BusinessLines) {
			return echo.NewHTTPError(http.StatusForbidden, "Insufficient scope to manage this user")
		}
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	params := db.UpdateUserParams{
		ID:        targetUserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Org:       db.UserOrg(req.Org),
		IsActive:  req.IsActive,
	}
	updatedUser, err := h.queries.UpdateUser(ctx, params)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to update user", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user")
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// --- Helper Functions ---

func hasPermission(permissions pq.StringArray, required string) bool {
	for _, p := range permissions {
		if p == required {
			return true
		}
	}
	return false
}

func isSubset(subset, superset pq.StringArray) bool {
	set := make(map[string]struct{}, len(superset))
	for _, s := range superset {
		set[s] = struct{}{}
	}
	for _, item := range subset {
		if _, found := set[item]; !found {
			return false
		}
	}
	return true
}
