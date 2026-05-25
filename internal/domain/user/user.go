package user

// User is the user entity
type User struct {
	name         Name
	passwordHash PasswordHash
	role         Role
}

// NewUser creates a new User
func NewUser(name Name, passwordHash PasswordHash, role Role) User {
	return User{
		name:         name,
		passwordHash: passwordHash,
		role:         role,
	}
}

// Name returns the name of the User
func (u User) Name() Name {
	return u.name
}

// PasswordHash returns the password hash of the User
func (u User) PasswordHash() PasswordHash {
	return u.passwordHash
}

// Role returns the role of the User
func (u User) Role() Role {
	return u.role
}

func RehydrateUser(name, hashedPassword, role string) User {
	return User{
		name:         RehydrateName(name),
		passwordHash: RehydratePasswordHash(hashedPassword),
		role:         RehydrateRole(role),
	}
}
