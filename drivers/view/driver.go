package view

import "errors"

// A Driver can initialize a Cache from a directory structure
type Driver interface {
	Init(root string) (*Cache, error)
}

var driverList = map[string]Driver{}

// DriverRegister registers a Driver with the package. A Driver must be
// registered before it can be used.
func DriverRegister(name string, d Driver) { driverList[name] = d }

// DriverAvailable returns true if the given driver has been registered.
func DriverAvailable(name string) bool { return driverList[name] != nil }

// ErrDriverUnavailable indicates that a given driver has not been registered,
// but was needed to complete some operation.
var ErrDriverUnavailable = errors.New("driver unavailable")

func DriverInit(driverName, root string) (*Cache, error) {
	driver := driverList[driverName]
	if driver == nil {
		return nil, ErrDriverUnavailable
	}
	return driver.Init(root)
}
