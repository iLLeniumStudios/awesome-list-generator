package utils

type StringList []string

func (l StringList) Contains(name string) bool {
	for _, item := range l {
		if item == name {
			return true
		}
	}
	return false
}
