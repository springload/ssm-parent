package ssm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandParameters(t *testing.T) {
	tests := []struct {
		name           string
		parameters     map[string]string
		expandAll      bool
		expandValues   []string
		envVars        map[string]string
		expectedParams map[string]string
		expectError    bool
	}{
		{
			name: "expand specific variable",
			parameters: map[string]string{
				"DATABASE_NAME": "DB_$ENVIRONMENT",
				"SOME_SECRET":   "abc$abc",
			},
			expandAll:    false,
			expandValues: []string{"DATABASE_NAME"},
			envVars: map[string]string{
				"ENVIRONMENT": "teststaging",
			},
			expectedParams: map[string]string{
				"DATABASE_NAME": "DB_teststaging",
				"SOME_SECRET":   "abc$abc",
			},
			expectError: false,
		},
		{
			name: "expand all variables",
			parameters: map[string]string{
				"DATABASE_NAME": "DB_$ENVIRONMENT",
				"SOME_SECRET":   "abc$abc",
			},
			expandAll:    true,
			expandValues: []string{},
			envVars: map[string]string{
				"ENVIRONMENT": "teststaging",
				"abc":         "def",
			},
			expectedParams: map[string]string{
				"DATABASE_NAME": "DB_teststaging",
				"SOME_SECRET":   "abcdef",
			},
			expectError: false,
		},
		{
			name: "missing variable in selective expansion",
			parameters: map[string]string{
				"SOME_SECRET": "abc$abc",
			},
			expandAll:    false,
			expandValues: []string{"NONEXISTENT"},
			envVars:      map[string]string{},
			expectedParams: map[string]string{
				"SOME_SECRET": "abc$abc",
			},
			expectError: true,
		},
		{
			name: "empty parameters",
			parameters: map[string]string{},
			expandAll:    false,
			expandValues: []string{"TEST"},
			envVars:      map[string]string{},
			expectedParams: map[string]string{},
			expectError: true,
		},
		{
			name: "complex variable expansion",
			parameters: map[string]string{
				"COMPLEX": "${VAR1}_${VAR2}_${VAR3}",
			},
			expandAll:    true,
			expandValues: []string{},
			envVars: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
				"VAR3": "value3",
			},
			expectedParams: map[string]string{
				"COMPLEX": "value1_value2_value3",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment variables
			for k, v := range tt.envVars {
				t.Setenv(k, v)
			}

			// Run the test
			err := expandParameters(tt.parameters, tt.expandAll, tt.expandValues)

			// Check error
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Check parameters
			for k, v := range tt.expectedParams {
				assert.Equal(t, v, tt.parameters[k], "parameter %s should be expanded correctly", k)
			}
		})
	}
}

func BenchmarkExpandParameters(b *testing.B) {
	parameters := map[string]string{
		"DATABASE_NAME": "DB_$ENVIRONMENT",
		"SOME_SECRET":   "abc$abc",
		"COMPLEX":       "${VAR1}_${VAR2}_${VAR3}",
	}

	b.Run("selective expansion", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = expandParameters(parameters, false, []string{"DATABASE_NAME"})
		}
	})

	b.Run("expand all", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = expandParameters(parameters, true, []string{})
		}
	})
}
