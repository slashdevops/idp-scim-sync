// Package repository persists the synchronization state.
//
// It provides two implementations of core.StateRepository:
//
//   - [S3Repository], used in production, which stores the state as a single
//     JSON object in an S3 bucket.
//   - [DiskRepository], which reads and writes any io.ReadWriter, for local runs
//     and tests.
//
// There is no database. The entire persisted state of the system is one JSON
// document, described in docs/State-File-example.md.
//
// # Absent state is not an error
//
// The first run of a fresh deployment finds no object. internal/core treats both
// s3types.NoSuchKey and [ErrStateFileEmpty] as "start from scratch" rather than
// as failures, so the errors this package returns are classified with
// errors.AsType and must stay distinguishable — which is why they are concrete
// types carrying an ErrorCode rather than opaque values.
package repository
