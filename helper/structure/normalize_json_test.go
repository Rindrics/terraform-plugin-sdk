// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package structure

import (
	"testing"
)

func TestNormalizeJsonString_valid(t *testing.T) {
	// Well formatted and valid.
	validJson := `{
   "abc": {
      "def": 123,
      "xyz": [
         {
            "a": "ホリネズミ"
         },
         {
            "b": "1\\n2"
         }
      ]
   }
}`
	expected := `{"abc":{"def":123,"xyz":[{"a":"ホリネズミ"},{"b":"1\\n2"}]}}`

	actual, err := NormalizeJsonString(validJson)
	if err != nil {
		t.Fatalf("Expected not to throw an error while parsing JSON, but got: %s", err)
	}

	if actual != expected {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, expected)
	}

	// Well formatted but not valid,
	// missing closing square bracket.
	invalidJson := `{
   "abc": {
      "def": 123,
      "xyz": [
         {
            "a": "1"
         }
      }
   }
}`
	actual, err = NormalizeJsonString(invalidJson)
	if err == nil {
		t.Fatalf("Expected to throw an error while parsing JSON, but got: %s", err)
	}

	// We expect the invalid JSON to be shown back to us again.
	if actual != invalidJson {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, invalidJson)
	}

	// Verify that it leaves strings alone
	testString := "2016-07-28t04:07:02z\nsomething else"
	expected = "2016-07-28t04:07:02z\nsomething else"
	actual, err = NormalizeJsonString(testString)
	if err == nil {
		t.Fatalf("Expected to throw an error while parsing JSON, but got: %s", err)
	}

	if actual != expected {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", actual, expected)
	}
}

func TestNormalizeJsonString_invalid(t *testing.T) {
	// Well formatted but not valid,
	// missing closing squre bracket.
	invalidJson := `{
   "abc": {
      "def": 123,
      "xyz": [
         {
            "a": "1"
         }
      }
   }
}`
	expected := `{"abc":{"def":123,"xyz":[{"a":"ホリネズミ"},{"b":"1\\n2"}]}}`
	actual, err := NormalizeJsonString(invalidJson)
	if err == nil {
		t.Fatalf("Expected to throw an error while parsing JSON, but got: %s", err)
	}

	// We expect the invalid JSON to be shown back to us again.
	if actual != invalidJson {
		t.Fatalf("Got:\n\n%s\n\nExpected:\n\n%s\n", expected, invalidJson)
	}
}

func TestNormalizeJsonString_arrayConversion(t *testing.T) {
	// Single-element array should be converted to string to conform to AWS behavior
	t.Run("SingleElementArrayToString", func(t *testing.T) {
		singleElementArrayJson := `{
    "Resource": ["arn:aws:iam::123456789012:root"]
}`
		expected := `{"Resource":"arn:aws:iam::123456789012:root"}`

		actual, err := NormalizeJsonString(singleElementArrayJson)
		if err != nil {
			t.Fatalf("Expected not to throw an error while parsing JSON, but got: %s", err)
		}

		if actual != expected {
			t.Fatalf("NormalizeJsonString did not convert single-element array to string. Got:\n\n%s\n\nExpected:\n\n%s\n", actual, expected)
		}
	})

	t.Run("MultipleElementArrayRemainsArray", func(t *testing.T) {
		multiElementArrayJson := `{
    "Resource": ["arn:aws:iam::123456789012:root", "arn:aws:iam::123456789012:user/user123"]
}`
		expectedMulti := `{"Resource":["arn:aws:iam::123456789012:root","arn:aws:iam::123456789012:user/user123"]}`

		actualMulti, errMulti := NormalizeJsonString(multiElementArrayJson)
		if errMulti != nil {
			t.Fatalf("Expected not to throw an error while parsing JSON, but got: %s", errMulti)
		}

		if actualMulti != expectedMulti {
			t.Fatalf("Multiple-element array should not be converted to string. Got:\n\n%s\n\nExpected:\n\n%s\n", actualMulti, expectedMulti)
		}
	})
}
