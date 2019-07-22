package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/g-airport/tool-box/email/config"
	"github.com/g-airport/tool-box/email/entity"
	"github.com/g-airport/tool-box/email/reader"
	"github.com/g-airport/tool-box/email/verifier"
	"github.com/g-airport/tool-box/email/writer"
)

// ipo 模式
type VerifierFactory struct {
	InData     []*entity.EmailInfo
	OutChannel chan *entity.EmailInfo
	Verifier   []*verifier.Verifier
}

var (
	failT int32
	trueF int32
	cnt   int32
)

type Verify func(email string) (*verifier.Lookup, error)

const dataCount = 500000

var Verifiers = &VerifierFactory{}

var workerData [][]*entity.EmailInfo

func main() {
	//step 1.get config
	config.InitConfig()
	fmt.Println("available ip", len(config.Config.IP))
	fmt.Println("available source email", len(config.Config.SourceEmail))
	//step 2.get verifier email
	inData := reader.InitSourceData()

	//map ip:source_email
	var (
		ip          = config.Config.IP
		sourceEmail = config.Config.SourceEmail
	)

	// todo use map ip:sourceEmail
	currentLen := len(config.Config.IP)
	verifiers := make([]*verifier.Verifier, 0, currentLen)
	for i := 0; i < currentLen; i++ {
		v := verifier.NewVerifier(ip[i], sourceEmail[i])
		verifiers = append(verifiers, v)
	}
	//step 3.init factory
	funcFactory := &VerifierFactory{
		InData:     make([]*entity.EmailInfo, 0),
		OutChannel: make(chan *entity.EmailInfo, dataCount),
		Verifier:   make([]*verifier.Verifier, 0),
	}
	Verifiers = funcFactory
	Verifiers.InData = inData
	Verifiers.Verifier = verifiers

	// log
	//fmt.Println(Verifiers.InData)
	//fmt.Println(Verifiers.Verifier)

	// start verify
	dataLen := len(Verifiers.InData)
	vLen := len(Verifiers.Verifier)
	nextWorker := dataLen / vLen

	fmt.Println(dataLen, nextWorker)

	//workerData
	workerData = make([][]*entity.EmailInfo, 0)
	//inData
	totalData := Verifiers.InData
	for i := 0; i < vLen+1; i++ {
		group := make([]*entity.EmailInfo, 0)
		if (i+1)*nextWorker < dataLen {
			group = totalData[i*nextWorker+1 : (i+1)*nextWorker]
			workerData = append(workerData, group)
		} else {
			// last group
			group = totalData[i*nextWorker+1:]
			for _, v := range group {
				workerData[vLen-1] = append(workerData[vLen-1], v)
			}
		}
	}
	//map reduce
	//fmt.Println(len(workerData),len(workerData[139]),workerData[139][3409])
	wg := new(sync.WaitGroup)
	for i := 0; i < vLen; i++ {
		wg.Add(1)
		go Do(Verifiers.Verifier[i].Verify, workerData[i], wg)
	}
	wg.Wait()
	close(Verifiers.OutChannel)

	done := make(chan struct{})
	go writer.Write2CSV(Verifiers.OutChannel, done)
	<-done
	fmt.Println(len(inData), trueF, failT)
}

func Do(f Verify, data []*entity.EmailInfo, wg *sync.WaitGroup) {

	for _, v := range data {
		e := &entity.EmailInfo{
			Email:     v.Email,
			SrcStatus: v.SrcStatus,
		}
		out, err := f(v.Email)
		fmt.Printf("go: %v ==> cnt: %d \n", out,atomic.AddInt32(&cnt,1))
		if out != nil {
			e.Err = err
			e.RetStatus = out.Deliverable
			if e.SrcStatus == "250" || e.SrcStatus == "1" {
				if !out.Deliverable {
					e.Extra = fmt.Sprintf("源数据状态正常，校验结果为无效邮箱")
					atomic.AddInt32(&failT, 1)
				}
			}
			if e.SrcStatus == "domain_invalid" || e.SrcStatus == "0" {
				if out.Deliverable {
					e.Extra = fmt.Sprintf("源数据状态不合法，校验结果证实该邮箱有效")
					atomic.AddInt32(&trueF, 1)
				}
			}
		}
		Verifiers.OutChannel <- e
	}
	wg.Done()
}
