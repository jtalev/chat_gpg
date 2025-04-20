package safety

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	application "github.com/jtalev/chat_gpg/application/services"
	models "github.com/jtalev/chat_gpg/domain/models"
	repo "github.com/jtalev/chat_gpg/infrastructure/repository"
)

type Swms struct {
	Swms     models.Swms
	SwmsArr  []models.Swms
	Errors   models.SwmsErrors
	UserRole string
	Db       *sql.DB
}

var defaultSwm = Swms{
	Swms: models.Swms{

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
	},
}

var p = application.Pdf{
	InPdfPath: "../ui/static/files/swm/safe_work_method_statement.pdf",

	OutPdfPath: "../ui/static/files/swm/safe_word_method_statement_output.pdf",
	OutPdfName: "safe_work_method_statement_output",

	S3StorageDir: "swms",

	JsonTemplatePath: "../ui/static/files/swm/swms_form_template.json",
	Data:             defaultSwm.Swms,
}

func (s *Swms) PostSwm(swms models.Swms) (models.SwmsErrors, error) {
	uuid := uuid.New().String()
	swms.UUID = uuid
	p.UUID = uuid
	errors := swms.Validate()
	if !errors.IsSuccessful {
		return errors, nil
	} else {
		_, err := repo.PostSwms(swms, s.Db)
		if err != nil {
			errors.IsSuccessful = false
			return errors, err
		}
		errors.SuccessMsg = "Swms submitted successfully."
		return errors, nil
	}
}

// GenerateSwmsPdf must be executed after PostSwm, p.UUID set during PostSwms execution
func (s *Swms) GenerateSwmsPdf(swms models.Swms) {
	p.Data = defaultSwm.Swms
	p.Data = p.WrapDataFieldText(p.Data, 32)
	err := p.ExecuteJsonTemplate()
	if err != nil {
		log.Println(err)
	}

	err = p.GeneratePdf()
	if err != nil {
		log.Println(err)
	}
}

func (s *Swms) GetSwms() error {
	s.SwmsArr = []models.Swms{}
	swmsArr, err := repo.GetSwms(s.Db)
	if err != nil {
		log.Printf("error getting swms: %v", err)
		return err
	}
	for _, swms := range swmsArr {
		temp := models.Swms{
			UUID:            swms.UUID,
			ProjectActivity: swms.ProjectActivity,
			ProjectNumber:   swms.ProjectNumber,
			SwmsDate:        swms.SwmsDate,
		}
		s.SwmsArr = append(s.SwmsArr, temp)
	}
	return nil
}
