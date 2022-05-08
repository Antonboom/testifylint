package main

import "testing"

func TestProduct(t *testing.T) {
	cases := []struct {
		a, b     []string
		expected [][]string
	}{
		{
			a:        nil,
			b:        nil,
			expected: [][]string{},
		},
		{
			a:        []string{},
			b:        []string{},
			expected: [][]string{},
		},
		{
			a:        []string{},
			b:        []string{"1", "2", "3"},
			expected: [][]string{},
		},
		{
			a: []string{"a"},
			b: []string{"1"},
			expected: [][]string{
				{"a", "1"},
			},
		},
		{
			a: []string{"a"},
			b: []string{"1", "2"},
			expected: [][]string{
				{"a", "1"}, {"a", "2"},
			},
		},
		{
			a: []string{"a"},
			b: []string{"1", "2", "3"},
			expected: [][]string{
				{"a", "1"}, {"a", "2"}, {"a", "3"},
			},
		},
		{
			a: []string{"a", "b"},
			b: []string{"1", "2", "3"},
			expected: [][]string{
				{"a", "1"}, {"a", "2"}, {"a", "3"},
				{"b", "1"}, {"b", "2"}, {"b", "3"},
			},
		},
		{
			a: []string{"a", "b", "c"},
			b: []string{"1", "2", "3"},
			expected: [][]string{
				{"a", "1"}, {"a", "2"}, {"a", "3"},
				{"b", "1"}, {"b", "2"}, {"b", "3"},
				{"c", "1"}, {"c", "2"}, {"c", "3"},
			},
		},
		{
			a: []string{"a", "b", "c"},
			b: []string{"1", "2"},
			expected: [][]string{
				{"a", "1"}, {"a", "2"},
				{"b", "1"}, {"b", "2"},
				{"c", "1"}, {"c", "2"},
			},
		},
		{
			a: []string{"a", "b", "c"},
			b: []string{"1"},
			expected: [][]string{
				{"a", "1"},
				{"b", "1"},
				{"c", "1"},
			},
		},
		{
			a:        []string{"a", "b", "c"},
			b:        []string{},
			expected: [][]string{},
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			res := Product(tt.a, tt.b)
			assertStringMatrixEqual(t, tt.expected, res)
		})
	}
}

func assertStringMatrixEqual(t *testing.T, expected, actual [][]string) {
	t.Helper()

	failed := func() bool {
		if len(actual) != len(expected) {
			return true
		}

		for i := range actual {
			if len(actual[i]) != len(expected[i]) {
				return true
			}

			for j := range actual[i] {
				if actual[i][j] != expected[i][j] {
					return true
				}
			}
		}

		return false
	}()

	if failed {
		t.Fatalf("actual %v != expected %v", actual, expected)
	}
}
