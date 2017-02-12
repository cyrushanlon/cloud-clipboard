package client

import (
	"net"
	"runtime"
	"strings"

	"fmt"

	log "github.com/Sirupsen/logrus"
)

//StringArrayContains returns if the slice contains a target string
func StringArrayContains(list []string, target string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == target {
			return true
		}
	}
	return false
}

//StringArrayRemove removes target string from slice
func StringArrayRemove(list *[]string, target string) {
	for i := 0; i < len(*list); i++ {
		if (*list)[i] == target {
			l := append((*list)[:i], (*list)[i+1:]...)
			list = &l
			return
		}
	}
}

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func getLogOutput(v ...interface{}) string {
	_, file, line, _ := runtime.Caller(2)

	split := strings.Split(file, "/")
	fileOut := split[len(split)-1]

	return fmt.Sprint("[", fileOut, " ", line, "] ", fmt.Sprint(v...))
}

//LogInfo outputs info with function and line number of where it was called
func LogInfo(v ...interface{}) {
	log.Info(getLogOutput(v))
}

//LogErr outputs error with function and line number of where it was called
func LogErr(v ...interface{}) {
	log.Error(getLogOutput(v))
}

//LogWarn outputs warning with function and line number of where it was called
func LogWarn(v ...interface{}) {
	log.Warn(getLogOutput(v))
}
