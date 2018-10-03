// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"time"
	"bytes"
	"strings"
	
	"github.com/sonjah/gosnmp"
)

func CheckHighAvailability(s *gosnmp.GoSNMP, u *Utilities) {
	service := "High-Availability"
	exitcode := OK
	message := ""
	perfdata := ""

	// Only proceed if check is configured
	if u.Args.HA == "" {
		return
	}
	
	// Retrieve configured serialnumbers
	serials := strings.Split(u.Args.HA, ",")
	
	// Errorhandling
	if len(serials) != 2 {
		exitcode = UNKNOWN
		message = fmt.Sprintf("Wrong input: %s Must be -ha=<primary-serial>,<secondary-serial>", u.Args.HA)
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}

	primary := serials[0]
	secondary := serials[1]

	// Errorhandling
	if primary == "" || secondary == "" {
		exitcode = UNKNOWN
		message = fmt.Sprintf("Wrong input: \"%s\" Must be -ha=<primary-serial>,<secondary-serial>", u.Args.HA)
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}

	// Fetch SNMP Data
	timeFetch := time.Now()
	
	oids := []string{OID_serialNumber}
	response, err := s.Get(oids)

	u.Metrics.TimeToFetch += time.Now().Sub(timeFetch)

	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}
	
	// Actual serialnumber
	serialnumber := strings.Trim(bytes.NewBuffer(response.Variables[0].Value.([]uint8)).String(), " ")
	
	// Check actual serialnumber agains configured ones
	switch {
	case serialnumber == primary:
		exitcode = OK
		message = fmt.Sprintf("Normal: Active-Server=Primary(%s)", serialnumber)
	case serialnumber == secondary:
		exitcode = CRITICAL
		message = fmt.Sprintf("Failover: Active-Server=Secondary(%s)", serialnumber)
	default:
		exitcode = CRITICAL
		message = fmt.Sprintf("Failover: Active-Server=Unknown(%s)", serialnumber)
	}

	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)

	return
}
