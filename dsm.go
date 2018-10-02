// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"time"

	"github.com/sonjah/gosnmp"
)

func CheckDSM(s *gosnmp.GoSNMP, u *Utilities) {
	// Required fields
	service := "DSM"
	exitcode := CRITICAL
	message := ""
	perfdata := ""

	// Fetch SNMP Data
	timeFetch := time.Now()

	oids := []string{OID_DSMVersion, OID_DSMUpgradeAvailable}
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
	dsmVersion := response.Variables[0].Value

	// Set response information
	if u.Args.UpgradeStatus == true {
		dsmUpgradeCode := response.Variables[1].Value.(int)
		dsmUpgrade     := ""
		
		switch dsmUpgradeCode {
		case 1:
			dsmUpgrade = "Available (New version ready for download)"
			exitcode  = WARNING
		
		case 2:
			dsmUpgrade = "Unavailable (DSM is up to date)"
			exitcode  = OK
		
		case 3:
			dsmUpgrade = "Connecting (Checking for the latest DSM)"
			exitcode  = OK
		
		case 4:
			dsmUpgrade = "Disconnected (Failed to connect to server)"
			exitcode  = CRITICAL
		
		case 5:
			dsmUpgrade = "Others (DSM is upgrading or downloading)"
			exitcode  = WARNING

		default:
			dsmUpgrade = "Unknown"
			exitcode  = UNKNOWN
			
		}
		message  = fmt.Sprintf("Version:%s UpgradeStatus:%s", dsmVersion, dsmUpgrade)
		perfdata = fmt.Sprintf("Upgrade_Status=%d", dsmUpgradeCode)
		
	} else {
		exitcode = OK
		message = fmt.Sprintf("Version:%s", response.Variables[0].Value)
	}
		
	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
	return
}
