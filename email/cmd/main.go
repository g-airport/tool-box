package main

import (
	"encoding/binary"
	"fmt"
	"github.com/g-airport/tool-box/email/client"
	"math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/g-airport/tool-box/email/config"
	"github.com/g-airport/tool-box/email/entity"
	"github.com/g-airport/tool-box/email/reader"
	"github.com/g-airport/tool-box/email/verifier"
	"github.com/g-airport/tool-box/email/writer"
)

// ipo 模式
type VerifierFactory struct {
	InData       []*entity.EmailInfo
	OutChannel   chan *entity.EmailInfo
	Verifier     []*verifier.Verifier
	ErrEmailChan chan *entity.EmailInfo
}

var (
	failT        int32
	trueF        int32
	cnt          int32
	errEmailCnt  int32
	nullEmailCnt int32
)

type Verify func(email string) (*verifier.Lookup, error)

const dataCount = 500000

var Verifiers = &VerifierFactory{}

var workerData [][]*entity.EmailInfo

func main() {
	//step 1.get config
	config.InitConfig()
	fmt.Println("available ip", len(config.Config.SourceIP))
	fmt.Println("available source email", len(config.Config.SourceEmail))
	//step 2.get verifier email
	inData := reader.InitSourceData()

	//map ip:source_email
	var (
		ip          = config.Config.SourceIP
		sourceEmail = config.Config.SourceEmail
	)

	// todo use map ip:sourceEmail
	currentLen := len(config.Config.SourceIP)
	verifiers := make([]*verifier.Verifier, 0, currentLen)
	for i := 0; i < currentLen; i++ {
		v := verifier.NewVerifier(ip[i], sourceEmail[i])
		verifiers = append(verifiers, v)
	}
	//step 3.init factory
	funcFactory := &VerifierFactory{
		InData:       make([]*entity.EmailInfo, 0),
		OutChannel:   make(chan *entity.EmailInfo, dataCount),
		Verifier:     make([]*verifier.Verifier, 0),
		ErrEmailChan: make(chan *entity.EmailInfo, dataCount),
	}
	Verifiers = funcFactory
	Verifiers.InData = inData
	Verifiers.Verifier = verifiers

	// log
	//fmt.Println(Verifiers.InData)
	//fmt.Println(Verifiers.Verifier)

	// start verify
	var currency int
	currency = 200
	dataLen := len(Verifiers.InData)
	//vLen := len(Verifiers.Verifier) // currency = vLen
	nextWorker := dataLen / currency

	fmt.Println(dataLen, nextWorker)

	//workerData
	workerData = make([][]*entity.EmailInfo, 0)
	//inData
	totalData := Verifiers.InData
	for i := 0; i < currency+1; i++ {
		group := make([]*entity.EmailInfo, 0)
		if (i+1)*nextWorker < dataLen {
			group = totalData[i*nextWorker+1 : (i+1)*nextWorker]
			workerData = append(workerData, group)
		} else {
			// last group
			group = totalData[i*nextWorker+1:]
			for _, v := range group {
				workerData[currency-1] = append(workerData[currency-1], v)
			}
		}
	}
	//map reduce
	//fmt.Println(len(workerData),len(workerData[139]),workerData[139][3409])
	//currency = vLen

	wg := new(sync.WaitGroup)
	for i := 0; i < currency; i++ {
		wg.Add(1)
		v := verifier.NewVerifier(GenPublicIP(), "insomnus@lovec.at")
		go Do(v.Verify, workerData[i], wg)
	}

	wg.Wait()
	close(Verifiers.OutChannel)

	done := make(chan struct{})
	go writer.Write2CSV(Verifiers.OutChannel, done)
	<-done
	fmt.Println(len(inData), trueF, failT, errEmailCnt, nullEmailCnt)
}

func Do(f Verify, data []*entity.EmailInfo, wg *sync.WaitGroup) {

	for _, v := range data {
		e := &entity.EmailInfo{
			Email:     v.Email,
			SrcStatus: v.SrcStatus,
		}
		out, err := verifier.NewVerifier(GenPublicIP(), "insomnus@lovec.at").Verify(v.Email)
		if err != nil {
			fmt.Printf("errEmailCnt: %d , err :%v , email : %s\n", atomic.AddInt32(&errEmailCnt, 1), err, v.Email)
			apiEmail := client.EmailProxyClientAPI(e.Email)
			if apiEmail.Email == "" {
				fmt.Printf("go: nullEmailCnt: %d  \n", atomic.AddInt32(&nullEmailCnt, 1))
			}
			//apiEmail.Email = e.Email
			Verifiers.OutChannel <- apiEmail
			continue
		}
		fmt.Printf("go: %v ==> cnt: %d \n", out, atomic.AddInt32(&cnt, 1))
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

func GenPublicIP() string {
	for {
		v := genIP()
		if !strings.HasPrefix(v, "10.") && !strings.HasPrefix(v, "192.") && !strings.HasPrefix(v, "172.") && !strings.HasPrefix(v, "127.") {
			return v

		}
	}
}

func genIP() string {
	ip := make(net.IP, net.IPv6len)
	copy(ip, net.IPv4zero)
	binary.BigEndian.PutUint32(ip.To4(), uint32(rand.Uint32()))
	return ip.To16().String()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
