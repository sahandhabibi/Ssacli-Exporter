# ssacli-exporter For HP RAID Controllers
This exporter for the Prometheus monitoring system calls into the ssacli utility to provide metrics for errors reported by HP RAID hardware. Under the hood, it invokes the following command when being scraped:

```
ssacli ctrl slot=0 physicaldrive all show detail
```

# Try it

```
ssacli_exporter -Port 9060
```
default port 9109 - /metrics

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

# Install HP Raid HPSA ssacli on Ubuntu

You just have to edit

```
/etc/apt/sources.list
```

Then add the source at the end of the file

```
  deb http://downloads.linux.hpe.com/SDR/repo/mcp xenial/current non-free

```

Enroll keys for DEB-based systems
Issue the following commands to enroll all keys on your deb-based system:

```
curl https://downloads.linux.hpe.com/SDR/hpPublicKey2048.pub | apt-key add -
curl https://downloads.linux.hpe.com/SDR/hpPublicKey2048_key1.pub | apt-key add -
curl https://downloads.linux.hpe.com/SDR/hpePublicKey2048_key1.pub | apt-key add -
```
Then install your package.

```
apt-get update && apt-get install ssacli

```

More Info 

https://downloads.linux.hpe.com/SDR/project/mcp/

https://downloads.linux.hpe.com/SDR/keys.html

https://support.hpe.com/hpsc/swd/public/detail?swItemId=MTX-f8f30da26d6749499adec36f8b#tab3
