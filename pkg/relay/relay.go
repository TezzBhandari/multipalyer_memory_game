package relay

import (
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
)

type ServerMessage struct {
	Counter int `json:"counter"`
}

type Relay struct {
	id      int
	players map[int]*Conn

	send int
	sync.Mutex
}

func NewRelay() *Relay {
	return &Relay{
		players: make(map[int]*Conn),
		send:    runtime.NumCPU(),
	}
}

func (r *Relay) relay(data []byte) {
	r.Lock()
	wg := &sync.WaitGroup{}

	msgCount := len(r.players)/r.send + 1

	curr := make([]*Conn, 0, msgCount)
	for _, conn := range r.players {
		if len(curr) == msgCount {
			wg.Add(1)
			go r.broadcast(curr, data, wg)
			curr = make([]*Conn, 0, msgCount)
		}
		curr = append(curr, conn)
	}

	if len(curr) > 0 {
		wg.Add(1)
		go r.broadcast(curr, data, wg)
	}

	wg.Wait()
	r.Unlock()
}

func (r *Relay) broadcast(curr []*Conn, msg []byte, s *sync.WaitGroup) {
	for _, conn := range curr {
		conn.msg(msg)
	}
	s.Done()
}

func (r *Relay) remove(id int) {
	r.Lock()
	delete(r.players, id)
	r.Unlock()
}

func (r *Relay) Add(w *websocket.Conn) {
	r.Lock()
	r.id++
	id := r.id
	r.Unlock()

	c := &Conn{
		id:    id,
		conn:  w,
		msgs:  make(chan []byte, 10),
		relay: r,
	}

	r.Lock()
	r.players[id] = c
	r.Unlock()

	go c.Read()
	go c.Write()
}
