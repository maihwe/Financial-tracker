// ================================
// ADMIN USER MANAGEMENT
// ================================

async function loadUsers() {

    const usersContainer =
        document.querySelector("#users-container");

    try {

        const response =
            await fetch("/admin/users", {
                cache: "no-store"
            });

        if (!response.ok) {

            usersContainer.innerHTML =
                "<p>Failed to load users.</p>";

            return;
        }

        const users =
            await response.json();

        if (users.length === 0) {

            usersContainer.innerHTML =
                "<p>No registered users found.</p>";

            return;
        }

        const table =
            document.createElement("table");

        table.innerHTML = `
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Role</th>
                    <th>Created</th>
                    <th>Action</th>
                </tr>
            </thead>

            <tbody></tbody>
        `;

        const tableBody =
            table.querySelector("tbody");

        users.forEach(function(user) {

            const row =
                document.createElement("tr");

            row.innerHTML = `
                <td>${user.id}</td>
                <td>${user.name || "—"}</td>
                <td>${user.email}</td>
            `;

            const roleCell =
                document.createElement("td");

            if (user.role === "super_admin") {

                roleCell.textContent =
                    "super_admin";

            } else {

                const roleSelect =
                    document.createElement("select");

                roleSelect.innerHTML = `
                    <option value="user">user</option>
                    <option value="admin">admin</option>
                `;

                roleSelect.value =
                    user.role;

                roleCell.appendChild(
                    roleSelect
                );

                const actionCell =
                    document.createElement("td");

                const saveButton =
                    document.createElement("button");

                saveButton.textContent =
                    "Save";

                saveButton.addEventListener(
                    "click",
                    function() {

                        updateUserRole(
                            user.id,
                            roleSelect.value
                        );
                    }
                );

                row.appendChild(roleCell);

                const createdCell =
                    document.createElement("td");

                createdCell.textContent =
                    new Date(
                        user.created_at
                    ).toLocaleString();

                row.appendChild(createdCell);

                actionCell.appendChild(
                    saveButton
                );

                row.appendChild(actionCell);

                tableBody.appendChild(row);

                return;
            }

            row.appendChild(roleCell);

            const createdCell =
                document.createElement("td");

            createdCell.textContent =
                new Date(
                    user.created_at
                ).toLocaleString();

            row.appendChild(createdCell);

            const actionCell =
                document.createElement("td");

            actionCell.textContent =
                "Protected";

            row.appendChild(actionCell);

            tableBody.appendChild(row);
        });

        usersContainer.innerHTML = "";

        usersContainer.appendChild(table);

    } catch (error) {

        console.error(
            "Error loading users:",
            error
        );

        usersContainer.innerHTML =
            "<p>Unable to load users.</p>";
    }
}


async function updateUserRole(
    userID,
    role
) {

    try {

        const response =
            await fetch(
                `/admin/users/${userID}/role`,
                {
                    method: "PUT",

                    headers: {
                        "Content-Type":
                            "application/json"
                    },

                    body: JSON.stringify({
                        role: role
                    })
                }
            );

        if (!response.ok) {

            const message =
                await response.text();

            alert(message);

            return;
        }

        alert(
            "User role updated successfully."
        );

        loadUsers();

    } catch (error) {

        console.error(
            "Error updating user role:",
            error
        );

        alert(
            "Unable to update user role."
        );
    }
}


loadUsers();