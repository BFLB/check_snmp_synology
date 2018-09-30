// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"os"
	"time"
)

var (
	//OID declarations
	OID_syno                = "1.3.6.1.4.1.6574"
	OID_model               = "1.3.6.1.4.1.6574.1.5.1.0"
	OID_serialNumber        = "1.3.6.1.4.1.6574.1.5.2.0"
	OID_DSMVersion          = "1.3.6.1.4.1.6574.1.5.3.0"
	OID_DSMUpgradeAvailable = "1.3.6.1.4.1.6574.1.5.4.0"
	OID_systemStatus        = "1.3.6.1.4.1.6574.1.1.0"
	OID_temperature         = "1.3.6.1.4.1.6574.1.2.0"
	OID_powerStatus         = "1.3.6.1.4.1.6574.1.3.0"
	OID_systemFanStatus     = "1.3.6.1.4.1.6574.1.4.1.0"
	OID_CPUFanStatus        = "1.3.6.1.4.1.6574.1.4.2.0"

	OID_disk       = ""
	OID_disk2      = ""
	OID_diskID     = "1.3.6.1.4.1.6574.2.1.1.2"
	OID_diskModel  = "1.3.6.1.4.1.6574.2.1.1.3"
	OID_diskStatus = "1.3.6.1.4.1.6574.2.1.1.5"
	OID_diskTemp   = "1.3.6.1.4.1.6574.2.1.1.6"

	OID_RAID       = ""
	OID_RAIDName   = "1.3.6.1.4.1.6574.3.1.1.2"
	OID_RAIDStatus = "1.3.6.1.4.1.6574.3.1.1.3"

	OID_Storage                = "1.3.6.1.2.1.25.2.3.1"
	OID_StorageDesc            = "1.3.6.1.2.1.25.2.3.1.3"
	OID_StorageAllocationUnits = "1.3.6.1.2.1.25.2.3.1.4"
	OID_StorageSize            = "1.3.6.1.2.1.25.2.3.1.5"
	OID_StorageSizeUsed        = "1.3.6.1.2.1.25.2.3.1.6"

	OID_UpsModel                = "1.3.6.1.4.1.6574.4.1.1.0"
	OID_UpsSN                   = "1.3.6.1.4.1.6574.4.1.3.0"
	OID_UpsStatus               = "1.3.6.1.4.1.6574.4.2.1.0"
	OID_UpsLoad                 = "1.3.6.1.4.1.6574.4.2.12.1.0"
	OID_UpsBatteryCharge        = "1.3.6.1.4.1.6574.4.3.1.1.0"
	OID_UpsBatteryChargeWarning = "1.3.6.1.4.1.6574.4.3.1.4.0"
)

type Args struct {
	Hostname    string
	Version     int
	Username    string
	Password    string
	Community   string
	Port        int
	Timeout     int
	Commandfile string
	TempWarn    int
	TempCrit    int
	StorageWarn int
	StorageCrit int
}

type Metrics struct {
	Checks        int
	TimeStart     time.Time
	TimeToConnect time.Duration
	TimeToProcess time.Duration
	TimeToPrint   time.Duration
	TimeTotal     time.Duration
}

type Utilities struct {
	Args    Args
	Metrics Metrics
}

// Nagios exit codes
const OK = 0
const WARNING = 1
const CRITICAL = 2
const UNKNOWN = 3

func Write(hostname string, service string, exitcode int, message string, perfdata string, utils *Utilities) {
	start := time.Now()
	// Parameters
	command := "PROCESS_SERVICE_CHECK_RESULT"
	timestamp := time.Now().Unix()
	var result string
	// Create the messaga
	if perfdata != "" {
		result = fmt.Sprintf("[%d] %s;%s;%s;%d;%s: %s | %s", timestamp, command, hostname, service, exitcode, NagiState(exitcode), message, perfdata)
	} else {
		result = fmt.Sprintf("[%d] %s;%s;%s;%d;%s: %s", timestamp, command, hostname, service, exitcode, NagiState(exitcode), message)
	}

	// TODO Errorhandling?

	if utils.Args.Commandfile == "stdout" {
		fmt.Println(result)

	} else {
		f, err := os.OpenFile(utils.Args.Commandfile, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(f, "%s\n", result)
		f.Close()
	}
	//  Write successful, finally update metrics
	utils.Metrics.Checks += 1
	utils.Metrics.TimeToPrint += time.Now().Sub(start)
}

func NagiState(exitcode int) (state string) {
	switch exitcode {
	case OK:
		return "OK"
	case WARNING:
		return "WARNING"
	case CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
