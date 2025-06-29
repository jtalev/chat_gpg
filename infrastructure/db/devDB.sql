PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
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
INSERT INTO employee_role VALUES('0cd836f9-394b-4968-b778-af1cbbf1d1bb','5972276','management','2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO employee_role VALUES('0cd836f9-394b-4968-b778-af1cbbf1d1bc','1234567','foreman','2025-03-23 06:53:39','2025-03-23 06:53:39');
INSERT INTO employee_role VALUES('0cd836f9-394b-4968-b778-af1cbbf1d1bd','7654321','employee','2025-03-23 06:53:39','2025-03-23 06:53:39');
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
CREATE TABLE note(
    uuid TEXT PRIMARY KEY,
    job_id INTEGER NOT NULL,
    note_type TEXT NOT NULL,
    note TEXT NOT NULL,
    is_archived INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO note VALUES('f83b9618-b73b-471c-beae-d1f574023055',1,'paint_note','{"note_uuid":"123abc","brand":"Haymes","product":"Expressions","colour":"Natural White","finish":"Low Sheen","area":"Living Room","coats":2,"surfaces":"Walls","notes":""}',0,'2025-06-10 06:41:39','2025-06-10 06:41:39');
INSERT INTO note VALUES('6b804426-281d-4d59-b6cc-e4de6c98ebc4',2,'paint_note','{"brand":"Haymes","product":"Ultratrim","colour":"Natural White","finish":"Semi-gloss","area":"Living Room","coats":2,"surfaces":"Woodwork","notes":""}',0,'2025-06-10 06:55:55','2025-06-10 06:55:55');
INSERT INTO note VALUES('667fa8d9-c50f-4d66-b694-52082b4956ca',2,'paint_note','{"note_uuid":"667fa8d9-c50f-4d66-b694-52082b4956ca","brand":"Haymes","product":"Ultratrim","colour":"Natural White","finish":"Semi-gloss","area":"Living Room","coats":2,"surfaces":"Woodwork","notes":""}',0,'2025-06-10 07:12:46','2025-06-10 07:12:46');
INSERT INTO note VALUES('9a09bc18-85e7-492f-91bd-386b31f531f8',24,'paint_note','{"note_uuid":"9a09bc18-85e7-492f-91bd-386b31f531f8","brand":"Haymes","product":"Ultratrim","colour":"Natural White","finish":"Semi-gloss","area":"Living Room","coats":2,"surfaces":"Woodwork","notes":""}',0,'2025-06-11 09:31:55','2025-06-11 09:31:55');
INSERT INTO note VALUES('043b9c51-7500-415a-8551-a614838c5f79',24,'paint_note','{"note_uuid":"043b9c51-7500-415a-8551-a614838c5f79","brand":"Haymes","product":"Ultratrim","colour":"Natural White","finish":"Semi-gloss","area":"Living Room","coats":2,"surfaces":"Woodwork","notes":""}',0,'2025-06-11 09:32:01','2025-06-11 09:32:01');
INSERT INTO note VALUES('6eb7f2f0-92d6-4955-a896-447b6ad74440',18,'paint_note','{"note_uuid":"6eb7f2f0-92d6-4955-a896-447b6ad74440","brand":"Haymes","product":"Expressions","colour":"Antique White","finish":"lkjg","area":"lkjgl","coats":0,"surfaces":"kjglkjg","notes":"lkjglk\n","JobId":0}',0,'2025-06-27 12:18:45','2025-06-27 12:18:45');
INSERT INTO note VALUES('087d6df0-e4a0-432d-8812-c60400d78916',18,'task_note','{"note_uuid":"087d6df0-e4a0-432d-8812-c60400d78916","title":"updated","description":"lkn;ln","status":"pending","priority":"low","notes":"ln;ln;ln;","JobId":0}',0,'2025-06-27 12:18:54','2025-06-27 12:18:54');
INSERT INTO note VALUES('ca6decbc-d309-4c06-afc0-148efed064bd',18,'task_note','{"note_uuid":"ca6decbc-d309-4c06-afc0-148efed064bd","title":"test","description":"updated description","status":"pending","priority":"n/a","notes":"test\n","JobId":0}',0,'2025-06-27 12:29:05','2025-06-27 12:29:05');
INSERT INTO note VALUES('7ba8fedf-ac50-4b8f-95f1-48327b695fa8',18,'task_note','{"note_uuid":"7ba8fedf-ac50-4b8f-95f1-48327b695fa8","title":"update","description":"","status":"pending","priority":"n/a","notes":"","JobId":0}',0,'2025-06-27 12:43:05','2025-06-27 12:43:05');
INSERT INTO note VALUES('5586903e-fd74-4a11-b3b8-0b9255911195',18,'paint_note','{"note_uuid":"5586903e-fd74-4a11-b3b8-0b9255911195","brand":"update","product":"update","colour":"update","finish":"lkh","area":"lkhj","coats":98,"surfaces":"kjlkjh","notes":"lkhk","JobId":0}',0,'2025-06-27 12:47:48','2025-06-27 12:47:48');
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
INSERT INTO job VALUES(17,'COGG - Dunes Toilet Block',12,'Ocean Rd','Ocean Grove','3214','Geelong',0,'2025-06-10 10:01:25','2025-06-10 10:01:25');
INSERT INTO job VALUES(18,'Will''s Inlaws',24,'Seventh Ave','Anglesea','3214','Anglesea',0,'2025-06-11 09:05:54','2025-06-21 02:09:53');
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
CREATE TABLE purchase_order(
    uuid TEXT PRIMARY KEY,
    employee_id TEXT NOT NULL,
    store_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    order_date TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO purchase_order VALUES('b6a10ca6-0380-4404-ab19-849bf8fdb305','5972276','8713290a-e683-4686-85da-4a8be332a157',1,'2025-06-03','2025-06-03 10:46:50','2025-06-03 10:46:50');
CREATE TABLE item_types(
    uuid TEXT PRIMARY KEY,
    type TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO item_types VALUES('b2b8650b-0400-403e-b972-e86f041db872','Expressions Low Sheen','Haymes','2025-06-03 08:09:42','2025-06-03 10:44:56');
INSERT INTO item_types VALUES('6b43e590-12b2-4930-bf82-ecb22db79b96','Brushes','Haymes','2025-06-03 10:46:19','2025-06-03 10:46:19');
CREATE TABLE item_size(
    uuid TEXT PRIMARY KEY,
    size TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO item_size VALUES('123','63mm','brushes','2025-06-03 07:00:29','2025-06-03 08:31:32');
INSERT INTO item_size VALUES('3ccfd644-d897-44bf-b037-56fa23a27d9f','15L','paint','2025-06-03 10:44:42','2025-06-03 10:44:42');
CREATE TABLE purchase_order_item(
    uuid TEXT PRIMARY KEY,
    purchase_order_id TEXT NOT NULL,
    item_name TEXT NOT NULL,
    item_type_id TEXT NOT NULL,
    item_size_id TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (purchase_order_id) REFERENCES purchase_order(uuid) ON DELETE CASCADE,
    FOREIGN KEY (item_type_id) REFERENCES item_types(uuid) ON DELETE CASCADE
);
INSERT INTO purchase_order_item VALUES('dafa02fe-a891-417c-8c3f-d30ed24644d7','b6a10ca6-0380-4404-ab19-849bf8fdb305','monument','b2b8650b-0400-403e-b972-e86f041db872','3ccfd644-d897-44bf-b037-56fa23a27d9f',1,'2025-06-03 10:46:50','2025-06-03 10:46:50');
INSERT INTO purchase_order_item VALUES('259279c5-48cb-4a7c-865a-f199640b6f1b','b6a10ca6-0380-4404-ab19-849bf8fdb305','Sash Cutter','6b43e590-12b2-4930-bf82-ecb22db79b96','123',6,'2025-06-03 10:46:50','2025-06-03 10:46:50');
CREATE TABLE stores(
    uuid TEXT PRIMARY KEY,
    business_name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT NOT NULL,
    address TEXT NOT NULL,
    suburb TEXT NOT NULL,
    city TEXT NOT NULL,
    account_code TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    modified_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO stores VALUES('8713290a-e683-4686-85da-4a8be332a157','Haymes Geelong Westfsafsd','jblkb','kblkjb','lkblk','jblkjbl','kjblkjbl','lkblkbklj','2025-06-03 07:03:17','2025-06-03 10:42:18');
CREATE TABLE tasks(
    uuid TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    handler TEXT NOT NULL,
    payload TEXT NOT NULL,
    status TEXT NOT NULL,
    retries INTEGER DEFAULT 0,
    max_retries INTEGET DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO tasks VALUES('e3109f9d-fad1-42cf-8df4-6e7e49cb614c','scheduled','db_backup',X'7b226e6578745f72756e5f6174223a22323032352d30362d30335431373a30313a34392e393232343033342b31303a3030227d','scheduled',0,3,'2025-06-03 07:01:49');
INSERT INTO tasks VALUES('0b4589f1-5038-4635-ba35-a10704b715c9','one_time','send_email',X'7b2273656e6465725f6e616d65223a224765656c6f6e67205061696e742047726f75702041646d696e222c2273656e6465725f656d61696c223a2261646d696e406765656c6f6e677061696e7467726f75702e636f6d2e6175222c22726563697069656e745f6e616d65223a224861796d6573204765656c6f6e672057657374222c22726563697069656e745f656d61696c223a226a2e74616c6576406f75746c6f6f6b2e636f6d222c227375626a656374223a225075726368617365204f72646572202d204765656c6f6e67205061696e742047726f7570222c22706c61696e5f746578745f636f6e74656e74223a225075726368617365204f72646572202d204765656c6f6e67205061696e742047726f75705c6e446174653a20323032352d30352d32375c6e53746f72653a204861796d6573204765656c6f6e6720576573745c6e4f7264657265642042793a20536c6964204b65737472656c5c6e5265666572656e63653a20547265776865656c61204176655c6e4163636f756e7420436f64653a206c6b626c6b626b6c6a5c6e5c6e4974656d733a5c6e3120782045787072657373696f6e73204c6f7720536865656e2031354c206d6f6e756d656e745c6e222c2268746d6c5f636f6e74656e74223a22227d','completed',0,3,'2025-06-03 07:20:57');
INSERT INTO tasks VALUES('9ff63189-f799-44e6-909e-377b56b50887','one_time','send_email',X'7b2273656e6465725f6e616d65223a224765656c6f6e67205061696e742047726f75702041646d696e222c2273656e6465725f656d61696c223a2261646d696e406765656c6f6e677061696e7467726f75702e636f6d2e6175222c22726563697069656e745f6e616d65223a224861796d6573204765656c6f6e672057657374222c22726563697069656e745f656d61696c223a226a2e74616c6576406f75746c6f6f6b2e636f6d222c227375626a656374223a225075726368617365204f72646572202d204765656c6f6e67205061696e742047726f7570222c22706c61696e5f746578745f636f6e74656e74223a225075726368617365204f72646572202d204765656c6f6e67205061696e742047726f75705c6e446174653a20323032352d30362d30325c6e53746f72653a204861796d6573204765656c6f6e6720576573745c6e4f7264657265642042793a20536c6964204b65737472656c5c6e5265666572656e63653a20547265776865656c61204176655c6e4163636f756e7420436f64653a206c6b626c6b626b6c6a5c6e5c6e4974656d733a5c6e3120782045787072657373696f6e73204c6f7720536865656e2031354c206d6f6e756d656e745c6e222c2268746d6c5f636f6e74656e74223a22227d','completed',0,3,'2025-06-03 07:23:48');
INSERT INTO tasks VALUES('9bf41c05-05f8-4c29-8034-c6f44d966402','one_time','send_email',X'7b2273656e6465725f6e616d65223a224765656c6f6e67205061696e742047726f75702041646d696e222c2273656e6465725f656d61696c223a2261646d696e406765656c6f6e677061696e7467726f75702e636f6d2e6175222c22726563697069656e745f6e616d65223a224861796d6573204765656c6f6e672057657374667361667364222c22726563697069656e745f656d61696c223a226a2e74616c6576406f75746c6f6f6b2e636f6d222c227375626a656374223a225075726368617365204f72646572202d204765656c6f6e67205061696e742047726f7570222c22706c61696e5f746578745f636f6e74656e74223a225075726368617365204f72646572202d204765656c6f6e67205061696e742047726f75705c6e446174653a20323032352d30362d30335c6e53746f72653a204861796d6573204765656c6f6e6720576573746673616673645c6e4f7264657265642042793a20536c6964204b65737472656c5c6e5265666572656e63653a20547265776865656c61204176655c6e4163636f756e7420436f64653a206c6b626c6b626b6c6a5c6e5c6e4974656d733a5c6e3120782031354c2045787072657373696f6e73204c6f7720536865656e206d6f6e756d656e745c6e3620782036336d6d20427275736865732053617368204375747465725c6e222c2268746d6c5f636f6e74656e74223a22227d','completed',0,3,'2025-06-03 10:46:50');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('employee',4);
INSERT INTO sqlite_sequence VALUES('employee_auth',4);
INSERT INTO sqlite_sequence VALUES('leave_request',11);
INSERT INTO sqlite_sequence VALUES('job',18);
INSERT INTO sqlite_sequence VALUES('timesheet_week',8);
INSERT INTO sqlite_sequence VALUES('timesheet',56);
COMMIT;
