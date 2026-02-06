package runtime

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// KubectlApply applies manifest YAML using kubectl.
func KubectlApply(ctx context.Context, manifest []byte) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(manifest)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl apply failed: %w: %s", err, string(out))
	}
	return nil
}
