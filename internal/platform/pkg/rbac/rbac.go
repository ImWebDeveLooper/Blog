package rbac

type Role string

const (
	RoleAdmin     Role = "Admin"
	RoleManager   Role = "Manager"
	RoleEditor    Role = "Editor"
	RoleAuthor    Role = "Author"
	RoleModerator Role = "Moderator"
	RoleMember    Role = "Member"
	RoleGuest     Role = "Guest"
)

func (r Role) String() string {
	return string(r)
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
