/*
Copyright © 2023 xiaoshangyangcode
*/
package cmd

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// 参数变量
var sender string
var receiver string
var subject string
var result string
var output string
var after string
var before string
var last string
var file string
var version bool

// 判断postfix配置中是否有主题的配置
var hasSubjectConf bool

// 判断日志中是否有年的配置
var hasYearConf bool

// 输出数量计数
var outputJishu int

// 判断是否要输出文件到csv
var outfile *os.File
var writer *csv.Writer

// 解析后的参数及可能用到的参数
var parseAfter time.Time
var parseBefore time.Time
var err error

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "postfixlogparse",
	Short: "A simple command line tool for postfix log parsing without dependencies",
	Long: `A simple command line tool for postfix log parsing without dependencies. 
Supports filtering specific time, sender, receiver, subject, and email status. There are multiple parameters available, and multiple parameters can be used together.
Author: https://github.com/xiaoshanyangcode  or  https://gitee.com/xiaoshanyangcode`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		// function main
		cmdMain(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.main.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// 添加一个名为 "name" 的标志
	pflag.StringVarP(&sender, "sender", "s", "", "Filter sender email addresses\n  Example：-s abc@163.com")
	pflag.StringVarP(&receiver, "receiver", "r", "", "Filter receiver email addresses\n  Example：-r abc@163.com")
	pflag.StringVarP(&subject, "theme", "t", "", "(Filter email subject) Filter the strings contained in the email subject (please add double quotes to the second side of the string when -t containing spaces 包含空格需加双引号)\n  Example：-t 会议  or  -t \"The subject\"")
	pflag.StringVarP(&result, "result", "m", "", "Filter email result status（sent 发送成功,bounced 退信,deferred 延迟,rejected 拒绝）, only supports：sent、bounced、rejected、deferred\n  Example：-m sent\n  Default: all")
	pflag.StringVarP(&output, "output", "o", "", "Output the results to CSV file, supporting relative paths\n  Example：-o /tmp/postfix_log_output.csv or -o postfix_log_output.csv")
	pflag.StringVarP(&file, "file", "f", "", "Which postifx log files to parse, supoort gz format, relative paths can be used,  * can be used to match  multiple log files(need to add \\ in front of *)\n  Example：-f /var/log/mail.\\* or mail.log\n  Default: Contains all postfix log files (comes from the system configuration)")
	pflag.StringVarP(&after, "after", "a", "", "Filter sending time after the parameter\n  The format only supports (year-month-day or monthday-hourminute or yearmonthday-hourminutesecond)\n  Example：-a 2023-12-25 or -a 1225-2250 or -a 20231225-205010\n  Default:The last 7 days when the three parameters -a -b -l are not used")
	pflag.StringVarP(&before, "before", "b", "", "Filter sending time before the parameter\n  The format only supports (year-month-day or monthday-hourminute or yearmonthday-hourminutesecond)\n  Example：-a 2023-12-25 or -a 1225-2250 or -a 20231225-205010\n  Default:The last 7 days when the three parameters -a -b -l are not used")
	pflag.StringVarP(&last, "last", "l", "", "Filter results by how recent they are（最近多久）\n  Format:Number + unit, the unit only supports (M minutes, h hours, d days, w weeks, m months, y years) \n  Example：-l 3d\n  Default:The last 7 days when the three parameters -a -b -l are not used")
	pflag.BoolVarP(&version, "version", "v", false, "软件版本(version)")

	// 将 pflag 绑定到 Cobra 命令
	rootCmd.Flags().AddFlagSet(pflag.CommandLine)
}
