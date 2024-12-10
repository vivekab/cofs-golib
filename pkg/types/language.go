package golibtypes

import (
	"github.com/vivekab/golib/protobuf/protoroot"
)

type Language string

const (
	LanguageEnglish   = Language("en")
	LanguageEnglishUS = Language("en_US")
	LanguageEnglishUK = Language("en_GB")
	LanguageRussian   = Language("ru_RU")
)

var LanguageFromProto = map[protoroot.Language]Language{
	protoroot.Language_LANGUAGE_EN:    LanguageEnglish,
	protoroot.Language_LANGUAGE_EN_US: LanguageEnglishUS,
}

var LanguageToProto = map[Language]protoroot.Language{
	LanguageEnglish:   protoroot.Language_LANGUAGE_EN,
	LanguageEnglishUS: protoroot.Language_LANGUAGE_EN_US,
}

func (v Language) IsValid(args map[string]string) *protoroot.GrpcError {
	return ValidateEnum(args, LanguageToProto, v)
}
