
## Go-Snowflake - The ID generator is based on the ideas of Twitter Snowflake.

Generates unique numeric 64-bit identifiers. The generator is based on ideas [Twitter Snowflake](https://github.com/twitter/snowflake/).

### Destination

Snowflake is designed to generate unique ordered identifiers on different machines. Supports 1024 machines (0-1023).

### Details

The 64-bit identifier consists of:

* 41bit contains a timestamp in milliseconds
* 10bit contains the number of the machine to generate
* 12bit contains the execution number in one millisecond

### Example

```go
package main

import (
	"fmt"
    "github.com/alexandergrom/go-snowflake"
)

func main() {
	// epoch is the starting point for all identifiers
	var epoch = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
	var gen, _ = snowflake.New(1, epoch)
	fmt.Printf("%d\n", gen.Generate()) // 1058897343283204096
}
```
