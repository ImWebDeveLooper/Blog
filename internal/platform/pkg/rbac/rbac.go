package rbac

type Role string

const (
	RoleAdmin     Role = "Admin"
	RoleManager   Role = "Manager"
	RoleEditor    Role = "Editor"
	RoleAuthor    Role = "Author"
	RoleModerator Role = "Moderator"
	RoleGuest     Role = "Guest"
)

func (s Role) String() string {
	return string(s)
}

type Object string

const (
	ObjectUser Object = "User"
	ObjectPost Object = "Post"
)

func (o Object) String() string {
	return string(o)
}

type Action string

const (
	ActionRead   Action = "Read"
	ActionWrite  Action = "Write"
	ActionModify Action = "Modify"
	ActionDelete Action = "Delete"
)

func (a Action) String() string {
	return string(a)
}
