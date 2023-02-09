package heartbeat

import "math/rand"

func ChooseRandomDataServer() string {
	dataServers := GetDataServer()
	count := len(dataServers)
	if count == 0 {
		return ""
	}
	return dataServers[rand.Intn(count)]
}
