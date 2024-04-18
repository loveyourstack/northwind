// sysrole represents enum system.role
package sysrole

const (
	Junior   string = "Junior"
	Senior   string = "Senior"
	Tech     string = "Tech"
	UserMgmt string = "User Mgmt"
)

var (
	GeneralActions    = [...]string{Junior, Senior, Tech}
	RestrictedActions = [...]string{Senior, Tech}
	UserMgmtActions   = [...]string{UserMgmt, Tech}
)
