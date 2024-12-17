DROP TABLE employee;
DROP TABLE employee_auth;

CREATE TABLE IF NOT EXISTS employee (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id TEXT NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT,
    is_admin INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employee_auth (
    auth_id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);

INSERT INTO employee (employee_id, first_name, last_name, email, phone_number, is_admin)
VALUES ('5972276', 'Slid', 'Kestrel', 'slid.kestrel@outlook.com', '0450579387', 1);

INSERT INTO employee_auth (employee_id, username, password_hash)
VALUES ('5972276', 'skestrel', '$2y$14$aiyzqIjN/Dyyuie6.mccdu8OC3GYB7XEPCdSU/P.UTlrRwR9ktIjq');