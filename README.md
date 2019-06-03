# GO Probe

# Description
This project is intended to build a stable and extensible SDK for instrumenting applications that are developed in GO language. As this is SDK based, users have to manually instrument the application by using the interfaces provided. The goprobe resides inside the application and records the client function entry/exit calls and then feed the data collected to APM collector.

## Short Description
An SDK based instrumentation for applications developed in GO language.

## APM version
APM 11.1.

## Supported third party versions
None

## License
[Apache License, Version 2.0](LICENSE)

# Installation Instructions

## Prerequisites
- Go 1.11+ - Working Go installed in the machine with all the required environment variables set like GOROOT, GOPATH
https://golang.org/dl/

- GoExtension - Download the latest extension from CA APM and deploy as part of Infrastructure agent process, we require credentials to login please contact Broadcom Support.
https://support.ca.com/

## Dependencies
- GOSDK depends on third party package for unique id generation, so install the package before instrumentation.<br/>
	```
	go get github.com/satori/go.uuid
	```

- Examples in the package are built using mux router, so install gorilla mux router before running the examples.<br/>
	```
	go get -u github.com/gorilla/mux
	```

## Installation
- Open command prompt and download the code from the repository.<br/>
	```
	go get github.com/CA-APM/goprobe
	```

## Configuration
- Navigate to goprobe and set the application configurations in the config.json
	- appName - Set the application name to be displayed in the ATC
	- hostName - IP address of the Infrastructure Agent
	- hostPort - Port to connect to Infrastructure Agent(ideally 5005)
	- logPath - Valid path to redirect the logs, this path should also contain the file name
	- debugEnabled - Set it to true if you need to enable the debug logs
	- interval - Set the time to collect the runtime metrics by the probe(seconds)

# Usage Instructions
Navigate to goprobe/cmd/ and run the below command, make sure the Infrastructure Agent is up and running.<br/>
	```
	go run cmd.go
	```

Open browser and hit any one of the below URL's to see the traces/metrics in ATC
- http://localhost:1110/home
- http://localhost:1110/about
- http://localhost:1110/info
- http://localhost:1110/print

## Application Instrumentation
As this is based on SDK, users need to instrument the application manually by using the interfaces provided. Follow the below steps to instrument the application.

### Interfaces exposed by SDK
These are the interfaces exposed by the SDK and these needs to be used to instrument and get the traces/metrics in ATC,
```
	InitGoProbe()
	HttpWrapper(f http.HandlerFunc) http.HandlerFunc
	StartTransaction(ctx context.Context, httpUrl string) Transaction
	EndTransaction()
	StartSegment(ctx context.Context, name string, optional ...string)
	EndSegment()
```

### Interface Usages
Inorder to use the interfaces provided by the SDK, users has to import the package in their application and then invoke the interfaces.
```
import "github.com/CA-APM/goprobe"
```

#### InitGoProbe
This is the base interface provided by the SDK and this needs to be invoked from the main entry point of the application,
```
Before Instrumentation:
func main() {
	handleRequests()
}
After Instrumentation:
import "github.com/CA-APM/goprobe"

func main() {
	goprobe.InitGoProbe()
	handleRequests()
}
```

#### HttpWrapper
This is used to get the transaction traces of http endpoints, inorder to intercept the http request users have to call the HttpWrapper for every incoming request, so that the SDK will help in retrieving the traces of http endpoints.
```
Before Instrumentation:
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/home", retrunHomeOnly)
	http.ListenAndServe(":1110", myRouter)
}
After Instrumentation:
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Handle("/home", goprobe.HttpWrapper(http.HandlerFunc(retrunHomeOnly)))
	http.ListenAndServe(":1110", myRouter)
}
```

#### StartTransaction, EndTransaction
This is used to collect the transaction information, these interfaces should be called at the function entry point so that the transaction information is collected and published to ATC.
```
Before Instrumentation:
func retrunHomeOnly(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
}
After Instrumentation:
func retrunHomeOnly(w http.ResponseWriter, r *http.Request) {
	txn := goprobe.StartTransaction(r.Context(), "/home")
	defer txn.EndTransaction()
	fmt.Fprintf(w, "home")
}
```

#### StartSegment, EndSegment
This is used to collect the traces inside a transaction, these interfaces should be called at the function entry point to identify the segments inside a transaction and publish metrics to ATC.
```
Before Instrumentation:
func retrunHomeOnly(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "home")
	fSegment1(r.Context())
}

func fSegment1(ctx context.Context) {
	fmt.Println("Inside Segment")
}
After Instrumentation:
func retrunHomeOnly(w http.ResponseWriter, r *http.Request) {
	txn := goprobe.StartTransaction(r.Context(), "/home")
	defer txn.EndTransaction()
	fmt.Fprintf(w, "home")
	fSegment1(r.Context())
}

func fSegment1(ctx context.Context) {
	seg := goprobe.StartSegment(ctx, "fSegment1")
	defer seg.EndSegment()
	fmt.Println("Inside Segment")
}
```
This is also used to collect the information from the application backends, this interface should be called at the function entry point where the application actually makes a database call with all the paramenters. In this case we need to pass two more additional paramenters like "dbConnectionString" and "actualQuery"
```
Before Instrumentation:
func sqlSegment(ctx context.Context) {
	//Actual Database call from the application
}
After Instrumentation:
func sqlSegment(ctx context.Context) {
	seg := goprobe.StartSegment(ctx, "query", "Server=lodsxvm58;Database=apm;Port=8080", "Select * from Booking where UserId = ?")
	defer seg.EndSegment()
	//Actual Database call from the application
}
```
## Metric description
- ProcessID - Represents the Process ID of the running application.
- GO Version - Represents the version of the Go language.
- Logical CPU Count - Represents the number of Logical CPUs that the current process currently uses.
- Host Name - Represents the Host name where the application is running.
- Total Heap Allocated - Represents the cumulative count of heap objects allocated.
- Total Heap Free - Represents the cumulative count of heap objects freed.
- Bytes In Use - Represents the Bytes of allocated heap objects.
- Bytes Total - Represents the cumulative bytes of allocated heap objects.
- Routines Total - Represents the number of running Go routines.

## Limitations
- Cross process Co-relation : Often transactions travel across multiple JVMs, CLRs or Node.js instances, or application services, depending on the environment. Collecting the full transaction path requires tracing synchronous and asynchronous calls across JVM, CLR, or Node.js instance boundaries. This ability lets you view details when transactions call methods on multiple JVMs or CLRs running on different servers.

- Attribute Decoration : Probes may have some useful information and it is critical to use as attributes. Information such as Environment details, static attributes and attributes from external properties which can be configured at the probe end needs to be propagated to ATC as custom attributes.

## Support
This document and associated tools are made available from CA Technologies as examples and provided at no charge as a courtesy to the CA APM Community at large. This resource may require modification for use in your environment. However, please note that this resource is not supported by CA Technologies, and inclusion in this site should not be construed to be an endorsement or recommendation by CA Technologies. These utilities are not covered by the CA Technologies software license agreement and there is no explicit or implied warranty from CA Technologies. They can be used and distributed freely amongst the CA APM Community, but not sold. As such, they are unsupported software, provided as is without warranty of any kind, express or implied, including but not limited to warranties of merchantability and fitness for a particular purpose. CA Technologies does not warrant that this resource will meet your requirements or that the operation of the resource will be uninterrupted or error free or that any defects will be corrected. The use of this resource implies that you understand and agree to the terms listed herein.

Although these utilities are unsupported, please let us know if you have any problems or questions by adding a comment to the CA APM Community Site area where the resource is located, so that the Author(s) may attempt to address the issue or question.

Unless explicitly stated otherwise this extension is only supported on the same platforms as the APM core agent. See [APM Compatibility Guide](http://www.ca.com/us/support/ca-support-online/product-content/status/compatibility-matrix/application-performance-management-compatibility-guide.aspx).

### Support URL
https://github.com/CA-APM/goprobe/issues

# Contributing
The [CA APM Community](https://communities.ca.com/community/ca-apm) is the primary means of interfacing with other users and with the CA APM product team.  The [developer subcommunity](https://communities.ca.com/community/ca-apm/ca-developer-apm) is where you can learn more about building APM-based assets, find code examples, and ask questions of other developers and the CA APM product team.

If you wish to contribute to this or any other project, please refer to [easy instructions](https://communities.ca.com/docs/DOC-231150910) available on the CA APM Developer Community.

# Change log
Changes for each version of the extension.

Version | Author | Comment
--------|--------|--------
1.0 | SrimanNarayana Vema | First version of the extension.
