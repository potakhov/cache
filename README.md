# Go cache

Package implements a basic cache that operates with generalized O(1) complexity for all possible (lookup, insert and expiration) operations.

Keys are always strings, values are arbitrary interface{} pointers.

Storage is not thread safe and should be properly guarded in case of shared usage.

Renew() and Store() methods update first and then expire so it is possible to renew a record that should be expired by now, this is done by design.