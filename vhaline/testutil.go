package vhaline

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func MakeAndMoveToTempDir() (origdir string, tmpdir string) {

	// make new temp dir
	var err error
	origdir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	tmpdir, err = ioutil.TempDir(origdir, "temp.test.dir")
	if err != nil {
		panic(err)
	}
	err = os.Chdir(tmpdir)
	if err != nil {
		panic(err)
	}

	return origdir, tmpdir
}

func TempDirCleanup(origdir string, tmpdir string) {
	// cleanup
	os.Chdir(origdir)
	err := os.RemoveAll(tmpdir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n TempDirCleanup of '%s' done.\n", tmpdir)
}

// GetAvailPort asks the OS for an unused port,
// returning a bound net.Listener and the port number
// to which it is bound. The caller should
// Close() the listener when it is done with
// the port.
func GetAvailPort() (net.Listener, int) {
	lsn, _ := net.Listen("tcp", ":0")
	r := lsn.Addr()
	return lsn, r.(*net.TCPAddr).Port
}

// waitUntilAddrAvailable returns -1 if the addr was
// alays unavailable after tries sleeps of dur time.
// Otherwise it returns the number of tries it took.
// Between attempts we wait 'dur' time before trying
// again.
func WaitUntilAddrAvailable(addr string, dur time.Duration, tries int) int {
	for i := 0; i < tries; i++ {
		var isbound bool
		isbound = IsAlreadyBound(addr)
		if isbound {
			time.Sleep(dur)
		} else {
			fmt.Printf("\n took %v %v sleeps for address '%v' to become available.\n", i, dur, addr)
			return i
		}
	}
	return -1
}

func sendViaTcp(hostport string, bts []byte) error {
	conn, err := net.Dial("tcp", hostport)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write(bts)
	return err
}
