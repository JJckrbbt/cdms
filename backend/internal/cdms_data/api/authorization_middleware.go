package api

import (
	"github.com/jjckrbbt/cdms/backend/internal/db"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user_context").(db.GetUserWithAuthorizationContextRow)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "User context not available")
			}

			permissions, ok := user.Permissions.([]string)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions (could not read permissions)")
			}

			for _, p := range permissions {
				if p == permission {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions")
		}
	}
}

func (h *UserHandler) LoadUserContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(db.CdmsUser)
		if !ok {
			return echo.NewHTTPError(http.StatusInternalServerError, "User not found in context")
		}

		fullUserContext, err := h.queries.GetUserWithAuthorizationContext(c.Request().Context(), user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to load user permissions")
		}

		c.Set("user_context", fullUserContext)
		return next(c)
	}
}
