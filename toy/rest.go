package toy

import (
	"github.com/gorilla/mux"
	"gopkg.in/pg.v5"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var concurrentConnections int64

type apiFunction func(db *pg.DB, w http.ResponseWriter, req *http.Request)

type RTE struct {
	Route  string
	Method string      //HTTP Method
	Fn     apiFunction //Function is always pass Cfg and Context as well, http response and request
}

type RouteTable []RTE

var routeTable RouteTable

func init() {
	routeTable = RouteTable{
		RTE{Route: "/users", Method: "GET", Fn: GetUsers},
	}
}

func logHandler(route string, f func(w http.ResponseWriter, req *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		atomic.AddInt64(&concurrentConnections, 1)
		defer atomic.AddInt64(&concurrentConnections, -1)

		start := time.Now()

		c := atomic.LoadInt64(&concurrentConnections)
		log.Printf("[%d] %s URL: %s - Start\n", c, r.Method, r.RequestURI)

		f(w, r)
		elapsed := time.Since(start)
		ms := int(elapsed.Nanoseconds() / int64(time.Millisecond))
		c = atomic.LoadInt64(&concurrentConnections)
		log.Printf("[%d] %s URL: %s - End Time: %d ms\n", c, r.Method, r.RequestURI, ms)

	})
}
func RegisteRouters(db *pg.DB) *mux.Router {
	router := mux.NewRouter()
	for _, entry := range routeTable {
		f := entry.Fn
		fn := logHandler(entry.Route, func(w http.ResponseWriter, req *http.Request) {
			f(db, w, req)
		})
		router.HandleFunc(entry.Route, fn).Methods(entry.Method, "OPTIONS")
	}
	return router
}

func GetUsers(db *pg.DB, w http.ResponseWriter, req *http.Request) {
	log.Printf("GetUsers - starting...n")
	vars := mux.Vars(req)
	userId := vars["userid"]
	//get user from postgre
	id, _ := strconv.ParseInt(userId, 10, 0)
	user, err := FetchUser(db, id)
	if err != nil {
		MarshalStringAndHttpWriteStatus(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user != nil {
		MarshalJsonAndHttpWriteStatus(w, user, http.StatusOK)
	}
}
