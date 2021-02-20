package extratypes

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

type intForTestingConvertion struct {
	I Int `json:"i" toml:"i"`
}

var (
	validInt = Int{
		Val: 10,
		Nil: false,
	}

	validMinusInt = Int{
		Val: -10,
		Nil: false,
	}

	nilInt = Int{
		Val: 0,
		Nil: true,
	}
)

var (
	testIntJSONError      = []byte("10")
	testIntMinusJSONError = []byte("-10")
	testIntNilJSONError   = []byte("null")
	testIntErrJSONError   = []byte("a")
)

func TestIntString(t *testing.T) {
	i := Int{
		Val: 10,
		Nil: false,
	}

	if i.String() != "10" {
		t.Errorf("i [%s] is not 10", i)
	}

	i = Int{
		Val: 0,
		Nil: true,
	}

	if i.String() != "nil" {
		t.Errorf("i [%s] is not nil", i)
	}
}

func TestIntScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("No error was expected but have: %s", err)
	}
	defer db.Close()

	t.Run("test scan valid", func(te *testing.T) {
		rows := mock.NewRows([]string{"i"}).
			AddRow(-1).
			AddRow(1).
			AddRow(1.1)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		rs, _ := db.Query("SELECT")
		defer rs.Close()
		for rs.Next() {
			var i Int
			err := rs.Scan(&i)
			if err != nil {
				t.Errorf("Unable to scan Int: %s", err)
			}
			if i.Nil {
				te.Errorf("i [%d] not expected to be nil", i.Val)
			}
		}

		if rs.Err() != nil {
			te.Errorf("got rows error: %s", rs.Err())
		}
	})

	t.Run("test nil", func(te *testing.T) {
		rows := mock.NewRows([]string{"i"}).
			AddRow(nil)

		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		rs, _ := db.Query("SELECT")
		defer rs.Close()
		for rs.Next() {
			var i Int
			err := rs.Scan(&i)
			if err != nil {
				t.Errorf("Unable to scan Int: %s", err)
			}
			if !i.Nil {
				te.Errorf("i [%d] expected to be nil", i.Val)
			}
		}

		if rs.Err() != nil {
			te.Errorf("got rows error: %s", rs.Err())
		}
	})
}

func TestIntJSONMarshal(t *testing.T) {
	t.Run("marshal int", func(te *testing.T) {
		result, err := validInt.MarshalJSON()
		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testIntJSONError)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testIntJSONError, result)
		}

	})

	t.Run("marshal minus", func(te *testing.T) {
		result, err := validMinusInt.MarshalJSON()

		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testIntMinusJSONError)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testIntNilJSONError, result)
		}

	})

	t.Run("marshal nil", func(te *testing.T) {
		result, err := nilInt.MarshalJSON()

		if err != nil {
			te.Errorf("Error marshaling Int to JSON: %s", err)
		}

		cmp := bytes.Compare(result, testIntNilJSONError)
		if cmp != 0 {
			te.Errorf("Expected '%s', got '%s'", testIntNilJSONError, result)
		}

	})
}

func TestIntJSONUnmarshal(t *testing.T) {
	t.Run("unmarshal int", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntJSONError)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, validInt)
		}
	})

	t.Run("unmarshal minus int", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntMinusJSONError)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, validMinusInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, validMinusInt)
		}
	})

	t.Run("unmarshal nil", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntNilJSONError)

		if err != nil {
			te.Errorf("Unexpected error: %s", err)
		}

		cmp := reflect.DeepEqual(result, nilInt)
		if !cmp {
			te.Errorf("result %v not equal to %v", result, nilInt)
		}

	})

	t.Run("unmarshal err", func(te *testing.T) {
		var result Int
		err := result.UnmarshalJSON(testIntErrJSONError)

		if err == nil {
			te.Errorf("expected error, but none given")
		}

	})
}
