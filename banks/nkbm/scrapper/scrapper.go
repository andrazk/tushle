package scrapper

import (
	"strings"
)

// const dateFormat = "02.01.2006"

// func scrapEntriesOsebniPast(from, to time.Time, buf interface{}) chromedp.Tasks {
// 	var (
// 		inputFrom = `//input[@name='dtm_min']`
// 		inputTo   = `//input[@name='dtm_max']`
// 		btn       = `//*[@id="inside"]/div/center/div[1]/div/form[1]//div[1]/table/tbody/tr[3]/td[5]/button`
// 		tbl       = `#zhtForm > table`
// 	)

// 	return chromedp.Tasks{
// 		chromedp.Navigate("https://bankanet.nkbm.si/bnk/Nkbm?action=prmRcn&rcnId=119833&vlt=EUR&error=stn"),
// 		chromedp.WaitVisible(`#bnk`),
// 		chromedp.Clear(inputFrom),
// 		chromedp.SendKeys(inputFrom, from.Format(dateFormat)),
// 		chromedp.Clear(inputTo),
// 		chromedp.SendKeys(inputTo, to.Format(dateFormat)),
// 		chromedp.Click(btn),
// 		chromedp.WaitVisible(tbl),
// 		chromedp.Click(`//*[@id="500"]`),
// 		chromedp.WaitVisible(tbl),
// 		chromedp.Evaluate(jsTraverseTablePast(), buf),
// 	}
// }

// func scrapEntriesOsebniFuture(buf interface{}) chromedp.Tasks {
// 	var waitFor = `//*[@id="straniZgoraj"][div or p]`

// 	return chromedp.Tasks{
// 		chromedp.Navigate("https://bankanet.nkbm.si/bnk/Nkbm?action=plcArh&plcSts=IZVD"),
// 		chromedp.WaitVisible(waitFor),
// 		chromedp.Evaluate(jsTraverseTableFuture(), buf),
// 	}
// }

// func scrapEntriesMastercard(from, to time.Time, buf interface{}) chromedp.Tasks {
// 	var (
// 		inputFrom = `//input[@name='dtm_min']`
// 		inputTo   = `//input[@name='dtm_max']`
// 		btn       = `//*[@id="inside"]/div/center/div[1]/div/form[1]//div[1]/table/tbody/tr[3]/td[5]/button`
// 		tbl       = `#zhtForm > table`
// 	)

// 	return chromedp.Tasks{
// 		chromedp.Navigate("https://bankanet.nkbm.si/bnk/Nkbm?action=pgdprm&pgdTip=KARTICA&pgdId=KARTICA-MAS-SI56044750239501094-X7562-EUR"),
// 		chromedp.WaitVisible(`#bnk`),
// 		chromedp.Clear(inputFrom),
// 		chromedp.SendKeys(inputFrom, from.Format(dateFormat)),
// 		chromedp.Clear(inputTo),
// 		chromedp.SendKeys(inputTo, to.Format(dateFormat)),
// 		chromedp.Click(btn),
// 		chromedp.WaitVisible(tbl),
// 		chromedp.Click(`//*[@id="500"]`),
// 		chromedp.WaitVisible(tbl),
// 		chromedp.Evaluate(jsTraverseTableCard(), buf),
// 	}
// }

// GetText in element and his children
func GetText(sel string) (js string) {
	const funcJS = `function getText(sel) {
				var text = [];
				var elements = document.body.querySelectorAll(sel);

				for(var i = 0; i < elements.length; i++) {
					var current = elements[i];
					if(current.children.length === 0 && current.textContent.replace(/ |\n/g,'') !== '') {
					// Check the element has no children && that it is not empty
						text.push(current.textContent);
					}
				}
				return text
			 };`

	invokeFuncJS := `var a = getText('` + sel + `'); a;`
	return strings.Join([]string{funcJS, invokeFuncJS}, " ")
}

// func jsTraverseTablePast() string {
// 	return `
// 		(function () {
// 			var data = []
// 			var rows = document.querySelectorAll('.template2 > tbody > tr')
// 			for (var i = 1; i < rows.length; i++) {
// 			var cells = rows[i].children
// 			data.push({
// 				exp: cells[1].textContent,
// 				inc: cells[2].textContent,
// 				desc: cells[4].textContent,
// 				payee: cells[5].textContent,
// 				acc: cells[6].textContent,
// 				ref: cells[7].textContent,
// 				book: cells[9].textContent,
// 				date: cells[10].textContent
// 			})
// 			}
// 			return data
// 		})()
// 	`
// }

// func jsTraverseTableFuture() string {
// 	return `
// 		(function () {
// 			var data = []
// 			var rows = document.querySelectorAll('.template2 > tbody > tr')
// 			for (var i = 1; i < rows.length; i++) {
// 				var cells = rows[i].children
// 				data.push({
// 					exp: cells[5].textContent,
// 					desc: cells[4].textContent,
// 					payee: cells[9].textContent,
// 					date: cells[8].textContent
// 				})
// 			}
// 			return data
// 		})()
// 	`
// }

// func jsTraverseTableCard() string {
// 	return `
// 		(function () {
// 			var data = []
// 			var rows = document.querySelectorAll('.template2 > tbody > tr')
// 			for (var i = 1; i < rows.length; i++) {
// 				var cells = rows[i].children
// 				if (cells[0].textContent === 'Skupaj:') {
// 					continue;
// 				}
// 				data.push({
// 					exp: cells[1].textContent,
// 					inc: cells[2].textContent,
// 					desc: cells[3].textContent + '; ' + cells[5].textContent,
// 					date: cells[4].textContent
// 				})
// 			}
// 			return data
// 		})()
// 	`
// }
