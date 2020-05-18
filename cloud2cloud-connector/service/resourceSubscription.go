package service

import (
	"context"
	"fmt"

	"github.com/go-ocf/go-coap"

	"github.com/go-ocf/cloud/cloud2cloud-connector/events"
	"github.com/go-ocf/cloud/cloud2cloud-connector/store"
	raCqrs "github.com/go-ocf/cloud/resource-aggregate/cqrs"
	pbCQRS "github.com/go-ocf/cloud/resource-aggregate/pb"
	pbRA "github.com/go-ocf/cloud/resource-aggregate/pb"
	kitNetGrpc "github.com/go-ocf/kit/net/grpc"
	kitHttp "github.com/go-ocf/kit/net/http"
)

func (s *SubscribeManager) subscribeToResource(ctx context.Context, l store.LinkedAccount, correlationID, signingSecret, deviceID, resourceHrefLink string) (string, error) {
	resp, err := subscribe(ctx, "/devices/"+deviceID+"/"+resourceHrefLink+"/subscriptions", correlationID, events.SubscriptionRequest{
		URL:           s.eventsURL,
		EventTypes:    []events.EventType{events.EventType_ResourceChanged},
		SigningSecret: signingSecret,
	}, l)
	if err != nil {
		return "", fmt.Errorf("cannot subscribe to device %v for %v: %v", deviceID, l.ID, err)
	}
	return resp.SubscriptionId, nil
}

func cancelResourceSubscription(ctx context.Context, l store.LinkedAccount, deviceID, resourceID, subscriptionID string) error {
	err := cancelSubscription(ctx, "/devices/"+deviceID+"/"+resourceID+"/subscriptions/"+subscriptionID, l)
	if err != nil {
		return fmt.Errorf("cannot cancel resource subscription for %v: %v", l.ID, err)
	}
	return nil
}

func notifyResourceChanged(ctx context.Context, raClient pbRA.ResourceAggregateClient, deviceID, href, userID string, contentType string, body []byte, cmdMeta pbCQRS.CommandMetadata) error {
	coapContentFormat := int32(-1)
	switch contentType {
	case coap.AppCBOR.String():
		coapContentFormat = int32(coap.AppCBOR)
	case coap.AppOcfCbor.String():
		coapContentFormat = int32(coap.AppOcfCbor)
	case coap.AppJSON.String():
		coapContentFormat = int32(coap.AppJSON)
	}

	_, err := raClient.NotifyResourceChanged(ctx, &pbRA.NotifyResourceChangedRequest{
		AuthorizationContext: &pbCQRS.AuthorizationContext{
			UserId:   userID,
			DeviceId: deviceID,
		},
		ResourceId:      raCqrs.MakeResourceId(deviceID, kitHttp.CanonicalHref(href)),
		CommandMetadata: &cmdMeta,
		Content: &pbRA.Content{
			Data:              body,
			ContentType:       contentType,
			CoapContentFormat: coapContentFormat,
		},
	})
	return err
}

func (s *SubscribeManager) HandleResourceChangedEvent(ctx context.Context, subscriptionData subscriptionData, header events.EventHeader, body []byte) error {
	userID, err := subscriptionData.linkedAccount.OriginCloud.AccessToken.GetSubject()
	if err != nil {
		return fmt.Errorf("cannot get userID for device (%v) resource (%v) content changed: %v", subscriptionData.subscription.DeviceID, subscriptionData.subscription.Href, err)
	}
	err = notifyResourceChanged(
		kitNetGrpc.CtxWithToken(ctx, subscriptionData.linkedAccount.OriginCloud.AccessToken.String()),
		s.raClient,
		subscriptionData.subscription.DeviceID,
		subscriptionData.subscription.Href,
		userID,
		header.ContentType,
		body,
		makeCommandMetadata(header.SequenceNumber),
	)
	if err != nil {
		return fmt.Errorf("cannot update resource aggregate (%v) resource (%v) content changed: %v", subscriptionData.subscription.DeviceID, subscriptionData.subscription.Href, err)
	}

	return nil

}

func (s *SubscribeManager) HandleResourceEvent(ctx context.Context, header events.EventHeader, body []byte, subscriptionData subscriptionData) error {
	switch header.EventType {
	case events.EventType_ResourceChanged:
		return s.HandleResourceChangedEvent(ctx, subscriptionData, header, body)
	}
	return fmt.Errorf("cannot handle resource event: unsupported Event-Type %v", header.EventType)
}
