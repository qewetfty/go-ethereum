package dpos

type Delegate struct{
	Address string
	Normal bool
	Vote int64 //投票数
	Nickname string // delegate name
}
