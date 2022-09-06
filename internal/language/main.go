package language

// Все языки должны иметь значения из этой таблицы:
// https://en.wikipedia.org/wiki/IETF_language_tag#List_of_common_primary_language_subtags
// TODO: Дозаполнить языки из спецификации.
const (
	// RU - русский.
	RU Lang = "ru"
	// EN - английский.
	EN Lang = "en"
	// Default - стандартный язык.
	Default = EN
)

type Lang string

// Known возвращает true в случае, если текущий язык является известным в
// проекте.
func (l Lang) Known() bool {
	switch l {
	case RU, EN:
		return true
	default:
		return false
	}
}
