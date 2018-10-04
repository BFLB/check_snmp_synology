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

type ebox struct {
	Index int
	Model string
	Power1 int
	Power2 int
	StatusName1	string
	StatusName2 string
	Exitcode int
}

func CheckEbox(s *gosnmp.GoSNMP, u *Utilities) {
	//Required fields
	service := "Disks"
	exitcode := CRITICAL
	message := ""
	perfdata := ""
	thresholdCritical := 2
	
	// Only proceed if check is configured
	if u.Args.Ebox == 0 {
		return
	}
	
	// Fetch SNMP Data
	timeFetch := time.Now()
	
	var err error
	var resultsEboxIndex  []gosnmp.SnmpPDU
	var resultsEboxModel  []gosnmp.SnmpPDU
	var resultsEboxPower1 []gosnmp.SnmpPDU
	var resultsEboxPower2 []gosnmp.SnmpPDU
	
	resultsEboxIndex, err  = s.BulkWalkAll(OID_eboxIndex)
	if err == nil {
		resultsEboxModel, err  = s.BulkWalkAll(OID_eboxModel)
	}
	if err == nil {
		resultsEboxPower1, err = s.BulkWalkAll(OID_eboxPower)
	}
	// Fetch SNMP Data Power2 (Only for models with redundant poewer supplies)
	if err == nil && u.Args.Ebox == 2 {
		resultsEboxPower2, err = s.BulkWalkAll(OID_eboxRedundantPower)	
	}
	
	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}
	
	u.Metrics.TimeToFetch += time.Now().Sub(timeFetch)
	
	
	// Create a ebox slice 
	eboxes := []ebox{}
	for i := 0; i < len(resultsEboxIndex); i++ {
		e := ebox{}
		e.Index = resultsEboxIndex[i].Value.(int)
		e.Model = strings.Trim(bytes.NewBuffer(resultsEboxModel[i].Value.([]uint8)).String(), " ")
		e.Power1 = resultsEboxPower1[i].Value.(int)
		if u.Args.Ebox == 2 {
			e.Power2 = resultsEboxPower2[i].Value.(int)
		}
		eboxes = append(eboxes, e)
    }

	// Set additional fields
	setupEboxes(eboxes)
	
	// If diskCheck set, create check for each disk
	for i := 0; i < len(eboxes); i++ {
		e := &eboxes[i]
			
		// Set servicename
		service = fmt.Sprintf("Expansion-Unit %02d", e.Index)
		
		// Set exitcode
		exitcode = e.Exitcode
			
		// Set message
		message = fmt.Sprintf("Model:%s Power1:%s Power2:%s", e.Model, e.StatusName1, e.StatusName2)
	
		// Set perfdata
		perfdata = fmt.Sprintf("Power1_status=%d;;%d Power2_status=%d;;%d", e.Power1, thresholdCritical, e.Power2, thresholdCritical)

		// Done. Write the check result
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)			
	}
	
	return
}

func setupEboxes (eboxes []ebox){
	for i := 0; i < len(eboxes); i++ {
		// Power1
		e := &eboxes[i]
		switch e.Power1 {
		case 1:
			e.StatusName1 = "Normal"
			e.Exitcode = OK
		case 2:
			e.StatusName1 = "Poor"
			e.Exitcode = CRITICAL
		case 3:
			e.StatusName1 = "Disconnected"
			e.Exitcode = CRITICAL
		default:
			e.StatusName1 = "Unknown"
			e.Exitcode = UNKNOWN
		}
		// Power2
		switch e.Power2{
		case 0:
			e.StatusName2 = "N/A"
		case 1:
			e.StatusName2 = "Normal"			
		case 2:
			e.StatusName2 = "Poor"
			e.Exitcode = CRITICAL
		case 3:
			e.StatusName2 = "Disconnected"
			e.Exitcode = CRITICAL
				// NOOP
		}		
	}
}
