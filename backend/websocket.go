package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"test-fullstack-loyalty/backend/consts"
	"test-fullstack-loyalty/backend/model"
	"test-fullstack-loyalty/backend/riderOps"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Proxy struct {
	wsConn         *websocket.Conn
	riderOps       *riderOps.RiderOperation
	riderPusher    chan model.Rider
	connectedRider string
}

func newProxy(riderOps *riderOps.RiderOperation, riderPusher chan model.Rider) *Proxy {
	if riderOps == nil {
		return nil
	}
	return &Proxy{nil, riderOps, riderPusher, ""}
}

// Run enalbes websocket listning on port 3000
func (p *Proxy) Run() {
	http.HandleFunc("/ws", p.webSocketHandler)
	err := http.ListenAndServe(*wsPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Proxy) webSocketHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("[Error] Websocket upgrader	err %v:", err)
		return
	}
	p.wsConn = wsConn
	go p.listenRiderRequest()
	go p.riderLiveUpdate()
}

// listenRiderRequest listens request from frontend
func (p *Proxy) listenRiderRequest() {
	for {
		_, data, err := p.wsConn.ReadMessage()
		//log.Printf("..... getting messages %v, %v\n", string(data), err)
		if err != nil {
			log.Printf("Connection closed")
			p.wsConn.Close()
			return
		}

		var request model.Request
		if json.Unmarshal(data, &request) != nil {
			log.Printf("[Error] Fail to marshal rquest %v", data)
			continue
		}

		if request.UserId != p.connectedRider {
			p.connectedRider = request.UserId
		}
		p.handleRequest(request)
	}
}

// handleRequest process request from frontend
func (p *Proxy) handleRequest(req model.Request) {
	switch string(req.Type) {
	case string(consts.RIDER_RECORD_REQUEST):
		rider, _ := p.riderOps.GetRider(req.UserId)
		p.riderPusher <- *rider
		break
	}
}

// riderLiveUpdate push live update to frontend
func (p *Proxy) riderLiveUpdate() {
	for rider := range p.riderPusher {
		var riderIdStr = strconv.FormatInt(rider.Id, 10)
		if riderIdStr == p.connectedRider {
			resp := &model.Response{string(consts.RIDER_RECORD_RESPONSE), rider}
			p.send(resp)
		}
	}
}

// send wraps websocket write
func (p *Proxy) send(data interface{}) {
	if err := p.wsConn.WriteJSON(data); err != nil {
		log.Print("[Error] Fail to send	data via websocket %v:", err)
		return
	}
}
