// salesregion represents enum sales.region
package salesregion

type Enum string

func (e Enum) String() string {
	return string(e)
}

const (
	Northern Enum = "Northern"
	Eastern  Enum = "Eastern"
	Southern Enum = "Southern"
	Western  Enum = "Western"
)
