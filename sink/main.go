package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)

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

type Nodes struct {
	nodes []*Node
	sync.Mutex
}

var (
	nodeman *Nodes
)

func main() {
	// Echo instance
	nodeman = new(Nodes)
	e := echo.New()

	e.POST("/node/register", register)
	e.POST("/data/upload", upload)

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}

type SenData struct {
	T string
	V string
}

type Ret struct {
	Code int
	Data interface{}
}

func upload(c echo.Context) error {
	d := new(SenData)
	if err := c.Bind(d); err != nil {
		return err
	}
	log.Infof("receive data:%v", d)
	return nil
}

// Handler
func register(c echo.Context) error {
	req := new(Node)
	if err := c.Bind(req); err != nil {
		return err
	}

	res := make([]*Node, 0)

	nodeman.Lock()
	defer nodeman.Unlock()

	exist := false
	for _, v := range nodeman.nodes {
		if v.ID == req.ID {
			exist = true
		}
		//if v.Loc.X < req.Loc.X && v.Loc.Y < req.Loc.Y && v.Loc.Z < req.Loc.Z {
		if isCloser(req, v) {
			res = append(res, v)
		}
	}

	//TODO: sort res
	for i := 0; i < len(res); i++ {
		for j := i + 1; j < len(res); j++ {
			if isSmaller(res[j], res[i], req) {
				res[i], res[j] = res[j], res[i]
			}
		}
	}
	if !exist {
		nodeman.nodes = append(nodeman.nodes, req)
		log.Infof("Register Sensor:%v", req)
	}
	return c.JSON(http.StatusOK, res)
}

func isSmaller(a, b, sensor *Node) bool {
	if distance(sensor, a) < distance(sensor, b) {
		return true
	}
	return false
}

func isCloser(sensor, relay *Node) bool {
	// 0.
	sink := &Node{Loc: Location{0, 0, 0}}
	if distance(relay, sink) >= distance(sensor, sink) {
		return false
	}
	// 1. replay-sink distance < sensor-sink distance
	/*if distance(relay, sink) < distance(sensor, sink) {
		return true
	}*/
	// 1. sensor-sink distance > sensor-relay distance
	if distance(sensor, sink) <= distance(sensor, relay) {
		return false
	}
	return true
}

func distance(a, b *Node) int {
	xd := a.Loc.X - b.Loc.X
	yd := a.Loc.Y - b.Loc.Y
	zd := a.Loc.Z - b.Loc.Z
	return xd*xd + yd*yd + zd*zd
}
