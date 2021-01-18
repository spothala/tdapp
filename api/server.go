package hook

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spothala/go-http-api/utils"
)

type Server struct {
	mux      *http.ServeMux
	appURL   string
	clientID string
	handlers map[string]bool
	debug    bool
}

// newServer - Create a new Web Server with supported handlers
func newServer(clientID string, debug bool) *Server {
	// Constructing Server Object
	s := &Server{
		mux:      http.NewServeMux(),
		clientID: clientID,
		appURL:   "https://auth.tdameritrade.com/auth?response_type=code&redirect_uri=http://localhost:8080&client_id=" + clientID + "%40AMER.OAUTHAP",
		handlers: make(map[string]bool),
		debug:    debug,
	}

	s.addHandler("/status", s.status)
	s.addHandler("/callback", s.callback)
	return s
}

// Start - Start the server to accept supported requests
func Start(port int, clientID string, debug bool) (*http.Server, *Server) {
	localServer := newServer(clientID, debug)
	h := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: localServer,
	}
	go func() {
		fmt.Println("Server started listening on Port: " + strconv.Itoa(port))
		if httpError := h.ListenAndServe(); httpError != nil {
			log.Println("While serving HTTP: ", httpError)
		}
	}()
	return h, localServer
}

func (s *Server) addHandler(pattern string, handler http.HandlerFunc) {
	s.handlers[pattern] = true
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) status(w http.ResponseWriter, req *http.Request) {
	utils.RespondJson(w, map[string]string{"health": "OK"})
}

func (s *Server) callback(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Println(req.URL.Query().Get("code"))
		if req.URL.Query().Get("code") != "" {
			headers := url.Values{}
			headers.Set("grant_type", "authorization_code")
			headers.Set("access_type", "offline")
			headers.Set("code", req.URL.Query().Get("code"))
			jsonOut, err := s.API("POST", headers, "/oauth2/token", nil, nil)
			if err != nil {
				return utils.RespondError(w, errors.New("Only GET method allowed on this endpoint"), http.StatusMethodNotAllowed)
			}
			//utils.WriteJsonToFile(jsonOut, instagram.AccessTokenFile)
		}
	default:
		utils.RespondError(w, errors.New("Only GET method allowed on this endpoint"), http.StatusMethodNotAllowed)
	}
}

func (s *Server) hook(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		s.ParseGithubRequest(w, req)
	default:
		utils.RespondError(w, errors.New("Only POST method allowed on this endpoint"), http.StatusMethodNotAllowed)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if enabled, ok := s.handlers[r.URL.Path]; !ok && !enabled {
		utils.RespondError(w, errors.New("This routing is not defined"), http.StatusNotImplemented)
		return
	}
	s.mux.ServeHTTP(w, r)
}

// GracefulShutdown - Shutdowns the server gracefully with timeout
func GracefulShutdown(hs *http.Server, timeout time.Duration) {
	// Created stop channel to wait on until it gets terminated signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Timeout to make sure all requests are completed
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Shutting down the server
	fmt.Printf("\nShutdown with timeout: %s\n", timeout)
	if err := hs.Shutdown(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Server stopped")
	}
}
