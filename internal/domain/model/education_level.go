package model

type EducationLevel string

const (
	EducationLevelBacholor EducationLevel = "学部生"
	EducationLevelMaster   EducationLevel = "大学院生"
	EducationLevelOther    EducationLevel = "その他"
)

func (m EducationLevel) Str() string {
	return string(m)
}

func (m EducationLevel) Validate() bool {
	if m != EducationLevelBacholor && m != EducationLevelMaster && m != EducationLevelOther {
		return false
	}
	return true
}
