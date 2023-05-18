package main

import (
        "errors"
        "fmt"
        "net/http"
        "os"
        "strings"
        "encoding/json"
)

func getBoot(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("got /boot request...\n")
        //get mac
		mac := getMAC(r)
        //build ipxe file using mac
        ipxe := getIpxe(mac)
        //respond with file
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Disposition", "attachment; filename=4a:2e:c4:0b:43:99.ipxe")
        w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
        w.Write([]byte(ipxe))
        return
}

func getISCSI(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /iscsi request...\n")
	//get mac
	mac := getMAC(r)
	//build ipxe file using mac
	pvc := getPVC(mac)
	//respond with file
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Disposition", "attachment; filename=4a:2e:c4:0b:43:99.ipxe")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Write([]byte(ipxe))
	return
}

func getMAC(r *http.Request) string {
	fmt.Printf("getting mac from request...\n")
	res := strings.ReplaceAll(r.URL.String(), "/boot/", "")
	res = strings.ReplaceAll(res, ".ipxe", "")
	mac := strings.ReplaceAll(res, "%3A", ":")
	return mac
}

func getHOST(r *http.Request) string {
	fmt.Printf("looking up hostname using mac...\n")
	return hostname
}

func getIpxe(mac string) string {

        m := `{"4a:2e:c4:0b:43:99":"testvm","4a:2e:c4:0b:43:98":"notreal"}`

        var data map[string]interface{}
        err := json.Unmarshal([]byte(m), &data)
        if err != nil {
                panic(err)
        }

        response := `!ipxe
set iscsi-target iscsi:%h-iscsi.stg.srvd.dev::::iqn.2016-09.com.openebs.cstor:pvc-%h
set gateway 10.0.1.1
set initiator-iqn iqn.2015-02.com.srvd.%h
set keep-san 1
sanboot ${iscsi-target}
boot
`

        resp := strings.ReplaceAll(response, "%h", data[mac].(string))
        return resp
}

func getPVC() string {

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