/*
This is a console web application which is used to test the GO SDK as SDK is a package
Navigate to this folder and run go run cmd.go
Open browser and hit http://localhost:1110/home
*/
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/CA-APM/goprobe"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("This folder contains files to test GO SDK functionality as SDK is a package...!!!")
	goprobe.InitGoProbe()
	handleRequests()
}

func retrunHomeOnly(w http.ResponseWriter, r *http.Request) {
	txn := goprobe.StartTransaction(r.Context(), "/home")
	defer txn.EndTransaction()
	fmt.Fprintf(w, "home")
	time.Sleep(1 * time.Second)
	fSegment1(r.Context())
	fmt.Println("Endpoint Hit Home...!!!")
}

func returnAboutOnly(w http.ResponseWriter, r *http.Request) {
	txn := goprobe.StartTransaction(r.Context(), "/about")
	defer txn.EndTransaction()
	fmt.Fprintf(w, "about")
	time.Sleep(1 * time.Second)
	sqlSegment(r.Context())
	fmt.Println("Endpoint Hit about...!!!")
}

func returnInfoOnly(w http.ResponseWriter, r *http.Request) {
	txn := goprobe.StartTransaction(r.Context(), "/info")
	defer txn.EndTransaction()
	fmt.Fprintf(w, "info")
	time.Sleep(1 * time.Second)
	fSegment1(r.Context())
	fmt.Println("Endpoint Hit info...!!!")
}

func returnPrintOnly(w http.ResponseWriter, r *http.Request) {
	txn := goprobe.StartTransaction(r.Context(), "/print")
	defer txn.EndTransaction()
	fmt.Fprintf(w, "print")
	time.Sleep(1 * time.Second)
	fSegment1(r.Context())
	fmt.Println("Endpoint Hit print...!!!")
}

func fSegment1(ctx context.Context) {
	seg := goprobe.StartSegment(ctx, "fSegment1")
	defer seg.EndSegment()
	time.Sleep(1 * time.Second)
	sqlSegment(ctx)
	fSegment2(ctx)
}

func fSegment2(ctx context.Context) {
	seg := goprobe.StartSegment(ctx, "fSegment2")
	defer seg.EndSegment()
	time.Sleep(1 * time.Second)
	fSegment3(ctx)
}
func fSegment3(ctx context.Context) {
	seg := goprobe.StartSegment(ctx, "fSegment3")
	defer seg.EndSegment()
	fSegment4(ctx)
	time.Sleep(1 * time.Second)
}

func fSegment4(ctx context.Context) {
	seg := goprobe.StartSegment(ctx, "fSegment4")
	defer seg.EndSegment()
	time.Sleep(1 * time.Second)
}

func sqlSegment(ctx context.Context) {
	seg := goprobe.StartSegment(ctx, "query", "Server=lodsxvm58;Database=apm;Port=8080", "Select * from Booking where UserId = ?")
	defer seg.EndSegment()
	time.Sleep(1 * time.Second)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Handle("/home", goprobe.HttpWrapper(http.HandlerFunc(retrunHomeOnly)))
	myRouter.Handle("/about", goprobe.HttpWrapper(http.HandlerFunc(returnAboutOnly)))
	myRouter.Handle("/info", goprobe.HttpWrapper(http.HandlerFunc(returnInfoOnly)))
	myRouter.Handle("/print", goprobe.HttpWrapper(http.HandlerFunc(returnPrintOnly)))
	http.ListenAndServe(":1110", myRouter)
}
