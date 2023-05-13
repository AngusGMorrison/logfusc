package logfusc

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSecret(t *testing.T) {
	t.Parallel()

	value := 42
	s := NewSecret(value)
	assert.Equal(t, value, s.value)
}

func Test_Secret_Expose(t *testing.T) {
	t.Parallel()

	value := 42
	s := NewSecret(value)
	assert.Equal(t, value, s.Expose())
}

func Test_Secret_String(t *testing.T) {
	t.Parallel()

	t.Run("when T is a string", func(t *testing.T) {
		t.Parallel()

		value := "foo"
		s := NewSecret(value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, value), s.String())
	})

	t.Run("when T is an int", func(t *testing.T) {
		t.Parallel()

		value := 42
		s := NewSecret(value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, value), s.String())
	})

	t.Run("when T is a struct", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Bar string
		}

		value := foo{Bar: "bar"}
		s := NewSecret(value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, value), s.String())
	})

	t.Run("when T is a pointer", func(t *testing.T) {
		t.Parallel()

		value := 42
		s := NewSecret(&value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, &value), s.String())
	})
}

func Test_Secret_GoString(t *testing.T) {
	t.Parallel()

	t.Run("when T is a string", func(t *testing.T) {
		t.Parallel()

		value := "foo"
		s := NewSecret(value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, value), s.GoString())
	})

	t.Run("when T is an int", func(t *testing.T) {
		t.Parallel()

		value := 42
		s := NewSecret(value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, value), s.GoString())
	})

	t.Run("when T is a struct", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Bar string
		}

		value := foo{Bar: "bar"}
		s := NewSecret(value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, value), s.GoString())
	})

	t.Run("when T is a pointer", func(t *testing.T) {
		t.Parallel()

		value := 42
		s := NewSecret(&value)
		assert.Equal(t, fmt.Sprintf(redactionTmpl, &value), s.GoString())
	})
}

func Test_string_formatting(t *testing.T) {
	t.Parallel()

	value := "password"
	secret := NewSecret(value)
	expected := fmt.Sprintf(redactionTmpl, value)

	testCases := []struct {
		verb     string
		expected string
	}{
		{"%s", expected},
		{"%q", fmt.Sprintf("\"%s\"", expected)},
		{"%v", expected},
		{"%+v", expected},
		{"%#v", expected},
		{"%x", fmt.Sprintf("%x", expected)},
		{"%X", fmt.Sprintf("%X", expected)},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.verb, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, fmt.Sprintf(tc.verb, secret))
		})
	}
}

func Test_Secret_MarshalJSON(t *testing.T) {
	t.Parallel()

	secretValue := "password"

	t.Run("can be marshaled standalone", func(t *testing.T) {
		t.Parallel()

		secret := NewSecret(secretValue)
		b, err := secret.MarshalJSON()

		assert.NoError(t, err)
		assert.JSONEq(t, fmt.Sprintf("%q", fmt.Sprintf(redactionTmpl, secretValue)), string(b))
	})

	t.Run("can be marshaled as part of a struct", func(t *testing.T) {
		t.Parallel()

		type foo struct {
			Secret Secret[string] `json:"secret"`
		}

		secret := NewSecret(secretValue)
		s := foo{Secret: secret}

		b, err := json.Marshal(s)

		assert.NoError(t, err)
		assert.JSONEq(t, fmt.Sprintf("{%q: %q}", "secret", fmt.Sprintf(redactionTmpl, secretValue)), string(b))
	})
}

func Test_Secret_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("when the the unmarshal target is valid", func(t *testing.T) {
		t.Parallel()

		t.Run("unmarshaling null is a no-op", func(t *testing.T) {
			t.Parallel()

			var secret Secret[string]
			err := json.Unmarshal([]byte(`null`), &secret)

			assert.NoError(t, err)
			assert.Empty(t, secret)
		})

		t.Run("Secret can be unmarshaled standalone", func(t *testing.T) {
			t.Parallel()

			var secret Secret[string]
			err := json.Unmarshal([]byte(`"password"`), &secret)

			assert.NoError(t, err)
			assert.Equal(t, "password", secret.value)
		})

		t.Run("Secret can be unmarshaled as part of a struct", func(t *testing.T) {
			t.Parallel()

			type foo struct {
				Secret Secret[string] `json:"secret"`
			}

			var secretStruct foo
			err := json.Unmarshal([]byte(`{"secret": "password"}`), &secretStruct)

			assert.NoError(t, err)
			assert.Equal(t, "password", secretStruct.Secret.value)
		})
	})

	t.Run("when the unmarshal target is invalid", func(t *testing.T) {
		t.Parallel()

		value := &unUnmarshalable{}
		secret := NewSecret(value)
		err := json.Unmarshal([]byte("0"), &secret)

		var unmarshalErr *UnmarshalError[*unUnmarshalable]
		assert.ErrorAs(t, err, &unmarshalErr)
	})
}

func Test_UnmarshalError_Error(t *testing.T) {
	t.Parallel()

	err := UnmarshalError[int]{
		intendedTarget: 0,
	}
	assert.Equal(t, fmt.Sprintf(unmarshalErrorTmpl, 0), err.Error())
}

type unUnmarshalable struct{}

func (u *unUnmarshalable) UnmarshalJSON([]byte) error {
	return errors.New("unmarshal error")
}
