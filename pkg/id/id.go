package golibid

import (
	"strings"

	"github.com/google/uuid"
)

type Id string

func (i Id) IsEmpty() bool {
	return i == ""
}

func (i Id) String() string {
	return string(i)
}

func NewId(prefix IdentifierPrefix) (Id, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return Id(""), err
	}
	if prefix == IdPrefixNone {
		return Id(strings.Replace(id.String(), "-", "", -1)), nil
	}
	return Id(prefix.String() + "_" + strings.Replace(id.String(), "-", "", -1)), nil
}

func (Id) GormDataType() string {
	return "varchar(50)"
}
