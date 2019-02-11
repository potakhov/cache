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

	deleteLine := CreateLine(time.Second * 60)
	if deleteLine.Store(1, 0) {
		t.Errorf("Unable to store a record")
	}
	if deleteLine.Store(2, 0) {
		t.Errorf("Unable to store a record")
	}
	if deleteLine.Store(3, 0) {
		t.Errorf("Unable to store a record")
	}

	deleteLine.Delete(2)

	if !deleteLine.Check(1) {
		t.Errorf("Unable to find a record after deletion")
	}
	if !deleteLine.Check(3) {
		t.Errorf("Unable to find a record after deletion")
	}
	if deleteLine.Check(2) {
		t.Errorf("Able to find a record after deletion")
	}

	deleteLine.Delete(3)

	if !deleteLine.Check(1) {
		t.Errorf("Unable to find a record after deletion")
	}
	if deleteLine.Check(3) {
		t.Errorf("Able to find a record after deletion")
	}

	if deleteLine.Store(4, 0) {
		t.Errorf("Unable to store a record")
	}

	deleteLine.Delete(1)
	if deleteLine.Check(1) {
		t.Errorf("Able to find a record after deletion")
	}
	if !deleteLine.Check(4) {
		t.Errorf("Unable to find a record after deletion")
	}

	if deleteLine.Delete(1) {
		t.Errorf("Able to delete a record after deletion")
	}

	deleteLine.Delete(4)
	if deleteLine.Check(4) {
		t.Errorf("Able to find a record after deletion")
	}
}
