package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Errors(t *testing.T) {
	testCases := []struct {
		name              string
		err               error
		errMsg            string
		errType           error
		presentableErrMsg string
	}{
		{
			name:              "header not found",
			err:               NewHeaderNotFound("example"),
			errMsg:            "header was not found with name: example",
			errType:           HeaderNotFoundType,
			presentableErrMsg: `Header "example" was not found`,
		},
		{
			name:              "multiple headers found",
			err:               NewMultipleHeadersFound("example"),
			errMsg:            `invalid input(s): header "example" cannot have multiple values`,
			errType:           MultipleHeadersFoundType,
			presentableErrMsg: `Header "example" cannot have multiple values`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.EqualError(t, tc.err, tc.errMsg)
			assert.IsType(t, tc.errType, tc.err)
		})
	}
}
