package main

import (
	"fmt"
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
	"github.com/paypal/gatt/examples/service"
)

//--------------------------------------------------------------------------
// gatt service
//--------------------------------------------------------------------------
func CubeService() *gatt.Service {
	s := gatt.NewService(gatt.MustParseUUID("19fc95c0-c111-11e3-9904-0002a5d5c51b"))
	s.AddCharacteristic(gatt.MustParseUUID("44fac9e0-c111-11e3-9246-0002a5d5c51b")).HandleReadFunc(
		func(rsp gatt.ResponseWriter, req *gatt.ReadRequest) {
			log.Println("read request")
			fmt.Fprintf(rsp, "123456")
		})
	s.AddCharacteristic(gatt.MustParseUUID("45fac9e0-c111-11e3-9246-0002a5d5c51b")).HandleWriteFunc(
		func(r gatt.Request, data []byte) (status byte) {
			log.Println("write request", data)
			return gatt.StatusSuccess
		})

	return s
}

//--------------------------------------------------------------------------
// gatt server
//--------------------------------------------------------------------------
func setupBtServer() {
	d, err := gatt.NewDevice(option.DefaultServerOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s", err)
	}

	// Register optional handlers.
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) { fmt.Println("Connect: ", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { fmt.Println("Disconnect: ", c.ID()) }),
	)

	// A mandatory handler for monitoring device state.
	onStateChanged := func(d gatt.Device, s gatt.State) {
		fmt.Printf("State: %s\n", s)
		switch s {
		case gatt.StatePoweredOn:
			// Setup GAP and GATT services for Linux implementation.
			// OS X doesn't export the access of these services.
			d.AddService(service.NewGapService("InfinityCube")) // no effect on OS X
			d.AddService(service.NewGattService())              // no effect on OS X

			// A simple count service for demo.
			s1 := CubeService()
			d.AddService(s1)

			// Advertise device name and service's UUIDs.
			d.AdvertiseNameAndServices("Gopher", []gatt.UUID{s1.UUID()})

			// Advertise as an OpenBeacon iBeacon
			d.AdvertiseIBeacon(gatt.MustParseUUID("AA6062F098CA42118EC4193EB73CCEB6"), 1, 2, -59)

		default:
		}
	}

	d.Init(onStateChanged)
	select {}
}
