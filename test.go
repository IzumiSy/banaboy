package main

import (
	config "pubsubroller/config"
	"pubsubroller/subscription"
	"pubsubroller/topic"
	"time"
)

// mock structures for testing

type fakeSubscriptionCallbacks struct {
	IsInitialized bool
	IsFinazlied   bool
	Calls         int
}

func (f *fakeSubscriptionCallbacks) Initialized() {
	f.IsInitialized = true
}

func (f *fakeSubscriptionCallbacks) Each(_ subscription.Subscription) {
	f.Calls++
}

func (f *fakeSubscriptionCallbacks) Finalized(counter countable) {
	f.IsFinazlied = true
}

type fakeTopicCallbacks struct {
	IsInitialized bool
	IsFinazlied   bool
	Calls         int
}

func (f *fakeTopicCallbacks) Initialized() {
	f.IsInitialized = true
}

func (f *fakeTopicCallbacks) Each(_ topic.Topic) {
	f.Calls++
}

func (f *fakeTopicCallbacks) Finalized(counter countable) {
	f.IsFinazlied = true
}

// all operations have 1s sleep in order to test whether goroutine is working

type fakeClient struct{}

func (_ fakeClient) CreateSubscription(_ subscription.Subscription) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (_ fakeClient) DeleteSubscription(_ subscription.Subscription) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (_ fakeClient) CreateTopic(_ topic.Topic) error {
	time.Sleep(1 * time.Second)
	return nil
}

func (_ fakeClient) DeleteTopic(_ topic.Topic) error {
	time.Sleep(1 * time.Second)
	return nil
}

var (
	mockConfig = config.Configuration{
		Internal_Topics_: map[string]config.Topic{
			"topic1": {
				Internal_Subscriptions_: []config.Subscription{
					{Name: "subscription11"},
					{Name: "subscription12"},
					{Name: "subscription13"},
				},
			},
			"topic2": {
				Internal_Subscriptions_: []config.Subscription{
					{Name: "subscription21"},
					{Name: "subscription22"},
					{Name: "subscription23"},
				},
			},
			"topic3": {
				Internal_Subscriptions_: []config.Subscription{
					{Name: "subscription31"},
					{Name: "subscription32"},
					{Name: "subscription33"},
				},
			},
		},
	}
)
