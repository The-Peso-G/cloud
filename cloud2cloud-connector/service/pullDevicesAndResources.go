package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/go-ocf/kit/log"

	"github.com/gofrs/uuid"

	pbAS "github.com/go-ocf/cloud/authorization/pb"
	"github.com/go-ocf/cloud/cloud2cloud-connector/store"
	pbCQRS "github.com/go-ocf/cloud/resource-aggregate/pb"
	pbRA "github.com/go-ocf/cloud/resource-aggregate/pb"
	"github.com/go-ocf/kit/codec/json"
	kitNetGrpc "github.com/go-ocf/kit/net/grpc"
	"github.com/go-ocf/sdk/schema"
)

type Device struct {
	Device schema.Device `json:"device"`
	Status string        `json:"status"`
}

type RetrieveDeviceWithLinksResponse struct {
	Device
	Links []schema.ResourceLink `json:"links"`
}

type pullDevicesHandler struct {
	s        store.Store
	asClient pbAS.AuthorizationServiceClient
	raClient pbRA.ResourceAggregateClient
}

func getUsersDevices(ctx context.Context, asClient pbAS.AuthorizationServiceClient) (map[string]bool, error) {
	getUserDevicesClient, err := asClient.GetUserDevices(ctx, &pbAS.GetUserDevicesRequest{})
	if err != nil {
		return nil, fmt.Errorf("cannot get users devices: %w", err)
	}
	defer getUserDevicesClient.CloseSend()
	userDevices := make(map[string]bool)
	for {
		userDevice, err := getUserDevicesClient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("cannot get users devices: %w", err)
		}
		if userDevice == nil {
			continue
		}

		userDevices[userDevice.DeviceId] = true
	}
	return userDevices, nil
}

func (p *pullDevicesHandler) getDevicesWithResourceLinks(ctx context.Context, account store.LinkedAccount) error {
	var errors []error
	connectionIDRand, err := uuid.NewV4()
	if err != nil {
		return err
	}
	connectionID := "c2c-connector-pull:/devices:" + connectionIDRand.String()

	client := NewHTTPClientWihoutVerifyServer()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, account.TargetURL+"/devices", nil)
	req.Header.Set("Authorization", "Bearer "+string(account.TargetCloud.AccessToken))
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var devices []RetrieveDeviceWithLinksResponse
	err = json.ReadFrom(resp.Body, devices)
	if err != nil {
		return err
	}
	registeredDevices, err := getUsersDevices(kitNetGrpc.CtxWithToken(ctx, account.OriginCloud.AccessToken.String()), p.asClient)
	if err != nil {
		return err
	}
	userID, err := account.OriginCloud.AccessToken.GetSubject()
	if err != nil {
		return fmt.Errorf("cannot get userID: %v", err)
	}

	for _, dev := range devices {
		deviceID := dev.Device.Device.ID
		ok := registeredDevices[deviceID]
		if !ok {
			_, err := p.asClient.AddDevice(kitNetGrpc.CtxWithToken(ctx, account.OriginCloud.AccessToken.String()), &pbAS.AddDeviceRequest{
				DeviceId: deviceID,
				UserId:   userID,
			})
			if err != nil {
				errors = append(errors, err)
				continue
			}

			err = publishCloudDeviceStatus(kitNetGrpc.CtxWithToken(ctx, account.OriginCloud.AccessToken.String()), p.raClient, userID, deviceID, pbCQRS.CommandMetadata{
				ConnectionId: connectionID,
			})
			if err != nil {
				errors = append(errors, err)
				continue
			}
		}
		delete(registeredDevices, deviceID)
		var online bool
		if strings.ToLower(dev.Status) == "online" {
			online = true
		}
		err = updateCloudStatus(kitNetGrpc.CtxWithToken(ctx, account.OriginCloud.AccessToken.String()), p.raClient, userID, deviceID, online, pbCQRS.CommandMetadata{
			ConnectionId: connectionID,
		})
		if err != nil {
			errors = append(errors, err)
			continue
		}
		for _, link := range dev.Links {
			err := publishResource(kitNetGrpc.CtxWithToken(ctx, account.OriginCloud.AccessToken.String()), p.raClient, userID, link, pbCQRS.CommandMetadata{
				ConnectionId: connectionID,
			})
			if err != nil {
				errors = append(errors, err)
				continue
			}
		}
	}
	for deviceID := range registeredDevices {
		_, err := p.asClient.RemoveDevice(kitNetGrpc.CtxWithToken(ctx, account.OriginCloud.AccessToken.String()), &pbAS.RemoveDeviceRequest{
			DeviceId: deviceID,
		})
		if err != nil {
			errors = append(errors, err)
			continue
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%+v", errors)
	}
	return nil
}

type Representation struct {
	Href           string      `json:"href"`
	Representation interface{} `json:"rep"`
}

type RetrieveDeviceContentAllResponse struct {
	Device
	Links []Representation `json:"links"`
}

func (p *pullDevicesHandler) getDevicesWithResourceValues(ctx context.Context, account store.LinkedAccount) error {
	var errors []error
	connectionIDRand, err := uuid.NewV4()
	if err != nil {
		return err
	}
	connectionID := "c2c-connector-pull:/devices?content=all:" + connectionIDRand.String()
	client := NewHTTPClientWihoutVerifyServer()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, account.TargetURL+"/devices?content=all", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+string(account.TargetCloud.AccessToken))
	req.Header.Set("Accept", "application/json")
	userID, err := account.OriginCloud.AccessToken.GetSubject()
	if err != nil {
		return fmt.Errorf("cannot get userID: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var devices []RetrieveDeviceContentAllResponse
	err = json.ReadFrom(resp.Body, devices)
	if err != nil {
		return err
	}
	for _, dev := range devices {
		deviceID := dev.Device.Device.ID
		for _, link := range dev.Links {
			body, err := json.Encode(link.Representation)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			err = notifyResourceChanged(
				kitNetGrpc.CtxWithToken(ctx, account.OriginCloud.AccessToken.String()),
				p.raClient,
				deviceID,
				link.Href,
				userID,
				"application/json",
				body,
				pbCQRS.CommandMetadata{
					ConnectionId: connectionID,
				},
			)
			if err != nil {
				errors = append(errors, err)
				continue
			}
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("%+v", errors)
	}

	return nil
}

func (p *pullDevicesHandler) pullDevicesFromAccount(ctx context.Context, account store.LinkedAccount) error {
	account, err := account.RefreshTokens(ctx, p.s, NewHTTPClientWihoutVerifyServer())
	if err != nil {
		return err
	}
	err = p.getDevicesWithResourceLinks(ctx, account)
	if err != nil {
		return err
	}
	err = p.getDevicesWithResourceValues(ctx, account)
	if err != nil {
		return err
	}
	return nil
}

func (p *pullDevicesHandler) Handle(ctx context.Context, iter store.LinkedAccountIter) error {
	var wg sync.WaitGroup
	for {
		var s store.LinkedAccount
		if !iter.Next(ctx, &s) {
			break
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := p.pullDevicesFromAccount(ctx, s)
			if err != nil {
				log.Errorf("cannot pull devices for linked account(%v): %v", s, err)
			}
		}()
	}
	wg.Wait()
	return iter.Err()
}

func pullDevices(ctx context.Context, s store.Store,
	asClient pbAS.AuthorizationServiceClient,
	raClient pbRA.ResourceAggregateClient) error {
	h := pullDevicesHandler{
		s: s,
	}
	return s.LoadLinkedAccounts(ctx, store.Query{}, &h)
}
