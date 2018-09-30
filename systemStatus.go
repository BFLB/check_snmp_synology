// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"strconv"

	"github.com/sonjah/gosnmp"
)

func CheckSystemStatus(s *gosnmp.GoSNMP, u *Utilities) {
	// Required fields
	service := "System-Status"
	exitcode := OK
	message := ""
	perfdata := ""

	// Fetch SNMP Data
	oids := []string{OID_systemStatus}
	result, err := s.Get(oids)

	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}

	// Get the result
	sysStat := result.Variables[0].Value.(int)

	// Set response information
	switch sysStat {
	case 1:
		exitcode = OK
		message = "Normal"
	case 2:
		exitcode = OK
		message = "Failed"
	default:
		exitcode = UNKNOWN
		message = "Unknown"
	}

	// Set perfdata
	perfdata = fmt.Sprintf("SystemStatus=%s", strconv.Itoa(sysStat))

	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
	return
}
