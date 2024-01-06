package cmd

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"mime"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/logrusorgru/aurora/v4"
)

type ReceiverCell struct {
	Receiver string
	Result   string
	Rtime    time.Time
	Message  string
}

type MailStruct struct {
	Subject      string
	Sender       string
	IsTrue       bool
	IsDone       bool
	ReceiverList *[]ReceiverCell
	SendTime     time.Time
}

var isDone MailStruct

// 邮件主题解码为人类可读
// 原始邮件的信头：Subject: =?utf-8?b?576O5Zu9?=
func dealWithFile(fileList []string, fileMap map[string]time.Time) {
	startTime := time.Now()
	// 输出为：美国
	// fmt.Println(decode_sub("=?utf-8?b?576O5Zu9?="))

	// 判断是否要输出文件到csv
	if len(output) > 0 {
		err = ensureDirectoryExists(output)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}
		outfile, err = os.Create(output)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}

		// 创建 CSV writer
		writer = csv.NewWriter(outfile)

		// 写入标题行
		header := []string{"Count", "EmailId", "Time", "Sender", "Subject", "Receiver", "Result", "Err log"}
		err = writer.Write(header)
		if err != nil {
			fmt.Println("Error writing header:", err)
			return
		}

		// 刷新缓冲区，确保标题行写入文件
		writer.Flush()

		// 检查错误
		if err := writer.Error(); err != nil {
			fmt.Println("Error flushing writer:", err)
			return
		}
	}

	resultMap := make(map[string]*MailStruct)
	mailCount := make(map[string]int, 4)
	mailCount["sent"] = 0
	mailCount["bounced"] = 0
	mailCount["rejected"] = 0
	mailCount["deferred"] = 0
	isDone = MailStruct{IsDone: true}
	outputJishu = 1
	for _, fileCell := range fileList {
		if fileMap[fileCell].Before(parseAfter) {
			// fmt.Println("continue", fileCell, fileMap[fileCell])
			continue
		} else {
			// fmt.Println("okok", fileCell, fileMap[fileCell])
			year := fileMap[fileCell].Year()
			if strings.HasSuffix(fileCell, ".gz") {
				analysisGZfile(fileCell, year, &resultMap, &mailCount)
			} else {
				analysisFile(fileCell, year, &resultMap, &mailCount)
			}
		}
	}

	for mailid, mailstruce := range resultMap {
		outputLine(mailid, mailstruce)
	}

	outputCount(&mailCount, startTime)

	if len(output) > 0 {
		if outfile != nil {
			defer outfile.Close()
		}

	}
}

func outputLine(mailid string, mailstruce *MailStruct) {
	// 如果没有收件人，就跳出本次循环
	if mailstruce == nil || !mailstruce.IsTrue || *mailstruce.ReceiverList == nil {
		return
	}
	for _, k := range *mailstruce.ReceiverList {
		if k.Result == "sent" {
			// fmt.Println("123")
			fmt.Printf("Count：%-4d EmailId：%-13s %-19s Sender：%-20s Subject：%-10s Receiver：%-20s Result：%-8s\n", outputJishu, mailid, k.Rtime.Format("2006-01-02 15:04:05"), mailstruce.Sender, mailstruce.Subject, k.Receiver, aurora.BgBlack(aurora.Green("success sent")))
			if writer != nil {
				err = appendToCSV(writer, []string{strconv.Itoa(outputJishu), mailid, k.Rtime.Format("2006-01-02 15:04:05"), mailstruce.Sender, mailstruce.Subject, k.Receiver, "success sent", ""})
				if err != nil {
					fmt.Println("Error appending to CSV:", err)
					return
				}
			}
		} else {
			fmt.Printf("Count：%-4d EmailId：%-13s %-19s Sender：%-20s Subject：%-10s Receiver：%-20s Result：%-8s Err log：%-20s\n", outputJishu, mailid, k.Rtime.Format("2006-01-02 15:04:05"), mailstruce.Sender, mailstruce.Subject, k.Receiver, aurora.BgBlack(aurora.Red(k.Result)), aurora.BgBlack(aurora.Yellow(k.Message)))
			if writer != nil {
				err = appendToCSV(writer, []string{strconv.Itoa(outputJishu), mailid, k.Rtime.Format("2006-01-02 15:04:05"), mailstruce.Sender, mailstruce.Subject, k.Receiver, k.Result, k.Message})
				if err != nil {
					fmt.Println("Error appending to CSV:", err)
					return
				}
			}
		}
		outputJishu += 1
	}
}

func outputCount(count *map[string]int, oldtime time.Time) {
	allCount := (*count)["sent"] + (*count)["bounced"] + (*count)["deferred"] + (*count)["rejected"]
	fmt.Println("\n************************* RESULTS *************************")
	fmt.Printf("%-20s    %-10d\n", aurora.BgBlack(aurora.Green("send（成功）:")), (*count)["sent"])
	fmt.Printf("%-20s    %-10d\n", aurora.BgBlack(aurora.Red("bounced（退信）:")), (*count)["bounced"])
	fmt.Printf("%-20s    %-10d\n", aurora.BgBlack(aurora.Red("deferred（延迟）:")), (*count)["deferred"])
	fmt.Printf("%-20s    %-10d\n\n", aurora.BgBlack(aurora.Red("rejected（拒绝）:")), (*count)["rejected"])
	fmt.Println("-----")
	fmt.Printf("%-10s    %-10d\n\n", "Total:", allCount)
	if writer != nil {
		fmt.Printf("Successfully written to the file：%s\n", aurora.BgBlack(aurora.Green(output)))
	}
	fmt.Println("***********************************************************")
	fmt.Printf("--- %-10.3f seconds ---\n", time.Since(oldtime).Seconds())
}

func analysisFile(fileCell string, year int, resultmap *map[string]*MailStruct, count *map[string]int) {
	// 打开文件
	file_tmp, err := os.Open(fileCell)
	if err != nil {
		fmt.Println("Unable to open file:", err)
		return
	}
	defer file_tmp.Close()

	// 使用 bufio.NewScanner 来逐行扫描文件
	scanner := bufio.NewScanner(file_tmp)

	// 逐行读取文件内容
	isSkipFile := false
	for scanner.Scan() {
		line := scanner.Text()
		AnalysizeLine(line, year, resultmap, count, &isSkipFile)
		if isSkipFile {
			break
		}
	}

	// 检查扫描过程中是否出现错误
	if err := scanner.Err(); err != nil {
		fmt.Println("File scan error:", err)
	}
}
func analysisGZfile(fileCell string, year int, resultmap *map[string]*MailStruct, count *map[string]int) {
	// 打开gz文件
	file_tmp, err := os.Open(fileCell)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file_tmp.Close()

	// 使用gzip.NewReader创建一个gzip.Reader
	gzipReader, err := gzip.NewReader(file_tmp)
	if err != nil {
		fmt.Println("Error creating gzip reader:", err)
		return
	}
	defer gzipReader.Close()

	// 使用bufio.NewScanner创建一个Scanner来逐行读取
	scanner := bufio.NewScanner(gzipReader)

	// 逐行读取文件内容
	isSkipFile := false
	for scanner.Scan() {
		line := scanner.Text()
		AnalysizeLine(line, year, resultmap, count, &isSkipFile)
		if isSkipFile {
			break
		}
	}

	// 检查扫描是否出错
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

// 解码主题
func decode_sub(encodedSubject string) string {
	decoder := new(mime.WordDecoder)
	decodedHeader, err := decoder.DecodeHeader(encodedSubject)
	if err != nil {
		print("err", err)
		return ""
	}
	return decodedHeader
}

// 判断字符串是否为邮件id
func isEmailid(s string) bool {
	n := len(s)

	// 检查字符串是否为空
	if n != 13 {
		return false
	}

	// 检查前 n-1 个字符是否只包含数字和大写字母
	for _, char := range s[:n-1] {
		if !(unicode.IsDigit(char) || unicode.IsUpper(char)) {
			return false
		}
	}

	// 检查最后一个字符是否是冒号
	return s[n-1] == ':'
}

func removeLastCharacter(s string) string {
	if len(s) == 0 {
		return s
	}
	return s[:len(s)-1]
}

func AnalysizeLine(line string, year int, resultmap *map[string]*MailStruct, count *map[string]int, isSkip *bool) {
	//  1.split line
	linelist := strings.Split(line, " ")
	if len(linelist) < 6 {
		return
	}

	//  2.parse time
	//  Jan  4 11:27:37  ==to=>  Jan 04 11:27:37
	if len(linelist[1]) == 0 {
		linelist = append(linelist[:1], linelist[2:]...)
	}
	if len(linelist[1]) == 1 {
		linelist[1] = "0" + linelist[1]
	}
	// 获取当前时区,并解析时间
	loc, _ := time.LoadLocation("Local")
	var parsedTime time.Time
	parsedTime, err = time.ParseInLocation("2006 Jan 02 15:04:05", strconv.Itoa(year)+" "+strings.Join(linelist[:3], " "), loc)
	if err != nil {
		parsedTime, err = time.ParseInLocation("2006 01-02 15:04:05", strings.Join(linelist[:3], " "), loc)
		if err != nil {
			fmt.Println("postfix the date format in the log only supports : Jan 02 15:04:05（default） or 2006 01-02 15:04:05")
			os.Exit(1)
		}
		// 判断日志中是否有年的配置
		hasYearConf = true
	}
	// 如果该行的时间超过了最大时间，就跳过该文件
	if parsedTime.After(parseBefore) {
		*isSkip = true
		return
	}
	if parsedTime.Before(parseAfter) {
		return
	}

	// 3.parse mailid
	// 如果行里不包含邮件id就退出
	if !isEmailid(linelist[5]) {
		return
	}
	mailId := removeLastCharacter(linelist[5])

	// 同一个id避免重复解析
	if (*resultmap)[mailId] == &isDone {
		return
	}

	if strings.HasPrefix(linelist[6], "client=") {
		if (*resultmap)[mailId] == nil {
			(*resultmap)[mailId] = &MailStruct{SendTime: parsedTime}
		}
		return
	}

	mailStruct := (*resultmap)[mailId]
	if mailStruct == nil {
		return
	}

	senderTmp := ""
	receiverTmp := ""
	subjectTmp := ""
	result_str := ""
	log_str := ""
	if len(linelist) > 8 && strings.Join(linelist[6:9], " ") == "warning: header Subject:" && len((*mailStruct).Subject) == 0 {
		// 判断是否有主题的配置
		hasSubjectConf = true
		tmp_sub_str := ""
		tmp_sub_num := 9
		for i := 9; i <= 20; i++ {
			if linelist[i+1] == "from" && strings.HasSuffix(linelist[i+2], "];") && strings.HasPrefix(linelist[i+3], "from=<") {
				tmp_sub_num = i
				tmp_sub_str = strings.Join(linelist[9:tmp_sub_num+1], " ")
				// 跳出循环
				i = 100
			}
		}

		subjectTmp = decode_sub(tmp_sub_str)
		// subjectTmp = decode_sub(linelist[9])
		if (len(subject) > 0 && strings.Contains(subjectTmp, subject)) || len(subject) == 0 {
			if (*mailStruct).ReceiverList == nil {
				(*mailStruct).ReceiverList = &[]ReceiverCell{}
			}
			(*mailStruct).Subject = subjectTmp
		} else {
			delete((*resultmap), mailId)
		}

		return
	}

	// 如有有主题的筛选，但是日志没有配置主题就删除该mailId
	if len(subject) > 0 && len((*mailStruct).Subject) == 0 {
		delete((*resultmap), mailId)
		return
	}

	if strings.HasPrefix(linelist[6], "from=<") && strings.HasSuffix(linelist[6], ">,") {
		senderTmp = linelist[6][6 : len(linelist[6])-2]
		if (len(sender) > 0 && senderTmp == sender) || (len(sender) == 0 && len(senderTmp) > 0) {
			(*mailStruct).Sender = senderTmp
		} else {
			delete((*resultmap), mailId)
		}
		return
	}

	if strings.HasPrefix(linelist[6], "to=<") && strings.HasSuffix(linelist[6], ">,") {
		receiverTmp = linelist[6][4 : len(linelist[6])-2]
		if (len(receiver) > 0 && receiverTmp == receiver) || (len(receiver) == 0 && len(receiverTmp) > 0) {

		} else {
			return
		}

		re := regexp.MustCompile(`status=([^\s,]+)(.*)`)
		match := re.FindStringSubmatch(line)
		if len(match) == 3 {
			result_str = match[1]
			// filter
			if len(result) > 0 && result != result_str {
				return
			}
			if result_str != "sent" {
				log_str = match[2][2 : len(match[2])-1]
			}
			if (*mailStruct).ReceiverList == nil {
				(*mailStruct).ReceiverList = &[]ReceiverCell{}
			}

			for _, j := range *(*mailStruct).ReceiverList {
				// when receiver appear more than once ,rewrite cell
				if j.Receiver == receiverTmp {
					j = ReceiverCell{Receiver: receiverTmp, Result: result_str, Rtime: parsedTime, Message: log_str}
					return
				}
			}

			*(*mailStruct).ReceiverList = append(*(*mailStruct).ReceiverList, ReceiverCell{Receiver: receiverTmp, Result: result_str, Rtime: parsedTime, Message: log_str})
			// 如果有收件人信息就设置为true
			(*mailStruct).IsTrue = true
			// result status count
			(*count)[result_str]++
		}
	}

	if len(linelist) == 7 && linelist[6] == "removed" {

		outputLine(mailId, mailStruct)
		// 标记已操作完成，避免同一个id重复操作
		(*resultmap)[mailId] = &isDone
	}
}
