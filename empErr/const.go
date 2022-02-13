package empErr

const (
	MustError                       Identifier = "MustError"
	NotAllowEmptyEnvError           Identifier = "NotAllowEmptyEnvError"
	CannotParseEnvStringToTypeError Identifier = "CannotParseEnvStringToTypeError"
	UnsupportedTypeError            Identifier = "UnsupportedTypeError"
	ArraySizeMismatchError          Identifier = "ArraySizeMismatchError"
)

var ErrorMap = map[Identifier]*Error{
	MustError: {
		Identifier: MustError,
	},
	NotAllowEmptyEnvError: {
		Identifier: NotAllowEmptyEnvError,
	},
	CannotParseEnvStringToTypeError: {
		Identifier: CannotParseEnvStringToTypeError,
	},
	UnsupportedTypeError: {
		Identifier: UnsupportedTypeError,
	},
	ArraySizeMismatchError: {
		Identifier: ArraySizeMismatchError,
	},
}
