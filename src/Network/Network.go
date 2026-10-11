package Network

import "fmt"

func NewApp() {
	network, _ := NetIOCountersStat_f()
	for _, v := range network {
		fmt.Println(v)
	}
}
