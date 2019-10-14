package subscription

import (
	"context"
	"pubsubroller/config"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
)

type Subscription struct {
	topic    *pubsub.Topic
	name     string
	endpoint string
	pull     bool
}

func New(name, endpoint string, pull bool, topic *pubsub.Topic) Subscription {
	return Subscription{
		topic:    topic,
		name:     name,
		endpoint: endpoint,
		pull:     pull,
	}
}

var (
	INTERNAL_ERR              error = errors.New("Internal error")
	SUBSCRIPTION_EXISTS_ERR   error = errors.New("Subscription already exists")
	NO_ENDPOINT_SPECIFIED_ERR error = errors.New("No endpoint specified")
)

func (subscription Subscription) Create(client *pubsub.Client, ctx context.Context) error {
	s := client.Subscription(subscription.name)
	exists, err := s.Exists(ctx)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	if exists {
		return SUBSCRIPTION_EXISTS_ERR
	}

	if subscription.endpoint == "" {
		return errors.WithMessage(NO_ENDPOINT_SPECIFIED_ERR, subscription.name)
	}

	var pushConfig pubsub.PushConfig
	if subscription.pull {
		pushConfig = pubsub.PushConfig{}
	} else {
		pushConfig = pubsub.PushConfig{Endpoint: subscription.endpoint}
	}

	_, err = client.CreateSubscription(
		ctx,
		subscription.name,
		pubsub.SubscriptionConfig{
			Topic:      subscription.topic,
			PushConfig: pushConfig,
		},
	)
	if err != nil {
		return errors.Wrap(err, INTERNAL_ERR.Error())
	}

	return nil
}

func FromConfig(conf config.Configuration, variables map[string]string, client *pubsub.Client) []Subscription {
	var subscriptions []Subscription

	for topicName, topic := range conf.Topics() {
		topicName := topicName
		topic := topic

		for _, sub := range topic.Subscriptions() {
			endpoint := sub.Endpoint
			for key, value := range variables {
				endpoint = strings.Replace(endpoint, "${"+key+"}", value, -1)
			}

			subscriptions =
				append(
					subscriptions,
					New(sub.Name, endpoint, sub.Pull, client.Topic(topicName)),
				)
		}
	}

	return subscriptions
}
