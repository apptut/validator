package validator

import "testing"

func TestNew(t *testing.T) {

	data := map[string][]string{
		"hello": {"2"},
	}

	rules := map[string]string{
		"hello": "gt:2",
	}

	_, err := New(data, rules)
	if err != nil {
		println(err)
	}
}
