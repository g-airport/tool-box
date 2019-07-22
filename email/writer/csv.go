package writer

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/g-airport/tool-box/email/entity"
)

var baseDir = "/Users/tqll/work/go/src/github.com/g-airport/tool-box/email/tmp/email/"

func CSVOutput(in chan *entity.EmailInfo) [][]string {
	out := make([][]string, 0)
	title := []string{"email", "valid", "check_result", "invalid_reason", "error"}
	out = append(out, title)
	for v := range in {
		str := make([]string, 0)
		str = append(str, v.Email)
		str = append(str, v.SrcStatus)
		str = append(str, fmt.Sprintf("%v", v.RetStatus))
		str = append(str, v.Extra)
		str = append(str, fmt.Sprintf("%v", v.Err))
		out = append(out, str)
	}
	return out

}

func Write2CSV(es chan *entity.EmailInfo,doneChan chan struct{}) string {
	fmt.Println("base dir",baseDir)
	fi, err := os.Stat(baseDir)
	fmt.Println("file ",fi.Name())
	if err != nil {
		_ = os.MkdirAll(baseDir, 0755)
	}

	filename := fmt.Sprintf(baseDir + "email_check_all.csv")
	f, err := os.Create(filename)
	fmt.Println("dst file",f.Name())
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// 写入UTF-8 BOM
	f.WriteString("\xEF\xBB\xBF")
	//创建一个新的写入文件流
	w := csv.NewWriter(f)
	data := CSVOutput(es)
	//写入数据
	w.WriteAll(data)

	w.Flush()
	doneChan <- struct{}{}
	return filename
}
