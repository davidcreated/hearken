package hearken

// Device is something that can be switched on and off.
type Device struct {
	ID string `json:"id"`
	// HardwareID is the device's own identifier, when it differs from ID.
	HardwareID string `json:"hardwareId,omitempty"`
	Name       string `json:"name"`
}

// Space is a room or area, and the devices in it.
type Space struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// DeviceIDs lists the devices in the space, by ID or HardwareID.
	DeviceIDs []string `json:"deviceIds,omitempty"`
}

// Mode is a scene: a saved set of device states, such as "Home" or
// "Movie night". A mode belongs to a space, or to the whole home when
// SpaceID is empty.
type Mode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// VoiceName is the name to say when it differs from Name, for example
	// "Office Boss Time" when several spaces have a "Boss Time".
	VoiceName string `json:"voiceName,omitempty"`
	// Draft modes are not finished and are never matched.
	Draft     bool   `json:"draft,omitempty"`
	SpaceID   string `json:"spaceId,omitempty"`
	SpaceName string `json:"spaceName,omitempty"`
}

// Home is everything a spoken name can refer to.
type Home struct {
	Devices []Device `json:"devices"`
	Spaces  []Space  `json:"spaces"`
	Modes   []Mode   `json:"modes"`
}

// DevicesIn returns the home's devices that are in space.
func (h Home) DevicesIn(space Space) []Device {
	ids := make(map[string]bool, len(space.DeviceIDs))
	for _, id := range space.DeviceIDs {
		ids[id] = true
	}
	var devices []Device
	for _, d := range h.Devices {
		if ids[d.ID] || (d.HardwareID != "" && ids[d.HardwareID]) {
			devices = append(devices, d)
		}
	}
	return devices
}
