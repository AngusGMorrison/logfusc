package logfusc

import (
	"bytes"
	"fmt"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func Test_jsoniter(t *testing.T) {
	t.Parallel()

	type foo struct {
		Secret Secret[string] `json:"secret"`
	}

	value := "bar"
	container := foo{Secret: NewSecret(value)}
	b, err := jsoniter.Marshal(container)

	assert.NoError(t, err)
	assert.JSONEq(t, fmt.Sprintf("{%q: %q}", "secret", fmt.Sprintf(redactionTmpl, value)), string(b))
}

func Test_zerolog(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	value := "bar"
	secret := NewSecret(value)
	logger.Print(secret)
	logged := buf.String()

	assert.Contains(t, logged, fmt.Sprintf(redactionTmpl, value))
	assert.NotContains(t, logged, value)
}

func Test_zap(t *testing.T) {
	t.Parallel()

	var buf zaptest.Buffer
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		&buf,
		zapcore.DebugLevel,
	)
	logger := zap.New(core)
	sugar := logger.Sugar()
	value := "bar"
	secret := NewSecret(value)
	sugar.Infow("sugar", secret)
	_ = logger.Sync()
	logged := buf.String()

	assert.Contains(t, logged, fmt.Sprintf(redactionTmpl, value))
	assert.NotContains(t, logged, value)
}

func Test_logrus(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	value := "bar"
	secret := NewSecret(value)
	logger.WithField("secret", secret).Info("")
	logged := buf.String()

	assert.Contains(t, logged, fmt.Sprintf(redactionTmpl, value))
	assert.NotContains(t, logged, value)
}
