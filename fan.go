// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"strings"

	"github.com/sonjah/gosnmp"
)

func CheckFanStatus(s *gosnmp.GoSNMP, u *Utilities) {
	// Required fields
	service := "Fan"
	exitcode := OK
	message := ""
	perfdata := ""
	stateOk := 1
	stateCritical := 2

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
	case systemFanStat == stateOk && cpuFanStat == stateOk:
		exitcode = OK
	case systemFanStat == stateCritical || cpuFanStat == stateCritical:
		exitcode = OK
	default:
		exitcode = UNKNOWN
	}

	// Set perfdata
	p := []string{
		fmt.Sprintf("SystemFan_Status=%d;;%d", systemFanStat, stateCritical),
		fmt.Sprintf("CPUFan_Status=%d;;%d", cpuFanStat, stateCritical),
	}
	perfdata = fmt.Sprint(strings.Join(p, " "))

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
