PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE employee (
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
INSERT INTO employee VALUES(1,'5972276','Slid','Kestrel','slid.kestrel@outlook.com','0450579387',1,'2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO employee VALUES(2,'5972279','Daddy','Doss','daddy.doss@outlook.com','0450579387',0,'2025-03-23 06:53:39','2025-03-26 05:54:43');
CREATE TABLE employee_auth (
    auth_id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);
INSERT INTO employee_auth VALUES(1,'5972276','skestrel','$2y$14$aiyzqIjN/Dyyuie6.mccdu8OC3GYB7XEPCdSU/P.UTlrRwR9ktIjq','2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO employee_auth VALUES(2,'5972279','ddoss','$2y$14$aiyzqIjN/Dyyuie6.mccdu8OC3GYB7XEPCdSU/P.UTlrRwR9ktIjq','2025-03-23 06:53:39','2025-03-26 05:54:43');
CREATE TABLE leave_request (
    request_id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id TEXT NOT NULL,
    leave_type TEXT NOT NULL,
    from_date TEXT NOT NULL,
    to_date TEXT NOT NULL,
    note TEXT,
    hours_per_day INTEGER NOT NULL DEFAULT 8,
    is_multi_day INTEGER,
    is_pending INTEGER NOT NULL DEFAULT 1,
    is_approved INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);
INSERT INTO leave_request VALUES(8,'5972276','annual','2025-03-25','2025-03-25','',4,0,0,1);
INSERT INTO leave_request VALUES(9,'5972276','annual','2025-03-27','2025-03-29','',8,1,0,1);
INSERT INTO leave_request VALUES(11,'5972276','annual','2025-04-05','2025-04-05','',8,0,1,0);
CREATE TABLE job (
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
INSERT INTO job VALUES(1,'Natts House',1,'Trewheela Ave','Manifold Heights','3218','Geelong',0,'2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO job VALUES(2,'saddsf',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:20:53','2025-03-26 04:20:53');
INSERT INTO job VALUES(3,'sdfadsfsd',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:20:57','2025-03-26 04:20:57');
INSERT INTO job VALUES(4,'sdfadsfsd',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:20:58','2025-03-26 04:20:58');
INSERT INTO job VALUES(5,'sdfadsfsd',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:20:59','2025-03-26 04:20:59');
INSERT INTO job VALUES(6,'sadfasdf',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:02','2025-03-26 04:21:02');
INSERT INTO job VALUES(7,'sadfasdf',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:03','2025-03-26 04:21:03');
INSERT INTO job VALUES(8,'sadfasdf',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:03','2025-03-26 04:21:03');
INSERT INTO job VALUES(9,'sadfasdf',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:03','2025-03-26 04:21:03');
INSERT INTO job VALUES(10,'sadfasdf',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:03','2025-03-26 04:21:03');
INSERT INTO job VALUES(11,'sadfasdf',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:03','2025-03-26 04:21:03');
INSERT INTO job VALUES(12,'dsafds',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:30','2025-03-26 04:21:30');
INSERT INTO job VALUES(13,'dsafds',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:31','2025-03-26 04:21:31');
INSERT INTO job VALUES(14,'dsafds',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:31','2025-03-26 04:21:31');
INSERT INTO job VALUES(15,'dsafds',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:31','2025-03-26 04:21:31');
INSERT INTO job VALUES(16,'dsafds',0,'n/a','n/a','n/a','n/a',0,'2025-03-26 04:21:31','2025-03-26 04:21:31');
CREATE TABLE timesheet_week (
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
INSERT INTO timesheet_week VALUES(1,'5972276',12,1,2,3,4,5,6,7,'2025-3-26','2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet_week VALUES(2,'5972276',1,8,9,10,11,12,13,14,'2025-3-26','2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet_week VALUES(3,'5972276',11,15,16,17,18,19,20,21,'2025-3-26','2025-03-26 05:12:18','2025-03-26 05:12:18');
CREATE TABLE timesheet (
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
INSERT INTO timesheet VALUES(1,1,'2025-3-26','wed',0,0,'2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet VALUES(2,1,'2025-3-27','thu',0,0,'2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet VALUES(3,1,'2025-3-28','fri',0,0,'2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet VALUES(4,1,'2025-3-29','sat',0,0,'2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet VALUES(5,1,'2025-3-30','sun',0,0,'2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet VALUES(6,1,'2025-3-31','mon',0,0,'2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet VALUES(7,1,'2025-4-1','tue',0,0,'2025-03-26 05:09:21','2025-03-26 05:09:21');
INSERT INTO timesheet VALUES(8,2,'2025-3-26','wed',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(9,2,'2025-3-27','thu',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(10,2,'2025-3-28','fri',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(11,2,'2025-3-29','sat',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(12,2,'2025-3-30','sun',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(13,2,'2025-3-31','mon',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(14,2,'2025-4-1','tue',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(15,3,'2025-3-26','wed',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(16,3,'2025-3-27','thu',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(17,3,'2025-3-28','fri',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(18,3,'2025-3-29','sat',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(19,3,'2025-3-30','sun',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(20,3,'2025-3-31','mon',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(21,3,'2025-4-1','tue',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('employee',2);
INSERT INTO sqlite_sequence VALUES('employee_auth',2);
INSERT INTO sqlite_sequence VALUES('leave_request',11);
INSERT INTO sqlite_sequence VALUES('job',16);
INSERT INTO sqlite_sequence VALUES('timesheet_week',3);
INSERT INTO sqlite_sequence VALUES('timesheet',21);

CREATE TABLE incident_report (
    uuid TEXT PRIMARY KEY,
    reporter_id TEXT NOT NULL,
    full_name TEXT NOT NULL,
    home_address TEXT NOT NULL,
    contact_number TEXT NOT NULL,
    incident_date TEXT NOT NULL,
    incident_time TEXT NOT NULL,
    police_notified TEXT NOT NULL,
    incident_location TEXT NOT NULL,
    incident_description TEXT NOT NULL,
    was_witnessed TEXT NOT NULL,
    was_injured TEXT NOT NULL,
    further_details TEXT NOT NULL,
    was_treated TEXT NOT NULL,
    treatment_location TEXT NOT NULL,
    inc_info_date_1 TEXT NOT NULL,
    inc_info_date_2 TEXT NOT NULL,
    inc_info_date_3 TEXT NOT NULL,
    inc_info_date_4 TEXT NOT NULL,
    inc_info_date_5 TEXT NOT NULL,
    inc_info_action_1 TEXT NOT NULL,
    inc_info_action_2 TEXT NOT NULL,
    inc_info_action_3 TEXT NOT NULL,
    inc_info_action_4 TEXT NOT NULL,
    inc_info_action_5 TEXT NOT NULL,
    inc_info_name_1 TEXT NOT NULL,
    inc_info_name_2 TEXT NOT NULL,
    inc_info_name_3 TEXT NOT NULL,
    inc_info_name_4 TEXT NOT NULL,
    inc_info_name_5 TEXT NOT NULL,
    reporter TEXT NOT NULL,
    signature TEXT NOT NULL,
    report_date TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reporter_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);

COMMIT;
