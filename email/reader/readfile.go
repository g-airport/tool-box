package reader

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/g-airport/tool-box/email/entity"
)

// this a read source file

var emailFile = "/Users/tqll/Downloads/email_sample.csv"

func InitSourceData() []*entity.EmailInfo {
	data, err := ioutil.ReadFile(emailFile)
	fmt.Println(emailFile)
	if err != nil {
		panic(fmt.Sprintf("read source file panic: %v", err))
	}
	csvReader := csv.NewReader(strings.NewReader(string(data[:])))
	out := make([]*entity.EmailInfo, 0)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, v := range record {
			r := strings.Split(v, "	")
			e := &entity.EmailInfo{
				Email:     r[0],
				SrcStatus: r[1],
			}
			out = append(out, e)
		}
	}
	return out
}
