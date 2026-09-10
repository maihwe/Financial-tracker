package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"financial-tracker/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AdminHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		if !RequireAdmin(pool, w, r) {
			return
		}

		switch r.URL.Path {

		case "/admin":

			if r.Method != http.MethodGet {
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			json.NewEncoder(w).Encode(
				map[string]string{
					"message": "Welcome to the admin area",
				},
			)

		case "/admin/users":

			if r.Method != http.MethodGet {
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			users, err := storage.ListUsersInDB(pool)

			if err != nil {
				http.Error(
					w,
					"Failed to retrieve users",
					http.StatusInternalServerError,
				)
				return
			}

			json.NewEncoder(w).Encode(users)

		case "/admin/statistics":

			if r.Method != http.MethodGet {
				http.Error(
					w,
					"Method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			statistics, err :=
				storage.GetAdminStatisticsFromDB(pool)

			if err != nil {
				http.Error(
					w,
					"Failed to retrieve statistics",
					http.StatusInternalServerError,
				)
				return
			}

			json.NewEncoder(w).Encode(statistics)

		default:

			// Role management endpoint:
			// PUT /admin/users/:id/role
			if strings.HasPrefix(
				r.URL.Path,
				"/admin/users/",
			) &&
				strings.HasSuffix(
					r.URL.Path,
					"/role",
				) {

				if r.Method != http.MethodPut {
					http.Error(
						w,
						"Method not allowed",
						http.StatusMethodNotAllowed,
					)
					return
				}

				path := strings.TrimPrefix(
					r.URL.Path,
					"/admin/users/",
				)

				path = strings.TrimSuffix(
					path,
					"/role",
				)

				userID, err :=
					strconv.Atoi(path)

				if err != nil || userID <= 0 {
					http.Error(
						w,
						"Invalid user ID",
						http.StatusBadRequest,
					)
					return
				}

				currentUserID, err :=
					GetAuthenticatedUserID(r)

				if err != nil {
					http.Error(
						w,
						"Authentication required",
						http.StatusUnauthorized,
					)
					return
				}

				if userID == currentUserID {
					http.Error(
						w,
						"You cannot change your own role",
						http.StatusForbidden,
					)
					return
				}

				var request struct {
					Role string `json:"role"`
				}

				err = json.NewDecoder(r.Body).Decode(
					&request,
				)

				if err != nil {
					http.Error(
						w,
						"Invalid request body",
						http.StatusBadRequest,
					)
					return
				}

				if request.Role != "user" &&
					request.Role != "admin" {

					http.Error(
						w,
						"Role must be either user or admin",
						http.StatusBadRequest,
					)
					return
				}

				user, err :=
					storage.UpdateUserRoleInDB(
						pool,
						userID,
						request.Role,
					)

				if err != nil {

					if err.Error() == "user not found" {
						http.Error(
							w,
							"User not found",
							http.StatusNotFound,
						)
						return
					}

					http.Error(
						w,
						"Failed to update user role",
						http.StatusInternalServerError,
					)
					return
				}

				json.NewEncoder(w).Encode(user)

				return
			}

			http.NotFound(w, r)
		}
	}
}
