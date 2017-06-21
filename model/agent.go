package model

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const agentsFile = "./agents.txt"
const separator = "||"

//Agent is an device registered in SNMP.
type Agent struct {
	IP        string `json:"ip" binding:"required"`
	Hostname  string `json:"hostname" binding:"required"`
	Community string `json:"community" binding:"required"`
	Name      string `json:"name"`
	SO        string `json:"so"`
	Uptime    string `json:"uptime"`
}

//GetAllAgents returns all agents registered in hosts file.
func GetAllAgents() ([]Agent, error) {
	var agents []Agent
	content, err := ioutil.ReadFile(agentsFile)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		data := strings.Split(line, separator)
		if len(data) > 1 {
			a := Agent{
				IP:        data[1],
				Hostname:  data[0],
				Community: data[2],
			}
			a.getBasic()
			agents = append(agents, a)
		}
	}
	return agents, nil
}

//Add registers a new Agent in hosts file.
func (a Agent) Add() error {
	file, err := os.OpenFile(agentsFile, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	if a.Get() == nil {
		return ErrDuplicate
	}
	str := (a.Hostname + separator + a.IP + separator + a.Community + "\n")
	_, err = file.WriteString(str)
	if err != nil {
		return err
	}
	return nil
}

//Get returns an agent using its IP.
func (a *Agent) Get() error {
	content, err := ioutil.ReadFile(agentsFile)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, a.IP) {
			data := strings.Split(line, separator)
			a.Hostname = data[0]
			a.Community = data[2]
			a.getBasic()
			return nil
		}
	}
	return ErrNotFound
}

func (a *Agent) getBasic() {
	o, err := a.GetOID("SO")
	if err != nil {
		a.SO = unknown
	} else {
		a.SO = o.Value
	}
	o, err = a.GetOID("NAME")
	if err != nil {
		a.Name = unknown
	} else {
		a.Name = o.Value
	}
	o, err = a.GetOID("UPTIME")
	if err != nil {
		a.Uptime = unknown
	} else {
		a.Uptime = o.Value
	}
}

//Delete removes an specific agent using its IP.
func (a *Agent) Delete() error {
	content, err := ioutil.ReadFile(agentsFile)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")

	idx := -1
	for i, line := range lines {
		if strings.Contains(line, a.IP) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return ErrNotFound
	}
	lines = append(lines[:idx], lines[idx+1:]...)
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(agentsFile, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

//GetOID returns value of specified OID.
func (a *Agent) GetOID(resource string) (SNMPObject, error) {
	object, err := NewSNMPObject(resource)
	if err != nil {
		return object, err
	}
	file := "snmpget.py"
	if object.Walk {
		file = "snmpwalk.py"
	}
	cmd := exec.Command("python", file, a.IP, a.Community, object.OID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return object, err
	}
	object.Value = string(out)
	if !object.Walk {
		if strings.Contains(object.Value, "timeout") {
			return object, ErrNotFound
		}
		object.Value = strings.TrimSuffix(strings.Split(object.Value, " = ")[1], "\n")
	} else {
		oID := ""
		rows := strings.Split(object.Value, "\n")
		for _, r := range rows {
			values := strings.Split(r, " = ")
			if oID != "" {
				if strings.HasSuffix(values[0], "6."+oID) {
					object.Value = values[1]
					break
				}
			} else {
				if values[1] == "/" {
					slice := strings.Split(values[0], ".")
					oID = slice[len(slice)-1]
				}
			}
		}
	}
	return object, nil
}
