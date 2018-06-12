package cache

import (
	"errors"
	"time"
)

// ListKeyType defines the list key type
type ListKeyType string

// ErrCacheRecordNotFound happens when Get() is unable to locate a record.
var ErrCacheRecordNotFound = errors.New("Record is not found")

type cacheList struct {
	prev       *cacheList
	next       *cacheList
	key        ListKeyType
	expiration time.Time
}

type cacheItem struct {
	value interface{}
	ptr   *cacheList
}

// Line is generalized object cache with expiration which generally works with O(1) complexity
type Line struct {
	defaultExpirationTime time.Duration
	store                 map[ListKeyType]*cacheItem
	first                 *cacheList
	last                  *cacheList
	nextCheckupTime       time.Time
}

// CreateLine creates a new line with specific default expiration duration for all objects
func CreateLine(defaultExpirationTime time.Duration) *Line {
	ret := &Line{defaultExpirationTime: defaultExpirationTime, nextCheckupTime: time.Now().Add(defaultExpirationTime)}
	ret.store = make(map[ListKeyType]*cacheItem)
	return ret
}

// StoreFor records a new hit for <key> (creating the record if it doesn't exist) shifting expiration further in time
// @returns true if record has been updated, false if new record was created
func (cl *Line) StoreFor(key ListKeyType, value interface{}, expires time.Duration) bool {
	curr, ok := cl.store[key]
	if !ok {
		el := &cacheList{key: key, expiration: time.Now().Add(expires)}
		cl.store[key] = &cacheItem{value: value, ptr: el}
		if cl.first == nil {
			cl.first = el
			cl.last = el
		} else {
			cl.last.next = el
			el.prev = cl.last
			cl.last = el
		}
	} else {
		curr.value = value
		curr.ptr.expiration = time.Now().Add(expires)

		if cl.last != curr.ptr {
			if curr.ptr.prev != nil {
				curr.ptr.prev.next = curr.ptr.next
			} else {
				cl.first = curr.ptr.next
			}

			curr.ptr.next.prev = curr.ptr.prev
			cl.last.next = curr.ptr
			curr.ptr.prev = cl.last
			curr.ptr.next = nil
			cl.last = curr.ptr
		}
	}
	cl.expire()
	return ok
}

// Store records a new hit for <key> (creating the record if it doesn't exist) shifting expiration further in time (using default expiration timeout)
// @returns true if record has been updated, false if new record was created
func (cl *Line) Store(key ListKeyType, value interface{}) bool {
	return cl.StoreFor(key, value, cl.defaultExpirationTime)
}

// RenewFor renews a <key> if it exists shifting expiration further in time without changing the value
// @returns true if record has been updated false otherwise
func (cl *Line) RenewFor(key ListKeyType, expires time.Duration) bool {
	curr, ok := cl.store[key]
	if ok {
		curr.ptr.expiration = time.Now().Add(expires)

		if cl.last != curr.ptr {
			if curr.ptr.prev != nil {
				curr.ptr.prev.next = curr.ptr.next
			} else {
				cl.first = curr.ptr.next
			}

			curr.ptr.next.prev = curr.ptr.prev
			cl.last.next = curr.ptr
			curr.ptr.prev = cl.last
			curr.ptr.next = nil
			cl.last = curr.ptr
		}
	}
	cl.expire()
	return ok
}

// Renew renews a <key> if it exists shifting expiration further in time without changing the value (using default expiration timeout)
// @returns true if record has been updated false otherwise
func (cl *Line) Renew(key ListKeyType) bool {
	return cl.RenewFor(key, cl.defaultExpirationTime)
}

// Check checks if the element is still present in the cache
func (cl *Line) Check(key ListKeyType) bool {
	cl.expire()
	_, ok := cl.store[key]
	return ok
}

// Get retrieves the element from the cache
func (cl *Line) Get(key ListKeyType) (interface{}, error) {
	cl.expire()
	val, ok := cl.store[key]

	if !ok {
		return nil, ErrCacheRecordNotFound
	}

	return val.value, nil
}

func (cl *Line) expire() {
	now := time.Now()
	if now.After(cl.nextCheckupTime) {
		cl.nextCheckupTime = now.Add(cl.defaultExpirationTime)

		if cl.first != nil {
			curr := cl.first
			for {
				if curr.expiration.After(now) {
					break
				}

				delete(cl.store, curr.key)

				if curr.next != nil {
					curr.next.prev = nil
					curr = curr.next
					cl.first = curr
				} else {
					cl.first = nil
					cl.last = nil
					break
				}
			}
		}
	}
}
