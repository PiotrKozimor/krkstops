package trie

type coord struct {
	x, y int
}

var (
	coordinates = map[rune]coord{
		'1': {0, -1},
		'2': {1, -1},
		'3': {2, -1},
		'4': {3, -1},
		'5': {4, -1},
		'6': {5, -1},
		'7': {6, -1},
		'8': {7, -1},
		'9': {8, -1},
		'0': {9, -1},

		'q': {0, 0},
		'w': {1, 0},
		'e': {2, 0},
		'r': {3, 0},
		't': {4, 0},
		'y': {5, 0},
		'u': {6, 0},
		'i': {7, 0},
		'o': {8, 0},
		'p': {9, 0},

		'a': {0, 1},
		's': {1, 1},
		'd': {2, 1},
		'f': {3, 1},
		'g': {4, 1},
		'h': {5, 1},
		'j': {6, 1},
		'k': {7, 1},
		'l': {8, 1},

		'z': {0, 2},
		'x': {1, 2},
		'c': {2, 2},
		'v': {3, 2},
		'b': {4, 2},
		'n': {5, 2},
		'm': {6, 2},
	}
)

func chebyshevDistance(a, b rune) int {
	aCoord := coordinates[a]
	bCoord := coordinates[b]
	return max(abs(aCoord.x-bCoord.x), abs(aCoord.y-bCoord.y))
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
