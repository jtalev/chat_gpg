package application

import (
	"os"
	"strings"
	"text/template"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func ExecuteJsonTemplate(inJsonPath string, data any) (string, error) {
	inJsonFile, err := os.ReadFile(inJsonPath)
	if err != nil {
		return "", err
	}

	outJsonName := strings.Split(strings.Split(inJsonPath, "/")[len(strings.Split(inJsonPath, "/"))-1], ".")[0]

	dotIdx := len(inJsonPath) - 5
	outJsonPath := inJsonPath[:dotIdx] + "-output" + inJsonPath[dotIdx:]
	outJsonFile, err := os.Create(outJsonPath)
	if err != nil {
		return "", err
	}
	defer outJsonFile.Close()

	t := template.Must(template.New(outJsonName).Parse(string(inJsonFile)))
	err = t.Execute(outJsonFile, data)
	if err != nil {
		return "", err
	}

	return outJsonFile.Name(), nil
}

func GeneratePdf(inPdfPath, inJsonPath, outPdfPath string, data any) error {
	templatedJsonPath, err := ExecuteJsonTemplate(inJsonPath, data)
	if err != nil {
		return err
	}

	err = api.FillFormFile(inPdfPath, templatedJsonPath, outPdfPath, nil)
	if err != nil {
		return err
	}
	return nil
}
