package main

import (
	"partyframe/logger"
)

func main() {
	helper()
}

func helper() {
	logger.Infof("what's %s", "name")
	logger.Debugf("what's %s", "name")
	logger.Tracef("what's %s", "name")
	logger.Errorf("what's %s", "name")
	logger.Fatalf("what's %s", "name")
}
