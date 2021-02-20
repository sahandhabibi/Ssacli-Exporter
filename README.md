# ssacli-exporter
ssacli Exporter For Prometheus

# Prometheus Exporter For HP RAID Controllers
This exporter for the Prometheus monitoring system calls into the ssacli utility to provide metrics for errors reported by HP RAID hardware. Under the hood, it invokes the following command when being scraped:

ssacli ctrl slot=0 physicaldrive all show detail


Example metrics output:

```
# HELP disk_current_temperature Disk Current Temperature
# TYPE disk_current_temperature gauge
disk_current_temperature{physicaldrive="box 3 bay 1 type SAS"} 46
disk_current_temperature{physicaldrive="box 3 bay 2 type SAS"} 51
disk_current_temperature{physicaldrive="box 3 bay 3 type SAS"} 48
disk_current_temperature{physicaldrive="box 3 bay 4 type SAS"} 44
# HELP disk_maximum_temperature Disk Maximum Temperature
# TYPE disk_maximum_temperature gauge
disk_maximum_temperature{physicaldrive="box 3 bay 1 type SAS"} 50
disk_maximum_temperature{physicaldrive="box 3 bay 2 type SAS"} 58
disk_maximum_temperature{physicaldrive="box 3 bay 3 type SAS"} 58
disk_maximum_temperature{physicaldrive="box 3 bay 4 type SAS"} 55
# HELP disk_status Disk Status (OK = 1)
# TYPE disk_status gauge
disk_status{physicaldrive="box 3 bay 1 type none"} 1
disk_status{physicaldrive="box 3 bay 2 type SAS"} 1
disk_status{physicaldrive="box 3 bay 3 type SAS"} 1
disk_status{physicaldrive="box 3 bay 4 type SAS"} 1
```
