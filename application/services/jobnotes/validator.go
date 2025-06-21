package jobnotes

type validatorFunc func() (interface{}, bool)

type paintnoteerrors struct {
	BrandErr   string
	ProductErr string
	ColourErr  string
	FinishErr  string

	IsSuccess  bool
	SuccessMsg string
}

func (p *paintnote) validate() (paintnoteerrors, bool) {
	return paintnoteerrors{BrandErr: "*required", IsSuccess: false}, false
}

type tasknoteerrors struct {
	TitleErr       string
	DescriptionErr string
	StatusErr      string

	IsSuccess  bool
	SuccessMsg string
}

func (t *tasknote) validate() (tasknoteerrors, bool) {
	return tasknoteerrors{}, true
}

type imagenoteerrors struct {
	IsSuccess  bool
	SuccessMsg string
}

func (i *imagenote) validate() (imagenoteerrors, bool) {
	errors := imagenoteerrors{
		IsSuccess:  true,
		SuccessMsg: "Image note submitted successfully",
	}
	return errors, errors.IsSuccess
}
