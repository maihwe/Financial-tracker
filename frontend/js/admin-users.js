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

            const roleCell =
                document.createElement("td");

            const roleSelect =
                document.createElement("select");

            roleSelect.innerHTML = `
                <option value="user">user</option>
                <option value="admin">admin</option>
            `;

            roleSelect.value = user.role;

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

            roleCell.appendChild(roleSelect);

            actionCell.appendChild(saveButton);

            row.innerHTML = `
                <td>${user.id}</td>
                <td>${user.name || "—"}</td>
                <td>${user.email}</td>
            `;

            row.appendChild(roleCell);

            const createdCell =
                document.createElement("td");

            createdCell.textContent =
            new Date(
                user.created_at
            ).toLocaleString();

            row.appendChild(createdCell);

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