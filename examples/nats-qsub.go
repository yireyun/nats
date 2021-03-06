// Copyright 2012-2015 Apcera Inc. All rights reserved.
// +build ignore

package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/nats-io/nats"
)

func usage() {
	log.Fatalf("Usage: nats-sub [-s server] [--ssl] [-t] <subject> <queue-group>\n")
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s] Queue[%s] Pid[%d]: '%s'\n", i, m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")
	var ssl = flag.Bool("ssl", false, "Use Secure Connection")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}

	opts := nats.DefaultOptions
	opts.Servers = strings.Split(*urls, ",")
	for i, s := range opts.Servers {
		opts.Servers[i] = strings.Trim(s, " ")
	}
	opts.Secure = *ssl

	nc, err := opts.Connect()
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	subj, queue, i := args[0], args[1], 0

	nc.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		i += 1
		printMsg(msg, i)
	})

	log.Printf("Listening on [%s]\n", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}
