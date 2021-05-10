package auth

import (
	"testing"
)

func TestHasPermissions(t *testing.T) {

	var permissionsTest = []struct {
		principal   *SecurityPrincipal // input
		permissions string             // arguments
		expected    bool               // expected result
	}{
		{&SecurityPrincipal{AuthProvider: ""}, "", true},
		{&SecurityPrincipal{AuthProvider: ""}, "*", false},
	}

	for _, tt := range permissionsTest {
		actual := tt.principal.HasPermissions(tt.permissions)
		if actual != tt.expected {
			t.Errorf("principal.HasPermissions(%s): expected %v, actual %v", tt.permissions,
				tt.expected, actual)
		}
	}
}
