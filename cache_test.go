package cache

import (
	"testing"
	"time"
)

const (
	testRecordName     = "test-record"
	testRecordValue    = "test-value"
	testRecordValueTwo = "test-value-two"
	testAnotherRecord  = "test-record-two"
)

func TestCache(t *testing.T) {
	line := CreateLine(time.Millisecond * 250)

	if line.Renew(testRecordName) {
		t.Error("Hit is able to reach non-existing record")
	}

	if line.Store(testRecordName, testRecordValue) {
		t.Error("Store returns true when adding new record on empty cache")
	}

	if !line.Check(testRecordName) {
		t.Error("Check is unable to find existing record")
	}

	if !line.Renew(testRecordName) {
		t.Error("Hit is not able to reach existing record")
	}

	if line.Store(testAnotherRecord, testRecordValue) {
		t.Error("Store returns true when adding another record")
	}

	val, err := line.Get(testRecordName)
	if err != nil {
		t.Error("Unable to read the record")
	} else {
		if val != testRecordValue {
			t.Errorf("Retrieved value doesn't match")
		}
	}

	if !line.Store(testRecordName, testRecordValueTwo) {
		t.Error("Store returns false when updating an existing record")
	}

	val, err = line.Get(testRecordName)
	if err != nil {
		t.Error("Unable to read the record after update")
	} else {
		if val != testRecordValueTwo {
			t.Errorf("Retrieved value doesn't match after update")
		}
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 100)
		if !line.Renew(testAnotherRecord) {
			t.Errorf("Unable to renew expiration time for the record")
		}
	}

	if line.Check(testRecordName) {
		t.Error("Record should've expired by now")
	}

	if !line.Check(testAnotherRecord) {
		t.Error("Record should have not expired by now")
	}

	_, err = line.Get(testRecordName)
	if err == nil {
		t.Error("Able to read expired record")
	}

	line.Store(testRecordName, testRecordValue)
	for i := 0; i < 5; i++ {
		time.Sleep(time.Millisecond * 100)
		if !line.Renew(testRecordName) {
			t.Errorf("Unable to renew expiration time for the record")
		}
	}

	if line.Check(testAnotherRecord) {
		t.Error("Record should've expired by now")
	}
}
