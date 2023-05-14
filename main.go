package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"encoding/json"
)

func getBoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /boot request\n")

	//get mac
	fmt.Println(r.URL)

	res := strings.ReplaceAll(r.URL.String(), "/boot/", "")
    fmt.Println(res)

	res = strings.ReplaceAll(res, ".ipxe", "")
    fmt.Println(res)

	mac := strings.ReplaceAll(res, "%3A", ":")
    fmt.Println(mac)

	getIpxe(mac)

	io.WriteString(w, "This is my website!\n")
}

func getIpxe(mac string) []byte {

	m := `{"4a:2e:c4:0b:43:99":"testvm","4a:2e:c4:0b:43:98":"notreal"}`

	var data map[string]interface{}
	err := json.Unmarshal([]byte(m), &data)
	if err != nil {
		panic(err)
	}

	hostname := data[mac]

	response := `
	!ipxe
	set iscsi-target iscsi:%h-iscsi.stg.srvd.dev::::iqn.2016-09.com.openebs.cstor:pvc-%h
	set gateway 10.0.1.1
	set initiator-iqn iqn.2015-02.com.srvd.%h
	set keep-san 1
	sanboot ${iscsi-target}
	boot
	`

	h := fmt.Sprintf(response2, hostname.(string))

	rawResponse := []byte(h)

	fmt.Println(response2)

	fmt.Println(h) 

    fmt.Println(rawResponse) // [65 66 67 226 130 172]

	return rawResponse

}


func main() {

	http.HandleFunc("/boot/", getBoot)

	err := http.ListenAndServe(":6666", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}