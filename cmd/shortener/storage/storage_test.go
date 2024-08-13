package storage

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLocalStore_GetData(t *testing.T) {
	type want struct {
		isError bool
		data    string
	}
	tests := []struct {
		name  string
		input string
		want  want
	}{
		{
			"simple test",
			"key",
			want{
				isError: false,
				data:    "value",
			},
		},
		{
			"empty test",
			"",
			want{
				isError: true,
				data:    "",
			},
		},
		{
			"hard test",
			"vdsdhhmggdsadcxvvfsdsaf",
			want{
				isError: false,
				data:    "fdsbhgkjmdfsaew341gfds",
			},
		},
	}

	if err := InitializeFileLocalStore("data.json"); err != nil {
		t.Fatal(err)
	}
	mockLocalStore := make(LocalStore)
	require.NoError(t, mockLocalStore.SaveData(nil, "key", "value"))
	require.NoError(t, mockLocalStore.SaveData(nil, "vdsdhhmggdsadcxvvfsdsaf", "fdsbhgkjmdfsaew341gfds"))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := mockLocalStore.GetData(nil, tt.input)
			if tt.want.isError {
				require.Error(t, err)
			}
			require.Equal(t, tt.want.data, value)
		})
	}
}

func TestLocalStore_SaveData(t *testing.T) {
	tests := []struct {
		name       string
		inputKey   string
		inputValue string
		isError    bool
	}{
		{
			"simple test",
			"key",
			"value",
			false,
		},
		{
			"empty test #1",
			"",
			"value",
			true,
		},
		{
			"empty test #2",
			"key",
			"",
			true,
		},
		{
			"empty test #3",
			"",
			"",
			true,
		},
		{
			"hard test",
			"vdsdhhmggdsadcxvvfsdsaf",
			"asdasd",
			false,
		},
	}

	if err := InitializeFileLocalStore("data.json"); err != nil {
		t.Fatal(err)
	}
	mockLocalStore := make(LocalStore)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mockLocalStore.SaveData(nil, tt.inputKey, tt.inputValue)
			if tt.isError {
				require.Error(t, err)
				return
			}
			value, err := mockLocalStore.GetData(nil, tt.inputKey)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tt.inputValue, value)
		})
	}
}
