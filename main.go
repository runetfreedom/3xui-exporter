package main

import (
	"context"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/publicsuffix"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var httpClient http.Client

var (
	panelUrl      = flag.String("endpoint", getEnv("PANEL_ENDPOINT", "http://localhost/"), "3x-ui endpoint. ")
	panelUserName = flag.String("username", getEnv("PANEL_USERNAME", "admin"), "3x-ui username")
	panelPassWord = flag.String("password", getEnv("PANEL_PASSWORD", "admin"), "3x-ui password")
	panelIp       = flag.String("ip", getEnv("PANEL_IP", ""), "3x-ui ip. Need to be set if the panel is behind a reverse proxy")
	metricsPath   = flag.String("metrics-path", getEnv("METRICS_PATH", "/metrics"), "Metrics path")
	listenAddress = flag.String("listen-address", getEnv("LISTEN_ADDRESS", ":9101"), "The address to listen on for HTTP requests.")
	httpTimeout   = flag.Int("http-timeout", 10, "Timeout for HTTP requests")
)

func main() {
	flag.Parse()

	*panelUrl = strings.TrimRight(*panelUrl, "/")

	log.Println("3x-ui endpoint: ", *panelUrl)

	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	httpClient = http.Client{
		Jar:     jar,
		Timeout: time.Duration(*httpTimeout) * time.Second,
	}

	if *panelIp != "" {

		dialer := &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}

		//default transport wth custom dialer
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				host := strings.Split(addr, ":")[0]
				port := strings.Split(addr, ":")[1]

				if url, err := url.Parse(*panelUrl); err == nil {
					if host == url.Hostname() {
						addr = *panelIp + ":" + port
					}
				}

				return dialer.DialContext(ctx, network, addr)
			},
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}

	err := login()
	if err != nil {
		log.Fatal(err)
	}

	collector := NewCollector()
	prometheus.MustRegister(collector)

	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
