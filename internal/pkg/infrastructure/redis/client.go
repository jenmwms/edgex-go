//
// Copyright (C) 2020-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"fmt"

	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	redisClient "github.com/edgexfoundry/edgex-go/internal/pkg/db/redis"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	model "github.com/edgexfoundry/go-mod-core-contracts/v2/models"

	"github.com/google/uuid"
)

type Client struct {
	*redisClient.Client
	loggingClient logger.LoggingClient
}

func NewClient(config db.Configuration, logger logger.LoggingClient) (*Client, errors.EdgeX) {
	var err error
	dc := &Client{}
	dc.Client, err = redisClient.NewClient(config, logger)
	dc.loggingClient = logger
	if err != nil {
		return nil, errors.NewCommonEdgeX(errors.KindDatabaseError, "redis client creation failed", err)
	}

	return dc, nil
}

// AddEvent adds a new event
func (c *Client) AddEvent(e model.Event) (model.Event, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if e.Id != "" {
		_, err := uuid.Parse(e.Id)
		if err != nil {
			return model.Event{}, errors.NewCommonEdgeX(errors.KindInvalidId, "uuid parsing failed", err)
		}
	}

	return addEvent(conn, e)
}

// EventById gets an event by id
func (c *Client) EventById(id string) (event model.Event, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	event, edgeXerr = eventById(conn, id)
	if edgeXerr != nil {
		return event, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// DeleteEventById removes an event by id
func (c *Client) DeleteEventById(id string) (edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr = deleteEventById(conn, id)
	if edgeXerr != nil {
		return errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// Add a new device profle
func (c *Client) AddDeviceProfile(dp model.DeviceProfile) (model.DeviceProfile, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if dp.Id != "" {
		_, err := uuid.Parse(dp.Id)
		if err != nil {
			return model.DeviceProfile{}, errors.NewCommonEdgeX(errors.KindInvalidId, "ID failed UUID parsing", err)
		}
	} else {
		dp.Id = uuid.New().String()
	}

	return addDeviceProfile(conn, dp)
}

// UpdateDeviceProfile updates a new device profile
func (c *Client) UpdateDeviceProfile(dp model.DeviceProfile) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateDeviceProfile(conn, dp)
}

// DeviceProfileNameExists checks the device profile exists by name
func (c *Client) DeviceProfileNameExists(name string) (bool, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()
	return deviceProfileNameExists(conn, name)
}

// AddDeviceService adds a new device service
func (c *Client) AddDeviceService(ds model.DeviceService) (model.DeviceService, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(ds.Id) == 0 {
		ds.Id = uuid.New().String()
	}

	return addDeviceService(conn, ds)
}

// DeviceServiceByName gets a device service by name
func (c *Client) DeviceServiceByName(name string) (deviceService model.DeviceService, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceService, edgeXerr = deviceServiceByName(conn, name)
	if edgeXerr != nil {
		return deviceService, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// DeviceServiceById gets a device service by id
func (c *Client) DeviceServiceById(id string) (deviceService model.DeviceService, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceService, edgeXerr = deviceServiceById(conn, id)
	if edgeXerr != nil {
		return deviceService, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// DeleteDeviceServiceById deletes a device service by id
func (c *Client) DeleteDeviceServiceById(id string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceServiceById(conn, id)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the device service with id %s", id), edgeXerr)
	}

	return nil
}

// DeleteDeviceServiceByName deletes a device service by name
func (c *Client) DeleteDeviceServiceByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceServiceByName(conn, name)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the device service with name %s", name), edgeXerr)
	}

	return nil
}

// DeviceServiceNameExists checks the device service exists by name
func (c *Client) DeviceServiceNameExists(name string) (bool, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()
	return deviceServiceNameExist(conn, name)
}

// UpdateDeviceService updates a device service
func (c *Client) UpdateDeviceService(ds model.DeviceService) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateDeviceService(conn, ds)
}

// DeviceProfileByName gets a device profile by name
func (c *Client) DeviceProfileByName(name string) (deviceProfile model.DeviceProfile, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceProfile, edgeXerr = deviceProfileByName(conn, name)
	if edgeXerr != nil {
		return deviceProfile, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return
}

// DeleteDeviceProfileById deletes a device profile by id
func (c *Client) DeleteDeviceProfileById(id string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceProfileById(conn, id)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the device profile with id %s", id), edgeXerr)
	}

	return nil
}

// DeleteDeviceProfileByName deletes a device profile by name
func (c *Client) DeleteDeviceProfileByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceProfileByName(conn, name)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the device profile with name %s", name), edgeXerr)
	}

	return nil
}

// AllDeviceProfiles query device profiles with offset and limit
func (c *Client) AllDeviceProfiles(offset int, limit int, labels []string) ([]model.DeviceProfile, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceProfiles, edgeXerr := deviceProfilesByLabels(conn, offset, limit, labels)
	if edgeXerr != nil {
		return deviceProfiles, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return deviceProfiles, nil
}

// DeviceProfilesByModel query device profiles with offset, limit and model
func (c *Client) DeviceProfilesByModel(offset int, limit int, model string) ([]model.DeviceProfile, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceProfiles, edgeXerr := deviceProfilesByModel(conn, offset, limit, model)
	if edgeXerr != nil {
		return deviceProfiles, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return deviceProfiles, nil
}

// DeviceProfilesByManufacturer query device profiles with offset, limit and manufacturer
func (c *Client) DeviceProfilesByManufacturer(offset int, limit int, manufacturer string) ([]model.DeviceProfile, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceProfiles, edgeXerr := deviceProfilesByManufacturer(conn, offset, limit, manufacturer)
	if edgeXerr != nil {
		return deviceProfiles, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return deviceProfiles, nil
}

// DeviceProfilesByManufacturerAndModel query device profiles with offset, limit, manufacturer and model
func (c *Client) DeviceProfilesByManufacturerAndModel(offset int, limit int, manufacturer string, model string) ([]model.DeviceProfile, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceProfiles, edgeXerr := deviceProfilesByManufacturerAndModel(conn, offset, limit, manufacturer, model)
	if edgeXerr != nil {
		return deviceProfiles, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return deviceProfiles, nil
}

// EventTotalCount returns the total count of Event from the database
func (c *Client) EventTotalCount() (uint32, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	count, edgeXerr := getMemberNumber(conn, ZCARD, EventsCollection)
	if edgeXerr != nil {
		return 0, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return count, nil
}

// EventCountByDevice returns the count of Event associated a specific Device from the database
func (c *Client) EventCountByDeviceName(deviceName string) (uint32, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	count, edgeXerr := getMemberNumber(conn, ZCARD, CreateKey(EventsCollectionDeviceName, deviceName))
	if edgeXerr != nil {
		return 0, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return count, nil
}

// AllDeviceServices returns multiple device services per query criteria, including
// offset: the number of items to skip before starting to collect the result set
// limit: The numbers of items to return
// labels: allows for querying a given object by associated user-defined labels
func (c *Client) AllDeviceServices(offset int, limit int, labels []string) (deviceServices []model.DeviceService, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	deviceServices, edgeXerr = deviceServicesByLabels(conn, offset, limit, labels)
	if edgeXerr != nil {
		return deviceServices, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return deviceServices, nil
}

// Add a new device
func (c *Client) AddDevice(d model.Device) (model.Device, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(d.Id) == 0 {
		d.Id = uuid.New().String()
	}

	return addDevice(conn, d)
}

// DeleteDeviceById deletes a device by id
func (c *Client) DeleteDeviceById(id string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceById(conn, id)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the device with id %s", id), edgeXerr)
	}

	return nil
}

// DeleteDeviceByName deletes a device by name
func (c *Client) DeleteDeviceByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteDeviceByName(conn, name)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the device with name %s", name), edgeXerr)
	}

	return nil
}

// DevicesByServiceName query devices by offset, limit and name
func (c *Client) DevicesByServiceName(offset int, limit int, name string) (devices []model.Device, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	devices, edgeXerr = devicesByServiceName(conn, offset, limit, name)
	if edgeXerr != nil {
		return devices, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query devices by offset %d, limit %d and name %s", offset, limit, name), edgeXerr)
	}
	return devices, nil
}

// DeviceIdExists checks the device existence by id
func (c *Client) DeviceIdExists(id string) (bool, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()
	exists, err := deviceIdExists(conn, id)
	if err != nil {
		return exists, errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to check the device existence by id %s", id), err)
	}
	return exists, nil
}

// DeviceNameExists checks the device existence by name
func (c *Client) DeviceNameExists(name string) (bool, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()
	exists, err := deviceNameExists(conn, name)
	if err != nil {
		return exists, errors.NewCommonEdgeX(errors.Kind(err), fmt.Sprintf("fail to check the device existence by name %s", name), err)
	}
	return exists, nil
}

// DeviceById gets a device by id
func (c *Client) DeviceById(id string) (device model.Device, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	device, edgeXerr = deviceById(conn, id)
	if edgeXerr != nil {
		return device, errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to query device by id %s", id), edgeXerr)
	}

	return
}

// DeviceByName gets a device by name
func (c *Client) DeviceByName(name string) (device model.Device, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	device, edgeXerr = deviceByName(conn, name)
	if edgeXerr != nil {
		return device, errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to query device by name %s", name), edgeXerr)
	}

	return
}

// DevicesByProfileName query devices by offset, limit and profile name
func (c *Client) DevicesByProfileName(offset int, limit int, profileName string) (devices []model.Device, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	devices, edgeXerr = devicesByProfileName(conn, offset, limit, profileName)
	if edgeXerr != nil {
		return devices, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query devices by offset %d, limit %d and name %s", offset, limit, profileName), edgeXerr)
	}
	return devices, nil
}

// Update a device
func (c *Client) UpdateDevice(d model.Device) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	return updateDevice(conn, d)
}

// AllEvents query events by offset and limit
func (c *Client) AllEvents(offset int, limit int) ([]model.Event, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	events, edgeXerr := c.allEvents(conn, offset, limit)
	if edgeXerr != nil {
		return events, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query events by offset %d and limit %d", offset, limit), edgeXerr)
	}
	return events, nil
}

// AllDevices query the devices with offset, limit, and labels
func (c *Client) AllDevices(offset int, limit int, labels []string) ([]model.Device, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	devices, edgeXerr := devicesByLabels(conn, offset, limit, labels)
	if edgeXerr != nil {
		return devices, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return devices, nil
}

// EventsByDeviceName query events by offset, limit and device name
func (c *Client) EventsByDeviceName(offset int, limit int, name string) (events []model.Event, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	events, edgeXerr = eventsByDeviceName(conn, offset, limit, name)
	if edgeXerr != nil {
		return events, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query events by offset %d, limit %d and name %s", offset, limit, name), edgeXerr)
	}
	return events, nil
}

// EventsByTimeRange query events by time range, offset, and limit
func (c *Client) EventsByTimeRange(startTime int, endTime int, offset int, limit int) (events []model.Event, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	events, edgeXerr = eventsByTimeRange(conn, startTime, endTime, offset, limit)
	if edgeXerr != nil {
		return events, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query events by time range %v ~ %v, offset %d, and limit %d", startTime, endTime, offset, limit), edgeXerr)
	}
	return events, nil
}

// ReadingTotalCount returns the total count of Event from the database
func (c *Client) ReadingTotalCount() (uint32, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	count, edgeXerr := getMemberNumber(conn, ZCARD, ReadingsCollection)
	if edgeXerr != nil {
		return 0, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return count, nil
}

// AllReadings query events by offset, limit, and labels
func (c *Client) AllReadings(offset int, limit int) ([]model.Reading, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	readings, edgeXerr := allReadings(conn, offset, limit)
	if edgeXerr != nil {
		return readings, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query readings by offset %d, and limit %d", offset, limit), edgeXerr)
	}
	return readings, nil
}

// ReadingsByTimeRange query readings by time range, offset, and limit
func (c *Client) ReadingsByTimeRange(start int, end int, offset int, limit int) (readings []model.Reading, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	readings, edgeXerr = readingsByTimeRange(conn, start, end, offset, limit)
	if edgeXerr != nil {
		return readings, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query readings by time range %v ~ %v, offset %d, and limit %d", start, end, offset, limit), edgeXerr)
	}
	return readings, nil
}

// ReadingsByResourceName query readings by offset, limit and resource name
func (c *Client) ReadingsByResourceName(offset int, limit int, resourceName string) (readings []model.Reading, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	readings, edgeXerr = readingsByResourceName(conn, offset, limit, resourceName)
	if edgeXerr != nil {
		return readings, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query readings by offset %d, limit %d and resourceName %s", offset, limit, resourceName), edgeXerr)
	}
	return readings, nil
}

// ReadingsByDeviceName query readings by offset, limit and device name
func (c *Client) ReadingsByDeviceName(offset int, limit int, name string) (readings []model.Reading, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	readings, edgeXerr = readingsByDeviceName(conn, offset, limit, name)
	if edgeXerr != nil {
		return readings, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query readings by offset %d, limit %d and name %s", offset, limit, name), edgeXerr)
	}
	return readings, nil
}

// ReadingCountByDeviceName returns the count of Readings associated a specific Device from the database
func (c *Client) ReadingCountByDeviceName(deviceName string) (uint32, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	count, edgeXerr := getMemberNumber(conn, ZCARD, CreateKey(ReadingsCollectionDeviceName, deviceName))
	if edgeXerr != nil {
		return 0, errors.NewCommonEdgeXWrapper(edgeXerr)
	}

	return count, nil
}

// ReadingsByResourceNameAndTimeRange query readings by resourceName and specified time range. Readings are sorted in descending order of origin time.
func (c *Client) ReadingsByResourceNameAndTimeRange(resourceName string, start int, end int, offset int, limit int) (readings []model.Reading, err errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	readings, err = readingsByResourceNameAndTimeRange(conn, resourceName, start, end, offset, limit)
	if err != nil {
		return readings, errors.NewCommonEdgeX(errors.Kind(err),
			fmt.Sprintf("fail to query readings by resourceName %s and time range %v ~ %v, offset %d, and limit %d", resourceName, start, end, offset, limit), err)
	}
	return readings, nil
}

func (c *Client) ReadingsByDeviceNameAndResourceName(deviceName string, resourceName string, offset int, limit int) (readings []model.Reading, err errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	readings, err = readingsByDeviceNameAndResourceName(conn, deviceName, resourceName, offset, limit)
	if err != nil {
		return readings, errors.NewCommonEdgeX(errors.Kind(err),
			fmt.Sprintf("fail to query readings by deviceName %s and resourceName %s", deviceName, resourceName), err)
	}

	return readings, nil
}

func (c *Client) ReadingsByDeviceNameAndResourceNameAndTimeRange(deviceName string, resourceName string, start int, end int, offset int, limit int) (readings []model.Reading, err errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	readings, err = readingsByDeviceNameAndResourceNameAndTimeRange(conn, deviceName, resourceName, start, end, offset, limit)
	if err != nil {
		return readings, errors.NewCommonEdgeX(errors.Kind(err),
			fmt.Sprintf("fail to query readings by deviceName %s, resourceName %s and time range %v ~ %v", deviceName, resourceName, start, end), err)
	}

	return readings, nil
}

// AddProvisionWatcher adds a new provision watcher
func (c *Client) AddProvisionWatcher(pw model.ProvisionWatcher) (model.ProvisionWatcher, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(pw.Id) == 0 {
		pw.Id = uuid.New().String()
	}

	return addProvisionWatcher(conn, pw)
}

// ProvisionWatcherById gets a provision watcher by id
func (c *Client) ProvisionWatcherById(id string) (provisionWatcher model.ProvisionWatcher, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	provisionWatcher, edgexErr = provisionWatcherById(conn, id)
	if edgexErr != nil {
		return provisionWatcher, errors.NewCommonEdgeX(errors.Kind(edgexErr), fmt.Sprintf("failed to query provision watcher by id %s", id), edgexErr)
	}

	return
}

// ProvisionWatcherByName gets a provision watcher by name
func (c *Client) ProvisionWatcherByName(name string) (provisionWatcher model.ProvisionWatcher, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	provisionWatcher, edgexErr = provisionWatcherByName(conn, name)
	if edgexErr != nil {
		return provisionWatcher, errors.NewCommonEdgeXWrapper(edgexErr)
	}

	return
}

//ProvisionWatchersByServiceName query provision watchers by offset, limit and service name
func (c *Client) ProvisionWatchersByServiceName(offset int, limit int, name string) (provisionWatchers []model.ProvisionWatcher, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	provisionWatchers, edgexErr = provisionWatchersByServiceName(conn, offset, limit, name)
	if edgexErr != nil {
		return provisionWatchers, errors.NewCommonEdgeX(errors.Kind(edgexErr),
			fmt.Sprintf("failed to query provision watcher by offset %d, limit %d and service name %s", offset, limit, name), edgexErr)
	}

	return
}

//ProvisionWatchersByProfileName query provision watchers by offset, limit and profile name
func (c *Client) ProvisionWatchersByProfileName(offset int, limit int, name string) (provisionWatchers []model.ProvisionWatcher, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	provisionWatchers, edgexErr = provisionWatchersByProfileName(conn, offset, limit, name)
	if edgexErr != nil {
		return provisionWatchers, errors.NewCommonEdgeX(errors.Kind(edgexErr),
			fmt.Sprintf("failed to query provision watcher by offset %d, limit %d and profile name %s", offset, limit, name), edgexErr)
	}

	return
}

// AllProvisionWatchers query provision watchers with offset, limit and labels
func (c *Client) AllProvisionWatchers(offset int, limit int, labels []string) (provisionWatchers []model.ProvisionWatcher, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	provisionWatchers, edgexErr = provisionWatchersByLabels(conn, offset, limit, labels)
	if edgexErr != nil {
		return provisionWatchers, errors.NewCommonEdgeXWrapper(edgexErr)
	}

	return
}

// DeleteProvisionWatcherByName deletes a provision watcher by name
func (c *Client) DeleteProvisionWatcherByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteProvisionWatcherByName(conn, name)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("failed to delete the provision watcher with name %s", name), edgeXerr)
	}

	return nil
}

// Update a provision watcher
func (c *Client) UpdateProvisionWatcher(pw model.ProvisionWatcher) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	return updateProvisionWatcher(conn, pw)
}

// AddInterval adds a new interval
func (c *Client) AddInterval(interval model.Interval) (model.Interval, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(interval.Id) == 0 {
		interval.Id = uuid.New().String()
	}

	return addInterval(conn, interval)
}

// IntervalByName gets a interval by name
func (c *Client) IntervalByName(name string) (interval model.Interval, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	interval, edgeXerr = intervalByName(conn, name)
	if edgeXerr != nil {
		return interval, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return
}

// IntervalById gets a interval by id
func (c *Client) IntervalById(id string) (interval model.Interval, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	interval, edgeXerr = intervalById(conn, id)
	if edgeXerr != nil {
		return interval, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return
}

// AllIntervals query intervals with offset and limit
func (c *Client) AllIntervals(offset int, limit int) (intervals []model.Interval, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	intervals, edgeXerr = allIntervals(conn, offset, limit)
	if edgeXerr != nil {
		return intervals, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return intervals, nil
}

// UpdateInterval updates a interval
func (c *Client) UpdateInterval(interval model.Interval) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateInterval(conn, interval)
}

// DeleteIntervalByName deletes the interval by name
func (c *Client) DeleteIntervalByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteIntervalByName(conn, name)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the interval with name %s", name), edgeXerr)
	}

	return nil
}

// AddIntervalAction adds a new intervalAction
func (c *Client) AddIntervalAction(action model.IntervalAction) (model.IntervalAction, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(action.Id) == 0 {
		action.Id = uuid.New().String()
	}
	return addIntervalAction(conn, action)
}

// AllIntervalActions query intervalActions with offset and limit
func (c *Client) AllIntervalActions(offset int, limit int) (intervalActions []model.IntervalAction, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	intervalActions, edgeXerr = allIntervalActions(conn, offset, limit)
	if edgeXerr != nil {
		return intervalActions, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return intervalActions, nil
}

// IntervalActionByName gets a intervalAction by name
func (c *Client) IntervalActionByName(name string) (action model.IntervalAction, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	action, edgeXerr = intervalActionByName(conn, name)
	if edgeXerr != nil {
		return action, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return
}

// IntervalActionsByIntervalName query intervalActions by offset, limit and intervalName
func (c *Client) IntervalActionsByIntervalName(offset int, limit int, intervalName string) (actions []model.IntervalAction, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	actions, edgeXerr = intervalActionsByIntervalName(conn, offset, limit, intervalName)
	if edgeXerr != nil {
		return actions, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query actions by offset %d, limit %d and intervalName %s", offset, limit, intervalName), edgeXerr)
	}
	return actions, nil
}

// DeleteIntervalActionByName deletes the intervalAction by name
func (c *Client) DeleteIntervalActionByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteIntervalActionByName(conn, name)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the intervalAction with name %s", name), edgeXerr)
	}

	return nil
}

// IntervalActionById gets a intervalAction by id
func (c *Client) IntervalActionById(id string) (action model.IntervalAction, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	action, edgeXerr = intervalActionById(conn, id)
	if edgeXerr != nil {
		return action, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return
}

// UpdateIntervalAction updates a intervalAction
func (c *Client) UpdateIntervalAction(action model.IntervalAction) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateIntervalAction(conn, action)
}

// AddSubscription adds a new subscription
func (c *Client) AddSubscription(subscription model.Subscription) (model.Subscription, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(subscription.Id) == 0 {
		subscription.Id = uuid.New().String()
	}

	return addSubscription(conn, subscription)
}

// AllSubscriptions returns multiple subscriptions per query criteria, including
// offset: The number of items to skip before starting to collect the result set.
// limit: The maximum number of items to return.
func (c *Client) AllSubscriptions(offset int, limit int) ([]model.Subscription, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	subscriptions, edgeXerr := allSubscriptions(conn, offset, limit)
	if edgeXerr != nil {
		return subscriptions, errors.NewCommonEdgeXWrapper(edgeXerr)
	}
	return subscriptions, nil
}

// SubscriptionsByCategory queries subscriptions by offset, limit and category
func (c *Client) SubscriptionsByCategory(offset int, limit int, category string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	subscriptions, edgeXerr = subscriptionsByCategory(conn, offset, limit, category)
	if edgeXerr != nil {
		return subscriptions, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query subscriptions by offset %d, limit %d and category %s", offset, limit, category), edgeXerr)
	}
	return subscriptions, nil
}

// SubscriptionsByLabel queries subscriptions by offset, limit and label
func (c *Client) SubscriptionsByLabel(offset int, limit int, label string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	subscriptions, edgeXerr = subscriptionsByLabel(conn, offset, limit, label)
	if edgeXerr != nil {
		return subscriptions, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query subscriptions by offset %d, limit %d and label %s", offset, limit, label), edgeXerr)
	}
	return subscriptions, nil
}

// SubscriptionsByReceiver queries subscriptions by offset, limit and receiver
func (c *Client) SubscriptionsByReceiver(offset int, limit int, receiver string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	subscriptions, edgeXerr = subscriptionsByReceiver(conn, offset, limit, receiver)
	if edgeXerr != nil {
		return subscriptions, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query subscriptions by offset %d, limit %d and receiver %s", offset, limit, receiver), edgeXerr)
	}
	return subscriptions, nil
}

// SubscriptionById gets a subscription by id
func (c *Client) SubscriptionById(id string) (subscription model.Subscription, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	subscription, edgexErr = subscriptionById(conn, id)
	if edgexErr != nil {
		return subscription, errors.NewCommonEdgeX(errors.Kind(edgexErr), fmt.Sprintf("failed to query subscription by id %s", id), edgexErr)
	}

	return
}

// SubscriptionByName queries subscription by name
func (c *Client) SubscriptionByName(name string) (subscription model.Subscription, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	subscription, edgeXerr = subscriptionByName(conn, name)
	if edgeXerr != nil {
		return subscription, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query subscription by name %s", name), edgeXerr)
	}
	return subscription, nil
}

// UpdateSubscription updates a new subscription
func (c *Client) UpdateSubscription(subscription model.Subscription) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateSubscription(conn, subscription)
}

// DeleteSubscriptionByName deletes a subscription by name
func (c *Client) DeleteSubscriptionByName(name string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteSubscriptionByName(conn, name)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the subscription with name %s", name), edgeXerr)
	}

	return nil
}

// SubscriptionsByCategoriesAndLabels queries subscriptions by offset, limit, categories and labels
func (c *Client) SubscriptionsByCategoriesAndLabels(offset int, limit int, categories []string, labels []string) (subscriptions []model.Subscription, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	subscriptions, edgeXerr = subscriptionsByCategoriesAndLabels(conn, offset, limit, categories, labels)
	if edgeXerr != nil {
		return subscriptions, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query subscriptions by offset %d, limit %d, categories %v and labels %v", offset, limit, categories, labels), edgeXerr)
	}
	return subscriptions, nil
}

// AddNotification adds a new notification
func (c *Client) AddNotification(notification model.Notification) (model.Notification, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(notification.Id) == 0 {
		notification.Id = uuid.New().String()
	}

	return addNotification(conn, notification)
}

// NotificationsByCategory queries notifications by offset, limit and category
func (c *Client) NotificationsByCategory(offset int, limit int, category string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	notifications, edgeXerr = notificationsByCategory(conn, offset, limit, category)
	if edgeXerr != nil {
		return notifications, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query notifications by offset %d, limit %d and category %s", offset, limit, category), edgeXerr)
	}
	return notifications, nil
}

// NotificationsByLabel queries notifications by offset, limit and label
func (c *Client) NotificationsByLabel(offset int, limit int, label string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	notifications, edgeXerr = notificationsByLabel(conn, offset, limit, label)
	if edgeXerr != nil {
		return notifications, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query notifications by offset %d, limit %d and label %s", offset, limit, label), edgeXerr)
	}
	return notifications, nil
}

// NotificationById gets a notification by id
func (c *Client) NotificationById(id string) (notification model.Notification, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	notification, edgexErr = notificationById(conn, id)
	if edgexErr != nil {
		return notification, errors.NewCommonEdgeX(errors.Kind(edgexErr), fmt.Sprintf("failed to query notification by id %s", id), edgexErr)
	}
	return
}

// NotificationsByStatus queries notifications by offset, limit and status
func (c *Client) NotificationsByStatus(offset int, limit int, status string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	notifications, edgeXerr = notificationsByStatus(conn, offset, limit, status)
	if edgeXerr != nil {
		return notifications, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query notifications by offset %d, limit %d and status %s", offset, limit, status), edgeXerr)
	}
	return notifications, nil
}

// NotificationsByTimeRange query notifications by time range, offset, and limit
func (c *Client) NotificationsByTimeRange(start int, end int, offset int, limit int) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	notifications, edgeXerr = notificationsByTimeRange(conn, start, end, offset, limit)
	if edgeXerr != nil {
		return notifications, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query notifications by time range %v ~ %v, offset %d, and limit %d", start, end, offset, limit), edgeXerr)
	}
	return notifications, nil
}

// NotificationsByCategoriesAndLabels queries notifications by offset, limit, categories and labels
func (c *Client) NotificationsByCategoriesAndLabels(offset int, limit int, categories []string, labels []string) (notifications []model.Notification, edgeXerr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	notifications, edgeXerr = notificationsByCategoriesAndLabels(conn, offset, limit, categories, labels)
	if edgeXerr != nil {
		return notifications, errors.NewCommonEdgeX(errors.Kind(edgeXerr),
			fmt.Sprintf("fail to query notifications by offset %d, limit %d, categories %v and labels %v", offset, limit, categories, labels), edgeXerr)
	}
	return notifications, nil
}

// DeleteNotificationById deletes a notification by id
func (c *Client) DeleteNotificationById(id string) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()

	edgeXerr := deleteNotificationById(conn, id)
	if edgeXerr != nil {
		return errors.NewCommonEdgeX(errors.Kind(edgeXerr), fmt.Sprintf("fail to delete the notification with id %s", id), edgeXerr)
	}

	return nil
}

// UpdateNotification updates a notification
func (c *Client) UpdateNotification(n model.Notification) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateNotification(conn, n)
}

// AddTransmission adds a new transmission
func (c *Client) AddTransmission(t model.Transmission) (model.Transmission, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	if len(t.Id) == 0 {
		t.Id = uuid.New().String()
	}

	return addTransmission(conn, t)
}

// UpdateTransmission updates a transmission
func (c *Client) UpdateTransmission(trans model.Transmission) errors.EdgeX {
	conn := c.Pool.Get()
	defer conn.Close()
	return updateTransmission(conn, trans)
}

// TransmissionById gets a transmission by id
func (c *Client) TransmissionById(id string) (trans model.Transmission, edgexErr errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	trans, edgexErr = transmissionById(conn, id)
	if edgexErr != nil {
		return trans, errors.NewCommonEdgeX(errors.Kind(edgexErr), fmt.Sprintf("failed to query transmission by id %s", id), edgexErr)
	}
	return
}

// TransmissionsByTimeRange query transmissions by time range, offset, and limit
func (c *Client) TransmissionsByTimeRange(start int, end int, offset int, limit int) (transmissions []model.Transmission, err errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	transmissions, err = transmissionsByTimeRange(conn, start, end, offset, limit)
	if err != nil {
		return transmissions, errors.NewCommonEdgeX(errors.Kind(err),
			fmt.Sprintf("fail to query transmissions by time range %v ~ %v, offset %d, and limit %d", start, end, offset, limit), err)
	}
	return transmissions, nil
}

// AllTransmissions returns multiple transmissions per query criteria, including
// offset: The number of items to skip before starting to collect the result set.
// limit: The maximum number of items to return.
func (c *Client) AllTransmissions(offset int, limit int) ([]model.Transmission, errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	transmission, err := allTransmissions(conn, offset, limit)
	if err != nil {
		return transmission, errors.NewCommonEdgeXWrapper(err)
	}
	return transmission, nil
}

// TransmissionsByStatus queries transmissions by offset, limit and status
func (c *Client) TransmissionsByStatus(offset int, limit int, status string) (transmissions []model.Transmission, err errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	transmissions, err = transmissionsByStatus(conn, offset, limit, status)
	if err != nil {
		return transmissions, errors.NewCommonEdgeX(errors.Kind(err),
			fmt.Sprintf("fail to query transmissions by offset %d, limit %d and status %s", offset, limit, status), err)
	}
	return transmissions, nil
}

// TransmissionsBySubscriptionName queries transmissions by offset, limit and subscription name
func (c *Client) TransmissionsBySubscriptionName(offset int, limit int, subscriptionName string) (transmissions []model.Transmission, err errors.EdgeX) {
	conn := c.Pool.Get()
	defer conn.Close()

	transmissions, err = transmissionsBySubscriptionName(conn, offset, limit, subscriptionName)
	if err != nil {
		return transmissions, errors.NewCommonEdgeX(errors.Kind(err),
			fmt.Sprintf("fail to query transmissions by offset %d, limit %d and subscription name %s", offset, limit, subscriptionName), err)
	}
	return transmissions, nil
}
