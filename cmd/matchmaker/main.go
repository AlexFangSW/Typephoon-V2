/*
The match maker works by providing 2 APIs, one for "queue in",
and another to poll the result.

It is designed with "Active-Passive" architecture,
a match maker becomes active once it aquires the lock in the external cache.

APIs:
- POST /api/v1/matchmaker/queue-in
- GET /api/v1/matchmaker/status/<ID>

External Services:
- Cache (Redis): For active lock
*/

package main

import "fmt"

func main() {
	fmt.Print("Hello from the match maker")
}
