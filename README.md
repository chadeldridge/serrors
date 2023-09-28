# serrors
Structured error stack using "log.slog".


## Contents
- [serrors](#serrors)
  - [Contents](#contents)
  - [Installation](#installation)
  - [Quick Start](#quick-start)

## Installation
To install serrors you must first have [Go](https://golang.org/) installed and setup.
1. Install serrors module.
```ssh
$ go get -u github.com/chadeldridge/serrors
```
2. Import serrors in your code:
```go
import "github.com/chadeldridge/serrors"
```

## Quick Start
```go
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/chadeldridge/serrors"
)

func main() {
	jErrs := struct {
		Errors serrors.SErrors `json:"errors"`
	}{
		Errors: serrors.New(os.Stdout, nil),
	}

	tErrs := struct {
		Errors serrors.SErrors
	}{
		Errors: serrors.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: serrors.UpperCaseKey}),
	}

	errs := doStuff()

	fmt.Println("JSON object")
	jErrs.Errors.Append(errs)
	j, _ := json.Marshal(jErrs)
	fmt.Println(string(j))

	fmt.Println("\nJSON Logging")
	jErrs.Errors.Log()

	fmt.Println("\nText Logging")
	tErrs.Errors.Append(errs)
	tErrs.Errors.Log()
}

func doStuff() serrors.SErrors {
	errs := doMore()
	// if !errs.IsEmpty() {
	if len(errs.Errors) > 0 {
		errs.WarnAny(time.Now(), "doMore failed to do more", "failed", true, "code", 500)
	}

	return errs
}

func doMore() serrors.SErrors {
	errs := serrors.New(os.Stdout, nil)
	errs.Add(time.Now(), slog.LevelError, "error was inevitable", slog.Bool("failed", true), slog.Int("code", 500))
	return errs
}
```

Output
```
JSON object
{"errors":[{"time":"2023-09-28T13:41:22.546518213-04:00","level":"ERROR","msg":"error was inevitable","failed":true,"code":500},{"time":"2023-09-28T13:41:22.546519365-04:00","level":"WARN","msg":"doMore failed to do more","failed":true,"code":500}]}

JSON Logging
{"time":"2023-09-28T13:41:22.546518213-04:00","level":"ERROR","msg":"error was inevitable","failed":true,"code":500}
{"time":"2023-09-28T13:41:22.546519365-04:00","level":"WARN","msg":"doMore failed to do more","failed":true,"code":500}

Text Logging
TIME=2023-09-28T13:41:22.546-04:00 LEVEL=ERROR MSG="error was inevitable" FAILED=true CODE=500
TIME=2023-09-28T13:41:22.546-04:00 LEVEL=WARN MSG="doMore failed to do more" FAILED=true CODE=500
```