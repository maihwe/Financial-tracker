// ================================
// AUTHENTICATION
// ================================


// Find the account sections.
const loginSection =
    document.querySelector("#login-section");

const registerSection =
    document.querySelector("#register-section");

const loggedInSection =
    document.querySelector("#logged-in-section");


// Find the account buttons.
const showRegisterButton =
    document.querySelector("#show-register-button");

const showLoginButton =
    document.querySelector("#show-login-button");

const logoutButton =
    document.querySelector("#logout-button");


// Find the login form.
const loginForm =
    document.querySelector("#login-form");


// Find login fields.
const loginEmail =
    document.querySelector("#login-email");

const loginPassword =
    document.querySelector("#login-password");

const toggleLoginPassword =
    document.querySelector("#toggle-login-password");


// Find login message area.
const loginMessage =
    document.querySelector("#login-message");


// Find registration form.
const registerForm =
    document.querySelector("#register-form");


// Find registration fields.
const registerName =
    document.querySelector("#register-name");

const registerEmail =
    document.querySelector("#register-email");

const registerPassword =
    document.querySelector("#register-password");

const toggleRegisterPassword =
    document.querySelector("#toggle-register-password");


// Find registration message area.
const registerMessage =
    document.querySelector("#register-message");


// Find logged-in user information.
const userEmail =
    document.querySelector("#user-email");


// Find the dashboard.
const dashboard =
    document.querySelector("#dashboard");


// Show the login section.
function showLogin() {

    loginSection.hidden = false;

    registerSection.hidden = true;

    loggedInSection.hidden = true;

    dashboard.hidden = true;
}


// Show the registration section.
function showRegister() {

    loginSection.hidden = true;

    registerSection.hidden = false;

    loggedInSection.hidden = true;

    dashboard.hidden = true;
}


// Show the dashboard.
function showDashboard() {

    loginSection.hidden = true;

    registerSection.hidden = true;

    loggedInSection.hidden = false;

    dashboard.hidden = false;
}


// ================================
// SWITCH TO REGISTER
// ================================

showRegisterButton.addEventListener(
    "click",
    function() {

        registerMessage.textContent = "";

        loginMessage.textContent = "";

        showRegister();
    }
);


// ================================
// SWITCH TO LOGIN
// ================================

showLoginButton.addEventListener(
    "click",
    function() {

        registerMessage.textContent = "";

        loginMessage.textContent = "";

        showLogin();
    }
);


// ================================
// SHOW / HIDE LOGIN PASSWORD
// ================================

toggleLoginPassword.addEventListener(
    "click",
    function() {

        if (
            loginPassword.type ===
            "password"
        ) {

            loginPassword.type =
                "text";

            toggleLoginPassword.textContent =
                "🙈";

            toggleLoginPassword.setAttribute(
                "aria-label",
                "Hide password"
            );

        } else {

            loginPassword.type =
                "password";

            toggleLoginPassword.textContent =
                "👁";

            toggleLoginPassword.setAttribute(
                "aria-label",
                "Show password"
            );
        }
    }
);


// ================================
// SHOW / HIDE REGISTER PASSWORD
// ================================

toggleRegisterPassword.addEventListener(
    "click",
    function() {

        if (
            registerPassword.type ===
            "password"
        ) {

            registerPassword.type =
                "text";

            toggleRegisterPassword.textContent =
                "🙈";

            toggleRegisterPassword.setAttribute(
                "aria-label",
                "Hide password"
            );

        } else {

            registerPassword.type =
                "password";

            toggleRegisterPassword.textContent =
                "👁";

            toggleRegisterPassword.setAttribute(
                "aria-label",
                "Show password"
            );
        }
    }
);


// ================================
// REGISTER
// ================================

registerForm.addEventListener(
    "submit",
    async function(event) {

        event.preventDefault();

        registerMessage.textContent = "";

        try {

            const response =
                await fetch(
                    "/register",
                    {
                        method: "POST",

                        credentials: "include",

                        headers: {
                            "Content-Type":
                                "application/json"
                        },

                        body: JSON.stringify({
                            name:
                                registerName.value,

                            email:
                                registerEmail.value,

                            password:
                                registerPassword.value
                        })
                    }
                );


            const data =
                await response.json();


            if (!response.ok) {

                registerMessage.textContent =
                    data;

                return;
            }


            // Save the registered email
            // before clearing the form.
            const registeredEmail =
                registerEmail.value;


            // Clear the registration form.
            registerForm.reset();


            // Put the registered email
            // into the login form.
            loginEmail.value =
                registeredEmail;


            // Move to the login screen.
            showLogin();


            // Show the success message
            // on the login screen.
            loginMessage.textContent =
                "Account created successfully. Please log in.";

        } catch (error) {

            console.error(
                "Registration error:",
                error
            );

            registerMessage.textContent =
                "Could not connect to the server.";
        }
    }
);


// ================================
// LOGIN
// ================================

loginForm.addEventListener(
    "submit",
    async function(event) {

        event.preventDefault();

        loginMessage.textContent = "";

        try {

            const response =
                await fetch(
                    "/login",
                    {
                        method: "POST",

                        credentials: "include",

                        headers: {
                            "Content-Type":
                                "application/json"
                        },

                        body: JSON.stringify({
                            email:
                                loginEmail.value,

                            password:
                                loginPassword.value
                        })
                    }
                );


            const data =
                await response.json();


            if (!response.ok) {

                loginMessage.textContent =
                    data;

                return;
            }
            // Display only the user's name.
            userEmail.textContent =
                data.name;


            // Send administrators to the admin dashboard.
            if (data.role === "admin") {

                 window.location.href =
                    "admin.html";

                return;
            }
            // Show normal user dashboard.
            showDashboard();


            // Clear login form.
            loginForm.reset();


            // Load this user's transactions.
            if (
                typeof loadTransactions ===
                "function"
            ) {

                loadTransactions();
            }

        } catch (error) {

            console.error(
                "Login error:",
                error
            );

            loginMessage.textContent =
                "Could not connect to the server.";
        }
    }
);


// ================================
// LOGOUT
// ================================

logoutButton.addEventListener(
    "click",
    async function() {

        try {

            const response =
                await fetch(
                    "/logout",
                    {
                        method: "POST",

                        credentials: "include"
                    }
                );


            if (!response.ok) {

                console.error(
                    "Logout failed."
                );

                return;
            }


            // Clear user information.
            userEmail.textContent = "";


            // Clear transactions.
            transactionList.innerHTML = "";


            // Clear financial summary.
            document.querySelector(
                "#total-income"
            ).textContent =
                formatCurrency(0);


            document.querySelector(
                "#total-expenses"
            ).textContent =
                formatCurrency(0);


            document.querySelector(
                "#balance"
            ).textContent =
                formatCurrency(0);


            // Return to login.
            showLogin();

        } catch (error) {

            console.error(
                "Logout error:",
                error
            );
        }
    }
);