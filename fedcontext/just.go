package fedcontext

// Return a pointer pointing to s. This is useful
// in cases where we want to fill pointer fields in
// structs from function results, i.e.
//
//   type Container struct {
//       X *string
//   }
//
//   c := Container{}
//   c.X = Just(fmt.Sprintf("name=%v", name))
func Just(s string) *string {
	return &s
}
