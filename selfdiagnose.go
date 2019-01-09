package main

import (
	"net/http"

	selfdiagnose "github.com/emicklei/go-selfdiagnose"
	"github.com/emicklei/go-selfdiagnose/task"
)

func addSelfdiagnose() {
	// add http handlers for /internal/selfdiagnose.(html|json|xml)
	selfdiagnose.AddInternalHandlers()
	selfdiagnose.Register(task.ReportHttpRequest{})
	selfdiagnose.Register(task.ReportHostname{})
	selfdiagnose.Register(task.ReportCPU())

	// start a HTTP server
	webServer := &http.Server{Addr: ":8080"}
	logInfo("HTTP api is listening on :8080")
	logInfo("open http://localhost:8080/internal/selfdiagnose.html")
	webServer.ListenAndServe()
}
