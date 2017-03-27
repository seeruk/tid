package parameters_test

import (
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/eidolon/console/assert"
	"github.com/eidolon/console/parameters"
)

func TestBoolValue(t *testing.T) {
	t.Run("NewBoolValue()", func(t *testing.T) {
		truthy := true
		falsey := false

		truthyValue := parameters.NewBoolValue(&truthy)
		falseyValue := parameters.NewBoolValue(&falsey)

		assert.Equal(t, truthyValue.String(), "true")
		assert.Equal(t, falseyValue.String(), "false")
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.BoolValue

			valid := []string{
				"1",
				"t",
				"T",
				"TRUE",
				"true",
				"True",
				"0",
				"f",
				"F",
				"FALSE",
				"false",
				"False",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should error for invalid values", func(t *testing.T) {
			var value parameters.BoolValue

			invalid := []string{
				"",
				"yes",
				"no",
				"Y",
				"N",
			}

			for _, item := range invalid {
				err := value.Set(item)
				assert.NotOK(t, err)
			}
		})

		t.Run("should modify the bool that it references", func(t *testing.T) {
			ref := true
			value := parameters.NewBoolValue(&ref)

			assert.Equal(t, true, ref)

			value.Set("false")
			assert.Equal(t, false, ref)

			value.Set("true")
			assert.Equal(t, true, ref)
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"1":     "true",
			"t":     "true",
			"T":     "true",
			"TRUE":  "true",
			"true":  "true",
			"True":  "true",
			"0":     "false",
			"f":     "false",
			"F":     "false",
			"FALSE": "false",
			"false": "false",
			"False": "false",
		}

		for in, expected := range inOut {
			value := new(parameters.BoolValue)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("FlagValue()", func(t *testing.T) {
		var value parameters.BoolValue

		assert.Equal(t, "true", value.FlagValue())
	})
}

func TestDateValue(t *testing.T) {
	t.Run("NewDateValue()", func(t *testing.T) {
		date, err := time.Parse("2006-01-02", "2017-02-27")
		if err != nil {
			t.Error(err)
		}

		dateValue := parameters.NewDateValue(&date)

		assert.Equal(t, "2017-02-27", dateValue.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.DateValue

			valid := []string{
				"2017-02-27",
				"2000-12-01",
				"0123-04-05",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should error for invalid values", func(t *testing.T) {
			var value parameters.DateValue

			invalid := []string{
				"",
				"hello",
				"00:14AM",
			}

			for _, item := range invalid {
				err := value.Set(item)
				assert.NotOK(t, err)
			}
		})

		t.Run("should modify the date that it references", func(t *testing.T) {
			ref, err := time.Parse("2006-01-02", "2017-02-27")
			if err != nil {
				t.Error(err)
			}

			value := parameters.NewDateValue(&ref)

			assert.Equal(t, "2017-02-27", ref.Format("2006-01-02"))

			value.Set("2000-12-01")
			assert.Equal(t, "2000-12-01", ref.Format("2006-01-02"))

			value.Set("0123-04-05")
			assert.Equal(t, "0123-04-05", ref.Format("2006-01-02"))
		})
	})

	t.Run("String()", func(t *testing.T) {
		results := []string{
			"2017-02-27",
			"2000-12-01",
			"0123-04-05",
		}

		for _, result := range results {
			value := new(parameters.DateValue)
			value.Set(result)

			actual := value.String()
			assert.Equal(t, result, actual)
		}
	})
}

func TestDurationValue(t *testing.T) {
	t.Run("NewDurationValue()", func(t *testing.T) {
		duration := time.Second
		durationValue := parameters.NewDurationValue(&duration)

		assert.Equal(t, "1s", durationValue.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.DurationValue

			valid := []string{
				"5us",
				"999ns",
				"1s",
				"3m",
				"1h33m2s",
				"5h",
				"365h",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should error for invalid values", func(t *testing.T) {
			var value parameters.DurationValue

			invalid := []string{
				"",
				"1d",
				"20y",
				"20 decades",
			}

			for _, item := range invalid {
				err := value.Set(item)
				assert.NotOK(t, err)
			}
		})

		t.Run("should modify the bool that it references", func(t *testing.T) {
			ref := time.Second
			value := parameters.NewDurationValue(&ref)

			assert.Equal(t, time.Second, ref)

			value.Set("1m")
			assert.Equal(t, time.Minute, ref)

			value.Set("1h")
			assert.Equal(t, ref, time.Hour)
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"1us": "1µs",
			"1ns": "1ns",
			"1s":  "1s",
			"1m":  "1m0s",
			"1h":  "1h0m0s",
			"5s":  "5s",
			"10h": "10h0m0s",
		}

		for in, expected := range inOut {
			value := new(parameters.DurationValue)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})
}

func TestFloat32Value(t *testing.T) {
	t.Run("NewFloat32Value()", func(t *testing.T) {
		float := float32(3.14)
		floatValue := parameters.NewFloat32Value(&float)

		assert.Equal(t, "3.14", floatValue.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.Float32Value

			valid := []string{
				"3",
				"3.1",
				"3.14",
				"314.159e-2",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should error for invalid values", func(t *testing.T) {
			var value parameters.Float32Value

			invalid := []string{
				"",
				"Hello, World!",
				"Three point one four",
			}

			for _, item := range invalid {
				err := value.Set(item)
				assert.NotOK(t, err)
			}
		})

		t.Run("should modify the float32 that it references", func(t *testing.T) {
			ref := float32(3.14)
			value := parameters.NewFloat32Value(&ref)

			assert.Equal(t, float32(3.14), ref)

			value.Set("3.14159")
			assert.Equal(t, float32(3.14159), ref)

			value.Set("10")
			assert.Equal(t, float32(10), ref)
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"3":          "3",
			"3.14":       "3.14",
			"3.14159":    "3.14159",
			"314.159e-2": "3.14159",
		}

		for in, expected := range inOut {
			value := new(parameters.Float32Value)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})
}

func TestFloat64Value(t *testing.T) {
	t.Run("NewFloat64Value()", func(t *testing.T) {
		float := float64(3.14)
		floatValue := parameters.NewFloat64Value(&float)

		assert.Equal(t, "3.14", floatValue.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.Float64Value

			valid := []string{
				"3",
				"3.1",
				"3.14",
				"314.159e-2",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should error for invalid values", func(t *testing.T) {
			var value parameters.Float64Value

			invalid := []string{
				"",
				"Hello, World!",
				"Three point one four",
			}

			for _, item := range invalid {
				err := value.Set(item)
				assert.NotOK(t, err)
			}
		})

		t.Run("should modify the float64 that it references", func(t *testing.T) {
			ref := float64(3.14)
			value := parameters.NewFloat64Value(&ref)

			assert.Equal(t, float64(3.14), ref)

			value.Set("3.14159")
			assert.Equal(t, float64(3.14159), ref)

			value.Set("10")
			assert.Equal(t, float64(10), ref)
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"3":          "3",
			"3.14":       "3.14",
			"3.14159":    "3.14159",
			"314.159e-2": "3.14159",
		}

		for in, expected := range inOut {
			value := new(parameters.Float64Value)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})
}

func TestIntValue(t *testing.T) {
	t.Run("NewIntValue()", func(t *testing.T) {
		intRef := 3
		intValue := parameters.NewIntValue(&intRef)

		assert.Equal(t, "3", intValue.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.Float64Value

			valid := []string{
				"3",
				"10",
				"25",
				"100000",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should error for invalid values", func(t *testing.T) {
			var value parameters.IntValue

			invalid := []string{
				"",
				"Hello, World!",
				"Three point one four",
				"92233720368547758070",
			}

			for _, item := range invalid {
				err := value.Set(item)
				assert.NotOK(t, err)
			}
		})

		t.Run("should modify the int that it references", func(t *testing.T) {
			ref := 5
			value := parameters.NewIntValue(&ref)

			assert.Equal(t, 5, ref)

			value.Set("10")
			assert.Equal(t, 10, ref)

			value.Set("25")
			assert.Equal(t, 25, ref)
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"3":    "3",
			"10":   "10",
			"1000": "1000",
		}

		for in, expected := range inOut {
			value := new(parameters.IntValue)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})
}

func TestIPValue(t *testing.T) {
	t.Run("NewIPValue()", func(t *testing.T) {
		ipRef := net.ParseIP("127.0.0.1")
		ipValue := parameters.NewIPValue(&ipRef)

		assert.Equal(t, "127.0.0.1", ipValue.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.IPValue

			valid := []string{
				"127.0.0.1",
				"192.168.0.1",
				"10.0.0.1",
				"255.255.255.0",
				"8.8.8.8",
				"8.8.4.4",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should error for invalid values", func(t *testing.T) {
			var value parameters.IPValue

			invalid := []string{
				"",
				"Not an IP adddress",
				"Hello, World!",
				"123 Fake Street",
				"127 0 0 1",
			}

			for _, item := range invalid {
				err := value.Set(item)
				assert.NotOK(t, err)
			}
		})

		t.Run("should modify the IP that it references", func(t *testing.T) {
			ref := net.ParseIP("127.0.0.1")
			value := parameters.NewIPValue(&ref)

			assert.Equal(t, value.String(), ref.String())

			value.Set("192.168.0.1")
			assert.Equal(t, value.String(), ref.String())

			value.Set("10.0.0.1")
			assert.Equal(t, value.String(), ref.String())
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"127.0.0.1":   "127.0.0.1",
			"192.168.0.1": "192.168.0.1",
			"10.0.0.1":    "10.0.0.1",
		}

		for in, expected := range inOut {
			value := new(parameters.IPValue)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})
}

func TestStringValue(t *testing.T) {
	t.Run("NewStringValue()", func(t *testing.T) {
		expected := "Hello, World!"
		actual := parameters.NewStringValue(&expected)

		assert.Equal(t, expected, actual.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.StringValue

			valid := []string{
				"Hello",
				"World",
				"Hello, World!",
				"3.14",
				"http://www.google.co.uk/",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should modify the string that it references", func(t *testing.T) {
			ref := "Hello"

			value := parameters.NewStringValue(&ref)
			assert.Equal(t, "Hello", ref)

			value.Set("World")
			assert.Equal(t, "World", ref)
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"Hello": "Hello",
			"World": "World",
			"hello": "hello",
			"world": "world",
		}

		for in, expected := range inOut {
			value := new(parameters.StringValue)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})
}

func TestUrlValue(t *testing.T) {
	t.Run("NewURLValue()", func(t *testing.T) {
		expected := "https://www.google.co.uk/"

		actual, err := url.Parse(expected)
		assert.OK(t, err)

		actualValue := parameters.NewURLValue(actual)
		assert.Equal(t, expected, actualValue.String())
	})

	t.Run("Set()", func(t *testing.T) {
		t.Run("should not error for valid values", func(t *testing.T) {
			var value parameters.URLValue

			valid := []string{
				"https://www.google.co.uk/",
				"ws://www.elliotdwright.com:9000/",
				"ftp://whouses.ftpanymore.com:21/",
			}

			for _, item := range valid {
				err := value.Set(item)
				assert.OK(t, err)
			}
		})

		t.Run("should modify the URL that it references", func(t *testing.T) {
			oldURL := "https://www.google.co.uk/"
			newURL := "https://www.elliotdwright.com/"

			ref, err := url.Parse(oldURL)
			assert.OK(t, err)

			value := parameters.NewURLValue(ref)
			assert.Equal(t, oldURL, ref.String())

			value.Set(newURL)
			assert.Equal(t, newURL, ref.String())
		})
	})

	t.Run("String()", func(t *testing.T) {
		inOut := map[string]string{
			"https://www.google.co.uk/":        "https://www.google.co.uk/",
			"ws://www.elliotdwright.com:9000/": "ws://www.elliotdwright.com:9000/",
			"ftp://whouses.ftpanymore.com:21/": "ftp://whouses.ftpanymore.com:21/",
		}

		for in, expected := range inOut {
			value := new(parameters.URLValue)
			value.Set(in)

			actual := value.String()
			assert.Equal(t, expected, actual)
		}
	})
}
