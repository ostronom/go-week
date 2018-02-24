package week

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	week, err := New(0, 1)

	require.NoError(t, err)
	assert.Equal(t, Week{year: 0, week: 1}, week)

	week, err = New(-1, 0)

	assert.Error(t, err)
}

func TestWeek_MarshalText(t *testing.T) {

	tests := []struct {
		Week     Week
		Expected string
		Error    bool
	}{
		{Week: Week{year: 1, week: 1}, Expected: `0001-W01`},
		{Week: Week{year: 2001, week: 22}, Expected: `2001-W22`},
		{Week: Week{year: 9999, week: 53}, Expected: `9999-W53`},
		{Week: Week{year: -100, week: 22}, Error: true},
		{Week: Week{year: 2001, week: 99}, Error: true},
	}

	for _, tt := range tests {
		result, err := tt.Week.MarshalText()

		if tt.Error {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.Expected, string(result))
		}
	}
}

func TestWeek_UnmarshalText(t *testing.T) {

	tests := []struct {
		Value    string
		Expected Week
		Error    bool
	}{
		{Value: `0001-W01`, Expected: Week{year: 1, week: 1}},
		{Value: `2001-W22`, Expected: Week{year: 2001, week: 22}},
		{Value: `9999-W53`, Expected: Week{year: 9999, week: 53}},
		{Value: `9999-W99`, Error: true},
	}

	for _, tt := range tests {
		var week Week
		err := week.UnmarshalText([]byte(tt.Value))

		if tt.Error {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.Expected, week)
		}
	}
}

func TestWeek_MarshalJSON(t *testing.T) {

	tests := []struct {
		Week     Week
		Expected string
		Error    bool
	}{
		{Week: Week{year: 1, week: 1}, Expected: `"0001-W01"`},
		{Week: Week{year: 2001, week: 22}, Expected: `"2001-W22"`},
		{Week: Week{year: 9999, week: 53}, Expected: `"9999-W53"`},
		{Week: Week{year: 2001, week: 99}, Error: true},
	}

	t.Run("method call", func(t *testing.T) {
		for _, tt := range tests {
			result, err := tt.Week.MarshalJSON()

			if tt.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, tt.Expected, string(result))
			}
		}
	})

	t.Run("marshal struct", func(t *testing.T) {
		const template = `{"Week":%s,"WeekPtr":%s}`

		type testType struct {
			Week    Week
			WeekPtr *Week
		}

		for _, tt := range tests {
			testStruct := testType{Week: tt.Week, WeekPtr: &tt.Week}
			result, err := json.Marshal(testStruct)

			if tt.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.JSONEq(t, fmt.Sprintf(template, tt.Expected, tt.Expected), string(result))
			}
		}
	})
}

func TestWeek_UnmarshalJSON(t *testing.T) {

	tests := []struct {
		Value    string
		Expected Week
		Error    bool
	}{
		{Value: `"0001-W01"`, Expected: Week{year: 1, week: 1}},
		{Value: `"2001-W22"`, Expected: Week{year: 2001, week: 22}},
		{Value: `"9999-W53"`, Expected: Week{year: 9999, week: 53}},
		{Value: `2001-W11`, Error: true},
		{Value: `"9999-W99"`, Error: true},
	}

	t.Run("method call", func(t *testing.T) {
		for _, tt := range tests {
			var week Week
			err := week.UnmarshalJSON([]byte(tt.Value))

			if tt.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.Expected, week)
			}
		}
	})

	t.Run("unmarshal struct", func(t *testing.T) {
		const template = `{"Week":%s,"WeekPtr":%s}`

		type testType struct {
			Week    Week
			WeekPtr *Week
		}

		for _, tt := range tests {
			value := fmt.Sprintf(template, tt.Value, tt.Value)

			var testStruct testType
			err := json.Unmarshal([]byte(value), &testStruct)

			if tt.Error {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, testType{Week: tt.Expected, WeekPtr: &tt.Expected}, testStruct)
			}
		}
	})
}