DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS employee_auth;
DROP TABLE IF EXISTS leave_request;

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

CREATE TABLE IF NOT EXISTS leave_request (
    request_id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id TEXT NOT NULL,
    leave_type TEXT NOT NULL,
    from_date TEXT NOT NULL,
    to_date TEXT NOT NULL,
    note TEXT,
    is_approved INTEGER DEFAULT 0,
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);

/** ---- employee inserts ---- **/
INSERT INTO employee (employee_id, first_name, last_name, email, phone_number, is_admin)
VALUES ('5972276', 'Slid', 'Kestrel', 'slid.kestrel@outlook.com', '0450579387', 1);
INSERT INTO employee (employee_id, first_name, last_name, email, phone_number, is_admin)
VALUES ('5972277', 'Daddy', 'Doss', 'daddy.doss@outlook.com', '0450579387', 0);

/** ---- employee_auth inserts ---- **/
INSERT INTO employee_auth (employee_id, username, password_hash)
VALUES ('5972276', 'skestrel', '$2y$14$aiyzqIjN/Dyyuie6.mccdu8OC3GYB7XEPCdSU/P.UTlrRwR9ktIjq');
INSERT INTO employee_auth (employee_id, username, password_hash)
VALUES ('5972277', 'ddoss', '$2y$14$aiyzqIjN/Dyyuie6.mccdu8OC3GYB7XEPCdSU/P.UTlrRwR9ktIjq');

/** ---- leave_request inserts ---- **/
INSERT INTO leave_request (employee_id, leave_type, from_date, to_date, note)
VALUES ('5972276', 'annual', '18/12/2024', '20/12/2024', 'gone fishin');