package jobnotes

type validatorFunc func() (interface{}, bool)

type paintnoteerrors struct {
	BrandErr    string
	ProductErr  string
	ColourErr   string
	FinishErr   string
	AreaErr     string
	SurfacesErr string

	IsSuccess  bool
	SuccessMsg string
}

func (p *Paintnote) validateBrand(errors paintnoteerrors) paintnoteerrors {
	if p.Brand == "" {
		errors.BrandErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (p *Paintnote) validateColour(errors paintnoteerrors) paintnoteerrors {
	if p.Colour == "" {
		errors.ColourErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (p *Paintnote) validateProduct(errors paintnoteerrors) paintnoteerrors {
	if p.Product == "" {
		errors.ProductErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (p *Paintnote) validateFinish(errors paintnoteerrors) paintnoteerrors {
	if p.Finish == "" {
		errors.FinishErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (p *Paintnote) validateArea(errors paintnoteerrors) paintnoteerrors {
	if p.Area == "" {
		errors.AreaErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (p *Paintnote) validateSurfaces(errors paintnoteerrors) paintnoteerrors {
	if p.Surfaces == "" {
		errors.SurfacesErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (p *Paintnote) validate() (paintnoteerrors, bool) {
	errors := paintnoteerrors{IsSuccess: true}

	errors = p.validateBrand(errors)
	errors = p.validateProduct(errors)
	errors = p.validateColour(errors)
	errors = p.validateFinish(errors)
	errors = p.validateArea(errors)
	errors = p.validateSurfaces(errors)

	return errors, errors.IsSuccess
}

type tasknoteerrors struct {
	TitleErr       string
	DescriptionErr string
	StatusErr      string

	IsSuccess  bool
	SuccessMsg string
}

func (t *Tasknote) validate() (tasknoteerrors, bool) {
	return tasknoteerrors{}, true
}

type imagenoteerrors struct {
	ImageErr   string
	CaptionErr string
	AreaErr    string

	IsSuccess  bool
	SuccessMsg string
}

func (i *Imagenote) validate() (imagenoteerrors, bool) {
	errors := imagenoteerrors{IsSuccess: true}

	errors = i.validateImage(errors)
	errors = i.validateCaption(errors)
	errors = i.validateArea(errors)

	return errors, errors.IsSuccess
}

func (i *Imagenote) validateImage(errors imagenoteerrors) imagenoteerrors {
	if i.ImageBase64 == "" {
		errors.ImageErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (i *Imagenote) validateCaption(errors imagenoteerrors) imagenoteerrors {
	if i.Caption == "" {
		errors.CaptionErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}

func (i *Imagenote) validateArea(errors imagenoteerrors) imagenoteerrors {
	if i.Area == "" {
		errors.AreaErr = "*required"
		errors.IsSuccess = false
		return errors
	}
	return errors
}
