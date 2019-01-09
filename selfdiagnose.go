package main

import (
	"log"
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

	webServer := &http.Server{Addr: ":8080"}
	log.Println("parzello HTTP api is listening on :8080")
	log.Println("open http://localhost:8080/internal/selfdiagnose.html")
	log.Fatal(webServer.ListenAndServe())
}
