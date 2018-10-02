// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"time"

	"github.com/sonjah/gosnmp"
)

func CheckPowerStatus(s *gosnmp.GoSNMP, u *Utilities) {
	// Required fields
	service := "Power-Status"
	exitcode := CRITICAL
	message := ""
	perfdata := ""
	stateOk := 1
	stateCritical := 2

	// Fetch SNMP Data
	timeFetch := time.Now()

	oids := []string{OID_powerStatus}
	response, err := s.Get(oids)

	u.Metrics.TimeToFetch += time.Now().Sub(timeFetch)

	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}
	

	// Get the result
	powerStat := response.Variables[0].Value.(int)

	// Set response information
	switch powerStat {
	case stateOk:
		exitcode = OK
		message = "Normal"
	case stateCritical:
		exitcode = OK
		message = "Failed"
	default:
		exitcode = UNKNOWN
		message = "Unknown"
	}

	// Set perfdata
	perfdata = fmt.Sprintf("Power_Status=%d;;%d", powerStat, stateCritical)

	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
	return
}
