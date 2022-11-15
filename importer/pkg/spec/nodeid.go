package spec

type (
	NodeID struct {
		Prop *Prop `yaml:",inline"`
	}
)

func (id *NodeID) Complete() {
	if id.Prop == nil {
		id.Prop = &Prop{}
	}
	id.Prop.Complete()
}

func (id *NodeID) Validate() error {
	//revive:disable-next-line:if-return
	if err := id.Prop.Validate(); err != nil {
		return err
	}
	return nil
}

func (id *NodeID) ValueStatement(record Record) (string, error) {
	props := Props{id.Prop}
	return props.ValueStatement(record)
}
