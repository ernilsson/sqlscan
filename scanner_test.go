package sqlscan

import (
	"errors"
	"testing"
)

type Entity struct {
	ID          int    `sql:"id"`
	Name        string `sql:"name"`
	Description string `sql:"description"`
}

type scannable struct {
	scan      func(...any) error
	scanCalls int
	columns   func() ([]string, error)
}

func (s *scannable) Scan(dest ...any) error {
	s.scanCalls++
	return s.scan(dest)
}

func (s *scannable) Columns() ([]string, error) {
	return s.columns()
}

func scanOk(expected int, t *testing.T) func(dest ...any) error {
	return func(dest ...any) error {
		if len(dest[0].([]any)) != expected {
			t.Error("Unexpected length", len(dest), "expected", expected)
		}
		return nil
	}
}

func scanErr(err error) func(dest ...any) error {
	return func(dest ...any) error {
		return err
	}
}

func columnsOk(cols ...string) func() ([]string, error) {
	return func() ([]string, error) {
		return cols, nil
	}
}

func columnsErr(err error) func() ([]string, error) {
	return func() ([]string, error) {
		return nil, err
	}
}

func TestStructScanner_ScanWithAllFields(t *testing.T) {
	s := scannable{
		scan:    scanOk(3, t),
		columns: columnsOk("id", "name", "description"),
	}
	scanner := New(&s)
	var e Entity
	err := scanner.Scan(&e)
	if err != nil {
		t.Error(err)
	}
	if s.scanCalls != 1 {
		t.Error("Expected 1 call to Scan, but got", s.scanCalls)
	}
}

func TestStructScanner_ScanWithSubsetOfFields(t *testing.T) {
	s := scannable{
		scan:    scanOk(2, t),
		columns: columnsOk("id", "description"),
	}
	scanner := New(&s)
	var e Entity
	err := scanner.Scan(&e)
	if err != nil {
		t.Error(err)
	}
	if s.scanCalls != 1 {
		t.Error("Expected 1 call to Scan, but got", s.scanCalls)
	}
}

func TestStructScanner_ScanWithNoFields(t *testing.T) {
	s := scannable{
		scan:    scanOk(0, t),
		columns: columnsOk(),
	}
	scanner := New(&s)
	var e Entity
	err := scanner.Scan(&e)
	if err != nil {
		t.Error(err)
	}
	if s.scanCalls != 1 {
		t.Error("Expected 1 call to Scan, but got", s.scanCalls)
	}
}

func TestStructScanner_ScanWithExtraneousFields(t *testing.T) {
	s := scannable{
		scan:    scanOk(3, t),
		columns: columnsOk("id", "name", "description", "amount"),
	}
	scanner := New(&s)
	var e Entity
	err := scanner.Scan(&e)
	if err != nil {
		t.Error(err)
	}
	if s.scanCalls != 1 {
		t.Error("Expected 1 call to Scan, but got", s.scanCalls)
	}
}

func TestStructScanner_ScanWithErroneousColumnsCall(t *testing.T) {
	s := scannable{
		columns: columnsErr(errors.New("mocked error")),
	}
	scanner := New(&s)
	var e Entity
	err := scanner.Scan(&e)
	if err == nil {
		t.Error("Expected error from Scan, but did not receive one")
	}
	if s.scanCalls != 0 {
		t.Error("Expected no calls to Scannable.Scan, but got", s.scanCalls)
	}
}

func TestStructScanner_ScanWithErroneousScanCall(t *testing.T) {
	s := scannable{
		scan:    scanErr(errors.New("mocked error")),
		columns: columnsOk("id", "name", "description", "amount"),
	}
	scanner := New(&s)
	var e Entity
	err := scanner.Scan(&e)
	if err == nil {
		t.Error("Expected error from Scan, but did not receive one")
	}
	if s.scanCalls != 1 {
		t.Error("Expected 1 call to Scannable.Scan, but got", s.scanCalls)
	}
}
