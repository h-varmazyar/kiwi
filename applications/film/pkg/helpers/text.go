package helpers

import "fmt"

func EscapeText(text string) string {
	escaped := ""
	for _, c := range text {
		switch c {
		case '_', '*', '[', ']', '(', ')', '~', '`', '>',
			'#', '+', '-', '=', '|', '{', '}', '.', '!':
			escaped = fmt.Sprintf("%v\\%c", escaped, c)
		default:
			escaped = fmt.Sprintf("%v%c", escaped, c)
		}
	}
	return escaped
}
