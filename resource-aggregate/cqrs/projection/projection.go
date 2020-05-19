package projection

import (
	"context"
	"fmt"

	"github.com/go-ocf/cqrs"
	"github.com/go-ocf/cqrs/eventbus"
	"github.com/go-ocf/cqrs/eventstore"
	"github.com/go-ocf/kit/log"

	raCqrsUtils "github.com/go-ocf/cloud/resource-aggregate/cqrs"
)

// Projection projects events from resource aggregate.
type Projection struct {
	projection *projection
}

// NewProjection creates new resource projection.
func NewProjection(ctx context.Context, name string, store eventstore.EventStore, subscriber eventbus.Subscriber, factoryModel eventstore.FactoryModelFunc) (*Projection, error) {
	projection, err := newProjection(ctx, name, store, subscriber, factoryModel, raCqrsUtils.GetTopics)
	if err != nil {
		return nil, fmt.Errorf("cannot create resource projection: %w", err)
	}
	return &Projection{projection: projection}, nil
}

// Register registers deviceId, loads events from eventstore and subscribe to eventbus.
// It can be called multiple times for same deviceId but after successful the a call Unregister
// must be called same times to free resources.
func (p *Projection) Register(ctx context.Context, deviceId string) (loaded bool, err error) {
	return p.projection.register(ctx, deviceId, []eventstore.SnapshotQuery{{GroupId: deviceId}})
}

// Unregister unregisters device and his resource from projection.
func (p *Projection) Unregister(deviceId string) error {
	return p.projection.unregister(deviceId)
}

// Models returns models for device, resource or nil for non exist.
func (p *Projection) Models(deviceId, resourceId string) []eventstore.Model {
	return p.projection.models([]eventstore.SnapshotQuery{{GroupId: deviceId, AggregateId: resourceId}})
}

// ForceUpdate invokes update registered resource model from evenstore.
func (p *Projection) ForceUpdate(ctx context.Context, deviceId, resourceId string) error {
	err := p.projection.forceUpdate(ctx, deviceId, []eventstore.SnapshotQuery{{GroupId: deviceId, AggregateId: resourceId}})
	if err != nil {
		return fmt.Errorf("cannot force update resource projection: %w", err)
	}
	return err
}

type projection struct {
	cqrsProjection *cqrs.Projection

	topicManager *TopicManager
	refCountMap  *RefCountMap
}

func newProjection(ctx context.Context, name string, store eventstore.EventStore, subscriber eventbus.Subscriber, factoryModel eventstore.FactoryModelFunc, getTopics GetTopicsFunc) (*projection, error) {
	cqrsProjection, err := cqrs.NewProjection(ctx, store, name, subscriber, factoryModel, func(string, ...interface{}) {})
	if err != nil {
		return nil, fmt.Errorf("cannot create Projection: %w", err)
	}
	return &projection{
		cqrsProjection: cqrsProjection,
		topicManager:   NewTopicManager(getTopics),
		refCountMap:    NewRefCountMap(),
	}, nil
}

// ForceUpdate invokes update registered resource model from evenstore.
func (p *projection) forceUpdate(ctx context.Context, registrationID string, query []eventstore.SnapshotQuery) error {
	_, err := p.refCountMap.Inc(registrationID, false)
	if err != nil {
		return fmt.Errorf("cannot force update projection: %w", err)
	}

	err = p.cqrsProjection.Project(ctx, query)
	if err != nil {
		return fmt.Errorf("cannot force update projection: %w", err)
	}
	_, err = p.refCountMap.Dec(registrationID)
	if err != nil {
		return fmt.Errorf("cannot force update projection: %w", err)
	}
	return nil
}

func (p *projection) models(query []eventstore.SnapshotQuery) []eventstore.Model {
	return p.cqrsProjection.Models(query)
}

func (p *projection) register(ctx context.Context, registrationID string, query []eventstore.SnapshotQuery) (loaded bool, err error) {
	created, err := p.refCountMap.Inc(registrationID, true)
	if err != nil {
		return false, fmt.Errorf("cannot register device: %w", err)
	}
	if !created {
		return false, nil
	}

	topics, updateSubscriber := p.topicManager.Add(registrationID)

	if updateSubscriber {
		err := p.cqrsProjection.SubscribeTo(topics)
		if err != nil {
			p.refCountMap.Dec(registrationID)
			return false, fmt.Errorf("cannot register device: %w", err)
		}
	}

	err = p.cqrsProjection.Project(ctx, query)
	if err != nil {
		return false, fmt.Errorf("cannot register device: %w", err)
	}

	return true, nil
}

func (p *projection) unregister(registrationID string) error {
	deleted, err := p.refCountMap.Dec(registrationID)
	if err != nil {
		return fmt.Errorf("cannot unregister device from projection: %w", err)
	}
	if !deleted {
		return nil
	}

	topics, updateSubscriber := p.topicManager.Remove(registrationID)

	if updateSubscriber {
		err := p.cqrsProjection.SubscribeTo(topics)
		if err != nil {
			log.Errorf("cannot change topics for projection: %w", err)
		}
	}
	return p.cqrsProjection.Forget([]eventstore.SnapshotQuery{
		{GroupId: registrationID},
	})
}
