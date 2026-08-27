package main

import (
	"log"
)

func main() {
	config := LoadConfig()
	deps := &AppDeps{
		wells:     NewWellService(NewWellStore()),
		readings:  newReadingService(newReadingStore(5000)),
		alerts:    newAlertService(newAlertStore()),
		tasks:     newTaskService(newTaskStore()),
		plans:     newPlanService(newPlanStore()),
		ops:       newOpsService(nil),
		backfill:  newBackfillRunner(newReadingStore(5000)),
		operators: newOperatorStore(defaultOperators()),
		sessions:  newSessionStore(),
	}
	log.Printf("geothermal well monitor listening on :%s", config.Port)
	log.Fatal(serveAddress(":"+config.Port, NewRouter(deps)))
}

func defaultOperators() []Operator {
	return []Operator{
		{Username: "liang", Name: "梁工", Role: "operator", Password: "well2026"},
		{Username: "zhao", Name: "赵工", Role: "supervisor", Password: "ops2026"},
	}
}
