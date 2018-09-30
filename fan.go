// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"strconv"

	"github.com/sonjah/gosnmp"
)

func CheckFanStatus(s *gosnmp.GoSNMP, u *Utilities) {
	// Required fields
	service := "Fan"
	exitcode := OK
	message := ""
	perfdata := ""

	// Fetch SNMP Data
	oids := []string{OID_systemFanStatus, OID_CPUFanStatus}
	result, err := s.Get(oids)

	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}

	// Get the result
	systemFanStat := result.Variables[0].Value.(int)
	cpuFanStat := result.Variables[1].Value.(int)

	// Set message
	message = fmt.Sprintf("System-Fan:%s, CPU-Fan:%s", fanStatusName(systemFanStat), fanStatusName(cpuFanStat))

	// Set ExitCode
	switch {
	case systemFanStat == 1 && cpuFanStat == 1:
		exitcode = OK
	case systemFanStat == 2 || cpuFanStat == 2:
		exitcode = OK
	default:
		exitcode = UNKNOWN
	}

	// Set perfdata
	perfdata = fmt.Sprintf("SystemFanStatus=%s CPUFanStatus=%s", strconv.Itoa(systemFanStat), strconv.Itoa(cpuFanStat))

	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
	return
}

func fanStatusName(s int) string {
	switch s {
	case 1:
		return "Normal"
	case 2:
		return "Failure"
	default:
		return "Unknown"
	}
}
