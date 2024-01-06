##  postfixlogparse

A simple command line tool for postfix log parsing without dependencies.
Supports filtering specific time, sender, receiver, subject, and email status. There are multiple parameters available, and multiple parameters can be used together.
Author: https://github.com/xiaoshanyangcode  or  https://gitee.com/xiaoshanyangcode

Such as:

```
$  ./postfixlogparse -a 20240105-130101 -b 20240105-230101 -m sent
---------------------------------------------------------------
Count：1    EmailId：AB6DE2099D1E  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：test@163.com  Result：success sent
......

************************* RESULTS *************************
send（成功）:               3
bounced（退信）:            0
......
```



###  Download

#####  Method1：

Download it from website：https://github.com/xiaoshanyangcode/postfixlogparse/releases



##### Method2：

This method relies on golang language.

Get  $GOPATH  : `go env`

Make sure $GOPATH/bin is in the system $PATH environment variable

```shell
go install github.com/xiaoshanyangcode/postfixlogparse@latest
```



###  Usage

```bash
$ ./postfixlogparse -h
```



```
Usage:
  postfixlogparse [flags]

Flags:
  -a, --after string      Filter sending time after the parameter
                            The format only supports (year-month-day or monthday-hourminute or yearmonthday-hourminutesecond)
                            Example：-a 2023-12-25 or -a 1225-2250 or -a 20231225-205010
                            Default:The last 7 days when the three parameters -a -b -l are not used
  -b, --before string     Filter sending time before the parameter
                            The format only supports (year-month-day or monthday-hourminute or yearmonthday-hourminutesecond)
                            Example：-a 2023-12-25 or -a 1225-2250 or -a 20231225-205010
                            Default:The last 7 days when the three parameters -a -b -l are not used
  -f, --file string       Which postifx log files to parse, supoort gz format, relative paths can be used,  * can be used to match  multiple log files(need to add \ in front of *)
                            Example：-f /var/log/mail.\* or mail.log
                            Default: Contains all postfix log files (comes from the system configuration)
  -h, --help              help for postfixlogparse
  -l, --last string       Filter results by how recent they are（最近多久）
                            Format:Number + unit, the unit only supports (M minutes, h hours, d days, w weeks, m months, y years)
                            Example：-l 3d
                            Default:The last 7 days when the three parameters -a -b -l are not used
  -o, --output string     Output the results to CSV file, supporting relative paths
                            Example：-o /tmp/postfix_log_output.csv or -o postfix_log_output.csv
  -r, --receiver string   Filter receiver email addresses
                            Example：-r abc@163.com
  -m, --result string     Filter email result status（sent 发送成功,bounced 退信,deferred 延迟,rejected 拒绝）, only supports：sent、bounced、rejected、deferred
                            Example：-m sent
                            Default: all
  -s, --sender string     Filter sender email addresses
                            Example：-s abc@163.com
  -t, --theme string      (Filter email subject) Filter the strings contained in the email subject (please add double quotes to the second side of the string when -t containing spaces 包含空格需加双引号)
                            Example：-t 会议  or  -t "The subject"
  -v, --version           软件版本(version)
```



####  Example 1

```shell
$ ./postfixlogparse -l 2h -o output.csv
---------------------------------------------------------------
Count：1    EmailId：AB6DE2099D1E  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：12345@qq.com         Result：bounced  Err log：host mx3.qq.com[183.47.111.94] said: 550 Domain may not exist or DNS check failed......
Count：2    EmailId：AB6DE2099D1E  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：test@163.com  Result：success sent
Count：3    EmailId：B18872099D25  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：test@163.com  Result：success sent
Count：4    EmailId：B18872099D25  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：12345@qq.com         Result：bounced  Err log：host mx3.qq.com[183.47.111.94] said: 550 Domain may not exist or DNS check failed ......
Count：5    EmailId：A8B202099D1F  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：12345@qq.com         Result：bounced  Err log：host mx3.qq.com[183.47.111.94] said: 550 Domain may not exist or DNS check failed ......
Count：6    EmailId：A8B202099D1F  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：test@163.com  Result：success sent

************************* RESULTS *************************
send（成功）:               3
bounced（退信）:            3
deferred（延迟）:           0
rejected（拒绝）:           0

-----
Total:        6

Successfully written to the file：output.csv
***********************************************************
--- 0.007      seconds ---
```

output.csv file such as：

| Count | EmailId      | Time             | Sender            | Subject | Receiver     | Result       | Err log |
| ----- | ------------ | ---------------- | ----------------- | ------- | ------------ | ------------ | ------- |
| 1     | 93E43218FB69 | 2023/12/26 17:30 | vm14@sun-site.com | test1   | test@163.com | success sent |         |
| 4     | 975F5218FB6E | 2023/12/26 17:30 | vm14@sun-site.com | test2   | test@163.com | success sent |         |

####  Example 2

```shell
$  ./postfixlogparse -a 20240105-130101 -b 20240105-230101 -m sent
---------------------------------------------------------------
Count：1    EmailId：AB6DE2099D1E  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：test@163.com  Result：success sent
Count：2    EmailId：B18872099D25  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：test@163.com  Result：success sent
Count：3    EmailId：A8B202099D1F  2024-01-05 15:10:45 Sender：vm14@sun-site.com    Subject：           Receiver：test@163.com  Result：success sent

************************* RESULTS *************************
send（成功）:               3
bounced（退信）:            0
deferred（延迟）:           0
rejected（拒绝）:           0

-----
Total:        3

***********************************************************
--- 0.008      seconds ---
```



####  Example 3

```
$ ./postfixlogparse -a 20240105-130101 -b 20240105-230101 -m sent -s vm14@sun-site.com -t testsubject

************************* RESULTS *************************
send（成功）:               0
bounced（退信）:            0
deferred（延迟）:           0
rejected（拒绝）:           0

-----
Total:        0

***********************************************************
--- 0.003      seconds ---
```



###  Recommended configuration

####  Configure subject message in log

This is a little trick for Postfix, it lets you log the subject, from and to of all the emails postfix sends (or which pass through it if you run it as a relay). It comes in handy when you need to debug an email issue and need to confirm your mailserver has sent the message.

First create the file and insert this into it:`/etc/postfix/header_checks`

```
/^subject:/      WARN
```

Now, in your postfix add the following to the end of the file:`/etc/postfix/main.cf`

```
header_checks = regexp:/etc/postfix/header_checks
```

Check whether the postfix configuration is correct. 

If there is no output after executing `postfix check`, the configuration is correct.

```
postfix check
```

And restart postfix:

```
systemctl restart postfix
```

You will hopefully now get log items like below, and if not you have a problem with your mailserver:

```
Dec  4 08:23:05 localhost postfix/cleanup[2278]: 90CA714: warning: header Subject: This is a testmail which gets logged from localhost[127.0.0.1]; from=<root@localhost> to=<root@localhost> proto=ESMTP helo=<localhost>
```



####  Configure year message in postfix log

This little trick is achieved by modifying the configuration of rsyslog.

#####  Step1: Modify  **/etc/rsyslog.conf**

First change the content of the line starting with mail.* to the following content

```
mail.*                       -/var/log/postfix/mail.log;FormatWithYear
```

Second insert a new line above the line starting with mail.*, as follows

```
$template FormatWithYear,"%timestamp:::date-year% %timestamp:::date-month%-%timestamp:::date-day% %timestamp:::date-hour%:%timestamp:::date-minute%:%timestamp:::date-second% %HOSTNAME% %syslogtag% %msg%\n"
```

**The final result is as follows**

```
$template FormatWithYear,"%timestamp:::date-year% %timestamp:::date-month%-%timestamp:::date-day% %timestamp:::date-hour%:%timestamp:::date-minute%:%timestamp:::date-second% %HOSTNAME% %syslogtag% %msg%\n"

mail.*                       -/var/log/postfix/mail.log;FormatWithYear
```

#####  Step2: Restart rsyslog

```
systemctl restart rsyslog
```

#####  Step3:  Check Result

```
tailf /var/log/postfix/mail.log
---------------------------------------------------------------
2024 01-05 17:06:37 vm13 postfix/cleanup[6669]: 41DAA218F382: message-id=<20240105090637.41DAA218F382@vm14>
2024 01-05 17:06:37 vm13 postfix/qmgr[1399]: 41DAA218F382: from=<>, size=2536, nrcpt=1 (queue active)
2024 01-05 17:06:37 vm13 postfix/bounce[6731]: DDD2B218F37C: sender non-delivery notification: 41DAA218F382
2024 01-05 17:06:37 vm13 postfix/qmgr[1399]: DDD2B218F37C: removed
```



#### Logrotate

Daily logrotate created in this format for postfix logs:

**Note**: /var/log/postfix/*.log is just an example, please follow the actual path.

```
/var/log/postfix/*.log {
    rotate 366
    dateext
    daily
    missingok
    notifempty
    compress
    delaycompress
}
```



Example of postfix log directory files:

```
mail.log
mail.log-20170724.gz
mail.log-20170725.gz
mail.log-20170726.gz
mail.log-20170727.gz
mail.log-20170728.gz
mail.log-20170729.gz
mail.log-20170730.gz
mail.log-20170731.gz
mail.log-20170801.gz
mail.log-20170802.gz
mail.log-20170803.gz
mail.log-20170804
```