/*Package goprobe ...*/
/*
This is used to handle the panics in the interfaces exposed so that the user application
will not crash because of exceptions in the SDK
*/
package goprobe

//handlePanic is called in all the exposed interfaces and this will catch the exceptions
//raised in the interfaces and will make the connection in inactive status
func handlePanic() {
	if r := recover(); r != nil {
		getAppInstance().dataConn = false
		getAppInstance().cmdConn = false
		getAppInstance().GetLogger().Error("Panicking...!!!", map[string]interface{}{
			"Exception": r,
		})
	}
}
