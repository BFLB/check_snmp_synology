// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"strconv"

	"github.com/sonjah/gosnmp"
)

func CheckTemperature(s *gosnmp.GoSNMP, u *Utilities) {
	// Required fields
	service := "Temperature"
	exitcode := OK
	message := ""
	perfdata := ""

	// Fetch SNMP Data
	oids := []string{OID_temperature}
	result, err := s.Get(oids)

	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}

	// Get the result
	temp := result.Variables[0].Value.(int)

	// Set message and perfdata
	message = fmt.Sprintf("%s\u00b0C", strconv.Itoa(temp))
	perfdata = strconv.Itoa(temp)
	perfdata = fmt.Sprintf("Temperature=%s", strconv.Itoa(temp))

	// Set exitcode
	switch {
	case temp < u.Args.TempWarn:
		exitcode = OK
	case temp < u.Args.TempCrit:
		exitcode = WARNING
	default:
		exitcode = CRITICAL
	}

	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
	return
}
