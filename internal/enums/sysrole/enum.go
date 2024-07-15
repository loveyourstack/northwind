// sysrole represents enum system.role
package sysrole

type Enum string

func (e Enum) String() string {
	return string(e)
}

const (
	Junior   Enum = "Junior"
	Senior   Enum = "Senior"
	Tech     Enum = "Tech"
	UserMgmt Enum = "User Mgmt"
)

var (
	GeneralActions    = [...]Enum{Junior, Senior, Tech}
	RestrictedActions = [...]Enum{Senior, Tech}
	UserMgmtActions   = [...]Enum{UserMgmt, Tech}
)
