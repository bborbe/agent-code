// Copyright (c) 2026 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package factory wires concrete dependencies for the agent-code binary.
//
// Pure-code agent — no Claude/Gemini/LLM dependencies, just deliverers.
package factory

import (
	"context"

	agentlib "github.com/bborbe/agent"
	delivery "github.com/bborbe/agent/delivery"
	healthcheck "github.com/bborbe/agent/healthcheck"
	"github.com/bborbe/cqrs/base"
	"github.com/bborbe/errors"
	libkafka "github.com/bborbe/kafka"
	libtime "github.com/bborbe/time"
	"github.com/bborbe/vault-cli/pkg/domain"

	"github.com/bborbe/agent-code/pkg/steps"
)

// CreateKafkaResultDeliverer creates a ResultDeliverer that publishes task
// updates to Kafka via CQRS commands. Uses the passthrough content generator
// — the agent framework's StepRunner already produces the full marshaled
// task in result.Output; the deliverer publishes it as-is and overrides
// status/phase frontmatter based on the result Status.
func CreateKafkaResultDeliverer(
	syncProducer libkafka.SyncProducer,
	topicPrefix base.TopicPrefix,
	taskID agentlib.TaskIdentifier,
	originalContent string,
	currentDateTime libtime.CurrentDateTimeGetter,
) agentlib.ResultDeliverer {
	return delivery.NewKafkaResultDeliverer(
		syncProducer,
		topicPrefix,
		taskID,
		originalContent,
		delivery.NewPassthroughContentGenerator(),
		currentDateTime,
	)
}

// CreateFileResultDeliverer creates a ResultDeliverer that writes the agent's
// output back to a markdown file (local CLI mode).
func CreateFileResultDeliverer(filePath string) agentlib.ResultDeliverer {
	return delivery.NewFileResultDeliverer(
		delivery.NewPassthroughContentGenerator(),
		filePath,
	)
}

// CreateAgent assembles the 3-phase pure-code agent — no LLM deps.
// PlanStep reads frontmatter, ExecuteStep computes, VerifyStep checks.
func CreateAgent() *agentlib.Agent {
	return agentlib.NewAgent(
		agentlib.NewPhase("planning", steps.NewPlanStep()),
		agentlib.NewPhase(domain.TaskPhaseExecution, steps.NewExecuteStep()),
		agentlib.NewPhase("ai_review", steps.NewVerifyStep()),
	)
}

// SyncProducerProvider defers Kafka sync-producer construction to call time.
//
// The construction can fail, but a Create* factory must not return an error
// (rule go-factory/no-error-return) and Kafka wiring must not live in main.go
// (rules go-factory/factory-moved, go-factory/main-holds-only-boot-lifecycle-config).
// A provider satisfies all three: the factory returns it without error, and the
// caller receives the error from Get at the point it can act on it.
type SyncProducerProvider interface {
	Get(ctx context.Context) (libkafka.SyncProducer, error)
}

type syncProducerProvider struct {
	brokers   libkafka.Brokers
	agentName string
}

func (s *syncProducerProvider) Get(ctx context.Context) (libkafka.SyncProducer, error) {
	syncProducer, err := libkafka.NewSyncProducerWithName(ctx, s.brokers, s.agentName)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create sync producer")
	}
	return syncProducer, nil
}

// CreateSyncProducerProvider returns a provider that builds the Kafka sync
// producer on demand. Pure composition — no I/O until Get is called.
func CreateSyncProducerProvider(brokers libkafka.Brokers, agentName string) SyncProducerProvider {
	return &syncProducerProvider{brokers: brokers, agentName: agentName}
}

// CreateAgentProvider wires the per-task-type dispatch for agent-code.
// Healthcheck-only binary, pure-Go (no LLM): TaskTypeHealthcheck routes to a
// Nop liveness agent — reaching it proves binary booted, envconfig parsed,
// Kafka client opened. Any other task_type hits the default-error branch.
func CreateAgentProvider() agentlib.AgentProvider {
	livenessAgent := healthcheck.NewAgent(healthcheck.NewNopStep())
	return agentlib.NewAgentProvider("agent-code", map[agentlib.TaskType]*agentlib.Agent{
		agentlib.TaskTypeHealthcheck: livenessAgent,
	})
}
