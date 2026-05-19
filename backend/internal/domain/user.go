package domain

// User is the pure-domain projection of a user, regardless of where
// the row physically lives (Better Auth's `user` table at the time of
// writing). Authentication and persistence are external concerns —
// this type carries the fields the application reads.
type User struct {
	ID    UserID
	Email string
	Name  string
}
