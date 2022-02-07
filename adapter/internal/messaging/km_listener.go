/*
 *  Copyright (c) 2021, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

// Package messaging holds the implementation for event listeners functions
package messaging

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wso2/product-microgateway/adapter/internal/discovery/xds"
	logger "github.com/wso2/product-microgateway/adapter/internal/loggers"
	eventhubTypes "github.com/wso2/product-microgateway/adapter/pkg/eventhub/types"
	"github.com/wso2/product-microgateway/adapter/pkg/logging"
	msg "github.com/wso2/product-microgateway/adapter/pkg/messaging"
)

// constants related to key manager events
const (
	keyManagerConfigEvent = "key_manager_configuration"
	actionAdd             = "add"
	actionUpdate          = "update"
	actionDelete          = "delete"
	superTenantDomain     = "carbon.super"
)

// handleKMEvent
func handleKMConfiguration() {
	var (
		indexOfKeymanager int
		isFound           bool
	)
	for d := range msg.KeyManagerChannel {
		var notification msg.EventKeyManagerNotification
		// var keyManagerConfig resourceTypes.KeymanagerConfig
		var kmConfigMap map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(string(d.Body)), &notification)
		if unmarshalErr != nil {
			logger.LoggerInternalMsg.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error occurred while unmarshalling key manager event data %v", unmarshalErr.Error()),
				Severity:  logging.MINOR,
				ErrorCode: 1240,
			})
			return
		}
		logger.LoggerInternalMsg.Infof("Event %s is received", notification.Event.PayloadData.EventType)
		for i := range xds.KeyManagerList {
			if strings.EqualFold(notification.Event.PayloadData.Name, xds.KeyManagerList[i].Name) {
				isFound = true
				indexOfKeymanager = i
				break
			}
		}

		var decodedByte, err = base64.StdEncoding.DecodeString(notification.Event.PayloadData.Value)

		if err != nil {
			if _, ok := err.(base64.CorruptInputError); ok {
				logger.LoggerInternalMsg.ErrorC(logging.ErrorDetails{
					Message:   "\nbase64 input is corrupt, check the provided key",
					Severity:  logging.TRIVIAL,
					ErrorCode: 1241,
				})
			}
			logger.LoggerInternalMsg.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error occurred while decoding the notification event %v", err.Error()),
				Severity:  logging.MAJOR,
				ErrorCode: 1242,
			})
			return
		}

		if strings.EqualFold(keyManagerConfigEvent, notification.Event.PayloadData.EventType) {
			if isFound && strings.EqualFold(actionDelete, notification.Event.PayloadData.Action) {
				logger.LoggerInternalMsg.Infof("Found KM %s to be deleted index %d", notification.Event.PayloadData.Name,
					indexOfKeymanager)
				if isFound {
					xds.KeyManagerList[indexOfKeymanager] = xds.KeyManagerList[len(xds.KeyManagerList)-1]
					xds.KeyManagerList = xds.KeyManagerList[:len(xds.KeyManagerList)-1]
				}
				xds.GenerateAndUpdateKeyManagerList()
			} else if decodedByte != nil {
				logger.LoggerInternalMsg.Infof("decoded stream %s", string(decodedByte))
				kmConfigMapErr := json.Unmarshal([]byte(string(decodedByte)), &kmConfigMap)
				if kmConfigMapErr != nil {
					logger.LoggerInternalMsg.ErrorC(logging.ErrorDetails{
						Message:   fmt.Sprintf("Error occurred while unmarshalling key manager config map %v", kmConfigMapErr.Error()),
						Severity:  logging.MINOR,
						ErrorCode: 1243,
					})
					return
				}

				if strings.EqualFold(actionAdd, notification.Event.PayloadData.Action) ||
					strings.EqualFold(actionUpdate, notification.Event.PayloadData.Action) {
					keyManager := eventhubTypes.KeyManager{Name: notification.Event.PayloadData.Name,
						Type: notification.Event.PayloadData.Type, Enabled: notification.Event.PayloadData.Enabled,
						TenantDomain: notification.Event.PayloadData.TenantDomain, Configuration: kmConfigMap}
					logger.LoggerInternalMsg.Infof("data %v", keyManager.Configuration)

					if isFound {
						xds.KeyManagerList[indexOfKeymanager] = keyManager
					} else {
						xds.KeyManagerList = append(xds.KeyManagerList, keyManager)
					}
					xds.GenerateAndUpdateKeyManagerList()
				}
			}
		}
		d.Ack(false)
	}
	logger.LoggerInternalMsg.Info("handle: deliveries channel closed")
}
