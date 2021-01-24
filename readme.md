# kl-logger

Standalone data logger for [TFA Klimalogg pro](https://www.tfa-dostmann.de/produkt/profi-thermo-hygrometer-mit-datenlogger-funktion-klimalogg-pro-30-3039/) weather station to make latest weather data accessible via API.

Inspired by https://github.com/matthewwall/weewx-klimalogg 

## Requirements

- obviously a [TFA Klimalogg pro](https://www.tfa-dostmann.de/produkt/profi-thermo-hygrometer-mit-datenlogger-funktion-klimalogg-pro-30-3039/) plus usb transceiver dongle
- [libusb](https://libusb.info/) installed
- [go](https://golang.org/dl/) installed (to build project)

## Usage

- checkout project
- either build executable using `go build` or run `go run main.go` to execute program
- to pair console press and hold USB button as described in manual
- execute kl-logger immediately after console switched into paring mode
- once console is paired, one can press USB button again to reconnect if connection got lost
- current weather data is available via http at `http://localhost:8088/weather` (see [Endpoints](#Endpoints) for more)

### Options

```
Usage of kl-logger:
  -apiPort int
        Port to serve http api requests (default 8088)
  -log string
        Logfile (default "stdout")
  -logLevel string
        Log level (e.g. error, info, debug, trace) (default "info")
  -usbTrace
        Trace usb control messages
```

### Endpoints

|URI|Description|
|---|-----------|
| GET `/weather` | Latest weather information received from paired klimalogg console |
| GET `/weather/{sensor:[0-8]}` | Latest weather information by sensor id |
| GET `/config` | Current console configuration |
| GET `/config/{sensor:[0-8]}` | Current console configuration by sensor id |
| GET `/debug/transceiver/trace?seconds=5` | Trace usb transceiver control messages |
| GET `/debug/pprof` | Profiling endpoints provided by go [net/http/pprof](https://pkg.go.dev/net/http/pprof) |

## Troubleshooting

### libusb usually needs superuser permissions to access usb devices on linux
That can be fixed using udev rules to change device permissions as described [here](https://github.com/libusb/libusb/wiki/FAQ#Can_I_run_libusb_applications_on_Linux_without_root_privilege).

1. create new rule file in `/etc/udev/rules.d` (e.g. `10-klimalogg.rules`)
2. define a rule to grant permissions (e.g. `SUBSYSTEM=="usb", ATTR{idVendor}=="6666", ATTR{idProduct}=="5555", GROUP="dialout"`) 