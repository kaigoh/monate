package scalar

import (
	"fmt"
	"io"

	"github.com/99designs/gqlgen/graphql"
)

// Long represents a 64-bit signed integer GraphQL scalar.
type Long int64

// MarshalGQL writes the scalar value.
func (l Long) MarshalGQL(w io.Writer) {
	graphql.MarshalInt64(int64(l)).MarshalGQL(w)
}

// UnmarshalGQL converts input values to Long.
func (l *Long) UnmarshalGQL(v any) error {
	if v == nil {
		*l = 0
		return nil
	}
	value, err := graphql.UnmarshalInt64(v)
	if err != nil {
		return fmt.Errorf("Long: %w", err)
	}
	*l = Long(value)
	return nil
}
