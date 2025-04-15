package lang

import (
	"errors"
	lang "github.com/KatePril/architecture-lab-3/painter/lang/errors"
	"strings"
	"testing"

	"github.com/KatePril/architecture-lab-3/painter"
	"github.com/stretchr/testify/assert"
)

func TestParser_EmptyInput(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader("")

	options, err := parser.Parse(input)

	assert.Nil(t, err)
	assert.Empty(t, options)
}

func TestParser_OneCommandLine(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader("figure 1 0.25")

	options, err := parser.Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(options))
}

func TestParser_CommandNameSensitivity(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader("FigUre 1 0.25")

	options, err := parser.Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(options))
}

func TestParser_CommandLine_WithSomeWhiteSpaces(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader("  figure    1 0.25 ")

	options, err := parser.Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(options))
}

func TestParser_MultipleCommands(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader(`white
bgrect 0.25 0.25 0.75 0.75
figure 0.5 0.5
green
figure 0.6 0.6
update`)

	options, err := parser.Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6, len(options))
}

func TestParser_TwoEqualCommandLines(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader(`figure 0.6 0.6
figure 0.6 0.6`)

	options, err := parser.Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(options))
}

func TestParser_MultipleCommands_WithEmptyLines(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader(`
white

bgrect 0.25 0.25 0.75 0.75
figure 0.5 0.5
green

figure 0.6 0.6
update`)

	options, err := parser.Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6, len(options))
}

func TestParser_InvalidCommandName(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader("invalid 1 0.25")

	_, err := parser.Parse(input)

	var cmdErr lang.InvalidCommandError
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &cmdErr))
	assert.Equal(t, "invalid", cmdErr.CommandName)
}

func TestParser_InvalidArguments(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader("figure some 0.25")

	_, err := parser.Parse(input)

	var cmdErr lang.InvalidArgumentsError
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &cmdErr))
}

func TestParser_InvalidAmountOfArguments(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader("figure 1 0.25 0")

	_, err := parser.Parse(input)

	var cmdErr lang.InvalidArgumentsError
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &cmdErr))
}

func TestParser_MultipleCommands_WithOneInvalid(t *testing.T) {
	state := painter.State{}
	var parser Parser = CreateOperationParser(&state)
	input := strings.NewReader(`
white

bgrect 0.25 0.25 0.75 0.75
figure 0.5 0.5
invalid

figure 0.6 0.6
update`)

	options, err := parser.Parse(input)

	var cmdErr lang.InvalidCommandError
	assert.Empty(t, options)
	assert.NotNil(t, err)
	assert.True(t, errors.As(err, &cmdErr))
	assert.Equal(t, "invalid", cmdErr.CommandName)
}
