// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

package check_snmp_synology

import (
	"fmt"

	"github.com/sonjah/gosnmp"
)

func CheckModel(s *gosnmp.GoSNMP, u *Utilities) {
	service := "Model"
	exitcode := OK
	perfdata := ""

	oids := []string{OID_model}
	result, err := s.Get(oids)
	if err != nil {
		// TODO Errorhandling
		return
	}
	modelname := result.Variables[0].Value
	var message string
	message = fmt.Sprintf("%s", modelname)
	Write(u.Args.Hostname, service, exitcode, message, perfdata, u)

	return
}
