PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;
DROP TABLE IF EXISTS employee;
DROP TABLE IF EXISTS employee_auth;
DROP TABLE IF EXISTS employee_role;
DROP TABLE IF EXISTS leave_request;
DROP TABLE IF EXISTS job;
DROP TABLE IF EXISTS timesheet_week;
DROP TABLE IF EXISTS timesheet;
DROP TABLE IF EXISTS incident_report;
DROP TABLE IF EXISTS swms;
DROP TABLE IF EXISTS purchase_order;
DROP TABLE IF EXISTS purchase_order_item;
DROP TABLE IF EXISTS item_types;
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS tasks;
CREATE TABLE employee (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone_number TEXT,
    is_admin INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO employee VALUES(1,'5972276','Slid','Kestrel','slid.kestrel@outlook.com','0450579387',1,'2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO employee VALUES(4,'1234567','Daddy','Doss','daddy@outlook.com','0452312655',0,'2025-03-31 05:17:16','2025-03-31 05:17:16');
INSERT INTO employee VALUES(5,'7654321','Big','Hoss','daddy@outlook.com','0452312655',0,'2025-03-31 05:17:16','2025-03-31 05:17:16');
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
INSERT INTO employee_auth VALUES(4,'1234567','ddoss','$2a$14$HXnQD4WebLNJd/FKwlmdJePfL4inmdLrNs6ogVPcMxOW86tfsKPXS','2025-03-31 05:17:17','2025-03-31 05:17:17');
INSERT INTO employee_auth VALUES(5,'7654321','bhoss','$2a$14$HXnQD4WebLNJd/FKwlmdJePfL4inmdLrNs6ogVPcMxOW86tfsKPXS','2025-03-31 05:17:17','2025-03-31 05:17:17');
CREATE TABLE employee_role (
    uuid TEXT PRIMARY KEY,
    employee_id TEXT UNIQUE NOT NULL,
    role TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id) ON DELETE CASCADE
);
INSERT INTO employee_role VALUES('0cd836f9-394b-4968-b778-af1cbbf1d1bb','5972276', 'management','2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO employee_role VALUES('0cd836f9-394b-4968-b778-af1cbbf1d1bc','1234567', 'foreman','2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO employee_role VALUES('0cd836f9-394b-4968-b778-af1cbbf1d1bd','7654321', 'employee','2025-03-23 06:53:39','2025-03-23 06:53:39');
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
INSERT INTO timesheet_week VALUES(4,'5972276',1,22,23,24,25,26,27,28,'2025-4-2','2025-04-02 07:26:04','2025-04-02 07:26:04');
INSERT INTO timesheet_week VALUES(5,'1234567',1,29,30,31,32,33,34,35,'2025-4-2','2025-04-04 05:23:35','2025-04-04 05:23:35');
INSERT INTO timesheet_week VALUES(7,'1234567',1,43,44,45,46,47,48,49,'2025-3-26','2025-04-04 05:47:30','2025-04-04 05:47:30');
INSERT INTO timesheet_week VALUES(8,'5972276',1,50,51,52,53,54,55,56,'2025-4-9','2025-04-04 05:50:10','2025-04-04 05:50:10');
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
INSERT INTO timesheet VALUES(9,2,'2025-3-27','thu',2,0,'2025-03-26 05:09:25','2025-04-04 05:49:24');
INSERT INTO timesheet VALUES(10,2,'2025-3-28','fri',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(11,2,'2025-3-29','sat',33,0,'2025-03-26 05:09:25','2025-04-04 05:49:25');
INSERT INTO timesheet VALUES(12,2,'2025-3-30','sun',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(13,2,'2025-3-31','mon',1,0,'2025-03-26 05:09:25','2025-04-04 05:49:26');
INSERT INTO timesheet VALUES(14,2,'2025-4-1','tue',0,0,'2025-03-26 05:09:25','2025-03-26 05:09:25');
INSERT INTO timesheet VALUES(15,3,'2025-3-26','wed',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(16,3,'2025-3-27','thu',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(17,3,'2025-3-28','fri',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(18,3,'2025-3-29','sat',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(19,3,'2025-3-30','sun',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(20,3,'2025-3-31','mon',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(21,3,'2025-4-1','tue',0,0,'2025-03-26 05:12:18','2025-03-26 05:12:18');
INSERT INTO timesheet VALUES(22,4,'2025-4-2','wed',0,0,'2025-04-02 07:26:04','2025-04-02 07:26:04');
INSERT INTO timesheet VALUES(23,4,'2025-4-3','thu',12,0,'2025-04-02 07:26:04','2025-04-08 07:10:33');
INSERT INTO timesheet VALUES(24,4,'2025-4-4','fri',0,0,'2025-04-02 07:26:04','2025-04-02 07:26:04');
INSERT INTO timesheet VALUES(25,4,'2025-4-5','sat',0,0,'2025-04-02 07:26:04','2025-04-02 07:26:04');
INSERT INTO timesheet VALUES(26,4,'2025-4-6','sun',0,0,'2025-04-02 07:26:04','2025-04-02 07:26:04');
INSERT INTO timesheet VALUES(27,4,'2025-4-7','mon',0,0,'2025-04-02 07:26:04','2025-04-02 07:26:04');
INSERT INTO timesheet VALUES(28,4,'2025-4-8','tue',0,0,'2025-04-02 07:26:04','2025-04-02 07:26:04');
INSERT INTO timesheet VALUES(29,5,'2025-4-2','wed',0,0,'2025-04-04 05:23:35','2025-04-04 05:23:35');
INSERT INTO timesheet VALUES(30,5,'2025-4-3','thu',8,0,'2025-04-04 05:23:35','2025-04-04 05:23:38');
INSERT INTO timesheet VALUES(31,5,'2025-4-4','fri',3,0,'2025-04-04 05:23:35','2025-04-04 05:23:39');
INSERT INTO timesheet VALUES(32,5,'2025-4-5','sat',4,0,'2025-04-04 05:23:35','2025-04-04 05:23:40');
INSERT INTO timesheet VALUES(33,5,'2025-4-6','sun',0,0,'2025-04-04 05:23:35','2025-04-04 05:23:35');
INSERT INTO timesheet VALUES(34,5,'2025-4-7','mon',5,0,'2025-04-04 05:23:35','2025-04-04 05:23:41');
INSERT INTO timesheet VALUES(35,5,'2025-4-8','tue',0,0,'2025-04-04 05:23:35','2025-04-04 05:23:35');
INSERT INTO timesheet VALUES(43,7,'2025-3-26','wed',0,0,'2025-04-04 05:47:30','2025-04-04 05:47:30');
INSERT INTO timesheet VALUES(44,7,'2025-3-27','thu',2,0,'2025-04-04 05:47:30','2025-04-04 05:47:33');
INSERT INTO timesheet VALUES(45,7,'2025-3-28','fri',0,0,'2025-04-04 05:47:30','2025-04-04 05:47:30');
INSERT INTO timesheet VALUES(46,7,'2025-3-29','sat',1,0,'2025-04-04 05:47:30','2025-04-04 05:47:34');
INSERT INTO timesheet VALUES(47,7,'2025-3-30','sun',2,0,'2025-04-04 05:47:30','2025-04-04 05:47:34');
INSERT INTO timesheet VALUES(48,7,'2025-3-31','mon',0,0,'2025-04-04 05:47:30','2025-04-04 05:47:30');
INSERT INTO timesheet VALUES(49,7,'2025-4-1','tue',0,0,'2025-04-04 05:47:30','2025-04-04 05:47:30');
INSERT INTO timesheet VALUES(50,8,'2025-4-9','wed',0,0,'2025-04-04 05:50:10','2025-04-04 05:50:10');
INSERT INTO timesheet VALUES(51,8,'2025-4-10','thu',1,0,'2025-04-04 05:50:10','2025-04-04 05:50:13');
INSERT INTO timesheet VALUES(52,8,'2025-4-11','fri',2,0,'2025-04-04 05:50:10','2025-04-04 05:50:14');
INSERT INTO timesheet VALUES(53,8,'2025-4-12','sat',0,0,'2025-04-04 05:50:10','2025-04-04 05:50:10');
INSERT INTO timesheet VALUES(54,8,'2025-4-13','sun',1,0,'2025-04-04 05:50:10','2025-04-04 05:50:15');
INSERT INTO timesheet VALUES(55,8,'2025-4-14','mon',3,0,'2025-04-04 05:50:10','2025-04-04 05:50:17');
INSERT INTO timesheet VALUES(56,8,'2025-4-15','tue',0,0,'2025-04-04 05:50:10','2025-04-04 05:50:10');
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
INSERT INTO incident_report VALUES('0cd836f9-394b-4968-b778-af1cbbf1d1bb','5972276','working','csdv','kjbkjbkj','2025-04-18','saasd','no','sadfdasd','dasfsda','no','no','sdafds','no','emergency room','','','','','','this is an action','this is a new action','','','','lolololol','','','','','Slid Kestrel','asfsd','2025-04-24','2025-04-01 05:21:39','2025-04-01 07:20:24');
INSERT INTO incident_report VALUES('199538a9-6fb1-4e22-80fa-978d3265cb63','5972276','dsdsafds','sadf','sadfs','2025-04-10','asdfsa','yes','dsaf','sdfas','yes','yes','','yes','onsite','','','','','','','','','','','','','','','','Slid Kestrel','safsd','2025-04-17','2025-04-15 07:13:50','2025-04-15 07:13:50');

CREATE TABLE swms (
    uuid TEXT PRIMARY KEY,
    job_id INTEGER NOT NULL,
    project_activity TEXT NOT NULL,
    project_number TEXT NOT NULL,
    site_address TEXT NOT NULL,
    contact_name TEXT NOT NULL,
    contact_number TEXT NOT NULL,
    email_address TEXT NOT NULL,
    swms_date TEXT NOT NULL,
    high_risk_works INTEGER DEFAULT 0,
    safety_gloves INTEGER DEFAULT 0,
    safety_boots INTEGER DEFAULT 0,
    safety_glasses INTEGER DEFAULT 0,
    protective_clothing INTEGER DEFAULT 0,
    respiratory_protection INTEGER DEFAULT 0,
    hi_vis_clothing INTEGER DEFAULT 0,
    safety_helmet INTEGER DEFAULT 0,
    fall_arrest INTEGER DEFAULT 0,
    other_1 TEXT,
    other_2 TEXT,
    step_below_2m INTEGER DEFAULT 0,
    step_above_2m INTEGER DEFAULT 0,
    scaffold INTEGER DEFAULT 0,
    pressure_washer_diesel INTEGER DEFAULT 0,
    roof_anchor_points INTEGER DEFAULT 0,
    extension_ladder INTEGER DEFAULT 0,
    electric_scissor_lift INTEGER DEFAULT 0,
    diesel_scissor_lift INTEGER DEFAULT 0,
    electric_knuckle_boom INTEGER DEFAULT 0,
    diesel_knuckle_boom INTEGER DEFAULT 0,
    airless_spray_gun INTEGER DEFAULT 0,
    angle_grinder INTEGER DEFAULT 0,
    step_1 TEXT,
    hazards_1 TEXT,
    risks_1 TEXT,
    initial_risk_1 TEXT,
    control_measures_1 TEXT,
    residual_risk_1 TEXT,
    control_responsibility_1 TEXT,
    step_2 TEXT,
    hazards_2 TEXT,
    risks_2 TEXT,
    initial_risk_2 TEXT,
    control_measures_2 TEXT,
    residual_risk_2 TEXT,
    control_responsibility_2 TEXT,
    step_3 TEXT,
    hazards_3 TEXT,
    risks_3 TEXT,
    initial_risk_3 TEXT,
    control_measures_3 TEXT,
    residual_risk_3 TEXT,
    control_responsibility_3 TEXT,
    step_4 TEXT,
    hazards_4 TEXT,
    risks_4 TEXT,
    initial_risk_4 TEXT,
    control_measures_4 TEXT,
    residual_risk_4 TEXT,
    control_responsibility_4 TEXT,
    step_5 TEXT,
    hazards_5 TEXT,
    risks_5 TEXT,
    initial_risk_5 TEXT,
    control_measures_5 TEXT,
    residual_risk_5 TEXT,
    control_responsibility_5 TEXT,
    step_6 TEXT,
    hazards_6 TEXT,
    risks_6 TEXT,
    initial_risk_6 TEXT,
    control_measures_6 TEXT,
    residual_risk_6 TEXT,
    control_responsibility_6 TEXT,
    step_7 TEXT,
    hazards_7 TEXT,
    risks_7 TEXT,
    initial_risk_7 TEXT,
    control_measures_7 TEXT,
    residual_risk_7 TEXT,
    control_responsibility_7 TEXT,
    step_8 TEXT,
    hazards_8 TEXT,
    risks_8 TEXT,
    initial_risk_8 TEXT,
    control_measures_8 TEXT,
    residual_risk_8 TEXT,
    control_responsibility_8 TEXT,
    step_9 TEXT,
    hazards_9 TEXT,
    risks_9 TEXT,
    initial_risk_9 TEXT,
    control_measures_9 TEXT,
    residual_risk_9 TEXT,
    control_responsibility_9 TEXT,
    step_10 TEXT,
    hazards_10 TEXT,
    risks_10 TEXT,
    initial_risk_10 TEXT,
    control_measures_10 TEXT,
    residual_risk_10 TEXT,
    control_responsibility_10 TEXT,
    step_11 TEXT,
    hazards_11 TEXT,
    risks_11 TEXT,
    initial_risk_11 TEXT,
    control_measures_11 TEXT,
    residual_risk_11 TEXT,
    control_responsibility_11 TEXT,
    step_12 TEXT,
    hazards_12 TEXT,
    risks_12 TEXT,
    initial_risk_12 TEXT,
    control_measures_12 TEXT,
    residual_risk_12 TEXT,
    control_responsibility_12 TEXT,
    sign_date_1 TEXT,
    sign_name_1 TEXT,
    sign_sig_1 TEXT,
    sign_date_2 TEXT,
    sign_name_2 TEXT,
    sign_sig_2 TEXT,
    sign_date_3 TEXT,
    sign_name_3 TEXT,
    sign_sig_3 TEXT,
    sign_date_4 TEXT,
    sign_name_4 TEXT,
    sign_sig_4 TEXT,
    sign_date_5 TEXT,
    sign_name_5 TEXT,
    sign_sig_5 TEXT,
    sign_date_6 TEXT,
    sign_name_6 TEXT,
    sign_sig_6 TEXT,
    sign_date_7 TEXT,
    sign_name_7 TEXT,
    sign_sig_7 TEXT,
    sign_date_8 TEXT,
    sign_name_8 TEXT,
    sign_sig_8 TEXT,
    sign_date_9 TEXT,
    sign_name_9 TEXT,
    sign_sig_9 TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at TEXT DATETIME DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS purchase_order;
CREATE TABLE purchase_order(
    uuid TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL,
    store_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    order_date TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS purchase_order_item;
CREATE TABLE purchase_order_item(
    uuid TEXT PRIMARY KEY,
    purchase_order_id TEXT NOT NULL,
    item_name TEXT NOT NULL,
    item_type_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_order(uuid) ON DELETE CASCADE,
    FOREIGN KEY (item_type_id) REFERENCES item_types(uuid) ON DELETE CASCADE
);

DROP TABLE IF EXISTS item_types;
CREATE TABLE item_types(
    uuid TEXT PRIMARY KEY,
    type TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS stores;
CREATE TABLE stores(
    uuid TEXT PRIMARY KEY,
    business_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    address TEXT NOT NULL,
    suburb TEXT NOT NULL,
    city TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS tasks;
CREATE TABLE tasks(
    uuid TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL,
    retries INTEGER DEFAULT 0,
    max_retries INTEGET DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('employee',4);
INSERT INTO sqlite_sequence VALUES('employee_auth',4);
INSERT INTO sqlite_sequence VALUES('leave_request',11);
INSERT INTO sqlite_sequence VALUES('job',16);
INSERT INTO sqlite_sequence VALUES('timesheet_week',8);
INSERT INTO sqlite_sequence VALUES('timesheet',56);
COMMIT;
