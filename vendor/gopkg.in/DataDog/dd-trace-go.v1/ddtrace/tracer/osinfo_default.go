// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016 Datadog, Inc.

// +build !windows,!linux,!darwin,!freebsd

package tracer

import (
	"runtime"
)

func osName() string {
	return runtime.GOOS
}

func osVersion() string {
	return unknown
}
