package deployer

import "testing"

func TestValidateRootDirInput(t *testing.T) {
	valid := []string{"", ".", "apps/api", "services\\bot"}
	for _, v := range valid {
		if err := ValidateRootDirInput(v); err != nil {
			t.Fatalf("expected %q to be valid, got %v", v, err)
		}
	}

	invalid := []string{"..", "../api", `C:\abs\path`}
	for _, v := range invalid {
		if err := ValidateRootDirInput(v); err == nil {
			t.Fatalf("expected %q to be invalid", v)
		}
	}
}

