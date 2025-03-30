package domain

type IncidentReport struct {
	UUID                string `json:"uuid"`
	ReporterId          string `json:"reporter_id"`
	FullName            string `json:"full_name"`
	HomeAddress         string `json:"home_address"`
	ContactNumber       string `json:"contact_number"`
	IncidentDate        string `json:"incident_date"`
	IncidentTime        string `json:"incident_time"`
	PoliceNotified      string `json:"police_notified"`
	IncidentLocation    string `json:"incident_location"`
	IncidentDescription string `json:"incident_description"`
	WasWitnessed        string `json:"was_witnessed"`
	WasInjured          string `json:"was_injured"`
	FurtherDetails      string `json:"further_details"`
	WasTreated          string `json:"was_treated"`
	TreatmentLocation   string `json:"treatment_location"`
	IncInfoDate1        string `json:"inc_info_date_1"`
	IncInfoDate2        string `json:"inc_info_date_2"`
	IncInfoDate3        string `json:"inc_info_date_3"`
	IncInfoDate4        string `json:"inc_info_date_4"`
	IncInfoDate5        string `json:"inc_info_date_5"`
	IncInfoAction1      string `json:"inc_info_action_1"`
	IncInfoAction2      string `json:"inc_info_action_2"`
	IncInfoAction3      string `json:"inc_info_action_3"`
	IncInfoAction4      string `json:"inc_info_action_4"`
	IncInfoAction5      string `json:"inc_info_action_5"`
	IncInfoName1        string `json:"inc_info_name_1"`
	IncInfoName2        string `json:"inc_info_name_2"`
	IncInfoName3        string `json:"inc_info_name_3"`
	IncInfoName4        string `json:"inc_info_name_4"`
	IncInfoName5        string `json:"inc_info_name_5"`
	Reporter            string `json:"reporter"`
	Signature           string `json:"signature"`
	ReportDate          string `json:"report_date"`
	CreatedAt           string `json:"created_at"`
	ModifiedAt          string `json:"modified_at"`
}

type IncidentReportErrors struct {
	UUIDErr                string
	ReporterIdErr          string
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
	IsSuccessful           bool
}

func (i *IncidentReport) Validate() IncidentReportErrors {
	errors := IncidentReportErrors{
		IsSuccessful: true,
	}
	errors = i.validateUUID(errors)
	errors = i.validateReporterId(errors)
	errors = i.validateFullName(errors)
	errors = i.validateHomeAddress(errors)
	errors = i.validateContactNumber(errors)
	errors = i.validateIncidentDate(errors)
	errors = i.validateTime(errors)
	errors = i.validatePoliceNotified(errors)
	errors = i.validateIncidentLocation(errors)
	errors = i.validateIncidentDescription(errors)
	errors = i.validateWasWitnessed(errors)
	errors = i.validateWasInjured(errors)
	errors = i.validateWasTreated(errors)
	errors = i.validateTreatmentLocation(errors)
	errors = i.validateReporter(errors)
	errors = i.validateSignature(errors)
	errors = i.validateReportDate(errors)
	return errors
}

func (i *IncidentReport) validateUUID(errors IncidentReportErrors) IncidentReportErrors {
	if i.UUID == "" {
		errors.IsSuccessful = false
		errors.UUIDErr = "*required"
		return errors
	}
	return errors
}

func (i *IncidentReport) validateReporterId(errors IncidentReportErrors) IncidentReportErrors {
	if i.ReporterId == "" {
		errors.IsSuccessful = false
		errors.ReporterIdErr = "*required"
		return errors
	}
	return errors
}

func (i *IncidentReport) validateFullName(errors IncidentReportErrors) IncidentReportErrors {
	if i.FullName == "" {
		errors.IsSuccessful = false
		errors.FullNameErr = "*required"
		return errors
	}
	return errors
}

func (i *IncidentReport) validateHomeAddress(errors IncidentReportErrors) IncidentReportErrors {
	if i.HomeAddress == "" {
		errors.IsSuccessful = false
		errors.HomeAddressErr = "*required"
		return errors
	}
	return errors
}

func (i *IncidentReport) validateContactNumber(errors IncidentReportErrors) IncidentReportErrors {
	if i.ContactNumber == "" {
		errors.IsSuccessful = false
		errors.ContactNumberErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateIncidentDate(errors IncidentReportErrors) IncidentReportErrors {
	if i.IncidentDate == "" {
		errors.IsSuccessful = false
		errors.IncidentDateErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateTime(errors IncidentReportErrors) IncidentReportErrors {
	if i.IncidentTime == "" {
		errors.IsSuccessful = false
		errors.IncidentTimeErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validatePoliceNotified(errors IncidentReportErrors) IncidentReportErrors {
	if i.PoliceNotified == "" {
		errors.IsSuccessful = false
		errors.PoliceNotifiedErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateIncidentLocation(errors IncidentReportErrors) IncidentReportErrors {
	if i.IncidentLocation == "" {
		errors.IsSuccessful = false
		errors.IncidentLocationErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateIncidentDescription(errors IncidentReportErrors) IncidentReportErrors {
	if i.IncidentDescription == "" {
		errors.IsSuccessful = false
		errors.IncidentDescriptionErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateWasWitnessed(errors IncidentReportErrors) IncidentReportErrors {
	if i.WasWitnessed == "" {
		errors.IsSuccessful = false
		errors.WasWitnessedErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateWasInjured(errors IncidentReportErrors) IncidentReportErrors {
	if i.WasInjured == "" {
		errors.IsSuccessful = false
		errors.WasInjuredErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateWasTreated(errors IncidentReportErrors) IncidentReportErrors {
	if i.WasTreated == "" {
		errors.IsSuccessful = false
		errors.WasTreatedErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateTreatmentLocation(errors IncidentReportErrors) IncidentReportErrors {
	if i.TreatmentLocation == "" {
		errors.IsSuccessful = false
		errors.TreatmentLocationErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateReporter(errors IncidentReportErrors) IncidentReportErrors {
	if i.Reporter == "" {
		errors.IsSuccessful = false
		errors.ReporterErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateSignature(errors IncidentReportErrors) IncidentReportErrors {
	if i.Signature == "" {
		errors.IsSuccessful = false
		errors.SignatureErr = "*required"
		return errors
	}
	return errors
}
func (i *IncidentReport) validateReportDate(errors IncidentReportErrors) IncidentReportErrors {
	if i.ReportDate == "" {
		errors.IsSuccessful = false
		errors.ReportDateErr = "*required"
		return errors
	}
	return errors
}
