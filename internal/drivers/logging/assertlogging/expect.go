package assertlogging

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alexZaicev/realm-mgr/internal/drivers/logging"
)

var (
	unexpectedDebugLogsMsg  = "unexpected expectations on debug logs present"
	unexpectedLogMsg        = "unexpected %s log at index %d"
	missingLogMsg           = "missing %s log at index %d"
	messageMsg              = "message of log at index %d of level %s does not match"
	unexpectedFieldMsg      = "unexpected field %q found in log at index %d of level %s"
	missingFieldMsg         = "missing field %q found in log at index %d of level %s"
	noFieldExpectationsMsg  = "no expectations found for field %q of log at index %d of level %s"
	unmatchedLogMsg         = "no expectation matched log at index %d of level %s: %s"
	unmatchedExpectationMsg = "no log matched expectation at index %d of level %s: %s"
)

// Logger is a in-memory mock implementation of the logging.Logger interface.
type Logger struct {
	options           *loggerOptions
	parent            *Logger
	mutex             sync.Mutex
	fields            logging.Fields
	errorLogs         []*assertableLog
	expectedErrorLogs []*ExpectedLog
	warnLogs          []*assertableLog
	expectedWarnLogs  []*ExpectedLog
	infoLogs          []*assertableLog
	expectedInfoLogs  []*ExpectedLog
	debugLogs         []*assertableLog
	expectedDebugLogs []*ExpectedLog
}

// NewLogger creates a new Logger instance.
//
// Any registered expectations are automatically verified when the test completes along with
// identification of unfulfilled expectations.
func NewLogger(t testing.TB, opts ...Option) *Logger {
	options := &loggerOptions{}

	for _, o := range opts {
		o(options)
	}

	logger := &Logger{options: options}

	if !options.skipAssert {
		t.Cleanup(func() {
			logger.assertExpectations(t)
		})
	}

	return logger
}

// ExpectedLog describes a single log statement that is expected to be generated by tested code.
type ExpectedLog struct {
	msg    string
	fields Fields
}

// newExpectedLog creates a new ExpectedLog instance.
func newExpectedLog(msg string) *ExpectedLog {
	return &ExpectedLog{
		msg:    msg,
		fields: make(map[string][]FieldValue),
	}
}

func (l *ExpectedLog) String() string {
	var fields []string
	for fieldName, fieldConstraints := range l.fields {
		constraintNames := make([]string, len(fieldConstraints))
		for i, fieldConstraint := range fieldConstraints {
			constraintNames[i] = fieldConstraint.String()
		}
		fields = append(fields, fmt.Sprintf("%s:%+v", fieldName, constraintNames))
	}
	sort.Strings(fields) // make output deteministic
	return fmt.Sprintf("message=%q, fields=map[%s]", l.msg, strings.Join(fields, " "))
}

// Field represents a log field with a key and number of expectations upon the log field value.
type Field struct {
	key   string
	value []FieldValue
}

// NewField creates a new Field instance.
func NewField(key string, value ...FieldValue) *Field {
	return &Field{
		key:   key,
		value: value,
	}
}

// Fields contains multiple log field expectations.
type Fields map[string][]FieldValue

// ExpectError adds an expectation for a log message at ERROR level. Expectations on log fields can
// be added using the returned ExpectedLog.
func (l *Logger) ExpectError(msg string) *ExpectedLog {
	log := newExpectedLog(msg)
	l.expectedErrorLogs = append(l.expectedErrorLogs, log)
	return log
}

// ExpectWarn adds an expectation for a log message at WARN level. Expectations on log fields can
// be added using the returned ExpectedLog.
func (l *Logger) ExpectWarn(msg string) *ExpectedLog {
	log := newExpectedLog(msg)
	l.expectedWarnLogs = append(l.expectedWarnLogs, log)
	return log
}

// ExpectInfo adds an expectation for a log message at INFO level. Expectations on log fields can
// be added using the returned ExpectedLog.
func (l *Logger) ExpectInfo(msg string) *ExpectedLog {
	log := newExpectedLog(msg)
	l.expectedInfoLogs = append(l.expectedInfoLogs, log)
	return log
}

// ExpectDebug adds an expectation for a log message at DEBUG level. If you wish to set expectations
// on DEBUG logs, you must also use the AssertDebugLogs option when calling NewLogger. Expectations
// on log fields can be added using the returned ExpectedLog.
func (l *Logger) ExpectDebug(msg string) *ExpectedLog {
	log := newExpectedLog(msg)
	l.expectedDebugLogs = append(l.expectedDebugLogs, log)
	return log
}

// assertExpectations ensures that the expected logs have all been called and that no unexpected
// logs have been produced.
//
// This is automatically registered as a test cleanup function so it is not necessary to call it
// manually.
func (l *Logger) assertExpectations(t assert.TestingT) bool {
	matched := true

	matched = l.assertExpectedLogs(t, l.expectedErrorLogs, l.errorLogs, "error") && matched
	matched = l.assertExpectedLogs(t, l.expectedWarnLogs, l.warnLogs, "warn") && matched
	matched = l.assertExpectedLogs(t, l.expectedInfoLogs, l.infoLogs, "info") && matched

	if l.options.assertDebugLogs {
		matched = l.assertExpectedLogs(t, l.expectedDebugLogs, l.debugLogs, "debug")
	} else if !assert.Empty(t, l.expectedDebugLogs, unexpectedDebugLogsMsg) {
		matched = false
	}

	return matched
}

// assertExpectedLogs checks that logs at the given level match and are in the same order.
func (l *Logger) assertExpectedLogs(
	t assert.TestingT,
	expected []*ExpectedLog,
	actual []*assertableLog,
	level string,
) bool {
	if l.options.assertWithoutOrder {
		return assertExpectedLogsUnordered(t, expected, actual, level)
	}

	matched := true
	length := max(len(expected), len(actual))

	for i := 0; i < length; i++ {
		// TODO: use assert.Greater when https://github.com/stretchr/testify/pull/1021 is merged
		if !assert.True(t, len(expected) > i, unexpectedLogMsg, level, i) {
			matched = false
			continue
		}

		// TODO: use assert.Greater when https://github.com/stretchr/testify/pull/1021 is merged
		if !assert.True(t, len(actual) > i, missingLogMsg, level, i) {
			matched = false
			continue
		}

		expectedLog := expected[i]
		actualLog := actual[i]

		matched = expectedLog.assertExpected(t, actualLog, i, level) && matched
	}

	return matched
}

// WithField adds an expectation upon a single log field. There must be at least one expectation
// on each log field.
func (l *ExpectedLog) WithField(key string, value ...FieldValue) *ExpectedLog {
	return l.WithFields(NewField(key, value...))
}

// WithField adds an expectation upon a log field added using WithError. There must be at least one
// expectation on each log field.
func (l *ExpectedLog) WithError(value ...FieldValue) *ExpectedLog {
	return l.WithField(logging.ErrKey, value...)
}

// WithFields adds expectations upon multiple log fields. There must be at least one expectation on
// each log field.
func (l *ExpectedLog) WithFields(fields ...*Field) *ExpectedLog {
	for _, field := range fields {
		l.fields[field.key] = field.value
	}

	return l
}

// assertExpected checks that an expected log matches a log that was produced.
func (l *ExpectedLog) assertExpected(
	t assert.TestingT,
	actual *assertableLog,
	index int,
	level string,
) bool {
	matched := assert.Equal(t, l.msg, actual.msg, messageMsg, index, level)
	keys := map[string]struct{}{}

	for key := range l.fields {
		keys[key] = struct{}{}
	}

	for key := range actual.fields {
		keys[key] = struct{}{}
	}

	for key := range keys {
		expectations, ok := l.fields[key]
		if !assert.True(t, ok, unexpectedFieldMsg, key, index, level) {
			matched = false
			continue
		}

		value, ok := actual.fields[key]
		if !assert.True(t, ok, missingFieldMsg, key, index, level) {
			matched = false
			continue
		}

		if !assert.True(t, len(expectations) > 0, noFieldExpectationsMsg, key, index, level) {
			matched = false
			continue
		}

		for _, expectation := range expectations {
			matched = expectation.assert(t, value, key, index, level) && matched
		}
	}

	return matched
}

// max compares to integers and returns the greatest.
func max(x, y int) int {
	if x > y {
		return x
	}

	return y
}

func assertExpectedLogsUnordered(
	t assert.TestingT,
	expected []*ExpectedLog,
	actual []*assertableLog,
	level string,
) bool {
	// innerT is used to capture failures from checking each expectation from the logs - since we only want to fail
	// the overall test if they all failed to match
	var innerT testing.TB = new(testing.T)

	// Need to make a set (map) so that we can remove the items from here once they have been used to match a log successfully
	// so that each expected log can only match once, and we can report exactly which expectations went unfulfilled
	expectedMap := make(map[int]*ExpectedLog, len(expected))
	for i, expectedLog := range expected {
		expectedMap[i] = expectedLog
	}

	matched := true

	for i, actualLog := range actual {
		logMatched := false
		for j, expectedLog := range expected {
			if _, ok := expectedMap[j]; ok {
				if logMatched = expectedLog.assertExpected(innerT, actualLog, i, level); logMatched {
					delete(expectedMap, j) // each expected log should only match one actual log
					break
				}
			}
		}
		// matching out of order so can only say that no expectation matched, not which expectation should have matched
		matched = assert.True(t, logMatched, unmatchedLogMsg, i, level, actualLog) && matched
	}

	for i, expectedLog := range expected {
		_, ok := expectedMap[i]
		// report unmatched expecations
		matched = assert.False(t, ok, unmatchedExpectationMsg, i, level, expectedLog) && matched
	}

	return matched
}
