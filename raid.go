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

type raid struct {
	Name string
	Status int
	PercentUsed int
	SizeFree int
	SizeTotal int
	StatusName	string
	Exitcode int
}

func CheckRaid(s *gosnmp.GoSNMP, u *Utilities) {
	//Required fields
	service := "Disks"
	exitcode := CRITICAL
	message := ""
	perfdata := ""
	
	// Only proceed if check is configured
	if u.Args.Ebox == 0 {
		return
	}
	
	// Fetch SNMP Data
	timeFetch := time.Now()
	
	var err error
	var resultsRaidName      []gosnmp.SnmpPDU
	var resultsRaidStatus    []gosnmp.SnmpPDU
	var resultsRaidSizeFree  []gosnmp.SnmpPDU
	var resultsRaidSizeTotal []gosnmp.SnmpPDU


	resultsRaidName, err  = s.BulkWalkAll(OID_raidName)
	if err == nil {
		resultsRaidStatus, err  = s.BulkWalkAll(OID_raidStatus)
	}
	if err == nil {
		resultsRaidSizeFree, err = s.BulkWalkAll(OID_raidFreeSize)
	}
	if err == nil {
		resultsRaidSizeTotal, err = s.BulkWalkAll(OID_raidTotalSize)	
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
	raids := []raid{}
	for i := 0; i < len(resultsRaidName); i++ {
		r := raid{}
		r.Name = strings.Trim(bytes.NewBuffer(resultsRaidName[i].Value.([]uint8)).String(), " ")
		r.Status = resultsRaidStatus[i].Value.(int)
		r.SizeFree = int(resultsRaidSizeFree[i].Value.(uint64)/1024/1024/1024) // Gigabytes
		r.SizeTotal = int(resultsRaidSizeTotal[i].Value.(uint64)/1024/1024/1024) // Gigabytes
		r.PercentUsed = int(100 * (r.SizeTotal - r.SizeFree) / r.SizeTotal)
		raids = append(raids, r)
    }

	// Set additional fields
	setupRaids(raids, u.Args.RaidWarn, u.Args.RaidCrit)
	
	// If diskCheck set, create check for each disk
	for i := 0; i < len(raids); i++ {
		r := &raids[i]

		// Set servicename
		service = fmt.Sprintf("RAID %d", i)
		
		// Set exitcode
		exitcode = r.Exitcode
			
		// Set message
		message = fmt.Sprintf("Name:%s Status:%s Used:%d%% Free:%dGB Total:%dGB", r.Name, r.StatusName, r.PercentUsed, r.SizeFree, r.SizeTotal)
	
		// Set perfdata
		perfdata = fmt.Sprintf("Status=%d Used=%d%% Size_free=%dG Size_total=%dG", r.Status, r.PercentUsed, r.SizeFree, r.SizeTotal)

		// Done. Write the check result
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)			
	}
	
	return
}

func setupRaids (raids []raid, warn int, crit int){
	for i := 0; i < len(raids); i++ {
		// Power1
		r := &raids[i]
		switch r.Status {
		case 1:
			r.StatusName = "Normal"
			r.Exitcode = OK
		case 2:
			r.StatusName = "Repairing"
			r.Exitcode = CRITICAL
		case 3:
			r.StatusName = "Migrating"
			r.Exitcode = CRITICAL
		case 4:
			r.StatusName = "Expanding"
			r.Exitcode = CRITICAL
		case 5:
			r.StatusName = "Deleting"
			r.Exitcode = CRITICAL
		case 6:
			r.StatusName = "Creating"
			r.Exitcode = CRITICAL
		case 7:
			r.StatusName = "Syncing"
			r.Exitcode = CRITICAL
		case 8:
			r.StatusName = "ParityChecking"
			r.Exitcode = CRITICAL
		case 9:
			r.StatusName = "Assembling"
			r.Exitcode = CRITICAL
		case 10:
			r.StatusName = "Canceling"
			r.Exitcode = CRITICAL
		case 11:
			r.StatusName = "Degrade"
			r.Exitcode = CRITICAL
		case 12:
			r.StatusName = "Crashed"
			r.Exitcode = CRITICAL
		case 13:
			r.StatusName = "DataScrubbing"
			r.Exitcode = CRITICAL
		case 14:
			r.StatusName = "Deploying"
			r.Exitcode = CRITICAL
		case 15:
			r.StatusName = "UnDeploying"
			r.Exitcode = CRITICAL
		case 16:
			r.StatusName = "MountCache"
			r.Exitcode = CRITICAL
		case 17:
			r.StatusName = "UnmountCache"
			r.Exitcode = CRITICAL
		case 18:
			r.StatusName = "ExpandingUnfinishedSHR"
			r.Exitcode = CRITICAL
		case 19:
			r.StatusName = "ConvertSHRToPool"
			r.Exitcode = CRITICAL
		case 20:
			r.StatusName = "ConvertSHRToSHR2"
			r.Exitcode = CRITICAL
		default:
			r.StatusName = "UnknownStatus"
			r.Exitcode = CRITICAL
		}
	}
}
