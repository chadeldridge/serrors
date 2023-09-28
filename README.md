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
    // TextHandler with keys to uppercase HandlerOption
	tErrs := serrors.NewTextHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: serrors.UpperCaseKey})

    // JSONHandler in a struct with no HandlerOptions
	jErrs := struct {
		Errors serrors.SErrors `json:"errors"`
	}{
		serrors.New(os.Stdout, nil),
	}

	errs := doStuff()

	fmt.Println("Text Logging")
	tErrs.Append(errs)
	tErrs.Log()

	fmt.Println("\nJSON Logging")
	jErrs.Errors.Append(errs)
	jErrs.Errors.Log()

	fmt.Println("\nJSON object")
	j, _ := json.Marshal(jErrs)
	fmt.Println(string(j))
}

func doStuff() serrors.SErrors {
	errs := doMore()
	if !errs.IsEmpty() {
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
Text Logging
TIME=2023-09-28T13:48:50.220-04:00 LEVEL=ERROR MSG="error was inevitable" FAILED=true CODE=500
TIME=2023-09-28T13:48:50.220-04:00 LEVEL=WARN MSG="doMore failed to do more" FAILED=true CODE=500

JSON Logging
{"time":"2023-09-28T13:48:50.220545992-04:00","level":"ERROR","msg":"error was inevitable","failed":true,"code":500}
{"time":"2023-09-28T13:48:50.220547194-04:00","level":"WARN","msg":"doMore failed to do more","failed":true,"code":500}

JSON object
{"errors":[{"time":"2023-09-28T13:48:50.220545992-04:00","level":"ERROR","msg":"error was inevitable","failed":true,"code":500},{"time":"2023-09-28T13:48:50.220547194-04:00","level":"WARN","msg":"doMore failed to do more","failed":true,"code":500}]}
```