package safety

import (
	"log"

	application "github.com/jtalev/chat_gpg/application/services"
)

type Swm struct {
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
}

var defaultSwm = Swm{
	HighRiskWorks:          true,
	SafetyBoots:            true,
	HiVisClothing:          true,
	StepBelow2m:            true,
	Step1:                  "Site Set Up",
	Hazards1:               "Muscle strains/ additional tripping hazards",
	Risks1:                 "Muscle strain and sprains from incorrect lifting technique/ tripping hazard over site storage position",
	InitialRisk1:           "B3",
	ControlMeasures1:       "Warm Up stretches, correct lifting technique, 2 person lifts when required/ preplanning to find suitable storage area for equipment",
	ResidualRisk1:          "A2",
	ControlResponsibility1: "Project Manager Foreman Painters",
	Step2:                  "Working on Ladders",
	Hazards2:               "Fall Hazard",
	Risks2:                 "Potential serious injury from falling from height of ladder",
	InitialRisk2:           "D3",
	ControlMeasures2:       "Use of platform steps whenever possible, no overreaching of ladder, where possible utilise mobile equipment or plant",
	ResidualRisk2:          "B2",
	ControlResponsibility2: "Project Manager Foreman Painters",
	Step3:                  "Working in Public Spaces",
	Hazards3:               "Public Safety",
	Risks3:                 "Risk of injury or harm to the general public",
	InitialRisk3:           "C2",
	ControlMeasures3:       "Isolate working area with the use of bollards, hard barriers and signage",
	ResidualRisk3:          "B1",
	ControlResponsibility3: "Project Manager Foreman Painters",
	Step4:                  "General Preparation works",
	Hazards4:               "Flying debris from sanding/ respiratory damage from sanding of paint and plaster",
	Risks4:                 "Eye damage from flying debris while undertaking preparation works/ Respiratory damage from paint/plaster dust",
	InitialRisk4:           "B4",
	ControlMeasures4:       "Use of correct PPE while undertaking preparation works including safety glasses/ P2 rated mask if required",
	ResidualRisk4:          "B1",
	ControlResponsibility4: "Project Manager Foreman Painters",
	Step5:                  "General Painting Works",
	Hazards5:               "Muscle strains, eye damage, fume inhalation",
	Risks5:                 "Muscle strains from overuse/ flying debris when working overhead/ fume inhalation",
	InitialRisk5:           "B4",
	ControlMeasures5:       "Take breaks when required and continue stretches throughout the day, use of safety glasses working overhead/ P2 mask when required",
	ResidualRisk5:          "B2",
	ControlResponsibility5: "Project Manager Foreman Painters",
	ContactNumber:          "054216133",
}

var p = application.Pdf{
	InPdfPath:        "../ui/static/files/swm/safe_work_method_statement.pdf",
	OutPdfPath:       "../ui/static/files/swm/safe_word_method_statement_output.pdf",
	JsonTemplatePath: "../ui/static/files/swm/swms_form_template.json",
	OutPdfName:       "safe_work_method_statement_output",
	Data:             defaultSwm,
}

func (s *Swm) GenerateSwmPdf() {
	err := p.GenerateJsonTemplate()
	if err != nil {
		log.Println(err)
	}
	err = p.FormatJsonTemplate()
	if err != nil {
		log.Println(err)
	}
	err = p.ExecuteJsonTemplate()
	if err != nil {
		log.Println(err)
	}

	err = p.GeneratePdf()
	if err != nil {
		log.Println(err)
	}
}
