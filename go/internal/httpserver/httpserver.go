package httpserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/acme/autocert"

	"context"

	"github.com/dimfeld/httptreemux"
	"vest.com/vest/go/internal/errutil"

	"vest.com/vest/go/internal/api"
	"vest.com/vest/go/internal/log"

	"github.com/rs/cors"
)

const (
	baseName  = "vest"
	codeParam = "code"
)

type handler interface {
	serve(http.ResponseWriter, *http.Request) error
}

func writeRespErr(err error, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	httpErr, ok := err.(httpError)
	if !ok {
		log.Error(ctx, err)
		httpErr = errInternalServerError
	}

	if k := contextRespKey(ctx); k != nil {
		k.key = httpErr.key()
	}

	for key := range w.Header() {
		if _, ok := httpErr.headers()[key]; ok {
			continue
		}
		delete(w.Header(), key)
	}

	if r.Method == http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpErr.code())
		return
	}

	response := struct {
		Key     string `json:"key"`
		Message string `json:"message"`
	}{
		httpErr.key(),
		httpErr.msg(),
	}

	b, err := marshal(ctx, response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.code())

	switch _, err := w.Write(b); {
	case err == nil:
	default:
		panic(errutil.Wrap(err))
	}
}

type commonParamHandler struct {
	innerH http.Handler
}

func (h commonParamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := contextWithPretty(r.Context(), false)

	for param, paramVals := range r.URL.Query() {
		if len(paramVals) > 1 {
			writeRespErr(msgError{errDuplicatedParam, fmt.Sprintf("Query parameter '%v' duplicated", param)}, w, r.WithContext(ctx))
			return
		}

		ctx = contextWithParam(ctx, param, paramVals[0])
	}

	h.innerH.ServeHTTP(w, r.WithContext(ctx))
}

type panicLogHandler struct {
	innerH http.Handler
}

func (h panicLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx = contextWithPublicParams(ctx)
	ctx = contextWithPublicPathParams(ctx)

	r = r.WithContext(ctx)

	defer func() {
		rec := recover()
		if rec == nil {
			return
		}
		if rec == http.ErrAbortHandler {
			panic(http.ErrAbortHandler)
		}
		if err, ok := rec.(error); ok {
			log.Error(errutil.Wrap(err))
		} else {
			log.Error(errutil.New(r))
		}
		panic(http.ErrAbortHandler)
	}()

	r = r.WithContext(ctx)

	h.innerH.ServeHTTP(w, r)
}

type erroringHandler struct {
	innerH handler
}

func (h erroringHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.innerH.serve(w, r)
	if err == nil {
		return
	}

	ctxErr := ctx.Err()

	switch ctxErr {
	case context.Canceled:
		writeRespErr(errCanceled, w, r)
	case context.DeadlineExceeded:
		writeRespErr(errTimeout, w, r)
	default:
		writeRespErr(err, w, r)
	}
}

type unknownPublicParamHandler struct {
	innerH       http.Handler
	params       map[string]struct{}
	paramsPrefix string
	puParams     map[string]struct{}
	puPathParams map[string]struct{}
}

func (h unknownPublicParamHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	prefixedParams := make(map[string]struct{})
	for param := range r.URL.Query() {
		if h.paramsPrefix != "" && strings.HasPrefix(param, h.paramsPrefix) {
			prefixedParams[param] = struct{}{}
			continue
		}

		var ok bool
		if h.params != nil {
			_, ok = h.params[param]
		}
		if !ok {
			writeRespErr(msgError{errUnknownParam, fmt.Sprintf("Query parameter '%v' unknown", param)}, w, r)
			return
		}
	}

	if pp := contextPublicParams(r.Context()); pp != nil {
		for k := range h.puParams {
			pp[k] = struct{}{}
		}
	}

	if ppp := contextPublicPathParams(r.Context()); ppp != nil {
		for k := range h.puPathParams {
			ppp[k] = struct{}{}
		}
	}

	if len(prefixedParams) > 0 {
		r = r.WithContext(contextWithPrefixedParams(r.Context(), h.paramsPrefix, prefixedParams))
	}

	h.innerH.ServeHTTP(w, r)
}

type methodNotAllowedHandler struct{}

func (methodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeRespErr(errMethodNotAllowed, w, r)
}

type notFoundHandler struct{}

func (notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeRespErr(errNotFound, w, r)
}

type HTTPServer struct {
	HTTPPort  int
	HTTPSPort int
	api.APIer

	shutdown       chan struct{}
	doneServeHTTP  chan struct{}
	doneServeHTTPS chan struct{}
}

type routeMethods map[string]http.Handler
type routePaths map[string]routeMethods
type router func(http.Handler) http.Handler

type routeBuilder struct {
	api    api.APIer
	routes func(http.Handler) http.Handler
	extra  routePaths
}

func (b routeBuilder) build() http.Handler {
	nopMiddle := func(h http.Handler) http.Handler {
		return h
	}
	if b.routes == nil {
		b.routes = nopMiddle
	}

	chartsGet := newChartsGetHandler(b.api)
	statsGet := newStatsGetHandler(b.api)
	navigationGet := newNavigationGetHandler(b.api)

	routes := func() routePaths {
		return routePaths{
			"v1/charts":     {http.MethodHead: chartsGet, http.MethodGet: chartsGet},
			"v1/stats":      {http.MethodHead: statsGet, http.MethodGet: statsGet},
			"v1/navigation": {http.MethodHead: navigationGet, http.MethodGet: navigationGet},
		}
	}

	r := routes()

	for p, methods := range b.extra {
		if _, ok := r[p]; !ok {
			r[p] = methods
			continue
		}
		for m, h := range methods {
			if _, ok := r[p][m]; ok {
				panic(fmt.Sprint("colliding route", p, m))
			}
			r[p][m] = h
		}
	}

	finalRoutes := make(routePaths)

	for p, methods := range r {
		p = "/" + baseName + "/" + p
		for m, h := range methods {
			if _, ok := finalRoutes[p]; !ok {
				finalRoutes[p] = make(routeMethods)
			}
			finalRoutes[p][m] = b.routes(h)

			finalRoutes[p][http.MethodOptions] = b.routes(notFoundHandler{})
		}
	}

	finalRoutes["/"] = make(routeMethods)
	for _, m := range []string{http.MethodOptions, http.MethodHead, http.MethodGet} {
		finalRoutes["/"][m] = b.routes(notFoundHandler{})
	}

	return makeRouter(finalRoutes, b.routes(methodNotAllowedHandler{}), b.routes(notFoundHandler{}))
}

func makeRouter(r routePaths, methodNotAllowed, notFound http.Handler) http.Handler {
	router := httptreemux.NewContextMux()
	router.NotFoundHandler = notFound.ServeHTTP
	router.MethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request, methods map[string]httptreemux.HandlerFunc) {
		for m := range methods {
			w.Header().Add("Allow", m)
		}
		methodNotAllowed.ServeHTTP(w, r)
	}
	router.RedirectCleanPath = false
	router.RedirectTrailingSlash = false
	router.PathSource = httptreemux.URLPath

	for p, methods := range r {
		for m, h := range methods {
			router.Handler(m, p, h)
		}
	}

	return router
}

func (s *HTTPServer) Start() {
	s.shutdown = make(chan struct{})
	s.doneServeHTTP = make(chan struct{})
	s.doneServeHTTPS = make(chan struct{})

	go s.serveHTTP()
	go s.serveHTTPS()
}

func (s *HTTPServer) Close() error {
	close(s.shutdown)
	<-s.doneServeHTTP
	<-s.doneServeHTTPS
	return nil
}

func (s *HTTPServer) serveHTTP() {
	defer close(s.doneServeHTTP)

	crs := corsHandler()

	ctx := context.Background()

	nonBatchRoutes := func() router {
		return func(h http.Handler) http.Handler {
			return panicLogHandler{
				crs.Handler(
					commonParamHandler{
						h,
					},
				),
			}
		}
	}

	builder := routeBuilder{
		api:    s.APIer,
		routes: nonBatchRoutes(),
	}

	h := builder.build()

	for {
		func() {
			addr := ":" + strconv.Itoa(s.HTTPPort)

			listener, err := tcpListen(addr)
			if err != nil {
				log.Error(err)
				return
			}

			server := &http.Server{Handler: h}

			serverClosed := make(chan struct{})
			go func() {
				defer close(serverClosed)

				log.Info("serving http", "address", addr)
				switch err := server.Serve(listener); err {
				case nil, http.ErrServerClosed:
				default:
					log.Error(err)
				}
			}()

			select {
			case <-s.shutdown:
			case <-serverClosed:
			}

			if err := server.Shutdown(ctx); err != nil {
				log.Error(err)
			}
		}()

		select {
		case <-s.shutdown:
			return
		case <-time.After(time.Second):
		}
	}
}

func (s *HTTPServer) serveHTTPS() {
	defer close(s.doneServeHTTPS)

	crs := corsHandler()

	ctx := context.Background()

	nonBatchRoutes := func() router {
		return func(h http.Handler) http.Handler {
			return panicLogHandler{
				crs.Handler(
					commonParamHandler{
						h,
					},
				),
			}
		}
	}

	builder := routeBuilder{
		api:    s.APIer,
		routes: nonBatchRoutes(),
	}

	h := builder.build()

	for {
		func() {
			addr := ":" + strconv.Itoa(s.HTTPSPort)

			m := &autocert.Manager{
				Cache:  autocert.DirCache(".cert"),
				Prompt: autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(
					"thevest.org",
					"www.thevest.org",
					"api.thevest.org",
					"ec2-54-84-182-191.compute-1.amazonaws.com"),
			}

			server := &http.Server{Handler: h, Addr: ":https", TLSConfig: m.TLSConfig()}

			serverClosed := make(chan struct{})
			go func() {
				defer close(serverClosed)

				log.Info("serving https", "address", addr)
				switch err := server.ListenAndServeTLS("", ""); err {
				case nil, http.ErrServerClosed:
				default:
					log.Error(err)
				}
			}()

			select {
			case <-s.shutdown:
			case <-serverClosed:
			}

			if err := server.Shutdown(ctx); err != nil {
				log.Error(err)
			}
		}()

		select {
		case <-s.shutdown:
			return
		case <-time.After(time.Second):
		}
	}
}

func marshal(ctx context.Context, v interface{}) ([]byte, error) {
	if contextPretty(ctx) {
		return json.MarshalIndent(v, "", "\t")
	}
	return json.Marshal(v)
}

func tcpListen(addr string) (net.Listener, error) {
	const network = "tcp"
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	tcpLn, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		return nil, err
	}
	return tcpLn, nil
}

func corsHandler() *cors.Cors {
	return cors.New(
		cors.Options{
			AllowedMethods: []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete},
			AllowedHeaders: []string{"*"},
			ExposedHeaders: []string{"Allow", "Content-Length", "Content-Encoding", "Content-Type", "Date", "ETag", "Vary", "Location"},
			MaxAge:         3600,
		},
	)
}
