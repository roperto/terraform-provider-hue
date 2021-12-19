package provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	// See: https://developers.meethue.com/develop/hue-api/error-messages/
	HueErrorUnauthorizedUser                 = 1
	HueErrorInvalidJSON                      = 2
	HueErrorResourceNotAvailable             = 3
	HueErrorMethodNotAvailable               = 4
	HueErrorMissingParameters                = 5
	HueErrorParameterNotAvailable            = 6
	HueErrorInvalidValueForParameter         = 7
	HueErrorParameterNotModifiable           = 8
	HueErrorTooManyItemsInList               = 11
	HueErrorPortalConnectionRequired         = 12
	HueErrorInternetlError                   = 901
	HueErrorLinkButtonNotPressed             = 101
	HueErrorCannotDisableDHCP                = 110
	HueErrorInvalidUpdateState               = 111
	HueErrorParameterNotModifiableDeviceOff  = 201
	HueErrorLightListFull                    = 203
	HueErrorGroupTableFull                   = 301
	HueErrorGroupTypeNotAllowed              = 305
	HueErrorLightAlreadyUsed                 = 306
	HueErrorSceneBufferFull                  = 402
	HueErrorSceneLocked                      = 403
	HueErrorSceneGroupEmpty                  = 404
	HueErrorCreateSensorTypeNotAllowed       = 501
	HueErrorSensorListFull                   = 502
	HueErrorSensorListFullZigBee             = 503
	HueErrorRuleEngineFull                   = 601
	HueErrorConditionError                   = 607
	HueErrorActionError                      = 608
	HueErrorUnableToActivate                 = 609
	HueErrorScheduleListFull                 = 701
	HueErrorScheduleTimeZoneInvalid          = 702
	HueErrorScheduleWithBothTimeAndLocalTime = 703
	HueErrorScheduleInvalidTag               = 704
	HueErrorScheduleExpired                  = 705
	HueErrorCommandError                     = 706
	HueErrorSourceModelInvalid               = 801
	HueErrorSourceFactoryNew                 = 802
	HueErrorInvalidState                     = 803
)

const (
	ResourceTypeGroup        = "groups"
	ResourceTypeLight        = "lights"
	ResourceTypeScene        = "scenes"
	ResourceTypeSensor       = "sensors"
	ResourceTypeSchedule     = "schedules"
	ResourceTypeResourcelink = "resourcelinks"
	ResourceTypeRule         = "rules"
)

type HueError struct {
	Type        int
	Address     string
	Description string
}

func (e *HueError) Error() string {
	return fmt.Sprintf("HueError #%d @ %s --> %s", e.Type, e.Address, e.Description)
}

func parseInt(i interface{}) int {
	switch i.(type) {
	case int:
		return i.(int)
	case float32:
		return int(i.(float32))
	case float64:
		return int(i.(float64))
	default:
		return 0
	}
}

type Hue struct {
	Hostname string
	Username string
}

func (hue *Hue) request(
	ctx context.Context,
	method string,
	resource_type string,
	resource_id string,
	data []byte,
) (interface{}, error) {
	if resource_id != "" {
		resource_id = fmt.Sprintf("/%s", resource_id)
	}
	url := fmt.Sprintf(
		"http://%s/api/%s/%s%s",
		hue.Hostname,
		hue.Username,
		resource_type,
		resource_id,
	)

	body := strings.NewReader(string(data))

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	jsonResultBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG REQUEST] %s %s --> %s --> %s", method, url, string(data), string(jsonResultBytes))

	var result interface{}
	err = json.Unmarshal(jsonResultBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (hue *Hue) parseSuccess(response interface{}, parseID bool) (string, error) {
	var lastId string

	if len(response.([]interface{})) == 0 {
		return "", errors.New(fmt.Sprintf("Empty response: %s", response))
	}

	for _, result := range response.([]interface{}) {
		result := result.(map[string]interface{})
		if len(result) != 1 {
			return "", errors.New(fmt.Sprintf("Only success expected but %d keys found.", len(result)))
		}

		for key, value := range result {
			if key != "success" {
				return "", errors.New(fmt.Sprintf("Non-success response: %s", result))
			}
			if parseID {
				lastId = value.(map[string]interface{})["id"].(string)
			}
		}
	}

	return lastId, nil
}

func (hue *Hue) parseError(res interface{}) error {
	resList, ok := res.([]interface{})
	if !ok {
		return nil
	}

	restMap, ok := resList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	error, ok := restMap["error"].(map[string]interface{})
	if !ok {
		return nil
	}

	errorTypeInt := parseInt(error["type"])

	return &HueError{
		Type:        errorTypeInt,
		Address:     error["address"].(string),
		Description: error["description"].(string),
	}
}

func (hue *Hue) Create(
	ctx context.Context,
	resource_type string,
	data map[string]interface{},
) (string, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	res, err := hue.request(ctx, http.MethodPost, resource_type, "", jsonBytes)
	if err != nil {
		return "", err
	}

	return hue.parseSuccess(res, true)
}

func (hue *Hue) List(
	ctx context.Context,
	resource_type string,
) (map[string]map[string]interface{}, error) {
	res, err := hue.request(ctx, http.MethodGet, resource_type, "", nil)
	if err != nil {
		return nil, err
	}

	err = hue.parseError(res)
	if err != nil {
		return nil, err
	}

	resMap := res.(map[string]interface{})
	retMap := map[string]map[string]interface{}{}
	for k, v := range resMap {
		retMap[k] = v.(map[string]interface{})
	}

	return retMap, nil
}

func (hue *Hue) Read(
	ctx context.Context,
	resource_type string,
	resource_id string,
) (map[string]interface{}, error) {
	res, err := hue.request(ctx, http.MethodGet, resource_type, resource_id, nil)
	if err != nil {
		return nil, err
	}

	err = hue.parseError(res)
	if err != nil {
		return nil, err
	}

	return res.(map[string]interface{}), nil
}

func (hue *Hue) Update(
	ctx context.Context,
	resource_type string,
	resource_id string,
	data map[string]interface{},
) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res, err := hue.request(ctx, http.MethodPut, resource_type, resource_id, jsonBytes)
	if err != nil {
		return err
	}

	_, err = hue.parseSuccess(res, false)
	return err
}

func (hue *Hue) Delete(
	ctx context.Context,
	resource_type string,
	resource_id string,
) error {
	res, err := hue.request(ctx, http.MethodDelete, resource_type, resource_id, nil)
	if err != nil {
		return err
	}

	_, err = hue.parseSuccess(res, false)
	return err
}
