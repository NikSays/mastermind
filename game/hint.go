package game

func GiveHint(input [4]int, real [4]int) (hit int, near int) {
	hit = 0
	near = 0
	for i := 0; i < 4; i++ {
		if input[i] == real[i] {
			hit++
			input[i] = -1
			real[i] = -1
		}
	}
	for i := 0; i < 4; i++ {
		if input[i] == -1 {
			continue
		}
		for j := 0; j < 4; j++ {
			if real[j] == -1 {
				continue
			}
			if input[i] == real[j] {
				near++
				input[i] = -1
				real[j] = -1
			}
		}
	}
	return
}
