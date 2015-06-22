// Package integrationtest is the integration test of the library.
// This package assumes the existance of a valid
// integrationtest/acd-token.json file with permissions 0600.
//
// The integration uses a real Amazon Cloud Drive account and manipulate files
// under the folder /acd_test_folder. Due to the nature of these tests, it is
// not recommended to run them against an account that has real data on it. The
// ACD team and contributors are not responsible for any data loss due to these
// tests.
package integrationtest
