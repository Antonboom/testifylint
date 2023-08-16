// Package analysisutil contains functions common for `analyzer` and `internal/checkers` packages.
// In addition, it is intended to "lighten" these packages.
//
// If the function is common to several packages, or it makes sense to test it separately without
// "polluting" the main package with tests of private functionality, then you should put function in this package.
package analysisutil
