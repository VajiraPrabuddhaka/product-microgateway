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

package messaging

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wso2/product-microgateway/adapter/config"
	"github.com/wso2/product-microgateway/adapter/internal/eventhub"
	logger "github.com/wso2/product-microgateway/adapter/internal/loggers"
	"github.com/wso2/product-microgateway/adapter/internal/synchronizer"
	"github.com/wso2/product-microgateway/adapter/pkg/logging"
	msg "github.com/wso2/product-microgateway/adapter/pkg/messaging"
)

func handleAzureOrganizationPurge() {
	for d := range msg.AzureOrganizationPurgeChannel {
		logger.LoggerInternalMsg.Info("message received for OrganizationPurgeChannel = " + string(d))
		var event msg.EventOrganizationPurge
		error := parseOrganizationPurgeJSONEvent(d, &event)

		if error != nil {
			logger.LoggerInternalMsg.ErrorC(logging.ErrorDetails{
				Message:   fmt.Sprintf("Error while processing the organization purge event %v. Hence dropping the event", error.Error()),
				Severity:  logging.MAJOR,
				ErrorCode: 1251,
			})
			continue
		}

		conf, errReadConfig := config.ReadConfigs()

		if errReadConfig != nil {
			log.Fatal("Error loading configuration. ", errReadConfig)
		}

		eventhub.LoadSubscriptionData(conf, nil)

		// clear existing Key Manager Data
		synchronizer.ClearKeyManagerData()
		// Pull Key Manager Data from APIM
		synchronizer.FetchKeyManagersOnStartUp(conf)

	}
}

func parseOrganizationPurgeJSONEvent(data []byte, event *msg.EventOrganizationPurge) error {
	unmarshalErr := json.Unmarshal(data, &event)
	if unmarshalErr != nil {
		logger.LoggerInternalMsg.ErrorC(logging.ErrorDetails{
			Message:   fmt.Sprintf("Error occurred while unmarshalling organization purge event data %v", unmarshalErr.Error()),
			Severity:  logging.MINOR,
			ErrorCode: 1252,
		})
	}
	return unmarshalErr
}
