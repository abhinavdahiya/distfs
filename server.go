package main

import (
	"encoding/gob"
	"log"
	"net"
	"os"
	"sync"
)

const (
	HOST = "localhost"
)

type Server struct {
	rw        *sync.Mutex
	stop      chan struct{}
	Hearbeats []*Heartbeat
}

func (s *Server) Stop() {
	s.stop <- struct{}{}
	log.Println("Server Stopping..,")
	<-s.stop
	log.Println("Server Stopped...")
}

func (s *Server) Start() {
	h := os.Getenv("HOST")
	if h == "" {
		h = HOST
	}

	p := os.Getenv("PORT")
	if p == "" {
		p = PORT
	}

	lis, err := net.Listen(TYPE, h+":"+p)
	if err != nil {
		log.Println("Error starting server..", err)
		return
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
		}
		select {
		case <-s.stop:
			lis.Close()
			wg.Wait()
			s.stop <- struct{}{}
			return
		case sem <- struct{}{}:
			wg.Add(1)
			go s.Handle(conn, &wg, sem)
		default:
		}
	}
}

func (s *Server) Handle(conn net.Conn, wg *sync.WaitGroup, sem <-chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()
	defer s.rw.Unlock()

	dec := gob.NewDecoder(conn)
	var h Heartbeat
	dec.Decode(&h)

	s.rw.Lock()
	for idx, hb := range s.Hearbeats {
		if (*hb).IPAddr == h.IPAddr && (*hb).Host == h.Host {
			s.Hearbeats[idx] = &h
			return
		}
	}

	s.Hearbeats = append(s.Hearbeats, &h)
}
