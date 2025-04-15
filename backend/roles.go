package backend

const (
    RoleGuest      = "guest"
    RoleUser       = "user"
    RoleModerator  = "moderator"
    RoleAdmin      = "admin"
)

func IsGuest(user *User) bool {
    return user == nil || user.Role == RoleGuest
}

func IsUser(user *User) bool {
    return user != nil && (user.Role == RoleUser || user.Role == RoleModerator || user.Role == RoleAdmin)
}

func IsModerator(user *User) bool {
    return user != nil && (user.Role == RoleModerator || user.Role == RoleAdmin)
}

func IsAdmin(user *User) bool {
    return user != nil && user.Role == RoleAdmin
}

func HasPermission(user *User, requiredRole string) bool {
    switch requiredRole {
    case RoleGuest:
        return true
    case RoleUser:
        return IsUser(user)
    case RoleModerator:
        return IsModerator(user)
    case RoleAdmin:
        return IsAdmin(user)
    default:
        return false
    }
}