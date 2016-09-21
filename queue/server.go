package queue

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"os"
	"os/signal"
	"syscall"
)

func Serve(listen, dir string) {

	doneChan := make(chan bool)
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	reQueue := regexp.MustCompile(`^/[a-zA-Z0-9\-_]+$`)
	reEnqueue := regexp.MustCompile(`^/[a-zA-Z0-9\-_]+/enqueue$`)
	reDequeue := regexp.MustCompile(`^/[a-zA-Z0-9\-_]+/dequeue$`)

	s := Service{newManager(dir)}

	handler := func (w http.ResponseWriter, r *http.Request) {
		path := r.RequestURI
		verb := r.Method
		fmt.Println(verb, path)
		switch {
		case reEnqueue.MatchString(path) && verb == http.MethodPost:
			s.enqueue(getName(path), w, r)
		case reDequeue.MatchString(path) && verb == http.MethodPost:
			s.dequeue(getName(path), w, r)
		case reQueue.MatchString(path):
			switch verb {
			case http.MethodPost:
				s.createQueue(path[1:], w)
			case http.MethodDelete:
				s.deleteQueue(path[1:], w)
			case http.MethodGet:
				s.getQueue(path[1:], w)
			default:
				http.NotFoundHandler()
			}
		case path == "/" && verb == http.MethodGet:
			s.getStatus(w)
		default:
			http.NotFoundHandler()
		}
	}

	go func() {
		<-signalChan
		fmt.Println("Stopping server ... ")
		s.shutdown()
		doneChan <- true
	}()

	http.HandleFunc("/", handler)

	go func(){
		log.Println("Listening on", listen)
		err := http.ListenAndServe(listen, nil)
		if err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	}()
	<- doneChan
}


func getName(path string) string {
	last := strings.LastIndex(path, "/")
	return path[1:last]
}

type Service struct {
	m *Manager
}

func (s *Service) getStatus(w http.ResponseWriter) {
	w.Write([]byte("Get Status\n"))
}

func (s *Service) enqueue(q string, w http.ResponseWriter, r *http.Request) {
	if data, err := ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		fmt.Println(string(data))
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) dequeue(q string, w http.ResponseWriter, r *http.Request) {

}

func (s *Service) getQueue(q string, w http.ResponseWriter) {
	w.Write([]byte("Get " + q + " \n"))
}

func (s *Service) createQueue(q string, w http.ResponseWriter) {
	w.Write([]byte("Create " + q + " \n"))
}

func (s *Service) deleteQueue(q string, w http.ResponseWriter) {
	w.Write([]byte("Delete " + q + " \n"))
}

func (s *Service) shutdown() {
	s.m.Shutdown()
}
