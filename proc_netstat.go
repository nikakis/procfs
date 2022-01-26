package procfs

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/nikakis/procfs/internal/util"
	"log"
	"strings"
)

type IpExt struct {
	InNoRoutes      string
	InTruncatedPkts string
	InMcastPkts     string
	OutMcastPkts    string
	InBcastPkts     string
	OutBcastPkts    string
	InOctets        string
	OutOctets       string
	InMcastOctets   string
	OutMcastOctets  string
	InBcastOctets   string
	OutBcastOctets  string
	InCsumErrors    string
	InNoECTPkts     string
	InECT1Pkts      string
	InECT0Pkts      string
	InCEPkts        string
	ReasmOverlaps   string
}

// ProcNetstat provides status information about the process,
// read from /proc/[pid]/net/netstat.
type ProcNetstat struct {
	IpExt
}

func (p Proc) NetStat() (ProcNetstat, error) {
	data, err := util.ReadFileNoStat(p.path("netstat"))

	if err != nil {
		log.Fatal(err)
		return ProcNetstat{}, err
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(data))

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	procNetstat := ProcNetstat{}

	for index, line := range lines {
		// If the line contains the header `IpExt` and is the last line
		// then this is line we need to parse
		if strings.Contains(line, "IpExt") && index == len(lines)-1 {
			ipExt := IpExt{}
			var header string

			_, err := fmt.Sscan(line,
				&header,
				&ipExt.InNoRoutes,
				&ipExt.InTruncatedPkts,
				&ipExt.InMcastPkts,
				&ipExt.OutMcastPkts,
				&ipExt.InBcastPkts,
				&ipExt.OutBcastPkts,
				&ipExt.InOctets,
				&ipExt.OutOctets,
				&ipExt.InMcastOctets,
				&ipExt.OutMcastOctets,
				&ipExt.InBcastOctets,
				&ipExt.OutBcastOctets,
				&ipExt.InCsumErrors,
				&ipExt.InNoECTPkts,
				&ipExt.InECT1Pkts,
				&ipExt.InECT0Pkts,
				&ipExt.InCEPkts,
				&ipExt.ReasmOverlaps,
			)
			if err != nil {
				log.Println(err)
				return ProcNetstat{}, err
			}
			procNetstat.IpExt = ipExt
		}
	}
	fmt.Printf("%+v", procNetstat)
	return procNetstat, nil

}
