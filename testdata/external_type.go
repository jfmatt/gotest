package testdata

// ExternalType is a type defined outside the test package.
// It has an unexported field that should NOT be compared when using Eq
// from another package.
type ExternalType struct {
	PublicField  string
	privateField string
}

// NewExternalType creates an ExternalType with both public and private fields set.
func NewExternalType(public, private string) ExternalType {
	return ExternalType{
		PublicField:  public,
		privateField: private,
	}
}
