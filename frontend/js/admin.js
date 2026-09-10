// ================================
// ADMIN DASHBOARD AUTHORIZATION
// ================================

async function checkAdminAccess() {
    const response = await fetch("/admin");

    if (response.status === 200) {
        console.log("Admin access confirmed");
        return;
    }

    if (response.status === 401) {
        window.location.href = "/";
        return;
    }

    if (response.status === 403) {
        window.location.href = "/";
        return;
    }

    console.log(
        "Unexpected response:",
        response.status
    );
}

checkAdminAccess();