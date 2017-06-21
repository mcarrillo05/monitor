package model

//SNMPObject is used for searching values of resource using SNMP.
type SNMPObject struct {
	Resource string `json:"resource"`
	OID      string `json:"OID"`
	Value    string `json:"value"`
	Walk     bool   `json:"-"`
}

const unknown = "Unknown"

//NewSNMPObject creates a new object using resource and search OID.
func NewSNMPObject(resource string) (SNMPObject, error) {
	object := SNMPObject{
		Resource: resource,
	}
	oID := getOID(resource)
	if oID == "" {
		return object, ErrInvalidResource
	}
	object.OID = oID
	// if resource == "DISK" {
	// 	object.Walk = true
	// }
	return object, nil
}

//getOID returns OID using specified resource.
func getOID(resource string) string {
	switch resource {
	case "SO":
		return "1.3.6.1.2.1.1.1.0"
	case "NAME":
		return "1.3.6.1.2.1.1.5.0"
	case "CPULOAD":
		return "1.3.6.1.4.1.2021.11.9.0"
	case "RAM":
		return "1.3.6.1.4.1.2021.4.5.0"
	case "RAMFREE":
		return "1.3.6.1.4.1.2021.4.6.0"
	case "DISK":
		return "1.3.6.1.2.1.25.2.3.1.5.1"
	case "UPTIME":
		return "1.3.6.1.2.1.25.1.1.0"
	case "SYSTEMDATE":
		return "1.3.6.1.2.1.25.1.2.0"
	case "PROCESSES":
		return "1.3.6.1.2.1.25.1.6.0"
	default:
		return ""
	}
}
