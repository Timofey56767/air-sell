package internal


var (
	RegistrText = "для того чтобы продолжить войдите в систему"
)

type Language string

type TransalteTexts struct {
	HelloText string
	SignInText string
	SignUpText string

}

var Translate = map[Language]TransalteTexts{
	"ru": {
		HelloText: "Добро Пожаловать",
		SignInText: "для того чтобы продолжить войдите в систему",
	},
}