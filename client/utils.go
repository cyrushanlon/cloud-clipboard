package client

import "net"

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
