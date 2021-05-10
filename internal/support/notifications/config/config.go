/*******************************************************************************
 * Copyright 2018 Dell Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *
 *******************************************************************************/

package config

import (
	bootstrapConfig "github.com/edgexfoundry/go-mod-bootstrap/v2/config"
)

type ConfigurationStruct struct {
	Writable    WritableInfo
	Clients     map[string]bootstrapConfig.ClientInfo
	Databases   map[string]bootstrapConfig.Database
	Registry    bootstrapConfig.RegistryInfo
	Service     bootstrapConfig.ServiceInfo
	Smtp        SmtpInfo
	SecretStore bootstrapConfig.SecretStoreInfo
}

type WritableInfo struct {
	ResendLimit     int
	LogLevel        string
	InsecureSecrets bootstrapConfig.InsecureSecrets
}

type SmtpInfo struct {
	Host                 string
	Username             string
	Password             string
	Port                 int
	Sender               string
	EnableSelfSignedCert bool
	Subject              string
}

// The earlier releases do not have Username field and are using Sender field where Usename will
// be used now, to make it backward compatible fallback to Sender, which is signified by the empty
// Username field.
func (s SmtpInfo) CheckUsername() string {
	if s.Username != "" {
		return s.Username
	}
	return s.Sender
}

// UpdateFromRaw converts configuration received from the registry to a service-specific configuration struct which is
// then used to overwrite the service's existing configuration struct.
func (c *ConfigurationStruct) UpdateFromRaw(rawConfig interface{}) bool {
	configuration, ok := rawConfig.(*ConfigurationStruct)
	if ok {
		// Check that information was successfully read from Registry
		if configuration.Service.Port == 0 {
			return false
		}
		*c = *configuration
	}
	return ok
}

// EmptyWritablePtr returns a pointer to a service-specific empty WritableInfo struct.  It is used by the bootstrap to
// provide the appropriate structure to registry.Client's WatchForChanges().
func (c *ConfigurationStruct) EmptyWritablePtr() interface{} {
	return &WritableInfo{}
}

// UpdateWritableFromRaw converts configuration received from the registry to a service-specific WritableInfo struct
// which is then used to overwrite the service's existing configuration's WritableInfo struct.
func (c *ConfigurationStruct) UpdateWritableFromRaw(rawWritable interface{}) bool {
	writable, ok := rawWritable.(*WritableInfo)
	if ok {
		c.Writable = *writable
	}
	return ok
}

// GetBootstrap returns the configuration elements required by the bootstrap.  Currently, a copy of the configuration
// data is returned.  This is intended to be temporary -- since ConfigurationStruct drives the configuration.toml's
// structure -- until we can make backwards-breaking configuration.toml changes (which would consolidate these fields
// into an bootstrapConfig.BootstrapConfiguration struct contained within ConfigurationStruct).
func (c *ConfigurationStruct) GetBootstrap() bootstrapConfig.BootstrapConfiguration {
	// temporary until we can make backwards-breaking configuration.toml change
	return bootstrapConfig.BootstrapConfiguration{
		Clients:     c.Clients,
		Service:     c.Service,
		Registry:    c.Registry,
		SecretStore: c.SecretStore,
	}
}

// GetLogLevel returns the current ConfigurationStruct's log level.
func (c *ConfigurationStruct) GetLogLevel() string {
	return c.Writable.LogLevel
}

// GetRegistryInfo returns the RegistryInfo from the ConfigurationStruct.
func (c *ConfigurationStruct) GetRegistryInfo() bootstrapConfig.RegistryInfo {
	return c.Registry
}

// GetDatabaseInfo returns a database information map.
func (c *ConfigurationStruct) GetDatabaseInfo() map[string]bootstrapConfig.Database {
	return c.Databases
}

// GetInsecureSecrets returns the service's InsecureSecrets.
func (c *ConfigurationStruct) GetInsecureSecrets() bootstrapConfig.InsecureSecrets {
	return c.Writable.InsecureSecrets
}

// GetMessageBusInfo returns empty bootstrapConfig.MessageBusInfo since this is not used.
func (c *ConfigurationStruct) GetMessageBusInfo() bootstrapConfig.MessageBusInfo {
	return bootstrapConfig.MessageBusInfo{}
}
