package service

import (
	"io"
	"os/exec"
	"regexp"
	"time"

	"github.com/isaquecsilva/diskstatus/unit"
)

type DiskSpaceService struct {
	cmd       []string
	interval  time.Duration
	byteUnit  unit.ByteUnit
	conversor unit.UnityConversionInterface
	re        *regexp.Regexp
}

func New(interval time.Duration, byteUnit unit.ByteUnit) *DiskSpaceService {
	dss := new(DiskSpaceService)
	dss.re = DefaultRegex
	dss.cmd = []string{"wmic", "logicaldisk", "get", "freespace"}
	dss.conversor = unit.Conversor(true)
	dss.interval = interval
	dss.byteUnit = byteUnit
	return dss
}

func (dss *DiskSpaceService) Execute(w io.Writer) {
	for {
		cmd := exec.Command(dss.cmd[0], dss.cmd[1:]...)
		out, err := cmd.CombinedOutput()

		if err != nil {
			println(err.Error())
			continue
		}

		if out = dss.getOnlyDigits(out); out == nil {
			println("no matches")
			continue
		}

		bytes_, err := dss.conversor.ToByte(string(out), dss.byteUnit)

		if err != nil {
			println("unity_conversion_error: ", err.Error())
			continue
		}

		w.Write([]byte(bytes_ + ","))
		time.Sleep(dss.interval)
	}
}

func (dss *DiskSpaceService) getOnlyDigits(b []byte) []byte {
	return dss.re.Find(b)
}
