package infrastructure

import (
	"database/sql"
	"log"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func GetSwms(db *sql.DB) ([]models.Swms, error) {
	q := `
	select * from swms;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outSwms := []models.Swms{}
	for rows.Next() {
		swms := models.Swms{}
		if err := rows.Scan(
			&swms.UUID, &swms.JobId, &swms.ProjectActivity, &swms.ProjectNumber, &swms.SiteAddress, &swms.ContactName, &swms.ContactNumber,
			&swms.EmailAddress, &swms.SwmsDate, &swms.HighRiskWorks, &swms.SafetyGloves, &swms.SafetyBoots, &swms.SafetyGlasses,
			&swms.ProtectiveClothing, &swms.RespiratoryProtection, &swms.HiVisClothing, &swms.SafetyHelmet, &swms.FallArrest,
			&swms.Other1, &swms.Other2, &swms.StepBelow2m, &swms.StepAbove2m, &swms.Scaffold, &swms.PressureWasherDiesel,
			&swms.RoofAnchorPoints, &swms.ExtensionLadder, &swms.ElectricScissorLift, &swms.DieselScissorLift,
			&swms.ElectricKnuckleBoom, &swms.DieselKnuckleBoom, &swms.AirlessSprayGun, &swms.AngleGrinder,
			&swms.Step1, &swms.Hazards1, &swms.Risks1, &swms.InitialRisk1, &swms.ControlMeasures1, &swms.ResidualRisk1, &swms.ControlResponsibility1,
			&swms.Step2, &swms.Hazards2, &swms.Risks2, &swms.InitialRisk2, &swms.ControlMeasures2, &swms.ResidualRisk2, &swms.ControlResponsibility2,
			&swms.Step3, &swms.Hazards3, &swms.Risks3, &swms.InitialRisk3, &swms.ControlMeasures3, &swms.ResidualRisk3, &swms.ControlResponsibility3,
			&swms.Step4, &swms.Hazards4, &swms.Risks4, &swms.InitialRisk4, &swms.ControlMeasures4, &swms.ResidualRisk4, &swms.ControlResponsibility4,
			&swms.Step5, &swms.Hazards5, &swms.Risks5, &swms.InitialRisk5, &swms.ControlMeasures5, &swms.ResidualRisk5, &swms.ControlResponsibility5,
			&swms.Step6, &swms.Hazards6, &swms.Risks6, &swms.InitialRisk6, &swms.ControlMeasures6, &swms.ResidualRisk6, &swms.ControlResponsibility6,
			&swms.Step7, &swms.Hazards7, &swms.Risks7, &swms.InitialRisk7, &swms.ControlMeasures7, &swms.ResidualRisk7, &swms.ControlResponsibility7,
			&swms.Step8, &swms.Hazards8, &swms.Risks8, &swms.InitialRisk8, &swms.ControlMeasures8, &swms.ResidualRisk8, &swms.ControlResponsibility8,
			&swms.Step9, &swms.Hazards9, &swms.Risks9, &swms.InitialRisk9, &swms.ControlMeasures9, &swms.ResidualRisk9, &swms.ControlResponsibility9,
			&swms.Step10, &swms.Hazards10, &swms.Risks10, &swms.InitialRisk10, &swms.ControlMeasures10, &swms.ResidualRisk10, &swms.ControlResponsibility10,
			&swms.Step11, &swms.Hazards11, &swms.Risks11, &swms.InitialRisk11, &swms.ControlMeasures11, &swms.ResidualRisk11, &swms.ControlResponsibility11,
			&swms.Step12, &swms.Hazards12, &swms.Risks12, &swms.InitialRisk12, &swms.ControlMeasures12, &swms.ResidualRisk12, &swms.ControlResponsibility12,
			&swms.SignDate1, &swms.SignName1, &swms.SignSig1,
			&swms.SignDate2, &swms.SignName2, &swms.SignSig2,
			&swms.SignDate3, &swms.SignName3, &swms.SignSig3,
			&swms.SignDate4, &swms.SignName4, &swms.SignSig4,
			&swms.SignDate5, &swms.SignName5, &swms.SignSig5,
			&swms.SignDate6, &swms.SignName6, &swms.SignSig6,
			&swms.SignDate7, &swms.SignName7, &swms.SignSig7,
			&swms.SignDate8, &swms.SignName8, &swms.SignSig8,
			&swms.SignDate9, &swms.SignName9, &swms.SignSig9,
			&swms.CreatedAt, &swms.ModifiedAt,
		); err != nil {
			return nil, err
		}
		outSwms = append(outSwms, swms)
	}
	return outSwms, nil
}

func GetSwmsByUUID(uuid string, db *sql.DB) (models.Swms, error) {
	q := `
	select * from swms where uuid = ?;
	`

	rows, err := db.Query(q, uuid)
	if err != nil {
		return models.Swms{}, err
	}
	defer rows.Close()

	swms := models.Swms{}
	for rows.Next() {
		if err := rows.Scan(
			&swms.UUID, &swms.JobId, &swms.ProjectActivity, &swms.ProjectNumber, &swms.SiteAddress, &swms.ContactName, &swms.ContactNumber,
			&swms.EmailAddress, &swms.SwmsDate, &swms.HighRiskWorks, &swms.SafetyGloves, &swms.SafetyBoots, &swms.SafetyGlasses,
			&swms.ProtectiveClothing, &swms.RespiratoryProtection, &swms.HiVisClothing, &swms.SafetyHelmet, &swms.FallArrest,
			&swms.Other1, &swms.Other2, &swms.StepBelow2m, &swms.StepAbove2m, &swms.Scaffold, &swms.PressureWasherDiesel,
			&swms.RoofAnchorPoints, &swms.ExtensionLadder, &swms.ElectricScissorLift, &swms.DieselScissorLift,
			&swms.ElectricKnuckleBoom, &swms.DieselKnuckleBoom, &swms.AirlessSprayGun, &swms.AngleGrinder,
			&swms.Step1, &swms.Hazards1, &swms.Risks1, &swms.InitialRisk1, &swms.ControlMeasures1, &swms.ResidualRisk1, &swms.ControlResponsibility1,
			&swms.Step2, &swms.Hazards2, &swms.Risks2, &swms.InitialRisk2, &swms.ControlMeasures2, &swms.ResidualRisk2, &swms.ControlResponsibility2,
			&swms.Step3, &swms.Hazards3, &swms.Risks3, &swms.InitialRisk3, &swms.ControlMeasures3, &swms.ResidualRisk3, &swms.ControlResponsibility3,
			&swms.Step4, &swms.Hazards4, &swms.Risks4, &swms.InitialRisk4, &swms.ControlMeasures4, &swms.ResidualRisk4, &swms.ControlResponsibility4,
			&swms.Step5, &swms.Hazards5, &swms.Risks5, &swms.InitialRisk5, &swms.ControlMeasures5, &swms.ResidualRisk5, &swms.ControlResponsibility5,
			&swms.Step6, &swms.Hazards6, &swms.Risks6, &swms.InitialRisk6, &swms.ControlMeasures6, &swms.ResidualRisk6, &swms.ControlResponsibility6,
			&swms.Step7, &swms.Hazards7, &swms.Risks7, &swms.InitialRisk7, &swms.ControlMeasures7, &swms.ResidualRisk7, &swms.ControlResponsibility7,
			&swms.Step8, &swms.Hazards8, &swms.Risks8, &swms.InitialRisk8, &swms.ControlMeasures8, &swms.ResidualRisk8, &swms.ControlResponsibility8,
			&swms.Step9, &swms.Hazards9, &swms.Risks9, &swms.InitialRisk9, &swms.ControlMeasures9, &swms.ResidualRisk9, &swms.ControlResponsibility9,
			&swms.Step10, &swms.Hazards10, &swms.Risks10, &swms.InitialRisk10, &swms.ControlMeasures10, &swms.ResidualRisk10, &swms.ControlResponsibility10,
			&swms.Step11, &swms.Hazards11, &swms.Risks11, &swms.InitialRisk11, &swms.ControlMeasures11, &swms.ResidualRisk11, &swms.ControlResponsibility11,
			&swms.Step12, &swms.Hazards12, &swms.Risks12, &swms.InitialRisk12, &swms.ControlMeasures12, &swms.ResidualRisk12, &swms.ControlResponsibility12,
			&swms.SignDate1, &swms.SignName1, &swms.SignSig1,
			&swms.SignDate2, &swms.SignName2, &swms.SignSig2,
			&swms.SignDate3, &swms.SignName3, &swms.SignSig3,
			&swms.SignDate4, &swms.SignName4, &swms.SignSig4,
			&swms.SignDate5, &swms.SignName5, &swms.SignSig5,
			&swms.SignDate6, &swms.SignName6, &swms.SignSig6,
			&swms.SignDate7, &swms.SignName7, &swms.SignSig7,
			&swms.SignDate8, &swms.SignName8, &swms.SignSig8,
			&swms.SignDate9, &swms.SignName9, &swms.SignSig9,
			&swms.CreatedAt, &swms.ModifiedAt,
		); err != nil {
			return swms, err
		}
	}
	return swms, nil
}

func PostSwms(swms models.Swms, db *sql.DB) (bool, error) {
	q := `
INSERT INTO swms (
	uuid, job_id, project_activity, project_number, site_address, contact_name, contact_number,
	email_address, swms_date, high_risk_works, safety_gloves, safety_boots, safety_glasses,
	protective_clothing, respiratory_protection, hi_vis_clothing, safety_helmet, fall_arrest,
	other_1, other_2, step_below_2m, step_above_2m, scaffold, pressure_washer_diesel,
	roof_anchor_points, extension_ladder, electric_scissor_lift, diesel_scissor_lift,
	electric_knuckle_boom, diesel_knuckle_boom, airless_spray_gun, angle_grinder,
	step_1, hazards_1, risks_1, initial_risk_1, control_measures_1, residual_risk_1, control_responsibility_1,
	step_2, hazards_2, risks_2, initial_risk_2, control_measures_2, residual_risk_2, control_responsibility_2,
	step_3, hazards_3, risks_3, initial_risk_3, control_measures_3, residual_risk_3, control_responsibility_3,
	step_4, hazards_4, risks_4, initial_risk_4, control_measures_4, residual_risk_4, control_responsibility_4,
	step_5, hazards_5, risks_5, initial_risk_5, control_measures_5, residual_risk_5, control_responsibility_5,
	step_6, hazards_6, risks_6, initial_risk_6, control_measures_6, residual_risk_6, control_responsibility_6,
	step_7, hazards_7, risks_7, initial_risk_7, control_measures_7, residual_risk_7, control_responsibility_7,
	step_8, hazards_8, risks_8, initial_risk_8, control_measures_8, residual_risk_8, control_responsibility_8,
	step_9, hazards_9, risks_9, initial_risk_9, control_measures_9, residual_risk_9, control_responsibility_9,
	step_10, hazards_10, risks_10, initial_risk_10, control_measures_10, residual_risk_10, control_responsibility_10,
	step_11, hazards_11, risks_11, initial_risk_11, control_measures_11, residual_risk_11, control_responsibility_11,
	step_12, hazards_12, risks_12, initial_risk_12, control_measures_12, residual_risk_12, control_responsibility_12,
	sign_date_1, sign_name_1, sign_sig_1,
	sign_date_2, sign_name_2, sign_sig_2,
	sign_date_3, sign_name_3, sign_sig_3,
	sign_date_4, sign_name_4, sign_sig_4,
	sign_date_5, sign_name_5, sign_sig_5,
	sign_date_6, sign_name_6, sign_sig_6,
	sign_date_7, sign_name_7, sign_sig_7,
	sign_date_8, sign_name_8, sign_sig_8,
	sign_date_9, sign_name_9, sign_sig_9
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
	$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37,
	$38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54,
	$55, $56, $57, $58, $59, $60, $61, $62, $63, $64, $65, $66, $67, $68, $69, $70, $71,
	$72, $73, $74, $75, $76, $77, $78, $79, $80, $81, $82, $83, $84, $85, $86, $87, $88,
	$89, $90, $91, $92, $93, $94, $95, $96, $97, $98, $99, $100, $101, $102, $103, $104,
	$105, $106, $107, $108, $109, $110, $111, $112, $113, $114, $115, $116, $117, $118,
	$119, $120, $121, $122, $123, $124, $125, $126, $127, $128, $129, $130, $131, $132,
	$133, $134, $135, $136, $137, $138, $139, $140, $141, $142, $143
);`

	_, err := db.Exec(q,
		swms.UUID, swms.JobId, swms.ProjectActivity, swms.ProjectNumber, swms.SiteAddress, swms.ContactName, swms.ContactNumber,
		swms.EmailAddress, swms.SwmsDate, swms.HighRiskWorks, swms.SafetyGloves, swms.SafetyBoots, swms.SafetyGlasses,
		swms.ProtectiveClothing, swms.RespiratoryProtection, swms.HiVisClothing, swms.SafetyHelmet, swms.FallArrest,
		swms.Other1, swms.Other2, swms.StepBelow2m, swms.StepAbove2m, swms.Scaffold, swms.PressureWasherDiesel,
		swms.RoofAnchorPoints, swms.ExtensionLadder, swms.ElectricScissorLift, swms.DieselScissorLift,
		swms.ElectricKnuckleBoom, swms.DieselKnuckleBoom, swms.AirlessSprayGun, swms.AngleGrinder,
		swms.Step1, swms.Hazards1, swms.Risks1, swms.InitialRisk1, swms.ControlMeasures1, swms.ResidualRisk1, swms.ControlResponsibility1,
		swms.Step2, swms.Hazards2, swms.Risks2, swms.InitialRisk2, swms.ControlMeasures2, swms.ResidualRisk2, swms.ControlResponsibility2,
		swms.Step3, swms.Hazards3, swms.Risks3, swms.InitialRisk3, swms.ControlMeasures3, swms.ResidualRisk3, swms.ControlResponsibility3,
		swms.Step4, swms.Hazards4, swms.Risks4, swms.InitialRisk4, swms.ControlMeasures4, swms.ResidualRisk4, swms.ControlResponsibility4,
		swms.Step5, swms.Hazards5, swms.Risks5, swms.InitialRisk5, swms.ControlMeasures5, swms.ResidualRisk5, swms.ControlResponsibility5,
		swms.Step6, swms.Hazards6, swms.Risks6, swms.InitialRisk6, swms.ControlMeasures6, swms.ResidualRisk6, swms.ControlResponsibility6,
		swms.Step7, swms.Hazards7, swms.Risks7, swms.InitialRisk7, swms.ControlMeasures7, swms.ResidualRisk7, swms.ControlResponsibility7,
		swms.Step8, swms.Hazards8, swms.Risks8, swms.InitialRisk8, swms.ControlMeasures8, swms.ResidualRisk8, swms.ControlResponsibility8,
		swms.Step9, swms.Hazards9, swms.Risks9, swms.InitialRisk9, swms.ControlMeasures9, swms.ResidualRisk9, swms.ControlResponsibility9,
		swms.Step10, swms.Hazards10, swms.Risks10, swms.InitialRisk10, swms.ControlMeasures10, swms.ResidualRisk10, swms.ControlResponsibility10,
		swms.Step11, swms.Hazards11, swms.Risks11, swms.InitialRisk11, swms.ControlMeasures11, swms.ResidualRisk11, swms.ControlResponsibility11,
		swms.Step12, swms.Hazards12, swms.Risks12, swms.InitialRisk12, swms.ControlMeasures12, swms.ResidualRisk12, swms.ControlResponsibility12,
		swms.SignDate1, swms.SignName1, swms.SignSig1,
		swms.SignDate2, swms.SignName2, swms.SignSig2,
		swms.SignDate3, swms.SignName3, swms.SignSig3,
		swms.SignDate4, swms.SignName4, swms.SignSig4,
		swms.SignDate5, swms.SignName5, swms.SignSig5,
		swms.SignDate6, swms.SignName6, swms.SignSig6,
		swms.SignDate7, swms.SignName7, swms.SignSig7,
		swms.SignDate8, swms.SignName8, swms.SignSig8,
		swms.SignDate9, swms.SignName9, swms.SignSig9,
	)
	if err != nil {
		log.Printf("failed to insert swms: %v", err)
		return false, err
	}
	return true, err
}

func PutSwms(swms models.Swms, db *sql.DB) (bool, error) {
	q := `
UPDATE swms SET
	job_id = $1, project_activity = $2, project_number = $3, site_address = $4, contact_name = $5,
	contact_number = $6, email_address = $7, swms_date = $8, high_risk_works = $9,
	safety_gloves = $10, safety_boots = $11, safety_glasses = $12, protective_clothing = $13,
	respiratory_protection = $14, hi_vis_clothing = $15, safety_helmet = $16, fall_arrest = $17,
	other_1 = $18, other_2 = $19, step_below_2m = $20, step_above_2m = $21, scaffold = $22,
	pressure_washer_diesel = $23, roof_anchor_points = $24, extension_ladder = $25,
	electric_scissor_lift = $26, diesel_scissor_lift = $27, electric_knuckle_boom = $28,
	diesel_knuckle_boom = $29, airless_spray_gun = $30, angle_grinder = $31,

	step_1 = $32, hazards_1 = $33, risks_1 = $34, initial_risk_1 = $35,
	control_measures_1 = $36, residual_risk_1 = $37, control_responsibility_1 = $38,
	step_2 = $39, hazards_2 = $40, risks_2 = $41, initial_risk_2 = $42,
	control_measures_2 = $43, residual_risk_2 = $44, control_responsibility_2 = $45,
	step_3 = $46, hazards_3 = $47, risks_3 = $48, initial_risk_3 = $49,
	control_measures_3 = $50, residual_risk_3 = $51, control_responsibility_3 = $52,
	step_4 = $53, hazards_4 = $54, risks_4 = $55, initial_risk_4 = $56,
	control_measures_4 = $57, residual_risk_4 = $58, control_responsibility_4 = $59,
	step_5 = $60, hazards_5 = $61, risks_5 = $62, initial_risk_5 = $63,
	control_measures_5 = $64, residual_risk_5 = $65, control_responsibility_5 = $66,
	step_6 = $67, hazards_6 = $68, risks_6 = $69, initial_risk_6 = $70,
	control_measures_6 = $71, residual_risk_6 = $72, control_responsibility_6 = $73,
	step_7 = $74, hazards_7 = $75, risks_7 = $76, initial_risk_7 = $77,
	control_measures_7 = $78, residual_risk_7 = $79, control_responsibility_7 = $80,
	step_8 = $81, hazards_8 = $82, risks_8 = $83, initial_risk_8 = $84,
	control_measures_8 = $85, residual_risk_8 = $86, control_responsibility_8 = $87,
	step_9 = $88, hazards_9 = $89, risks_9 = $90, initial_risk_9 = $91,
	control_measures_9 = $92, residual_risk_9 = $93, control_responsibility_9 = $94,
	step_10 = $95, hazards_10 = $96, risks_10 = $97, initial_risk_10 = $98,
	control_measures_10 = $99, residual_risk_10 = $100, control_responsibility_10 = $101,
	step_11 = $102, hazards_11 = $103, risks_11 = $104, initial_risk_11 = $105,
	control_measures_11 = $106, residual_risk_11 = $107, control_responsibility_11 = $108,
	step_12 = $109, hazards_12 = $110, risks_12 = $111, initial_risk_12 = $112,
	control_measures_12 = $113, residual_risk_12 = $114, control_responsibility_12 = $115,

	sign_date_1 = $116, sign_name_1 = $117, sign_sig_1 = $118,
	sign_date_2 = $119, sign_name_2 = $120, sign_sig_2 = $121,
	sign_date_3 = $122, sign_name_3 = $123, sign_sig_3 = $124,
	sign_date_4 = $125, sign_name_4 = $126, sign_sig_4 = $127,
	sign_date_5 = $128, sign_name_5 = $129, sign_sig_5 = $130,
	sign_date_6 = $131, sign_name_6 = $132, sign_sig_6 = $133,
	sign_date_7 = $134, sign_name_7 = $135, sign_sig_7 = $136,
	sign_date_8 = $137, sign_name_8 = $138, sign_sig_8 = $139,
	sign_date_9 = $140, sign_name_9 = $141, sign_sig_9 = $142

WHERE uuid = $143;
`

	_, err := db.Exec(q,
		swms.JobId, swms.ProjectActivity, swms.ProjectNumber, swms.SiteAddress, swms.ContactName,
		swms.ContactNumber, swms.EmailAddress, swms.SwmsDate, swms.HighRiskWorks,
		swms.SafetyGloves, swms.SafetyBoots, swms.SafetyGlasses, swms.ProtectiveClothing,
		swms.RespiratoryProtection, swms.HiVisClothing, swms.SafetyHelmet, swms.FallArrest,
		swms.Other1, swms.Other2, swms.StepBelow2m, swms.StepAbove2m, swms.Scaffold,
		swms.PressureWasherDiesel, swms.RoofAnchorPoints, swms.ExtensionLadder,
		swms.ElectricScissorLift, swms.DieselScissorLift, swms.ElectricKnuckleBoom,
		swms.DieselKnuckleBoom, swms.AirlessSprayGun, swms.AngleGrinder,

		swms.Step1, swms.Hazards1, swms.Risks1, swms.InitialRisk1,
		swms.ControlMeasures1, swms.ResidualRisk1, swms.ControlResponsibility1,
		swms.Step2, swms.Hazards2, swms.Risks2, swms.InitialRisk2,
		swms.ControlMeasures2, swms.ResidualRisk2, swms.ControlResponsibility2,
		swms.Step3, swms.Hazards3, swms.Risks3, swms.InitialRisk3,
		swms.ControlMeasures3, swms.ResidualRisk3, swms.ControlResponsibility3,
		swms.Step4, swms.Hazards4, swms.Risks4, swms.InitialRisk4,
		swms.ControlMeasures4, swms.ResidualRisk4, swms.ControlResponsibility4,
		swms.Step5, swms.Hazards5, swms.Risks5, swms.InitialRisk5,
		swms.ControlMeasures5, swms.ResidualRisk5, swms.ControlResponsibility5,
		swms.Step6, swms.Hazards6, swms.Risks6, swms.InitialRisk6,
		swms.ControlMeasures6, swms.ResidualRisk6, swms.ControlResponsibility6,
		swms.Step7, swms.Hazards7, swms.Risks7, swms.InitialRisk7,
		swms.ControlMeasures7, swms.ResidualRisk7, swms.ControlResponsibility7,
		swms.Step8, swms.Hazards8, swms.Risks8, swms.InitialRisk8,
		swms.ControlMeasures8, swms.ResidualRisk8, swms.ControlResponsibility8,
		swms.Step9, swms.Hazards9, swms.Risks9, swms.InitialRisk9,
		swms.ControlMeasures9, swms.ResidualRisk9, swms.ControlResponsibility9,
		swms.Step10, swms.Hazards10, swms.Risks10, swms.InitialRisk10,
		swms.ControlMeasures10, swms.ResidualRisk10, swms.ControlResponsibility10,
		swms.Step11, swms.Hazards11, swms.Risks11, swms.InitialRisk11,
		swms.ControlMeasures11, swms.ResidualRisk11, swms.ControlResponsibility11,
		swms.Step12, swms.Hazards12, swms.Risks12, swms.InitialRisk12,
		swms.ControlMeasures12, swms.ResidualRisk12, swms.ControlResponsibility12,

		swms.SignDate1, swms.SignName1, swms.SignSig1,
		swms.SignDate2, swms.SignName2, swms.SignSig2,
		swms.SignDate3, swms.SignName3, swms.SignSig3,
		swms.SignDate4, swms.SignName4, swms.SignSig4,
		swms.SignDate5, swms.SignName5, swms.SignSig5,
		swms.SignDate6, swms.SignName6, swms.SignSig6,
		swms.SignDate7, swms.SignName7, swms.SignSig7,
		swms.SignDate8, swms.SignName8, swms.SignSig8,
		swms.SignDate9, swms.SignName9, swms.SignSig9,
		swms.UUID,
	)

	if err != nil {
		log.Printf("failed to update swms: %v", err)
		return false, err
	}
	return true, nil
}

func DeleteSwms(uuid string, db *sql.DB) error {
	q := `
	delete from swms where uuid = ?;
	`
	_, err := db.Exec(q, uuid)
	if err != nil {
		return err
	}

	return nil
}
