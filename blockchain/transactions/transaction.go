package transactions

type transaction struct {
	ID      []byte
	inputs  []txInput
	outputs []txOutput
}

func createTransaction(outputs []txOutput)
