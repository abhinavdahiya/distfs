package jobs

import "log"

//SendHeartbeats is a slave node job which sends a heartbeat to the master node
type SendHeartbeat struct {
	Master string
}

func (j SendHeartbeat) Run() {
	c := &heartbeat.Client{
		Dest: j.Master,
	}

	done, err := c.Send()
	if err != nil || !done {
		log.Error(err)
	} else {
		log.Info("Heartbeat sent successfully")
	}
}
