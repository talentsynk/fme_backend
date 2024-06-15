package myuser

// check if the user1 has general access
func CanSuspendActivate(user1, instance *User) bool {
	if (user1.Role == 1) && (instance.Role == 2 || instance.Role == 3 || instance.Role == 4 || instance.Role == 5) {
		return true
	} else {
		return false
	}

}
