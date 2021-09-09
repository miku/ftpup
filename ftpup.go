package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

var (
	hostPort    = flag.String("l", "localhost:15201", "hostport to listen on")
	ftpHostPort = flag.String("p", "ftp.ncbi.nlm.nih.gov:21", "ftp host to proxy to")
	ftpTimeout  = flag.Duration("T", 10*time.Second, "ftp timeout")
	maxInflight = flag.Int("X", 3, "max requests in flight at the same time")
)

// UserPassword allows to pass in user:password in flags.
type UserPassword struct {
	User     string
	Password string
}

// String rendering of username and password.
func (u *UserPassword) String() string {
	return fmt.Sprintf("%s:%s", u.User, u.Password)
}

// Set parses credentials string into username and password.
func (u *UserPassword) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return fmt.Errorf("user:password required")
	}
	u.User = parts[0]
	u.Password = parts[1]
	return nil
}

// server proxies FTP requests to a FTP host
type server struct {
	ftpHostPort string
	ftpTimeout  time.Duration
	ftpUsername string
	ftpPassword string

	sem chan bool // bounded concurrency, imposing a minimal limit
}

// ServeHTTP proxies requests to FTP, only single path supported.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.sem <- true
	defer func() {
		<-s.sem
	}()
	path := r.URL.Path
	log.Printf("retrieving %v", path)
	c, err := ftp.Dial(s.ftpHostPort, ftp.DialWithTimeout(s.ftpTimeout))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = c.Login(s.ftpUsername, s.ftpPassword); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	defer func() {
		if err := c.Quit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}()
	resp, err := c.Retr(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Close()
	n, err := io.Copy(w, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("[ok] copied %d bytes from ftp://%s%s to %s", n, s.ftpHostPort, path, r.RemoteAddr)
}

func main() {
	var ftpUserPass = UserPassword{
		User:     "anonymous",
		Password: "anonymous",
	}
	flag.Var(&ftpUserPass, "u", "username and password")
	flag.Parse()
	srv := &server{
		ftpHostPort: *ftpHostPort,
		ftpTimeout:  *ftpTimeout,
		ftpUsername: ftpUserPass.User,
		ftpPassword: ftpUserPass.Password,
		sem:         make(chan bool, *maxInflight),
	}
	log.Printf("starting ftpup on http://%v", *hostPort)
	log.Fatal(http.ListenAndServe(*hostPort, srv))
}
