package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/go-chi/chi/v5"
	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//go:embed posts
var postsFS embed.FS

//go:embed static
var staticFS embed.FS

//go:embed static/img/logo.webp
var logo []byte

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)

	conf := flag.String("conf", "/etc/sbs.conf", "app config file")
	flag.Parse()

	cfg, err := ReadConfigFile(*conf)
	if err != nil {
		logger.Fatalln(err)
	}

	var mailer Mailer
	if cfg.SendGridAPIKey != "" {
		mailer = NewSendGridMailer(cfg.SendGridAPIKey)
	} else {
		mailer = NewLogMailer()
	}

	app := NewApplication(mailer, logger)

	// setup http.Handler for blog posts
	posts, _ := fs.Sub(postsFS, "posts")
	postsServer := http.FileServer(http.FS(posts))
	gzipPostsServer := gzhttp.GzipHandler(postsServer)

	// setup http.Handler for static files
	static, _ := fs.Sub(staticFS, "static")
	staticServer := http.FileServer(http.FS(static))
	gzipStaticServer := gzhttp.GzipHandler(staticServer)

	r := chi.NewRouter()
	r.Mount("/", app.Router())
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/posts/*", http.StripPrefix("/posts", gzipPostsServer))
	r.Handle("/static/*", http.StripPrefix("/static", gzipStaticServer))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write(logo)
	})
	r.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	addr := fmt.Sprintf("127.0.0.1:%s", cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,

		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// open up the socket listener
	l, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalln(err)
	}

	// let systemd know that we are good to go (no-op if not using systemd)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	logger.Printf("started server on %s\n", addr)

	err = srv.Serve(l)
	if err != nil {
		logger.Fatalln(err)
	}
}
