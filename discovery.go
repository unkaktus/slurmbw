package slurmbw

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetHostname() (string, error) {
	hostname := os.Getenv("SLURMD_NODENAME")
	if hostname == "" {
		return "", fmt.Errorf("No SLURMD_NODENAME set")
	}
	return hostname, nil
}

type NodeList []string

func GetNodeList() (NodeList, error) {
	nodelist_string := os.Getenv("SLURM_JOB_NODELIST")
	if nodelist_string == "" {
		return nil, fmt.Errorf("No SLURM_JOB_NODELIST set")
	}
	sp := strings.Split(nodelist_string, "[")
	prefix := sp[0]
	range_string := strings.Split(sp[1], "]")[0]

	range_split := strings.Split(range_string, "-")
	start, err := strconv.Atoi(range_split[0])
	if err != nil {
		return nil, fmt.Errorf("parse range start: %w", err)
	}
	stop, err := strconv.Atoi(range_split[1])
	if err != nil {
		return nil, fmt.Errorf("parse range start: %w", err)
	}

	nodelist := NodeList{}
	for i := start; i <= stop; i++ {
		nodename := fmt.Sprintf("%s%02d", prefix, i)
		nodelist = append(nodelist, nodename)
	}

	return nodelist, nil
}

func (nodelist NodeList) Others(hostname string) NodeList {
	others := NodeList{}
	for _, nodename := range nodelist {
		if nodename == hostname {
			continue
		}
		others = append(others, nodename)
	}
	return others
}

func Rank() (int, error) {
	hostname, err := GetHostname()
	if err != nil {
		return -1, fmt.Errorf("getting hostname: %w", err)
	}
	log.Printf("hostname: %s", hostname)

	nodelist, err := GetNodeList()
	if err != nil {
		return -1, fmt.Errorf("GetNodeList: %w", err)
	}

	for rank, nodename := range nodelist {
		if nodename == hostname {
			return rank, nil
		}
	}
	return -1, fmt.Errorf("hostname not found")
}

func NodenameByRank(rank int) (string, error) {
	nodelist, err := GetNodeList()
	if err != nil {
		return "", fmt.Errorf("GetNodeList: %w", err)
	}

	if rank >= len(nodelist) {
		return "", fmt.Errorf("rank is to high")
	}

	return nodelist[rank], nil

}

func Others() (NodeList, error) {
	hostname, err := GetHostname()
	if err != nil {
		return nil, fmt.Errorf("getting hostname: %w", err)
	}
	log.Printf("hostname: %s", hostname)

	nodelist, err := GetNodeList()
	if err != nil {
		return nil, fmt.Errorf("GetNodeList: %w", err)
	}

	others := nodelist.Others(hostname)
	return others, nil
}
