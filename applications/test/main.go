package main

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"strings"
)

// CustomRenderer is a simple custom renderer to avoid escaping Markdown characters
type CustomRenderer struct {
	blackfriday.Renderer
}

func (r *CustomRenderer) NormalText(w *strings.Builder, text []byte) {
	w.WriteString(string(text))
}

var tt = `
📹 دانلود فیلم چیزی در انبار وجود دارد
🇮🇷 چیزی در انبار وجود دارد
🏴󠁧󠁢󠁥󠁮󠁧󠁿 There’s Something in the Barn
📝 داستان یک خانواده آمریکایی که کلبه‌ای دورافتاده در نروژ به ارث می‌برند و...
📽 2023
🗣  انگلیسی
💬 زیر نویس چسبیده دارد

\#️⃣  \#چیزی_در_انبار \#فیلم_خارجی
`

type Pa struct {
	name string
}

type TmpPa struct {
	Payload []*Pa
}

func main() {
	// متن Markdown شما

	tmp := new(TmpPa)

	tmp.Payload = nil

	if len(tmp.Payload) == 0 {
		fmt.Println("zero")
	}

	return

	output := escapeText(tt)

	// نمایش نتیجه
	fmt.Println(output)
}

func escapeText(text string) string {
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
