package sggj

const (
	Castle = iota //城堡
	SMITHY        //铁匠铺
	FARM          //农场
	TRAP          //陷阱
)

type City struct {
	Name     string
	Gold     uint32
	Food     uint32
	Wood     uint32
	Building [4]uint32
}

func (c *City) get() {

}

func (c *City) build() {

}
