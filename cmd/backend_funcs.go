package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v4"
)

func isEmailValid(email string) bool {
	// 定义邮箱格式的正则表达式
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// 编译正则表达式
	re := regexp.MustCompile(emailRegex)
	// 使用正则表达式检查字符串是否匹配
	return re.MatchString(email)
}

func paseLast(last string) (befo time.Time, aft time.Time, err error) {
	numstr := last[:len(last)-1]
	unitbyte := last[len(last)-1]
	num, err := strconv.Atoi(numstr)
	if err != nil {
		return befo, aft, err
	}
	unitMap := map[byte]time.Duration{
		'M': time.Minute,
		'h': time.Hour,
		'd': 24 * time.Hour,
		'w': 7 * 24 * time.Hour,
		'm': 30 * 24 * time.Hour,
		'y': 365 * 24 * time.Hour,
	}
	tmp_val, ok := unitMap[unitbyte]
	if !ok {
		return befo, aft, err
	}

	tmp_val *= time.Duration(num)

	befo = time.Now()
	aft = befo.Add(-tmp_val)
	return befo, aft, nil
}

func parseRsyslog() []string {
	file, err := os.Open("/etc/rsyslog.conf")
	if err != nil {
		fmt.Println("Can't open /etc/rsyslog.conf, err msg:", err)
		os.Exit(1)
		return nil
	}
	defer file.Close()
	var rsysloglist []string
	// 避免重复
	uniqueMap := make(map[string]bool)
	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tmp_str := strings.TrimSpace(line)

		if strings.HasPrefix(tmp_str, "mail.*") {
			tmp_str_mail1 := strings.TrimSpace(tmp_str[6:])[1:]
			tmp_str_mail1 = strings.Split(tmp_str_mail1, ";")[0]
			_, ok := uniqueMap[tmp_str_mail1]
			if !ok {
				rsysloglist = append(rsysloglist, tmp_str_mail1)
				uniqueMap[tmp_str_mail1] = true
			}
		}
		if strings.HasPrefix(tmp_str, "mail.info") {
			tmp_str_mail1 := strings.TrimSpace(tmp_str[9:])[1:]
			tmp_str_mail1 = strings.Split(tmp_str_mail1, ";")[0]
			_, ok := uniqueMap[tmp_str_mail1]
			if !ok {
				rsysloglist = append(rsysloglist, tmp_str_mail1)
				uniqueMap[tmp_str_mail1] = true
			}
		}
	}

	// Check scan for errors
	if err := scanner.Err(); err != nil {
		fmt.Println("/etc/rsyslog.conf ,the file scan error:", err)
		os.Exit(1)
		return nil
	}
	return rsysloglist
}

func parseFile() ([]string, map[string]time.Time) {
	var fileListTmp []string
	fileListOut := []string{}
	if len(file) == 0 {
		fileListTmp = parseRsyslog()
	} else {
		fileListTmp = parseFileflag(file)
	}
	tmp_file_str := ""
	for _, str := range fileListTmp {
		if len(str) > 0 {
			tmp_file_str += str + " " + str + "-* "
		}

	}
	command := "ls -lrt " + tmp_file_str + "  | awk '{print $9}'"
	// fmt.Println(command)
	// run shell cmd
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
	}
	fileList := strings.Split(string(output), "\n")
	// fmt.Println("fileList", fileList, command)
	if len(fileList) == 0 {
		fmt.Printf("%s file does not exist\n", aurora.BgBlack(aurora.Red(fileListTmp)))
		os.Exit(1)
	}
	newFileList := []string{}
	for _, file := range fileList {
		if !strings.HasPrefix(file, "ls:") {
			newFileList = append(newFileList, file)
		}
	}
	if len(fileList) == 0 {
		fmt.Printf("%s file does not exist\n", aurora.BgBlack(aurora.Red(fileListTmp)))
		os.Exit(1)
	}
	fileMtime := make(map[string]time.Time)
	for _, file := range newFileList {
		if len(file) > 0 {
			// fmt.Println("file", file)
			fileMtime[file] = getMtime(file)
			fileListOut = append(fileListOut, file)
		}
	}
	return fileListOut, fileMtime

}

func parseFileflag(path string) []string {
	command := "ls -lrt " + path + "  | awk '{print $9}'"
	// run shell command
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err", err)
		os.Exit(1)
		return nil
	}
	fileList := strings.Split(string(output), "\n")
	if len(fileList) == 0 || strings.HasPrefix(fileList[0], "ls:") {
		fmt.Printf("%s file does not exist\n", aurora.BgBlack(aurora.Red(path)))
		os.Exit(1)
	}
	// 将输出字符串按行分割
	return fileList
}

// 提示：必须是已经存在的文件才能用此函数
func getMtime(filePathInfo string) time.Time {
	// filePath := "/path/to/your/file.txt"

	// 获取文件信息
	// filepa, _ := filepath.Abs(filePathInfo)
	// fmt.Println("aaaa", filepa, "asdfasdfasdf", filePathInfo)
	fileInfo, _ := os.Stat(filePathInfo)
	// if err != nil {
	// 	fmt.Println("Error--:", err.Error())

	// }

	// fmt.Println("zzz", filePathInfo, fileInfo)
	// 获取修改时间（mtime）
	return fileInfo.ModTime()

	// // 打印修改时间
	// fmt.Printf("File: %s\n", filePath)
	// fmt.Printf("Last Modified Time: %s\n", modTime.Format(time.RFC3339))
}

func ensureDirectoryExists(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, os.ModePerm)
}

// appendToCSV 用于在现有文件上追加写入数据行
func appendToCSV(writer *csv.Writer, data []string) error {
	// 写入数据行
	err = writer.Write(data)
	if err != nil {
		return err
	}

	// 刷新缓冲区，确保数据行写入文件
	writer.Flush()

	// 检查错误
	return writer.Error()
}
