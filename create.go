package main

import (
	"context"
	"fmt"
	"pubsubroller/client"
	config "pubsubroller/config"
	subscription "pubsubroller/subscription"
	topic "pubsubroller/topic"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func createSubscriptions(c client.PubsubClient, ctx context.Context, conf config.Configuration, opts Options) {
	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	fmt.Println("Start creating subscriptions...")

	for _, sub := range subscription.FromConfig(conf, opts.Variables, c) {
		sub := sub

		egSubscriptions.Go(func() error {
			if !opts.IsDryRun {
				if err := sub.Create(c, ctx); err != nil {
					if errors.Cause(err) == subscription.SUBSCRIPTION_EXISTS_ERR {
						subscriptionSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			subscriptionCreatedCount += 1
			fmt.Printf("Subscription created: %s\n", sub.Name())
			return nil
		})
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("Subscriptions created: %d, skipped: %d\n", subscriptionCreatedCount, subscriptionSkippedCount)
}

func createTopics(c client.PubsubClient, ctx context.Context, conf config.Configuration, opts Options) {
	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	fmt.Println("Start creating topics...")

	for _, tp := range topic.FromConfig(conf, opts.Variables, c) {
		tp := tp

		egTopics.Go(func() error {
			if !opts.IsDryRun {
				if err := tp.Create(c, ctx); err != nil {
					if errors.Cause(err) == topic.TOPIC_EXISTS_ERR {
						topicSkippedCount += 1
						return nil
					} else {
						return err
					}
				}
			}

			topicCreatedCount += 1
			fmt.Printf("Topic created: %s\n", tp.Name())
			return nil
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("Topics created: %d, skipped: %d\n", topicCreatedCount, topicSkippedCount)
}
