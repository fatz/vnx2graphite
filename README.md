vnx2graphite
===============

vnx tool for pushing server_stats to graphite

# build and install #
build your crosscompile environment [http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go](http://dave.cheney.net/2012/09/08/an-introduction-to-cross-compilation-with-go)

then do a 
    
    go-linux-386 build vnx2graphite.go

because the controlstations are 32-bit OS you have to do a 386 build.

after this just copy the conf, the binary and the shell script to 
    
    nasadmin@your-control-station:~/vnx2graphite/

# configure #
you have to configure four options in the vnx2graphite.conf:
    
    host=<hostname or ipaddress of graphite>
    port=<carbon port>
    timeout=<timeout in seconds>
    basename=<graphite namespace e.g. com.company.your.vnx>

# run #
if your config and binary are at `$HOME/vnx2graphite/` you can use the vnx2graphite wrapper script. It expects that the nasdb is mounted at /nas. Otherwise take it as an example and write your own script.

## run the wrapper ##

    $HOME/vnx2graphite/vnx2graphite.sh <stats e.g. nfs> <server name e.g. server_2>

## periodicly run the wrapper ##
Use cron to ensure the script will run every minute. If you need it faster just do a loop with sleep but be aware of timeouts. The default timeout for vnx2graphite is 10 seconds.

### example `/etc/cron.d/vnx2graphite` ###

    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh nfs server_2
    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh fs server_2
    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh kernel server_2
    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh mpfs server_2
    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh store server_2
    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh rep server_2
    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh rpc server_2
    * * * * *      nasadmin     /home/nasadmin/vnx2graphite/vnx2graphite.sh snap server_2

# environment #
I tested this script with the 5300 model and the 5500 model