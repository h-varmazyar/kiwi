package handlers

import "fmt"

func respProxyCaption(proxyLinks []string) string {
	caption := `
📌سرور پرسرعت کانال⭐️🎁

🔺مخصوص تمام نت‌ها ☄️
🔺مناسب موبایل و لپ‌تاپ📱🖥
🔺مخصوص اختلالات اخیر⚡️🔐

🔸🔹🔸🔹🔸🔹🔸🔹🔸🔹🔸🔹
🛜 پروکسی ضد فیلتر و پر سرعت

%v

🔸🔹🔸🔹🔸🔹🔸🔹🔸🔹🔸🔹

🌐 @kiwi_proxy
`

	proxyLinksText := ""
	for _, link := range proxyLinks {
		proxyLinksText = fmt.Sprintf("%v *[اتصال به پروکسی ✅](%v)*\n", proxyLinksText, link)

	}

	caption = fmt.Sprintf(caption, proxyLinksText)

	return caption
}

var (
	responseContentSaved = "محتوا با موفقیت افزوده شد"
)
