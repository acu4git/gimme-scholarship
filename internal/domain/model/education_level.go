package model

type EducationLevel string

func (m EducationLevel) String() string {
	return string(m)
}

func (m EducationLevel) Validate() bool {
	var res bool
	switch m.String() {
	case "学部生", "大学院生", "その他":
		res = true
	default:
		res = false
	}
	return res
}
