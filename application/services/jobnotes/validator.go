package jobnotes

type validatorFunc func() (interface{}, bool)

type paintnoteerrors struct {
	brandErr   string
	productErr string
	colourErr  string
	finishErr  string

	isSuccess  bool
	successMsg string
}

func (p *paintnote) validate() (paintnoteerrors, bool) {
	return paintnoteerrors{}, true
}

type tasknoteerrors struct {
	titleErr       string
	descriptionErr string
	statusErr      string

	isSuccess  bool
	successMsg string
}

func (t *tasknote) validate() (tasknoteerrors, bool) {
	return tasknoteerrors{}, true
}

type imagenoteerrors struct {
	isSuccess  bool
	successMsg string
}

func (i *imagenote) validate() (imagenoteerrors, bool) {
	errors := imagenoteerrors{
		isSuccess:  true,
		successMsg: "Image note submitted successfully",
	}
	return errors, errors.isSuccess
}
