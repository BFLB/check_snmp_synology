// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"
	"time"
	
	"github.com/sonjah/gosnmp"
)

func CheckModel(s *gosnmp.GoSNMP, u *Utilities) {
	service := "Model"
	exitcode := OK
	message := ""
	perfdata := ""

	// Fetch SNMP Data
	timeFetch := time.Now()
		
	oids := []string{OID_model,OID_serialNumber}
	response, err := s.Get(oids)

	u.Metrics.TimeToFetch += time.Now().Sub(timeFetch)

	// Errorhandling
	if err != nil {
		exitcode = UNKNOWN
		message = err.Error()
		Write(u.Args.Hostname, service, exitcode, message, perfdata, u)
		return
	}
	
	modelname    := response.Variables[0].Value
	serialnumber := response.Variables[1].Value
	
	message = fmt.Sprintf("%s (S/N:%s)", modelname, serialnumber)
	// Done. Write the check result
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)

	return
}
