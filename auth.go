package client

import "strings"

//IsAuthedTCP checks a TCP clipboard change for authorisation
func IsAuthedTCP(msg string) (bool, string) {
	split := strings.Split(msg, ":")
	if len(split) < 2 && Conf.Username != split[0] {
		return false, ""
	}

	msgClip := ""
	for i := 1; i < len(split)-1; i++ {
		msgClip += split[i] + ":"
	}

	msgClip += split[len(split)-1]

	return true, msgClip
}

//AddAuthTCP adds TCP clipboard change authorisation before sending
func AddAuthTCP(msg string) []byte {
	return []byte(Conf.Username + ":" + msg)
}

//IsAuthedUDP checks a UDP packet for authorisation
func IsAuthedUDP(msg string) (bool, string) {
	split := strings.Split(msg, ":")
	if len(split) < 2 && Conf.Username != split[0] {
		return false, ""
	}

	msgClip := ""
	for i := 1; i < len(split); i++ {
		msgClip += split[i]
	}

	return true, msgClip
}

//AddAuthUDP checks a UDP packet for authorisation
func AddAuthUDP(msg string) []byte {
	return []byte(Conf.Username + ":" + msg)
}