package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/unkaktus/slurmbw"
)

// run with srun --pty -t 00:30:00 -p test -N 2 ./ib-speedtest
func run() error {
	opa := true

	others, err := slurmbw.Others()
	if err != nil {
		return fmt.Errorf("get others: %w", err)
	}
	for _, nodename := range others {
		log.Printf("other: %s", nodename)
	}

	rank, err := slurmbw.Rank()
	if err != nil {
		return fmt.Errorf("get rank: %w", err)
	}

	log.Printf("my rank: %v", rank)
	hostname, _ := slurmbw.GetHostname()

	if opa {
		hostname += "opa"
	}

	if rank == 0 {
		err = slurmbw.Listen(hostname + ":4343")
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
	} else {
		time.Sleep(10 * time.Second)
		rank0_nodename, err := slurmbw.NodenameByRank(0)
		if err != nil {
			return fmt.Errorf("find node with rank 0: %w", err)
		}
		if opa {
			rank0_nodename += "opa"
		}
		conn, err := net.Dial("tcp4", rank0_nodename+":4343")
		if err != nil {
			return fmt.Errorf("dial: %w", err)
		}
		_, err = io.Copy(ioutil.Discard, conn)
		if err != nil {
			return fmt.Errorf("copy: %w", err)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
