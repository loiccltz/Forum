package backend

func IsAdmin(user *User) bool {
	return user.Role == "admin"
}

func IsModerator(user *User) bool {
	return user.Role == "moderator"
}
