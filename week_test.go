package week

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

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

func TestWeek_Next(t *testing.T) {

	tests := []struct {
		Curr  Week
		Next  Week
		Error bool
	}{
		{Curr: Week{year: 2003, week: 51}, Next: Week{year: 2003, week: 52}},
		{Curr: Week{year: 2003, week: 52}, Next: Week{year: 2004, week: 1}},
		{Curr: Week{year: 2004, week: 01}, Next: Week{year: 2004, week: 2}},
		{Curr: Week{year: 2004, week: 52}, Next: Week{year: 2004, week: 53}},
		{Curr: Week{year: 2004, week: 53}, Next: Week{year: 2005, week: 1}},
		{Curr: Week{year: 9999, week: 52}, Error: true},
	}

	for _, tt := range tests {
		prev, err := tt.Curr.Next()
		if tt.Error {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.Next, prev)
		}
	}
}

func TestWeek_Previous(t *testing.T) {

	tests := []struct {
		Curr  Week
		Prev  Week
		Error bool
	}{
		{Curr: Week{year: 2004, week: 01}, Prev: Week{year: 2003, week: 52}},
		{Curr: Week{year: 2003, week: 52}, Prev: Week{year: 2003, week: 51}},
		{Curr: Week{year: 2005, week: 01}, Prev: Week{year: 2004, week: 53}},
		{Curr: Week{year: 2004, week: 53}, Prev: Week{year: 2004, week: 52}},
		{Curr: Week{year: 0, week: 01}, Error: true},
	}

	for _, tt := range tests {
		prev, err := tt.Curr.Previous()
		if tt.Error {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.Prev, prev)
		}
	}
}

func TestWeek_MarshalText(t *testing.T) {

	tests := []struct {
		Week     Week
		Expected string
		Error    bool
	}{
		{Week: Week{year: 1, week: 1}, Expected: "0001-W01"},
		{Week: Week{year: 2001, week: 22}, Expected: "2001-W22"},
		{Week: Week{year: 9999, week: 52}, Expected: "9999-W52"},
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
		{Value: "0001-W01", Expected: Week{year: 1, week: 1}},
		{Value: "2001-W22", Expected: Week{year: 2001, week: 22}},
		{Value: "9999-W52", Expected: Week{year: 9999, week: 52}},
		{Value: "9999-W99", Error: true},
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
		{Week: Week{year: 9999, week: 52}, Expected: `"9999-W52"`},
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
		{Value: `"9999-W52"`, Expected: Week{year: 9999, week: 52}},
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

func TestWeek_Value(t *testing.T) {

	tests := []struct {
		Week     Week
		Expected string
		Error    bool
	}{
		{Week: Week{year: 1, week: 1}, Expected: "0001-W01"},
		{Week: Week{year: 2001, week: 22}, Expected: "2001-W22"},
		{Week: Week{year: 9999, week: 52}, Expected: "9999-W52"},
		{Week: Week{year: -100, week: 22}, Error: true},
		{Week: Week{year: 2001, week: 99}, Error: true},
	}

	for _, tt := range tests {
		result, err := tt.Week.Value()

		if tt.Error {
			assert.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, []byte(tt.Expected), result)
		}
	}
}

func TestWeek_Scan(t *testing.T) {
	const query = `SELECT week FROM test_table ORDER BY week LIMIT 1`

	tests := []struct {
		Value    driver.Value
		Expected Week
		Error    bool
	}{
		{Value: "0001-W01", Expected: Week{year: 1, week: 1}},
		{Value: "2001-W22", Expected: Week{year: 2001, week: 22}},
		{Value: "9999-W52", Expected: Week{year: 9999, week: 52}},
		{Value: "9999-W99", Error: true},
		{Value: 500, Error: true},
	}

	for _, tt := range tests {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"week"}).AddRow(tt.Value))

		row := db.QueryRow(query)

		var week Week
		err = row.Scan(&week)
		if tt.Error {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			assert.Equal(t, tt.Expected, week)
		}
	}
}
