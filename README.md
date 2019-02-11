# Go cache

Package implements a basic cache that operates with generalized O(1) complexity for all possible (lookup, insert and expiration) operations.

Keys and values are arbitrary interface{} pointers.

Storage is not thread safe and should be properly guarded in case of shared usage.

Expiration is being run periodically so it is possible to still pull a record that should be expired by now if the expiration period did not happen yet.

The container itself doesn't run an expiration goroutine in background so actual expiration happens during Store(), Renew(), Check() or Get() operations. If required periodic expiration could be triggered from outside by checking for an arbitrary (non-existent) key presence.
