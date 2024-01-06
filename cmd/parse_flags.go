package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"
)

func cmdMain(cmd *cobra.Command, args []string) {
	parseFlag(cmd, args)
	fileList, fileMap := parseFile()

	dealWithFile(fileList, fileMap)

	if !hasSubjectConf {
		fmt.Println(aurora.BgBlack(aurora.Yellow("\nThis postfix is not configured to write subject into the log. For configuration, please refer to: https://github.com/xiaoshanyangcode/postfixlogparse or https://gitee.com/xiaoshanyangcode/postfixlogparse  ，Chapter 'Configure postfix subject message in log' in the README.md file ")))
	}
	if !hasYearConf {
		fmt.Println(aurora.BgBlack(aurora.Yellow("\nThis postfix is not configured to write years into the log. For configuration, please refer to: https://github.com/xiaoshanyangcode/postfixlogparse or https://gitee.com/xiaoshanyangcode/postfixlogparse  ，Chapter 'Configure year message in postfix log' in the README.md file ")))
	}

}

func parseFlag(cmd *cobra.Command, args []string) {
	// output version
	if version {
		fmt.Println("1.0.1")
		os.Exit(0)
	}

	if len(sender) != 0 && !isEmailValid(sender) {
		fmt.Printf("The 'sender' parameter value error：%s , Not a standard email format！\n", aurora.BgBlack(aurora.Red(sender)))
		os.Exit(1)
	}
	if len(receiver) != 0 && !isEmailValid(receiver) {
		fmt.Printf("The 'receiver' parameter value error：%s , Not a standard email format！\n", aurora.BgBlack(aurora.Red(receiver)))
		os.Exit(1)
	}
	if result != "sent" && result != "bounced" && result != "rejected" && result != "deferred" && result != "" {
		fmt.Printf("The 'result' parameter value error：%s , Only support: send、bounced、rejected、deferred\n", aurora.BgBlack(aurora.Red(result)))
		os.Exit(1)
	}
	if len(output) > 0 && !strings.HasSuffix(output, ".csv") {
		fmt.Printf("The 'output' parameter value error：%s ,Don't end with .csv\n", aurora.BgBlack(aurora.Red(output)))
		os.Exit(1)
	}

	if len(after) != len(before) {
		fmt.Printf("Parameter value error：%s \nThe format only supports (year-month-day or monthday-hourminute or yearmonthday-hourminutesecond)  Example：-a 2023-12-25 or -a 1225-2250 or -a 20231225-205010\n", aurora.BgBlack(aurora.Red("'before' and 'after' need to be used in combination, and the format must be consistent, ")))
		os.Exit(1)
	}

	if len(last) > 0 && (len(after) > 0 || len(before) > 0) {
		fmt.Printf("Parameter value error：%s \n", aurora.BgBlack(aurora.Red("'last' cannot be used together with 'before' and 'after'")))
		os.Exit(1)
	}
	loc, _ := time.LoadLocation("Local")
	if len(after) > 0 && len(before) > 0 {
		parseAfter, err = time.ParseInLocation("20060102-150405", after, loc)
		if err != nil {
			parseAfter, err = time.ParseInLocation("20060102-1504", strconv.Itoa(time.Now().Local().Year())+after, loc)
			if err != nil {
				parseAfter, err = time.ParseInLocation("2006-01-02", after, loc)
				if err != nil {
					fmt.Printf("The 'after' parameter value error：%s The format only supports (year-month-day or monthday-hourminute or yearmonthday-hourminutesecond)\n  Example：-a 2023-12-25 or -a 1225-2250 or -a 20231225-205010\n", aurora.BgBlack(aurora.Red(after)))
					os.Exit(1)
				}
			}
		}

		parseBefore, err = time.ParseInLocation("20060102-150405", before, loc)
		if err != nil {
			parseBefore, err = time.ParseInLocation("20060102-150405", strconv.Itoa(time.Now().Local().Year())+before+"59", loc)
			if err != nil {
				parseBefore, err = time.ParseInLocation("2006-01-02150405", before+"235959", loc)
				if err != nil {
					fmt.Printf("The 'before' parameter value error：%s The format only supports (year-month-day or monthday-hourminute or yearmonthday-hourminutesecond)\n  Example：-a 2023-12-25 or -a 1225-2250 or -a 20231225-205010\n", aurora.BgBlack(aurora.Red(after)))
					os.Exit(1)
				}
			}
		}
		if parseBefore.Before(parseAfter) {
			fmt.Printf("Parameter value error：%s\n", aurora.BgBlack(aurora.Red("The value of 'after' should be smaller than the value of 'before'")))
			os.Exit(1)
		}
	}

	if len(after) == 0 && len(before) == 0 && len(last) == 0 {
		if len(file) == 0 {
			last = "1m"
		} else {
			last = "10y"
		}

	}

	if len(after) == 0 && len(before) == 0 {
		parseBefore, parseAfter, err = paseLast(last)
		if err != nil {
			fmt.Printf("The 'last' parameter value error：%s Format:Number + unit, the unit only supports (M minutes, h hours, d days, w weeks, m months, y years)  Example：-l 3d\n", aurora.BgBlack(aurora.Red(last)))
			os.Exit(1)
		}
	}

	// 在这里放置你的应用程序逻辑
	// fmt.Println("sender", sender)
	// fmt.Println("receiver", receiver)
	// fmt.Println("subject", subject)
	// fmt.Println("result", result)
	// fmt.Println("output", output)
	// fmt.Println("file", file)
	// fmt.Println("after", after)
	// fmt.Println("before", before)
	// fmt.Println("last", last)
	// fmt.Println("parseBefore", parseBefore)
	// fmt.Println("parseAfter", parseAfter)

}
