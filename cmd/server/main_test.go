package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthEndpointReturnsOK(t *testing.T) {
	testServer := httptest.NewServer(newHandler())
	t.Cleanup(testServer.Close)

	response, err := testServer.Client().Get(testServer.URL + "/healthz")
	if err != nil {
		t.Fatalf("get health endpoint: %v", err)
	}
	t.Cleanup(func() { response.Body.Close() })

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected health endpoint status %d, got %d", http.StatusOK, response.StatusCode)
	}
}

func TestNewServerConfiguresHTTPServer(t *testing.T) {
	const address = "127.0.0.1:8080"

	server := newServer(address, newHandler())

	if server.Addr != address {
		t.Fatalf("expected server address %q, got %q", address, server.Addr)
	}
	if server.ReadHeaderTimeout <= 0 {
		t.Fatal("expected a positive read header timeout")
	}
	if server.ReadTimeout <= 0 {
		t.Fatal("expected a positive read timeout")
	}
	if server.WriteTimeout <= 0 {
		t.Fatal("expected a positive write timeout")
	}
	if server.IdleTimeout <= 0 {
		t.Fatal("expected a positive idle timeout")
	}
	if server.MaxHeaderBytes <= 0 {
		t.Fatal("expected a positive maximum header size")
	}
}

func TestServeShutsDownWhenContextIsCanceled(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen for test server: %v", err)
	}

	server := newServer(listener.Addr().String(), newHandler())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- serve(ctx, server, listener, time.Second)
	}()

	waitForHealth(t, "http://"+listener.Addr().String()+"/healthz")
	cancel()

	select {
	case err := <-serveErrors:
		if err != nil {
			t.Fatalf("serve after cancellation: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("serve did not shut down")
	}
}

func waitForHealth(t *testing.T, url string) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		response, err := http.Get(url)
		if err == nil {
			response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("health endpoint did not become ready: %s", url)
}

func TestServeReturnsUnexpectedError(t *testing.T) {
	wantErr := errors.New("listener failed")

	err := serve(context.Background(), newServer("", newHandler()), errorListener{err: wantErr}, time.Second)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected serve error %v, got %v", wantErr, err)
	}
}

func TestServeIgnoresServerClosed(t *testing.T) {
	err := serve(context.Background(), newServer("", newHandler()), errorListener{err: http.ErrServerClosed}, time.Second)
	if err != nil {
		t.Fatalf("expected normal server closure to return nil, got %v", err)
	}
}

func TestServerAddress(t *testing.T) {
	tests := []struct {
		name       string
		configured string
		expected   string
	}{
		{name: "default", expected: ":8080"},
		{name: "configured", configured: "127.0.0.1:9090", expected: "127.0.0.1:9090"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := serverAddress(test.configured); got != test.expected {
				t.Fatalf("expected address %q, got %q", test.expected, got)
			}
		})
	}
}

func TestRunWithContextReturnsListenError(t *testing.T) {
	wantErr := errors.New("listen failed")

	err := runWithContext(context.Background(), ":0", newHandler(), func(_, _ string) (net.Listener, error) {
		return nil, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected listen error %v, got %v", wantErr, err)
	}
}

func TestServeReturnsShutdownTimeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen for test server: %v", err)
	}

	requestStarted := make(chan struct{})
	requestRelease := make(chan struct{})
	server := newServer(listener.Addr().String(), newHandler())
	server.Handler = http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		close(requestStarted)
		<-requestRelease
		response.WriteHeader(http.StatusOK)
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- serve(ctx, server, listener, 10*time.Millisecond)
	}()

	requestErrors := make(chan error, 1)
	go func() {
		response, err := http.Get("http://" + listener.Addr().String() + "/healthz")
		if response != nil {
			response.Body.Close()
		}
		requestErrors <- err
	}()

	select {
	case <-requestStarted:
	case <-time.After(time.Second):
		t.Fatal("request did not start")
	}
	cancel()

	select {
	case err := <-serveErrors:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("expected shutdown timeout, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("serve did not return after shutdown timeout")
	}

	close(requestRelease)
	select {
	case <-requestErrors:
	case <-time.After(time.Second):
		t.Fatal("in-flight request did not finish after release")
	}
	server.Close()
}

type errorListener struct {
	err error
}

func (listener errorListener) Accept() (net.Conn, error) {
	return nil, listener.err
}

func (listener errorListener) Close() error {
	return nil
}

func (listener errorListener) Addr() net.Addr {
	return errorListenerAddr{}
}

type errorListenerAddr struct{}

func (errorListenerAddr) Network() string {
	return "test"
}

func (errorListenerAddr) String() string {
	return "test"
}
