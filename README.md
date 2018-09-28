# check_snmp_synology
Golang Nagios Plugin for Synology

This plugin can be executed by any Nagios compatible Monitoring System as an actice check.
During the execution, the various results are sent to the Monitoring System as independent passive checks.
At the end, the plugin returns with some information of the overall execution results, e.g. Execution time,
Number of checks created etc.
The benefit of this approach is that a lot of information can be pulled out of the system in an easy and
efficient way.
