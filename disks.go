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

type disk struct {
	ID string
	Model string
	DiskType string
	Status int
	Temperature int
	StatusName string
	Exitcode int
}

func CheckDisks(s *gosnmp.GoSNMP, u *Utilities) {
	//Required fields
	service := "Disks"
	exitcode := CRITICAL
	message := ""
	perfdata := ""

	// Fetch SNMP Data
	timeFetch := time.Now()

	resultsDiskId, err := s.BulkWalkAll(OID_diskID)
	resultsDiskModel, err := s.BulkWalkAll(OID_diskModel)	
	resultsDiskType, err := s.BulkWalkAll(OID_diskType)	
	resultsDiskStatus, err := s.BulkWalkAll(OID_diskStatus)	
	resultsDiskTemperature, err := s.BulkWalkAll(OID_diskTemperature)	

	u.Metrics.TimeToFetch += time.Now().Sub(timeFetch)
	
	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}
	
	// Create a disk slice 
	disks := []disk{}
	for i := 0; i < len(resultsDiskId); i++ {
		d := disk{}
		d.ID = strings.Trim(bytes.NewBuffer(resultsDiskId[i].Value.([]uint8)).String(), " ")
		d.Model = strings.Trim(bytes.NewBuffer(resultsDiskModel[i].Value.([]uint8)).String(), " ")
		d.DiskType = strings.Trim(bytes.NewBuffer(resultsDiskType[i].Value.([]uint8)).String(), " ")
		d.Status = resultsDiskStatus[i].Value.(int)
		d.Temperature = resultsDiskTemperature[i].Value.(int)
		disks = append(disks, d)
    }

	// Set additional fields
	setStatusAndExitcodes(disks, u.Args.TempWarn, u.Args.TempCrit)
	
	// Exitcode
	exitcode = getExitcode(disks)
	
	// Create message and perfdata
	countDisks := len(disks)
	countCritical := countExitcode(disks, CRITICAL)
	countWarning  := countExitcode(disks, WARNING)
	countUnknown  := countExitcode(disks, UNKNOWN)
	countOk       := countExitcode(disks, OK)
	tMin          := minTemp(disks)
	tAvg          := avgTemp(disks)
	tMax          := maxTemp(disks)
	switch exitcode {
	case OK:
		message = fmt.Sprintf("Total:%d All disks Ok (Max. Temp.%d\u00b0C)", countDisks, tMax)
	default:
		message = fmt.Sprintf("Total:%d Critical:%d Warning:%d Unknown:%d Ok:%d (Max. Temp.%d\u00b0C)", countDisks, countCritical, countWarning, countUnknown, countOk, tMax)
	}
	
	perfdata = fmt.Sprintf("Disks_Total=%d Disks_OK=%d Disks_WARNING=%d Disks_CRITICAL=%d Disks_UNKNOWN=%d Temp_Min=%d Temp_Avg=%d Temp_Max=%d", countDisks, countOk, countWarning, countCritical, countUnknown, tMin, tAvg, tMax)

	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
	
	
	// If diskCheck set, create check for each disk
	if u.Args.DiskChecks == true {
		for i := 0; i < len(disks); i++ {
			d := disks[i]
			
			// Set servicename
			service = d.ID
			
			// Set message
			message = fmt.Sprintf("Type:%s Model:%s Status%s Temperature:%d\u00b0C", d.DiskType, d.Model, d.StatusName, d.Temperature)
	
			// Set perfdata
			perfdata = fmt.Sprintf("Disk_Status=%d Disk_Temperature=%d;%d;%d", d.Status, d.Temperature, u.Args.TempWarn, u.Args.TempCrit)

			// Done. Write the check result
			Write(u.Args.Hostname, service, exitcode, message, perfdata, u)			
		}
	}
	
	return
}


func setStatusAndExitcodes (disks []disk, tempWarn int, tempCrit int){
	for i := 0; i < len(disks); i++ {
		d := disks[i]
		// StatusCode
		switch d.Status {
		case 1:
			d.StatusName = "Normal"
			d.Exitcode = OK
		case 2:
			d.StatusName = "Initialized"
			d.Exitcode = WARNING
		case 3:
			d.StatusName = "NotInitialized"
			d.Exitcode = WARNING
		case 4:
			d.StatusName = "SystemPartitionFailed"
			d.Exitcode = CRITICAL
		case 5:
			d.StatusName = "Crashed"
			d.Exitcode = CRITICAL
		default:
			d.StatusName = "Unknown"
			d.Exitcode = UNKNOWN
		}
		// Temp
		switch d.Exitcode{
		case CRITICAL:
			// NOOP
		case WARNING:
			switch {
			case d.Temperature >= tempCrit:
				d.StatusName = "Overheating"
				d.Exitcode = OK
			default:
				//NOOP
			}
		default: // OK
			switch {
			case d.Temperature >= tempCrit:
				d.StatusName = "Overheating"
				d.Exitcode = OK
			case d.Temperature >= tempWarn:
				d.StatusName = "Temperature Warning"
				d.Exitcode = WARNING
			default:
				// NOOP
			}
		}		
	}
}

func countExitcode (disks []disk, exitcode int) (c int) {
	for i := 0; i < len(disks); i++ {
		if disks[i].Exitcode == exitcode {
			c += 1
		}
	}
	return c
}

func countStatuscode (disks []disk, statuscode int) (c int){
	for i := 0; i < len(disks); i++ {
		if disks[i].Status == statuscode {
			c += 1
		}
	}
	return c
}

func getExitcode (disks []disk) (c int){
	switch {
	case countExitcode(disks, CRITICAL) > 0:
		c = CRITICAL
	case countExitcode(disks, WARNING) > 0:
		c = WARNING
	case countExitcode(disks, UNKNOWN) > 0:
		c = UNKNOWN
	default:
		c = OK
	}
	return c
}

func minTemp (disks []disk) (t int){
	t = disks[0].Temperature
	for i := 0; i < len(disks); i++ {
		if disks[i].Temperature < t {
			t = disks[i].Temperature
		}
	}
	return t
}

func maxTemp (disks []disk) (t int){
	t = 0
	for i := 0; i < len(disks); i++ {
		if disks[i].Temperature > t {
			t = disks[i].Temperature
		}
	}
	return t
}

func avgTemp (disks []disk) (t int){
	t = 0
	for i := 0; i < len(disks); i++ {
		t += disks[i].Temperature
	}
	return int(t / len(disks))
}



