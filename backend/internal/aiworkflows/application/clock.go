package application

import "time"

// nowFn is overridable in tests so deterministic timestamps can be
// asserted. Production callers leave it pointing at time.Now.
var nowFn = time.Now
