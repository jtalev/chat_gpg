package domain

type Swms struct {
	UUID                    string `json:"uuid"`
	JobId                   string `json:"job_id"`
	ProjectActivity         string `json:"project_activity"`
	ProjectNumber           string `json:"project_number"`
	SiteAddress             string `json:"site_address"`
	ContactName             string `json:"contact_name"`
	ContactNumber           string `json:"contact_number"`
	EmailAddress            string `json:"email_address"`
	SwmsDate                string `json:"swms_date"`
	HighRiskWorks           bool   `json:"high_risk_works"`
	SafetyGloves            bool   `json:"safety_gloves"`
	SafetyBoots             bool   `json:"safety_boots"`
	SafetyGlasses           bool   `json:"safety_glasses"`
	ProtectiveClothing      bool   `json:"protective_clothing"`
	RespiratoryProtection   bool   `json:"respiratory_protection"`
	HiVisClothing           bool   `json:"hi_vis_clothing"`
	SafetyHelmet            bool   `json:"safety_helmet"`
	FallArrest              bool   `json:"fall_arrest"`
	Other1                  string `json:"other_1"`
	Other2                  string `json:"other_2"`
	StepBelow2m             bool   `json:"step_below_2m"`
	StepAbove2m             bool   `json:"step_above_2m"`
	Scaffold                bool   `json:"scaffold"`
	PressureWasherDiesel    bool   `json:"pressure_washer_diesel"`
	RoofAnchorPoints        bool   `json:"roof_anchor_points"`
	ExtensionLadder         bool   `json:"extension_ladder"`
	ElectricScissorLift     bool   `json:"electric_scissor_lift"`
	DieselScissorLift       bool   `json:"diesel_scissor_lift"`
	ElectricKnuckleBoom     bool   `json:"electric_knuckle_boom"`
	DieselKnuckleBoom       bool   `json:"diesel_knuckle_boom"`
	AirlessSprayGun         bool   `json:"airless_spray_gun"`
	AngleGrinder            bool   `json:"angle_grinder"`
	Step1                   string `json:"step_1"`
	Hazards1                string `json:"hazards_1"`
	Risks1                  string `json:"risks_1"`
	InitialRisk1            string `json:"initial_risk_1"`
	ControlMeasures1        string `json:"control_measures_1"`
	ResidualRisk1           string `json:"residual_risk_1"`
	ControlResponsibility1  string `json:"control_responsibility_1"`
	Step2                   string `json:"step_2"`
	Hazards2                string `json:"hazards_2"`
	Risks2                  string `json:"risks_2"`
	InitialRisk2            string `json:"initial_risk_2"`
	ControlMeasures2        string `json:"control_measures_2"`
	ResidualRisk2           string `json:"residual_risk_2"`
	ControlResponsibility2  string `json:"control_responsibility_2"`
	Step3                   string `json:"step_3"`
	Hazards3                string `json:"hazards_3"`
	Risks3                  string `json:"risks_3"`
	InitialRisk3            string `json:"initial_risk_3"`
	ControlMeasures3        string `json:"control_measures_3"`
	ResidualRisk3           string `json:"residual_risk_3"`
	ControlResponsibility3  string `json:"control_responsibility_3"`
	Step4                   string `json:"step_4"`
	Hazards4                string `json:"hazards_4"`
	Risks4                  string `json:"risks_4"`
	InitialRisk4            string `json:"initial_risk_4"`
	ControlMeasures4        string `json:"control_measures_4"`
	ResidualRisk4           string `json:"residual_risk_4"`
	ControlResponsibility4  string `json:"control_responsibility_4"`
	Step5                   string `json:"step_5"`
	Hazards5                string `json:"hazards_5"`
	Risks5                  string `json:"risks_5"`
	InitialRisk5            string `json:"initial_risk_5"`
	ControlMeasures5        string `json:"control_measures_5"`
	ResidualRisk5           string `json:"residual_risk_5"`
	ControlResponsibility5  string `json:"control_responsibility_5"`
	Step6                   string `json:"step_6"`
	Hazards6                string `json:"hazards_6"`
	Risks6                  string `json:"risks_6"`
	InitialRisk6            string `json:"initial_risk_6"`
	ControlMeasures6        string `json:"control_measures_6"`
	ResidualRisk6           string `json:"residual_risk_6"`
	ControlResponsibility6  string `json:"control_responsibility_6"`
	Step7                   string `json:"step_7"`
	Hazards7                string `json:"hazards_7"`
	Risks7                  string `json:"risks_7"`
	InitialRisk7            string `json:"initial_risk_7"`
	ControlMeasures7        string `json:"control_measures_7"`
	ResidualRisk7           string `json:"residual_risk_7"`
	ControlResponsibility7  string `json:"control_responsibility_7"`
	Step8                   string `json:"step_8"`
	Hazards8                string `json:"hazards_8"`
	Risks8                  string `json:"risks_8"`
	InitialRisk8            string `json:"initial_risk_8"`
	ControlMeasures8        string `json:"control_measures_8"`
	ResidualRisk8           string `json:"residual_risk_8"`
	ControlResponsibility8  string `json:"control_responsibility_8"`
	Step9                   string `json:"step_9"`
	Hazards9                string `json:"hazards_9"`
	Risks9                  string `json:"risks_9"`
	InitialRisk9            string `json:"initial_risk_9"`
	ControlMeasures9        string `json:"control_measures_9"`
	ResidualRisk9           string `json:"residual_risk_9"`
	ControlResponsibility9  string `json:"control_responsibility_9"`
	Step10                  string `json:"step_10"`
	Hazards10               string `json:"hazards_10"`
	Risks10                 string `json:"risks_10"`
	InitialRisk10           string `json:"initial_risk_10"`
	ControlMeasures10       string `json:"control_measures_10"`
	ResidualRisk10          string `json:"residual_risk_10"`
	ControlResponsibility10 string `json:"control_responsibility_10"`
	Step11                  string `json:"step_11"`
	Hazards11               string `json:"hazards_11"`
	Risks11                 string `json:"risks_11"`
	InitialRisk11           string `json:"initial_risk_11"`
	ControlMeasures11       string `json:"control_measures_11"`
	ResidualRisk11          string `json:"residual_risk_11"`
	ControlResponsibility11 string `json:"control_responsibility_11"`
	Step12                  string `json:"step_12"`
	Hazards12               string `json:"hazards_12"`
	Risks12                 string `json:"risks_12"`
	InitialRisk12           string `json:"initial_risk_12"`
	ControlMeasures12       string `json:"control_measures_12"`
	ResidualRisk12          string `json:"residual_risk_12"`
	ControlResponsibility12 string `json:"control_responsibility_12"`
	SignDate1               string `json:"sign_date_1"`
	SignName1               string `json:"sign_name_1"`
	SignSig1                string `json:"sign_sig_1"`
	SignDate2               string `json:"sign_date_2"`
	SignName2               string `json:"sign_name_2"`
	SignSig2                string `json:"sign_sig_2"`
	SignDate3               string `json:"sign_date_3"`
	SignName3               string `json:"sign_name_3"`
	SignSig3                string `json:"sign_sig_3"`
	SignDate4               string `json:"sign_date_4"`
	SignName4               string `json:"sign_name_4"`
	SignSig4                string `json:"sign_sig_4"`
	SignDate5               string `json:"sign_date_5"`
	SignName5               string `json:"sign_name_5"`
	SignSig5                string `json:"sign_sig_5"`
	SignDate6               string `json:"sign_date_6"`
	SignName6               string `json:"sign_name_6"`
	SignSig6                string `json:"sign_sig_6"`
	SignDate7               string `json:"sign_date_7"`
	SignName7               string `json:"sign_name_7"`
	SignSig7                string `json:"sign_sig_7"`
	SignDate8               string `json:"sign_date_8"`
	SignName8               string `json:"sign_name_8"`
	SignSig8                string `json:"sign_sig_8"`
	SignDate9               string `json:"sign_date_9"`
	SignName9               string `json:"sign_name_9"`
	SignSig9                string `json:"sign_sig_9"`
	CreatedAt               string `json:"created_at"`
	ModifiedAt              string `json:"modified_at"`
}

type SwmsErrors struct {
	JobIdErr           string
	ProjectActivityErr string
	ProjectNumberErr   string
	SiteAddressErr     string
	ContactNameErr     string
	ContactNumberErr   string
	EmailAddressErr    string
	SwmsDateErr        string
	IsSuccessful       bool
	SuccessMsg         string
}

func (s *Swms) Validate() SwmsErrors {
	errors := SwmsErrors{IsSuccessful: true}
	errors = s.ValidateJobId(errors)
	errors = s.ValidateProjectActivity(errors)
	errors = s.ValidateProjectNumber(errors)
	errors = s.ValidateSiteAddress(errors)
	errors = s.ValidateContactName(errors)
	errors = s.ValidateContactNumber(errors)
	errors = s.ValidateEmailAddress(errors)
	errors = s.ValidateSwmsDate(errors)
	return errors
}

func (s *Swms) ValidateJobId(errors SwmsErrors) SwmsErrors {
	if s.JobId == "" {
		errors.JobIdErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (s *Swms) ValidateProjectActivity(errors SwmsErrors) SwmsErrors {
	if s.ProjectActivity == "" {
		errors.ProjectActivityErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (s *Swms) ValidateProjectNumber(errors SwmsErrors) SwmsErrors {
	if s.ProjectNumber == "" {
		errors.ProjectNumberErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (s *Swms) ValidateSiteAddress(errors SwmsErrors) SwmsErrors {
	if s.SiteAddress == "" {
		errors.SiteAddressErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (s *Swms) ValidateContactName(errors SwmsErrors) SwmsErrors {
	if s.ContactName == "" {
		errors.ContactNameErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (s *Swms) ValidateContactNumber(errors SwmsErrors) SwmsErrors {
	if s.ContactNumber == "" {
		errors.ContactNumberErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (s *Swms) ValidateEmailAddress(errors SwmsErrors) SwmsErrors {
	if s.EmailAddress == "" {
		errors.EmailAddressErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (s *Swms) ValidateSwmsDate(errors SwmsErrors) SwmsErrors {
	if s.SwmsDate == "" {
		errors.SwmsDateErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}
