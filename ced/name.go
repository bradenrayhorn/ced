package ced

type Name ValidatableString

func (n Name) Validate() error {
	s := ValidatableString(n)
	return validate(Required(s), Max(s, 255, "characters"))
}
