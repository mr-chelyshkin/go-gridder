package config

type FileConfigOpt func(*fileConfig)

// WithFileType rewrite FileParse file type.
func WithFileType(e string) FileConfigOpt {
	return func(fc *fileConfig) {
		fc.ext = e
	}
}

// WithValidator set object for check data after parse.
func WithValidator(v validator) FileConfigOpt {
	return func(fc *fileConfig) {
		fc.validator = v
	}
}
