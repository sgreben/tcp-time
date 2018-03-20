package main

import (
	"context"
	"log"
	"net"
	"net/http/httptrace"
	"time"
)

func dialDuration(url string) (time.Duration, error) {
	log.Println("dialing", url)
	start := time.Now()
	conn, err := net.Dial("tcp", url)
	end := time.Now()
	if err == nil {
		conn.Close()
	} else {
		log.Println("error", err)
	}
	delta := end.Sub(start)
	log.Println("duration", delta)
	return delta, err
}

func connectDuration(addr string) (time.Duration, error) {
	var start, end time.Time
	valid := false
	trace := httptrace.ClientTrace{
		DNSDone: func(_ httptrace.DNSDoneInfo) { start = time.Now() },
		GotConn: func(_ httptrace.GotConnInfo) { end = time.Now(); valid = true },
	}
	ctx := httptrace.WithClientTrace(context.Background(), &trace)
	var dialer net.Dialer
	log.Println("connecting", addr)
	start = time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if !valid {
		end = time.Now()
	}
	if err == nil {
		conn.Close()
	} else {
		log.Println("error", err)
	}
	delta := end.Sub(start)
	log.Println("duration", delta)
	return delta, err
}
