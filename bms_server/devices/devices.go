package devices

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/karalabe/hid"
	"github.com/ronnie-w/bms-server/database"
	"github.com/ronnie-w/bms-server/utils"
)

func GetDevices(vendorId, deviceId uint16) []hid.DeviceInfo {
	return hid.Enumerate(vendorId, deviceId)
}

func ConnectQRScanner(rw http.ResponseWriter, r *http.Request) {
	var vendorIdStr, deviceIdStr string
	if err := database.Conn().QueryRow("SELECT qr_code_scanner_vendor_id, qr_code_scanner_device_id FROM settings").Scan(&vendorIdStr, &deviceIdStr); err != nil {
		utils.ErrResponse(rw, nil, "Could not find vendor id or device id")
		return
	}

	vendorId, _ := strconv.ParseUint(vendorIdStr, 16, 16)
	deviceId, _ := strconv.ParseUint(deviceIdStr, 16, 16)

	qrCodeDevices := GetDevices(uint16(vendorId), uint16(deviceId))
	if len(qrCodeDevices) == 0 {
		utils.ErrResponse(rw, nil, "No devices with the given IDs found")
		return
	}

	qrCodeDevice, err := qrCodeDevices[0].Open()
	if err != nil {
		utils.ErrResponse(rw, nil, "Failed to open device")
		return
	}

	dataChan := make(chan []byte)
	go func() {
		for {
			data := make([]byte, 64)
			n, err := qrCodeDevice.Read(data)
			if err != nil {
				fmt.Println("Error reading from device")
			}

			if n > 0 {
				dataChan <- data[:n]
			}
		}
	}()

	for {
		select {
		case qrData := <-dataChan:
			log.Println(string(qrData))
		}
	}
}
