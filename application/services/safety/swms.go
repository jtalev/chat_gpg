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
	Jobs     []models.Job
	Db       *sql.DB
}

var p = application.Pdf{
	InPdfPath: "../ui/static/files/swm/safe_work_method_statement.pdf",

	OutPdfPath: "../ui/static/files/swm/safe_word_method_statement_output.pdf",
	OutPdfName: "safe_work_method_statement_output",

	S3StorageDir: "swms",

	JsonTemplatePath: "../ui/static/files/swm/swms_form_template.json",
}

func (s *Swms) PostSwm(swms models.Swms) (models.SwmsErrors, error) {
	uuid := uuid.New().String()
	swms.UUID = uuid
	p.UUID = uuid
	errors := swms.Validate()
	if !errors.IsSuccessful {
		return errors, nil
	} else {
		log.Println("posting swms")
		_, err := repo.PostSwms(swms, s.Db)
		if err != nil {
			errors.IsSuccessful = false
			return errors, err
		}
		errors.SuccessMsg = "Swms submitted successfully."
		return errors, nil
	}
}

func (s *Swms) PutSwms(swms models.Swms) (models.SwmsErrors, error) {
	p.UUID = swms.UUID
	errors := swms.Validate()
	if !errors.IsSuccessful {
		return errors, nil
	} else {
		log.Println("putting swms")
		_, err := repo.PutSwms(swms, s.Db)
		if err != nil {
			log.Printf("error updating swms: %v", err)
			return errors, err
		}
		err = p.Delete()
		if err != nil {
			log.Printf("error deleting swms pdf from s3: %v", err)
			return errors, err
		}
		errors.SuccessMsg = "Swms updated successfully."
		return errors, nil
	}
}

// GenerateSwmsPdf must be executed after PostSwm, p.UUID set during PostSwms execution
func (s *Swms) GenerateSwmsPdf(swms models.Swms) {
	p.Data = swms
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

func (s *Swms) GetUserRole(employeeId string) error {
	userRole, err := repo.GetEmployeeRole(employeeId, s.Db)
	if err != nil {
		log.Printf("error getting user role: %v", err)
		return err
	}
	s.UserRole = userRole.Role
	return nil
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

func (s *Swms) GetSwmsPdfUrl(uuid string) (string, error) {
	p.UUID = uuid
	p.S3FileName = uuid + "_" + p.OutPdfName + ".pdf"

	url, err := p.GetPresignedURL(100)
	if err != nil {
		log.Printf("error getting presigned pdf url: %v", err)
		return "", err
	}

	return url, nil
}

func (s *Swms) DeleteSwms() error {
	err := repo.DeleteSwms(s.Swms.UUID, s.Db)
	if err != nil {
		log.Printf("error deleting swms: %v", err)
		return err
	}
	p.UUID = s.Swms.UUID
	p.S3FileName = p.UUID + "_" + p.OutPdfName + ".pdf"
	err = p.Delete()
	if err != nil {
		log.Printf("error deleting swms pdf: %v", err)
		return err
	}
	return nil
}

func (s *Swms) GetSwmsByUUID() error {
	swms, err := repo.GetSwmsByUUID(s.Swms.UUID, s.Db)
	if err != nil {
		log.Printf("error fetching swms by UUID: %v", err)
		return err
	}
	s.Swms = swms
	return nil
}
