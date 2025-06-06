package safety

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	services "github.com/jtalev/chat_gpg/application/services"
	domain "github.com/jtalev/chat_gpg/domain/models"
	infrastructure "github.com/jtalev/chat_gpg/infrastructure/repository"
)

type IncidentReportValues struct {
	UUID                   string `json:"uuid"`
	ReporterId             string `json:"reporter_id"`
	FullName               string `json:"full_name"`
	HomeAddress            string `json:"home_address"`
	ContactNumber          string `json:"contact_number"`
	IncidentDate           string `json:"incident_date"`
	IncidentTime           string `json:"incident_time"`
	PoliceNotified         string `json:"police_notified"`
	IncidentLocation       string `json:"incident_location"`
	IncidentDescription    string `json:"incident_description"`
	WasWitnessed           string `json:"was_witnessed"`
	WasInjured             string `json:"was_injured"`
	FurtherDetails         string `json:"further_details"`
	WasTreated             string `json:"was_treated"`
	TreatmentLocation      string `json:"treatment_location"`
	IncInfoDate1           string `json:"inc_info_date_1"`
	IncInfoDate2           string `json:"inc_info_date_2"`
	IncInfoDate3           string `json:"inc_info_date_3"`
	IncInfoDate4           string `json:"inc_info_date_4"`
	IncInfoDate5           string `json:"inc_info_date_5"`
	IncInfoAction1         string `json:"inc_info_action_1"`
	IncInfoAction2         string `json:"inc_info_action_2"`
	IncInfoAction3         string `json:"inc_info_action_3"`
	IncInfoAction4         string `json:"inc_info_action_4"`
	IncInfoAction5         string `json:"inc_info_action_5"`
	IncInfoName1           string `json:"inc_info_name_1"`
	IncInfoName2           string `json:"inc_info_name_2"`
	IncInfoName3           string `json:"inc_info_name_3"`
	IncInfoName4           string `json:"inc_info_name_4"`
	IncInfoName5           string `json:"inc_info_name_5"`
	Reporter               string `json:"reporter"`
	Signature              string `json:"signature"`
	ReportDate             string `json:"report_date"`
	FullNameErr            string
	HomeAddressErr         string
	ContactNumberErr       string
	IncidentDateErr        string
	IncidentTimeErr        string
	PoliceNotifiedErr      string
	IncidentLocationErr    string
	IncidentDescriptionErr string
	WasWitnessedErr        string
	WasInjuredErr          string
	WasTreatedErr          string
	TreatmentLocationErr   string
	ReporterErr            string
	SignatureErr           string
	ReportDateErr          string
	SuccessMsg             string
}

const (
	inIncidentReportPath       = "../ui/static/files/incident_report/incident_report_forms.pdf"
	inIncidentReportJsonPath   = "../ui/static/files/incident_report/incident_form_data.json"
	outIncidentReportPath      = "../ui/static/files/incident_report/incident_report_output.pdf"
	outIncidentReportPdfName   = "incident_report_output"
	incidentReportS3StorageDir = "incident_report"
)

func mapIncidentReportValuesToIncidentReport(incidentReportValues IncidentReportValues) domain.IncidentReport {
	return domain.IncidentReport{
		UUID:                incidentReportValues.UUID,
		ReporterId:          incidentReportValues.ReporterId,
		FullName:            incidentReportValues.FullName,
		HomeAddress:         incidentReportValues.HomeAddress,
		ContactNumber:       incidentReportValues.ContactNumber,
		IncidentDate:        incidentReportValues.IncidentDate,
		IncidentTime:        incidentReportValues.IncidentTime,
		PoliceNotified:      incidentReportValues.PoliceNotified,
		IncidentLocation:    incidentReportValues.IncidentLocation,
		IncidentDescription: incidentReportValues.IncidentDescription,
		WasWitnessed:        incidentReportValues.WasWitnessed,
		WasInjured:          incidentReportValues.WasInjured,
		FurtherDetails:      incidentReportValues.FurtherDetails,
		WasTreated:          incidentReportValues.WasTreated,
		TreatmentLocation:   incidentReportValues.TreatmentLocation,
		IncInfoDate1:        incidentReportValues.IncInfoDate1,
		IncInfoDate2:        incidentReportValues.IncInfoDate2,
		IncInfoDate3:        incidentReportValues.IncInfoDate3,
		IncInfoDate4:        incidentReportValues.IncInfoDate4,
		IncInfoDate5:        incidentReportValues.IncInfoDate5,
		IncInfoAction1:      incidentReportValues.IncInfoAction1,
		IncInfoAction2:      incidentReportValues.IncInfoAction2,
		IncInfoAction3:      incidentReportValues.IncInfoAction3,
		IncInfoAction4:      incidentReportValues.IncInfoAction4,
		IncInfoAction5:      incidentReportValues.IncInfoAction5,
		IncInfoName1:        incidentReportValues.IncInfoName1,
		IncInfoName2:        incidentReportValues.IncInfoName2,
		IncInfoName3:        incidentReportValues.IncInfoName3,
		IncInfoName4:        incidentReportValues.IncInfoName4,
		IncInfoName5:        incidentReportValues.IncInfoName5,
		Reporter:            incidentReportValues.Reporter,
		Signature:           incidentReportValues.Signature,
		ReportDate:          incidentReportValues.ReportDate,
	}
}

func mapIncidentReportToIncidentReportValues(incidentReport domain.IncidentReport) IncidentReportValues {
	return IncidentReportValues{
		UUID:                incidentReport.UUID,
		ReporterId:          incidentReport.ReporterId,
		FullName:            incidentReport.FullName,
		HomeAddress:         incidentReport.HomeAddress,
		ContactNumber:       incidentReport.ContactNumber,
		IncidentDate:        incidentReport.IncidentDate,
		IncidentTime:        incidentReport.IncidentTime,
		PoliceNotified:      incidentReport.PoliceNotified,
		IncidentLocation:    incidentReport.IncidentLocation,
		IncidentDescription: incidentReport.IncidentDescription,
		WasWitnessed:        incidentReport.WasWitnessed,
		WasInjured:          incidentReport.WasInjured,
		FurtherDetails:      incidentReport.FurtherDetails,
		WasTreated:          incidentReport.WasTreated,
		TreatmentLocation:   incidentReport.TreatmentLocation,
		IncInfoDate1:        incidentReport.IncInfoDate1,
		IncInfoDate2:        incidentReport.IncInfoDate2,
		IncInfoDate3:        incidentReport.IncInfoDate3,
		IncInfoDate4:        incidentReport.IncInfoDate4,
		IncInfoDate5:        incidentReport.IncInfoDate5,
		IncInfoAction1:      incidentReport.IncInfoAction1,
		IncInfoAction2:      incidentReport.IncInfoAction2,
		IncInfoAction3:      incidentReport.IncInfoAction3,
		IncInfoAction4:      incidentReport.IncInfoAction4,
		IncInfoAction5:      incidentReport.IncInfoAction5,
		IncInfoName1:        incidentReport.IncInfoName1,
		IncInfoName2:        incidentReport.IncInfoName2,
		IncInfoName3:        incidentReport.IncInfoName3,
		IncInfoName4:        incidentReport.IncInfoName4,
		IncInfoName5:        incidentReport.IncInfoName5,
		Reporter:            incidentReport.Reporter,
		Signature:           incidentReport.Signature,
		ReportDate:          incidentReport.ReportDate,
	}
}

func mapToPdf(inPdfPath, inJsonPath, outPdfPath, outPdfName, s3StorageDir, uuid string, data any) services.Pdf {
	return services.Pdf{
		InPdfPath:        inPdfPath,
		JsonTemplatePath: inJsonPath,
		OutPdfPath:       outPdfPath,
		OutPdfName:       outPdfName,
		S3StorageDir:     s3StorageDir,
		UUID:             uuid,
		Data:             data,
	}
}

func postIncidentReport(incidentReport domain.IncidentReport, db *sql.DB) (sql.Result, error) {
	return infrastructure.PostIncidentReport(incidentReport, db)
}

func mapErrorsToIncidentReportValues(incidentReportValues IncidentReportValues, errors domain.IncidentReportErrors) IncidentReportValues {
	incidentReportValues.FullNameErr = errors.FullNameErr
	incidentReportValues.HomeAddressErr = errors.HomeAddressErr
	incidentReportValues.ContactNumberErr = errors.ContactNumberErr
	incidentReportValues.IncidentDateErr = errors.IncidentDateErr
	incidentReportValues.IncidentTimeErr = errors.IncidentTimeErr
	incidentReportValues.PoliceNotifiedErr = errors.PoliceNotifiedErr
	incidentReportValues.IncidentLocationErr = errors.IncidentLocationErr
	incidentReportValues.IncidentDescriptionErr = errors.IncidentDescriptionErr
	incidentReportValues.WasWitnessedErr = errors.WasWitnessedErr
	incidentReportValues.WasInjuredErr = errors.WasInjuredErr
	incidentReportValues.WasTreatedErr = errors.WasTreatedErr
	incidentReportValues.TreatmentLocationErr = errors.TreatmentLocationErr
	incidentReportValues.ReporterErr = errors.ReporterErr
	incidentReportValues.SignatureErr = errors.SignatureErr
	incidentReportValues.ReportDateErr = errors.ReportDateErr
	return incidentReportValues
}

func wrapText(text string, limit int) string {
	var result string
	words := strings.Fields(text)
	lineLen := 0

	for _, word := range words {
		if lineLen+len(word) > limit {
			result += `\n`
			lineLen = 0
		}
		result += word + " "
		lineLen += len(word) + 1
	}
	return strings.TrimSpace(result)
}

func GenerateIncidentReportPdf(incidentReportValues IncidentReportValues, db *sql.DB) (IncidentReportValues, error) {
	uuid := uuid.New().String()

	incidentReport := mapIncidentReportValuesToIncidentReport(incidentReportValues)
	incidentReport.UUID = uuid
	errors := incidentReport.Validate()
	if !errors.IsSuccessful {
		incidentReportValues = mapErrorsToIncidentReportValues(incidentReportValues, errors)
		return incidentReportValues, nil
	}
	incidentReport.FurtherDetails = wrapText(incidentReport.FurtherDetails, 100)
	incidentReport.IncidentDescription = wrapText(incidentReport.IncidentDescription, 100)
	result, err := postIncidentReport(incidentReport, db)
	if err != nil {
		log.Printf("error posting incident report: %v", result)
		return incidentReportValues, err
	}

	p := mapToPdf(
		inIncidentReportPath,
		inIncidentReportJsonPath,
		outIncidentReportPath,
		outIncidentReportPdfName,
		incidentReportS3StorageDir,
		uuid,
		incidentReport,
	)

	err = p.GeneratePdf()
	if err != nil {
		log.Printf("Error generating incident report pdf: %v", err)
		return incidentReportValues, err
	}
	incidentReportValues.SuccessMsg = "Incident report submitted successfully."
	return incidentReportValues, nil
}

func GetIncidentReportUrl(uuid string, db *sql.DB) (string, error) {
	p := services.Pdf{
		UUID:         uuid,
		S3FileName:   uuid + "_" + outIncidentReportPdfName + ".pdf",
		S3StorageDir: incidentReportS3StorageDir,
	}

	url, err := p.GetPresignedURL(100)
	if err != nil {
		return "", err
	}

	return url, nil
}

func DeleteIncidentReport(uuid string, db *sql.DB) error {
	err := infrastructure.DeleteIncidentReport(uuid, db)
	if err != nil {
		log.Printf("error deleting incident report from db: %v", err)
		return err
	}

	p := services.Pdf{
		UUID:         uuid,
		S3FileName:   uuid + "_" + outIncidentReportPdfName + ".pdf",
		S3StorageDir: incidentReportS3StorageDir,
	}

	err = p.Delete()
	if err != nil {
		log.Printf("error deleting incident report pdf from s3 bucket: %v", err)
		return err
	}

	return nil
}

func GetIncidentReport(uuid string, db *sql.DB) (IncidentReportValues, error) {
	incidentReport, err := infrastructure.GetIncidentReport(uuid, db)
	if err != nil {
		return IncidentReportValues{}, err
	}

	incidentReportValues := mapIncidentReportToIncidentReportValues(incidentReport)

	return incidentReportValues, nil
}

func PutIncidentReport(incidentReportValues IncidentReportValues, db *sql.DB) (IncidentReportValues, error) {
	incidentReport := mapIncidentReportValuesToIncidentReport(incidentReportValues)
	errors := incidentReport.Validate()
	if errors.IsSuccessful == false {
		incidentReportValues = mapErrorsToIncidentReportValues(incidentReportValues, errors)
		return incidentReportValues, nil
	}

	_, err := infrastructure.PutIncidentReport(incidentReport, db)
	if err != nil {
		log.Printf("error updating incident report at db: %v", err)
		return incidentReportValues, err
	}

	p := mapToPdf(
		inIncidentReportPath,
		inIncidentReportJsonPath,
		outIncidentReportPath,
		outIncidentReportPdfName,
		incidentReportS3StorageDir,
		incidentReport.UUID,
		incidentReport,
	)
	p.S3FileName = fmt.Sprintf("%s_%s.pdf", incidentReport.UUID, p.OutPdfName)

	err = p.Delete()
	if err != nil {
		log.Printf("error deleting incident report from s3 bucket: %v", err)
		return incidentReportValues, err
	}
	err = p.GeneratePdf()
	if err != nil {
		log.Printf("error generating pdf from updated incident report: %v", err)
		return incidentReportValues, err
	}

	incidentReportValues.SuccessMsg = "Successfully updated incident report."
	return incidentReportValues, nil
}
