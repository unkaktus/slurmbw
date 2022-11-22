package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/unkaktus/slurmbw"
)

// Build with `env GOOS=linux go build -v`
// Run with `srun --pty -t 00:30:00 -p test -N 2 ./slurmbw“
func run() error {
	network := flag.String("network", "udp4", "Network to use (tcp4, udp4, ...)")
	omnipath := flag.Bool("omnipath", false, "Use OmniPath network (IPoIB)")
	flag.Parse()

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

	if *omnipath {
		hostname += "opa"
	}

	if rank == 0 {
		err = slurmbw.Listen(*network, hostname+":4343")
		if err != nil {
			return fmt.Errorf("listen: %w", err)
		}
	} else {
		time.Sleep(10 * time.Second)
		rank0_nodename, err := slurmbw.NodenameByRank(0)
		if err != nil {
			return fmt.Errorf("find node with rank 0: %w", err)
		}
		if *omnipath {
			rank0_nodename += "opa"
		}
		err = slurmbw.Dial(*network, rank0_nodename+":4343")
		if err != nil {
			return fmt.Errorf("dial: %w", err)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
