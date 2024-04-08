//
// Copyright (c) 2021 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
// 	// "github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/transforms"
// )

import (
	"errors"
	"fmt"
	"os"

	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg"
	"github.com/edgexfoundry/app-functions-sdk-go/v2/pkg/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"

	"encoding/json"
)

// // Function to print event data
// func printEventData(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
// 	// Convert the data interface to a byte array
// 	if data == nil {
// 		return false, errors.New("no event received")
// 	}
// 	ctx.LoggingClient().Debug("THIS IS A TEST ")
// 	fmt.Println("###################### THIS IS DATA ####################", data)

// 	// event, ok := data.(models.Event)
// 	// if !ok {
// 	// 	return false, errors.New("unexpected type received")
// 	// }
// 	// if len(event.Readings) == 0 {
// 	// 	return false, errors.New("no event readings to transform")
// 	// }

// 	// // print data here
// 	// fmt.Println("###################### THIS IS EVENT ####################", event)

// 	return false, nil
// }

func printEventData1(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Convert the data interface to a byte array
	if data == nil {
		return false, errors.New("no event received")
	}
	ctx.LoggingClient().Debug("THIS IS A TEST ")
	fmt.Println("###################### THIS IS DATA ####################", data)

	// Convert the incoming data to the desired format
	rawData := data.([]interface{})
	if len(rawData) != 7 {
		return false, errors.New("unexpected data format")
	}

	device := rawData[2].(string)
	profile := rawData[3].(string)
	source := rawData[4].(string)
	origin := rawData[5].(int64)

	readingData := rawData[6].([]interface{})
	if len(readingData) != 1 {
		return false, errors.New("unexpected reading data format")
	}

	reading := readingData[0].(map[string]interface{})
	readingID := reading["id"].(string)
	readingOrigin := reading["origin"].(int64)
	readingDevice := reading["device"].(string)
	readingResource := reading["resource"].(string)
	readingProfile := reading["profile"].(string)
	readingValueType := reading["valueType"].(string)
	readingUnits := reading["units"].(string)
	// readingValue := int(reading["value"].(float64))
	// readingValue := reading["value"]

	// Create the Event object
	event := dtos.Event{
		DeviceName:  device,
		ProfileName: profile,
		SourceName:  source,
		Origin:      origin,
		Readings: []dtos.BaseReading{
			{
				Id:           readingID,
				Origin:       readingOrigin,
				DeviceName:   readingDevice,
				ResourceName: readingResource,
				ProfileName:  readingProfile,
				ValueType:    readingValueType,
				Units:        readingUnits,
				// Value:     readingValue,
				// SimpleReading: {
				// 	Value: readingValue,
				// },

			},
		},
	}

	// Convert the Event object to JSON
	eventJSON, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return false, err
	}

	fmt.Println("###################### THIS IS EVENT ####################")
	fmt.Println(string(eventJSON))

	return false, nil
}

func printEventData2(ctx interfaces.AppFunctionContext, data interface{}) (bool, interface{}) {
	// Convert the data interface to a byte array
	if data == nil {
		return false, errors.New("no event received")
	}
	ctx.LoggingClient().Debug("THIS IS A TEST ")
	fmt.Println("###################### THIS IS DATA ####################", data)
	if event, ok := data.(dtos.Event); ok {
		fmt.Println("###################### THIS IS EVENT ####################", event.Readings)

		for _, reading_raw := range event.Readings {
			fmt.Println(reading_raw.Value)
			fmt.Println(reading_raw.DeviceName)
			fmt.Println(reading_raw.ProfileName)

		}

	}

	return false, nil
}

const (
	// serviceKey = "app-service-demo"
	serviceKey = "app-service-demo"
)

func main() {
	// turn off secure mode for examples. Not recommended for production
	_ = os.Setenv("EDGEX_SECURITY_SECRET_STORE", "false")

	// 1) First thing to do is to create an new instance of an EdgeX Application Service.
	service, ok := pkg.NewAppService(serviceKey)
	if !ok {
		os.Exit(-1)
	}

	// Leverage the built in logging service in EdgeX
	// lc := service.LoggingClient()

	// 2) shows how to access the application's specific configuration settings.
	deviceNames, err := service.GetAppSettingStrings("DeviceNames")

	fmt.Println("Filtering for devices", deviceNames)

	if err != nil {
		// lc.Error(err.Error())
		os.Exit(-1)
	}
	// lc.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))

	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	if err := service.SetFunctionsPipeline(
		// transforms.NewFilterFor(deviceNames).FilterByDeviceName,
		printEventData2,

	// Place code to print here

	); err != nil {
		// lc.Errorf("SetFunctionsPipeline returned error: %s", err.Error())
		os.Exit(-1)
	}

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = service.MakeItRun()
	if err != nil {
		// lc.Errorf("MakeItRun returned error: %s", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}
