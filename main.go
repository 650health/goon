package main

import (
	"flag"
	"fmt"
	"os"
)

type protoImplT func(bool, bool, string, string, []string) error

const usage = "Usage: goon -proto <version> [options] -- <program> [<arg>...]"

var protoFlag = flag.String("proto", "", "protocol version (one of: 0.0)")
var ackFlag = flag.String("ack", "", "arbitrary data used during handshake")
var inFlag  = flag.Bool("in", false, "enable reading from stdin")
var outFlag = flag.Bool("out", false, "output program's stdout")
var errFlag = flag.String("err", "nil", "output or redirect stderr")
var dirFlag = flag.String("dir", ".", "working directory for the spawned process")
var logFlag = flag.String("log", "", "enable logging")

func main() {
	flag.Parse()
	args := flag.Args()

	if *protoFlag == "" {
		die_usage("Please specify the protocol version.")
	}

	if *ackFlag != "" {
		os.Exit(performInitialHandshake(*protoFlag, *ackFlag))
	}

	initLogger(*logFlag)
	validateArgs(args)

	/* Run external program and block until it terminates */
	protoImpl := findProtocolImpl(*protoFlag)
	if protoImpl == nil {
		reason := fmt.Sprintf("Unsupported protocol version: %v", *protoFlag)
		die(reason)
	}
	err := findProtocolImpl(*protoFlag)(*inFlag, *outFlag, *errFlag, *dirFlag, args)
	if err != nil {
		os.Exit(getExitStatus(err))
	}
}

func performInitialHandshake(protoFlag, ackstr string) int {
	if findProtocolImpl(protoFlag) == nil {
		return 1
	}
	os.Stdout.WriteString(ackstr)
	return 0
}

func validateArgs(args []string) {
	if len(args) < 1 {
		die_usage("Not enough arguments.")
	}

	logger.Printf("Flag values:\n  proto: %v\n  in: %v\n  out: %v\n  err: %v\n  dir: %v\nArgs: %v\n",
				  *protoFlag, *inFlag, *outFlag, *errFlag, *dirFlag, args)
}

func findProtocolImpl(flag string) (impl protoImplT) {
	switch flag {
	case "0.0":
		impl = proto_0_0
	}
	return
}
