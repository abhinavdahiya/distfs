package file

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	rw    *sync.Mutex
	stop  chan struct{}
	Files []*File
}

func (s *Server) Stop() {
	s.stop <- struct{}{}
	log.Println("Server stopping...")
	<-s.stop
	log.Println("Server stopped...")
}

func (s *Server) Start() {
	h := os.Getenv("HOST")
	if h == "" {
		h = HOST
	}
	p := os.Getenv("PORT_FILE")
	if p == "" {
		p = PORT
	}
	lis, err := net.Listen(TYPE, r+":"+p)
	if err != nil {
		log.Println("Error starting server...", err)
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
			go s.ReceiveFile(conn, &wg, sem)
		}
	}
}

func (s *Server) ReceiveFile(conn net.Conn, wg *sync.WaitGroup, sem <-chan struct{}) {
	defer wg.Done()
	defer func() { <-sem }()

	bname := make([]byte, 64)
	bsize := make([]byte, 10)

	conn.Read(bname)
	conn.Read(bsize)

	size, err := strconv.ParseInt(strings.Trim(string(bsize), ":"), 10, 64)
	if err != nil {
		log.Println("Error parsing buffer size", err)
		return
	}
	name := strings.Trim(string(bname), ":")

	var file []byte
	var rbyt int64
	line := &File{
		Name: name,
		Size: size,
	}

	for rbyt = 0; size-rbyts > BUFFER_SIZE; rbyt = rbyt + BUFFER_SIZE {
		conn.Read(file)
	}
	conn.Read(file)
	s.AddFile(line)

}

func (s *Server) AddFile(f File) {
	defer s.rw.Unlock()
	s.rw.Lock()
	f.AddedOn = time.Now()
	s.Files = append(s.Files, &f)
}

type File struct {
	AddedOn time.Time
	Name    string
	Size    string
	Dests   []string
}
