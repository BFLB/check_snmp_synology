// Copyright (c) 2018 the check_snmp_authors. All rights reserved.
// Use of this source code is governed by ISC-style license
// that can be found in the LICENSE file.

// FIXME desctription

package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	. "github.com/BFLB/check_snmp_synology"
	"github.com/sonjah/gosnmp"
)

var (
	pluginVersion = "0.1"
	// TODO Fix and improve
	host        = flag.String("H", "", "Synology hostname")
	version     = flag.String("v", "2c", "SNMP version")
	community   = flag.String("C", "public", "Community")
	port        = flag.String("P", "161", "Port")
	timeout     = flag.Int("T", 10, "Timeout")
	commandfile = flag.String("cmd", "stdout", "Commandfile")
	tempWarn    = flag.Int("tempWarn", 50, "Warning Temperatur")
	tempCrit    = flag.Int("tempCrit", 50, "Critical Temperature")
	//storageWarn = flag.Int("storWarn", 50, "Warning Storage")
	//storageCrit = flag.Int("storCrit", 50, "Critical Storage")
)



func main() {

	timeStart := time.Now()
		
	flag.Usage = func() {
		// TODO Fix and improve
		fmt.Printf("Monitoring-Plugin check_snmp_synology\n")
		fmt.Printf("License: ISC\n")
		fmt.Printf("Source copyright and information: https://github.com/BFLB/check_snmp_synology\n\n")
		fmt.Printf("Usage:\n")

		flag.PrintDefaults()
	}

	flag.Parse()
	
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, ' ', 0)
	defer w.Flush()

	// Init Utilities
	var u Utilities
	u.Args.Commandfile = *commandfile
	u.Args.TempWarn = *tempWarn
	u.Args.TempCrit = *tempCrit
	u.Args.Timeout = *timeout
	u.Metrics.TimeStart = timeStart

	exitcode := CRITICAL
	execTimeCrit := u.Args.Timeout
	execTimeWarn := int((execTimeCrit / 10) * 8)

	// Init SNMP
	gosnmp.Default.Target = *host
	gosnmp.Default.Community = *community
	gosnmp.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking FIXME variable

	err := gosnmp.Default.Connect()
	if err != nil {
		exitcode = CRITICAL
		fmt.Printf("%s: Plugin version: %s - %s\n", NagiState(exitcode), pluginVersion, err.Error())
		os.Exit(exitcode)
	}
	defer gosnmp.Default.Conn.Close()

	u.Metrics.TimeToConnect = time.Now().Sub(u.Metrics.TimeStart)
	
	// Common checks
	CheckModel(gosnmp.Default, &u)
	CheckSystemStatus(gosnmp.Default, &u)
	CheckTemperature(gosnmp.Default, &u)
	CheckPowerStatus(gosnmp.Default, &u)
	CheckFanStatus(gosnmp.Default, &u)

	u.Metrics.TimeTotal = time.Now().Sub(u.Metrics.TimeStart)
	
	// Prepare exit information
	// Set exitcode
	if u.Metrics.TimeTotal.Seconds() > float64(u.Args.Timeout) {
		exitcode = CRITICAL

	} else if u.Metrics.TimeTotal.Seconds() >= float64(u.Args.Timeout) {
		exitcode = WARNING
	} else {
		exitcode = OK
	}

	timeTotal := float64(u.Metrics.TimeTotal.Seconds())
	timeConnect := float64(u.Metrics.TimeToConnect.Seconds())
	timeFetch := float64(u.Metrics.TimeToFetch.Seconds())
	timeProcess := float64((u.Metrics.TimeToProcess - u.Metrics.TimeToPrint).Seconds())
	timePrint := float64(u.Metrics.TimeToPrint.Seconds())

	message := fmt.Sprintf("%d passive-check(s) generated in %.3f seconds (t_conn:%.3fs t_fetch:%.3fs, t_proc:%.3fs, t_print:%.3fs)", u.Metrics.Checks, timeTotal, timeConnect, timeFetch, timeProcess, timePrint)

	perfdata := fmt.Sprintf("ExecTime=%3.3fs;%d;%d t_conn=%3.3fs t_load=%3.3fs t_proc=%3.3fs t_print=%3.3fs StatusCode=%d ChecksCreated=%d", timeTotal, execTimeWarn, execTimeCrit, timeConnect, timeFetch, timeProcess, timePrint, exitcode, u.Metrics.Checks)

	// TODO Make exit function in utils
	// Print exit information
	fmt.Printf("%s: Plugin version: %s - %s | %s\n", NagiState(exitcode), pluginVersion, message, perfdata)

	// Done. Exit with exitcode
	os.Exit(exitcode)

}
