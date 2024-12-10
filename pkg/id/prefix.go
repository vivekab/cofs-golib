package golibid

/*
 * User prefix + uuid will be used internally when passed between services and
 * removed when passing to/from the database.
 */

type IdentifierPrefix string

const (
	IdPrefixNone    = IdentifierPrefix("")
	IdPrefixRequest = IdentifierPrefix("req")
)

func (pr IdentifierPrefix) String() string {
	return string(pr)
}
