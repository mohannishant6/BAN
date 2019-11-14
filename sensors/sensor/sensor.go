package sensor

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type Sensor interface {
	ForwardData(d SenData) error
	GenerateData() error
	StartDutyCycle()
	Register() error
	ID() string
}

type Location struct {
	X int
	Y int
	Z int
}

type Node struct {
	ID   string
	Addr string
	Port string
	Loc  Location
}

type RealSensor struct {
	self        Node
	Sink        string
	Interval    int
	Duration    int
	CloserNodes []Node
	DataSet     string
	alive       bool
}

var (
	httpclient = createHTTPClient()
	ErrSleep   = errors.New("sensor is not alive")
	ErrForward = errors.New("data forwarding failed")
)

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(3) * time.Second,
	}

	return client
}

func NewRealCensor(self Node, sink, dataset string, interv int, dur int) *RealSensor {
	return &RealSensor{
		self:     self,
		Sink:     sink,
		Interval: interv,
		Duration: dur,
		DataSet:  dataset,
		alive:    true,
	}
}

func (s *RealSensor) ID() string {
	return s.self.ID
}

func (s *RealSensor) Register() error {
	payload, _ := json.Marshal(s.self)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/node/register", s.Sink), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpclient.Do(req)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	nodes := make([]Node, 0)
	err = json.Unmarshal(body, &nodes)
	if err != nil {
		return err
	}
	s.CloserNodes = nodes
	log.Infof("found closer sensor:%v", s.CloserNodes)
	return nil
}

func (s *RealSensor) ForwardData(d SenData) error {
	if !s.alive {
		return ErrSleep
	}
	//origint := d.T
	var err error
	for _, v := range s.CloserNodes {
		log.Infof("forward data:%v to %v", d, v)
		//d.T = origint + "->" + v.ID
		err = s.senddata(d, fmt.Sprintf("%s:%s", v.Addr, v.Port))
		if err == nil {
			log.Info("forward success")
			return nil
		} else {
			log.Errorf("forward failed:%v", err)
		}
	}
	log.Infof("no forward node avaliable, send data:%v to the sink:%v", d, s.Sink)
	err = s.senddata(d, s.Sink)
	if err != nil {
		log.Errorf("send data to sink failed:%v", err)
	} else {
		log.Infof("send data to sink successfully")
	}
	return nil
}

func (s *RealSensor) senddata(d SenData, addr string) error {
	payload, _ := json.Marshal(d)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/data/upload", addr), bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpclient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ErrForward
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Infof("send data result:%v", string(body))
	return nil
}

func (s *RealSensor) GenerateData() error {
	dat, err := ioutil.ReadFile(s.DataSet)
	if err != nil {
		return err
	}
	dataset := strings.Split(string(dat), "\n")
	idx := 0
	go func() {
		for true {
			//TODO: add critial logic
			for !s.alive {
				time.Sleep(time.Second)
			}
			dvalue := ""
			//for dvalue == "" {
			dvalue = strings.Trim(dataset[idx%len(dataset)], "\t \n")
			//	idx++
			//}
			report := SenData{T: s.self.ID, V: dvalue}
			log.Infof("report:%v", report)
			s.ForwardData(report)
			time.Sleep(time.Second)
			idx++
		}
	}()
	return nil
}

func (s *RealSensor) StartDutyCycle() {
	s.alive = true
	duration := time.Duration(s.Duration) * time.Second
	interval := time.Duration(s.Interval) * time.Second

	go func() {
		for true {
			time.Sleep(duration)
			s.alive = false
			log.Warn("sensor goes to sleep")
			time.Sleep(interval)
			s.alive = true
			log.Warn("sensor wakes up")
		}
	}()
}
