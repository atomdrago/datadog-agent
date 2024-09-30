// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

// Package ddflareextensionimpl defines the OpenTelemetry Extension implementation.
package ddflareextensionimpl

import (
	"context"
	"fmt"

	"github.com/DataDog/datadog-agent/comp/otelcol/ddflareextension/impl/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
)

const (
	defaultHTTPPort = 7777
)

type ddExtensionFactory struct {
	extension.Factory

	factories              *otelcol.Factories
	configProviderSettings otelcol.ConfigProviderSettings
}

// NewFactory creates a factory for HealthCheck extension.
func NewFactory(factories *otelcol.Factories, configProviderSettings otelcol.ConfigProviderSettings) extension.Factory {
	return &ddExtensionFactory{
		factories:              factories,
		configProviderSettings: configProviderSettings,
	}
}

func (f *ddExtensionFactory) CreateExtension(ctx context.Context, set extension.Settings, cfg component.Config) (extension.Extension, error) {
	config := &Config{
		factories:              f.factories,
		configProviderSettings: f.configProviderSettings,
	}
	config.HTTPConfig = cfg.(*Config).HTTPConfig
	return NewExtension(ctx, config, set.TelemetrySettings, set.BuildInfo)
}

func (f *ddExtensionFactory) CreateDefaultConfig() component.Config {
	return &Config{
		HTTPConfig: &confighttp.ServerConfig{
			Endpoint: fmt.Sprintf("localhost:%d", defaultHTTPPort),
		},
	}
}

func (f *ddExtensionFactory) Type() component.Type {
	return metadata.Type
}

func (f *ddExtensionFactory) ExtensionStability() component.StabilityLevel {
	return metadata.ExtensionStability
}
