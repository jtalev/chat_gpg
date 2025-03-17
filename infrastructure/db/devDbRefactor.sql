DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS employee_auth;
DROP TABLE IF EXISTS leave_request;
DROP TABLE IF EXISTS job;
DROP TABLE IF EXISTS timesheet;
DROP TABLE IF EXISTS timesheet_week;

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
    is_pending INTEGER DEFAULT 1,
    is_approved INTEGER DEFAULT 0,
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS job (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    number INTEGER NOT NULL,
    address TEXT NOT NULL,
    suburb TEXT NOT NULL,
    post_code TEXT NOT NULL,
    city TEXT NOT NULL,
    is_complete INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS timesheet_week (
    timesheet_week_id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    wed_timesheet_id INTEGER,
    thu_timesheet_id INTEGER,
    fri_timesheet_id INTEGER,
    sat_timesheet_id INTEGER,
    sun_timesheet_id INTEGER,
    mon_timesheet_id INTEGER,
    tue_timesheet_id INTEGER,
    week_start_date TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS timesheet (
    timesheet_id INTEGER PRIMARY KEY AUTOINCREMENT,
    timesheet_week_id INTEGER NOT NULL,
    timesheet_date TEXT NOT NULL,
    day TEXT NOT NULL,
    hours INTEGER NOT NULL,
    minutes INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (timesheet_week_id) REFERENCES timesheet_week(timesheet_week_id) ON DELETE CASCADE
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
VALUES ('5972276', 'annual', '2024-12-18', '2024-12-20', 'gone fishin');
INSERT INTO leave_request (employee_id, leave_type, from_date, to_date, note, is_pending, is_approved)
VALUES ('5972276', 'sick', '2024-12-18', '2024-12-20', 'gone fishin', 0, 0);
INSERT INTO leave_request (employee_id, leave_type, from_date, to_date, note, is_pending, is_approved)
VALUES ('5972276', 'annual', '2024-12-18', '2024-12-20', 'gone fishin', 0, 1);

/** ---- job inserts ---- **/
INSERT INTO job (name, number, address, suburb, post_code, city)
VALUES ('Natts House', 1, 'Trewheela Ave', 'Manifold Heights', '3218', 'Geelong');
INSERT INTO job (name, number, address, suburb, post_code, city)
VALUES ('Matts House', 1, 'Trewheela Ave', 'Manifold Heights', '3218', 'Geelong');

/** ---- timesheet_week inserts ---- **/
INSERT INTO timesheet_week 
(employee_id, job_id, wed_timesheet_id, thu_timesheet_id, fri_timesheet_id, sat_timesheet_id, sun_timesheet_id, mon_timesheet_id, tue_timesheet_id, week_start_date)
VALUES ('5972276', 2, 1, 2, 3, 4, 5, 6, 7, '2025-1-1');

/** ---- timesheet inserts ---- **/
INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
VALUES (1, '2025-01-01', 'wed', 4, 30);
INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
VALUES (1, '2025-01-02', 'thu', 5, 30);
INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
VALUES (1, '2025-01-03', 'fri', 1, 30);
INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
VALUES (1, '2025-01-04', 'sat', 0, 0);
INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
VALUES (1, '2025-01-05', 'sun', 0, 0);
INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
VALUES (1, '2025-01-06', 'mon', 8, 30);
INSERT INTO timesheet (timesheet_week_id, timesheet_date, day, hours, minutes)
VALUES (1, '2025-01-07', 'tue', 8, 0);