// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package factory_test

import (
	"context"

	"github.com/IBM/sarama"
	agentlib "github.com/bborbe/agent"
	"github.com/bborbe/cqrs/base"
	libkafka "github.com/bborbe/kafka"
	kafkamocks "github.com/bborbe/kafka/mocks"
	libtime "github.com/bborbe/time"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bborbe/agent-code/pkg/factory"
)

// This exercises the REAL delivery.NewKafkaResultDeliverer constructor (via
// factory.CreateKafkaResultDeliverer, not a fake sender) with a fake
// libkafka.SyncProducer, so a fat-fingered topicPrefix would be caught: the
// resulting Kafka message topic must match the exact prefix wiring.
var _ = Describe("CreateKafkaResultDeliverer (topic prefix wiring)", func() {
	var (
		ctx             context.Context
		syncProducer    *kafkamocks.KafkaSyncProducer
		taskID          agentlib.TaskIdentifier
		originalContent string
	)

	BeforeEach(func() {
		ctx = context.Background()
		syncProducer = &kafkamocks.KafkaSyncProducer{}
		syncProducer.SendMessageReturns(int32(0), int64(123), nil)
		taskID = agentlib.TaskIdentifier("task-abc-123")
		originalContent = "---\ntitle: My Task\nstatus: in_progress\n---\n\nBody.\n"
	})

	deliver := func(topicPrefix base.TopicPrefix) *sarama.ProducerMessage {
		deliverer := factory.CreateKafkaResultDeliverer(
			syncProducer,
			topicPrefix,
			taskID,
			originalContent,
			libtime.NewCurrentDateTime(),
		)
		err := deliverer.DeliverResult(ctx, agentlib.AgentResultInfo{
			Status: agentlib.AgentStatusDone,
			Output: "---\ntitle: My Task\n---\n\nBody.\n",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(syncProducer.SendMessageCallCount()).To(Equal(1))
		_, msg := syncProducer.SendMessageArgsForCall(0)
		return msg
	}

	It("publishes to the develop-prefixed request topic for TopicPrefix(\"develop\")", func() {
		msg := deliver(base.TopicPrefix("develop"))
		Expect(msg.Topic).To(Equal(libkafka.Topic("develop-agent-task-v1-request").String()))
	})

	It("publishes to the master-prefixed request topic for TopicPrefix(\"master\")", func() {
		msg := deliver(base.TopicPrefix("master"))
		Expect(msg.Topic).To(Equal(libkafka.Topic("master-agent-task-v1-request").String()))
	})

	It("publishes to the unprefixed request topic for TopicPrefix(\"\")", func() {
		msg := deliver(base.TopicPrefix(""))
		Expect(msg.Topic).To(Equal(libkafka.Topic("agent-task-v1-request").String()))
	})
})
