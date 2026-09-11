package main

import (
	"fmt"
	"net/http"

	"financial-tracker/database"
	"financial-tracker/handlers"
)

func main() {

	pool, err := database.Connect()

	if err != nil {

		fmt.Println(
			"Database connection failed:",
			err,
		)

		return
	}

	defer pool.Close()

	transactionHandler :=
		handlers.TransactionHandler(pool)

	http.Handle(
		"/",
		http.FileServer(
			http.Dir("./frontend"),
		),
	)

	http.HandleFunc(
		"/transactions",
		transactionHandler,
	)

	http.HandleFunc(
		"/transactions/",
		transactionHandler,
	)

	registerHandler :=
		handlers.RegisterHandler(pool)

	http.HandleFunc(
		"/register",
		registerHandler,
	)

	loginHandler :=
		handlers.LoginHandler(pool)

	http.HandleFunc(
		"/login",
		loginHandler,
	)
	forgotPasswordHandler :=
		handlers.ForgotPasswordHandler(pool)

	http.HandleFunc(
		"/forgot-password",
		forgotPasswordHandler,
	)

	resetPasswordHandler :=
		handlers.ResetPasswordRouteHandler(pool)

	http.HandleFunc(
		"/reset-password",
		resetPasswordHandler,
	)

	http.HandleFunc(
		"/logout",
		handlers.LogoutHandler,
	)

	adminHandler :=
		handlers.AdminHandler(pool)

	http.HandleFunc(
		"/admin",
		adminHandler,
	)

	http.HandleFunc(
		"/admin/",
		adminHandler,
	)

	fmt.Println(
		"Server running on http://localhost:8080",
	)

	err = http.ListenAndServe(
		":8080",
		nil,
	)

	if err != nil {

		fmt.Println(
			"Server error:",
			err,
		)
	}
}
