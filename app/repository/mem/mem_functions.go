package mem

var idSeq int = 0

func generateID() int {
	idSeq++
	return idSeq
}
