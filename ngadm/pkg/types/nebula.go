package types

type NebulaServiceComponent string

const (
	Metad            NebulaServiceComponent = "metad"
	Graphd           NebulaServiceComponent = "graphd"
	Storaged         NebulaServiceComponent = "storaged"
	AllNebulaSerivce NebulaServiceComponent = "all"
)

func (n NebulaServiceComponent) String() string {
	return string(n)
}

var NebulaServiceComponentMap = map[string]NebulaServiceComponent{
	"metad":       Metad,
	"graphd":      Graphd,
	"storaged":    Storaged,
	"nebulagraph": AllNebulaSerivce,
	"all":         AllNebulaSerivce,
}
