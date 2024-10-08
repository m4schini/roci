package oci

import (
	"context"
	"github.com/opencontainers/runtime-spec/specs-go"
	"os/exec"
	"time"
)

type LifecycleHook uint8

const (
	HookCreateRuntime LifecycleHook = iota
	HookCreateContainer
	HookStartContainer
	HookPostStart
	HookPostStop
	// HookPreStart is deprecated
	HookPreStart
)

// RunHook executes a single OCI hook with the specified context.
func RunHook(ctx context.Context, hook specs.Hook) error {
	var cancel context.CancelFunc = func() {}
	if hook.Timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*hook.Timeout)*time.Second)
	}
	defer cancel()

	cmd := exec.CommandContext(ctx, hook.Path, hook.Args...)
	cmd.Env = hook.Env
	return cmd.Run()
}

// RunHooks executes a sequence of hooks in the order they are provided.
func RunHooks(ctx context.Context, hooks []specs.Hook) (err error) {
	for _, hook := range hooks {
		err = RunHook(ctx, hook)
		if err != nil {
			return err
		}
	}

	return nil
}

// InvokeHooks runs the hooks corresponding to the specified lifecycle stage.
func InvokeHooks(hooks *specs.Hooks, hook LifecycleHook) (err error) {
	return RunHooks(context.Background(), HooksFromSpec(hooks, hook))
}

// HooksFromSpec retrieves the hooks corresponding to the specified lifecycle stage from the specification.
func HooksFromSpec(spec *specs.Hooks, hook LifecycleHook) (hooks []specs.Hook) {
	hooks = make([]specs.Hook, 0)
	if spec == nil {
		return hooks
	}

	// Helper function to safely return the hooks if they are not nil.
	must := func(h []specs.Hook) []specs.Hook {
		if h == nil {
			return hooks
		}
		return h
	}

	// Switch case to select the correct hooks based on the lifecycle stage.
	switch hook {
	case HookCreateRuntime:
		return must(spec.CreateRuntime)
	case HookCreateContainer:
		return must(spec.CreateContainer)
	case HookStartContainer:
		return must(spec.StartContainer)
	case HookPostStart:
		return must(spec.Poststart)
	case HookPostStop:
		return must(spec.Poststop)
	default:
		return must(nil)
	}
}
