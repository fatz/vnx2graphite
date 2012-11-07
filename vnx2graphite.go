package main

import (
	"os"
	"bufio"
	"bytes"
	"io"
	"fmt"
	"strings"
	"time"
	"net"
)

import goopt "github.com/droundy/goopt"

import config "github.com/kless/goconfig/config"

var version = "0.1"

///////////////////////////////////
//
//		Taken form goopt example
//
///////////////////////////////////

// The Flag function creates a boolean flag, possibly with a negating
// alternative.  Note that you can specify either long or short flags
// naturally in the same list.
var amVerbose = goopt.Flag([]string{"-v", "--verbose"}, []string{"--quiet"},
	"output verbosely this will also show the graphite data", "be quiet, instead")

// This is just a logging function that uses the verbosity flags to
// decide whether or not to log anything.
func log(x ...interface{}) {
	if *amVerbose {
		fmt.Println(x...)
	}
}

///////////////////////////////////
///////////////////////////////////

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

//replace slashes and whitespaces with underscore
func stringify(tstring string) (stringified string) {
	str := strings.Replace(tstring, "\"", "", -1)
	str = strings.Replace(str, "/", "_", -1)
	stringified = strings.Replace(str, " ", "_", -1)
	return
}

var opt_conf = goopt.String([]string{"-c", "--config"}, "config file", "path to config file")
var opt_data = goopt.String([]string{"-d", "--data"}, "data csv", "path to data csv file")
var opt_statsname = goopt.String([]string{"-s", "--statsname"}, "nfs", "extending name for the bucket: $basename.nfs")
var opt_mover = goopt.String([]string{"-m", "--datamover"}, "server_2", "extending name for the bucket: $basename.$movername.nfs")

func main() {
	goopt.Version = version
	goopt.Summary = "send emc vnx performance data to graphite"
	goopt.Parse(nil)

	if f, _ := exists(*opt_conf); f == false {
		fmt.Print(goopt.Help())
		fmt.Println("ERROR: config file " + *opt_conf + " doesn't exist")
		return
	}
	c, _ := config.ReadDefault(*opt_conf)

	host, _ := c.String("graphite", "host")
	port, _ := c.Int("graphite", "port")
	timeout, _ := c.Int("graphite", "timeout")
	basename, _ := c.String("graphite", "basename")
	//todo: we also can parse this from the output... 
	timestamp := time.Now().Unix()

	log(fmt.Sprintf("using graphite with host %s:%d", host, port))

	if f, _ := exists(*opt_data); f == false {
		fmt.Print(goopt.Help())
		fmt.Println("ERROR: data file " + *opt_data + " doesn't exist")
		return
	}
	lines, err := readLines(*opt_data)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	//we only want to use the head and the first line of data
	head := strings.Split(lines[0], ",")
	data := strings.Split(lines[1], ",")

	if len(head) != len(data) {
		fmt.Println("ERROR: malformed csv (length of head != length of data")
		return
	}

	//create the graphite connection.
	// Todo: export this in a small pkg

	cstr := fmt.Sprintf("%s:%d", host, port)
	log("trying to connect to " + cstr)
	conn, err := net.DialTimeout("tcp", cstr, time.Duration(timeout)*time.Second)

	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	for n, val := range head {
		key := stringify(val)
		value := stringify(data[n])
		if key != "Timestamp" {
			msg := fmt.Sprintf("%s.%s.%s.%s %s %d", basename, *opt_mover, *opt_statsname, key, value, timestamp)
			log("sending: " + msg)
			fmt.Fprint(conn, "\n"+msg+"\n")
		} else {
			log("Timestamp: ... next...")
		}
	}
	conn.Close()
}
