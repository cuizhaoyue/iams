package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions_Validate(t *testing.T) {
	opts := &Options{
		Level:            "test",
		Format:           "test",
		EnableColor:      true,
		DisableCaller:    false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	errs := opts.Validate()
	expected := `[unrecognized level: "test" not a valid log format: "test"]`
	fmt.Printf("%s", errs)
	assert.Equal(t, expected, fmt.Sprintf("%s", errs))
}
