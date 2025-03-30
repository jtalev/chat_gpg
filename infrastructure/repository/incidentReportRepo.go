package infrastructure

import (
	"database/sql"

	domain "github.com/jtalev/chat_gpg/domain/models"
)

func GetIncidentReports(db *sql.DB) ([]domain.IncidentReport, error) {
	q := `
	select * from incident_report;
	`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	outIncidentReports := []domain.IncidentReport{}
	for rows.Next() {
		incidentReport := domain.IncidentReport{}
		if err := rows.Scan(
			&incidentReport.UUID,
			&incidentReport.ReporterId,
			&incidentReport.FullName,
			&incidentReport.HomeAddress,
			&incidentReport.ContactNumber,
			&incidentReport.IncidentDate,
			&incidentReport.IncidentTime,
			&incidentReport.PoliceNotified,
			&incidentReport.IncidentLocation,
			&incidentReport.IncidentDescription,
			&incidentReport.WasWitnessed,
			&incidentReport.WasInjured,
			&incidentReport.FurtherDetails,
			&incidentReport.WasTreated,
			&incidentReport.TreatmentLocation,
			&incidentReport.IncInfoDate1,
			&incidentReport.IncInfoDate2,
			&incidentReport.IncInfoDate3,
			&incidentReport.IncInfoDate4,
			&incidentReport.IncInfoDate5,
			&incidentReport.IncInfoAction1,
			&incidentReport.IncInfoAction2,
			&incidentReport.IncInfoAction3,
			&incidentReport.IncInfoAction4,
			&incidentReport.IncInfoAction5,
			&incidentReport.IncInfoName1,
			&incidentReport.IncInfoName2,
			&incidentReport.IncInfoName3,
			&incidentReport.IncInfoName4,
			&incidentReport.IncInfoName5,
			&incidentReport.Reporter,
			&incidentReport.Signature,
			&incidentReport.ReportDate,
			&incidentReport.CreatedAt,
			&incidentReport.ModifiedAt,
		); err != nil {
			return nil, err
		}
		outIncidentReports = append(outIncidentReports, incidentReport)
	}
	return outIncidentReports, nil
}

func PostIncidentReport(incidentReport domain.IncidentReport, db *sql.DB) (sql.Result, error) {
	q := `
	INSERT INTO incident_report (uuid, reporter_id, full_name, home_address, contact_number, incident_date,
		incident_time, police_notified, incident_location, incident_description, was_witnessed, 
		was_injured, further_details, was_treated, treatment_location, inc_info_date_1, inc_info_date_2, 
		inc_info_date_3, inc_info_date_4, inc_info_date_5, inc_info_action_1, inc_info_action_2, inc_info_action_3, 
		inc_info_action_4, inc_info_action_5, inc_info_name_1, inc_info_name_2, inc_info_name_3, 
		inc_info_name_4, inc_info_name_5, reporter, signature, report_date)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, 
		$19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33);
	`

	result, err := db.Exec(
		q,
		incidentReport.UUID,
		incidentReport.ReporterId,
		incidentReport.FullName,
		incidentReport.HomeAddress,
		incidentReport.ContactNumber,
		incidentReport.IncidentDate,
		incidentReport.IncidentTime,
		incidentReport.PoliceNotified,
		incidentReport.IncidentLocation,
		incidentReport.IncidentDescription,
		incidentReport.WasWitnessed,
		incidentReport.WasInjured,
		incidentReport.FurtherDetails,
		incidentReport.WasTreated,
		incidentReport.TreatmentLocation,
		incidentReport.IncInfoDate1,
		incidentReport.IncInfoDate2,
		incidentReport.IncInfoDate3,
		incidentReport.IncInfoDate4,
		incidentReport.IncInfoDate5,
		incidentReport.IncInfoAction1,
		incidentReport.IncInfoAction2,
		incidentReport.IncInfoAction3,
		incidentReport.IncInfoAction4,
		incidentReport.IncInfoAction5,
		incidentReport.IncInfoName1,
		incidentReport.IncInfoName2,
		incidentReport.IncInfoName3,
		incidentReport.IncInfoName4,
		incidentReport.IncInfoName5,
		incidentReport.Reporter,
		incidentReport.Signature,
		incidentReport.ReportDate,
	)
	if err != nil {
		return result, err
	}
	return result, nil
}

func DeleteIncidentReport(uuid string, db *sql.DB) error {
	q := `
	delete from incident_report where uuid = ?;
	`

	_, err := db.Exec(q, uuid)
	if err != nil {
		return err
	}

	return nil
}
