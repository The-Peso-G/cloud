package service

import (
	"fmt"
	"net/http"

	"github.com/go-ocf/sdk/schema"

	"github.com/go-ocf/cloud/grpc-gateway/client"
	"github.com/go-ocf/cloud/grpc-gateway/pb"
	"github.com/go-ocf/cloud/http-gateway/uri"
	"github.com/gorilla/mux"
)

func (requestHandler *RequestHandler) getDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx, cancel := requestHandler.makeCtx(r)
	defer cancel()

	sdkDevice, err := requestHandler.client.GetDevice(ctx, vars[uri.DeviceIDKey])
	if err != nil {
		writeError(w, fmt.Errorf("cannot get device: %w", err))
		return
	}

	jsonResponseWriter(w, mapToDevice(sdkDevice))
}

type Status string

const Status_ONLINE Status = "online"
const Status_OFFLINE Status = "offline"

func toStatus(isOnline bool) Status {
	if isOnline {
		return "online"
	}
	return "offline"
}

type Device struct {
	Device schema.Device `json:"device"`
	Status Status        `json:"status"`
}

type RetrieveDeviceWithLinksResponse struct {
	Device
	Links []schema.ResourceLink `json:"links"`
}

func toLocalizedString(s *pb.LocalizedString) schema.LocalizedString {
	return schema.LocalizedString{
		Value:    s.Value,
		Language: s.Language,
	}
}

func toLocalizedStrings(s []*pb.LocalizedString) []schema.LocalizedString {
	r := make([]schema.LocalizedString, 0, 16)
	for _, v := range s {
		r = append(r, toLocalizedString(v))
	}
	return r
}

func toResourceLink(s pb.ResourceLink) schema.ResourceLink {
	return schema.ResourceLink{
		ResourceTypes: s.GetTypes(),
		Interfaces:    s.GetInterfaces(),
		Href:          s.GetHref(),
		DeviceID:      s.GetDeviceId(),
	}
}

func toResourceLinks(s []pb.ResourceLink) []schema.ResourceLink {
	r := make([]schema.ResourceLink, 0, 16)
	for _, v := range s {
		r = append(r, toResourceLink(v))
	}
	return r
}

func mapToDevice(d client.DeviceDetails) RetrieveDeviceWithLinksResponse {
	return RetrieveDeviceWithLinksResponse{
		Device: Device{
			Device: schema.Device{
				ResourceTypes:    d.Device.GetTypes(),
				ID:               d.ID,
				ManufacturerName: toLocalizedStrings(d.Device.GetManufacturerName()),
				ModelNumber:      d.Device.GetModelNumber(),
				Name:             d.Device.GetName(),
			},
			Status: toStatus(d.Device.GetIsOnline()),
		},
		Links: toResourceLinks(d.Resources),
	}
}
