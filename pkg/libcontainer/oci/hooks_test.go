package oci

import (
	"testing"

	"github.com/opencontainers/runtime-spec/specs-go"
)

func TestHooksFromSpec(t *testing.T) {
	tests := []struct {
		name     string
		spec     *specs.Hooks
		hook     LifecycleHook
		expected []specs.Hook
	}{
		{
			name:     "Nil spec",
			spec:     nil,
			hook:     HookCreateRuntime,
			expected: []specs.Hook{},
		},
		{
			name: "Nil CreateRuntime hooks",
			spec: &specs.Hooks{
				CreateRuntime: nil,
			},
			hook:     HookCreateRuntime,
			expected: []specs.Hook{},
		},
		{
			name: "Non-nil CreateRuntime hooks",
			spec: &specs.Hooks{
				CreateRuntime: []specs.Hook{
					{
						Path: "path/to/hook1",
					},
					{
						Path: "path/to/hook2",
					},
				},
			},
			hook: HookCreateRuntime,
			expected: []specs.Hook{
				{
					Path: "path/to/hook1",
				},
				{
					Path: "path/to/hook2",
				},
			},
		},
		{
			name: "Non-nil CreateContainer hooks",
			spec: &specs.Hooks{
				CreateContainer: []specs.Hook{
					{
						Path: "path/to/hook1",
					},
				},
			},
			hook: HookCreateContainer,
			expected: []specs.Hook{
				{
					Path: "path/to/hook1",
				},
			},
		},
		{
			name: "Nil StartContainer hooks",
			spec: &specs.Hooks{
				StartContainer: nil,
			},
			hook:     HookStartContainer,
			expected: []specs.Hook{},
		},
		{
			name: "Non-nil PostStart hooks",
			spec: &specs.Hooks{
				Poststart: []specs.Hook{
					{
						Path: "path/to/hook1",
					},
				},
			},
			hook: HookPostStart,
			expected: []specs.Hook{
				{
					Path: "path/to/hook1",
				},
			},
		},
		{
			name: "Non-nil PostStop hooks",
			spec: &specs.Hooks{
				Poststop: []specs.Hook{
					{
						Path: "path/to/hook1",
					},
				},
			},
			hook: HookPostStop,
			expected: []specs.Hook{
				{
					Path: "path/to/hook1",
				},
			},
		},
		{
			name: "Unknown lifecycle hook",
			spec: &specs.Hooks{
				CreateRuntime: []specs.Hook{
					{
						Path: "path/to/hook1",
					},
				},
			},
			hook:     LifecycleHook(99), // Invalid lifecycle hook
			expected: []specs.Hook{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HooksFromSpec(tt.spec, tt.hook)
			if len(got) != len(tt.expected) {
				t.Errorf("expected %d hooks, got %d", len(tt.expected), len(got))
			}
			for i, hook := range got {
				if hook.Path != tt.expected[i].Path {
					t.Errorf("expected hook path %s, got %s", tt.expected[i].Path, hook.Path)
				}
			}
		})
	}
}
