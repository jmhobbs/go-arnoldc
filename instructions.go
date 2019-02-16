package arnoldc

// This is a map of valid instructions for ArnoldC to their token.
var instructions map[string]int

func init() {
	instructions = map[string]int{
		"IT'S SHOWTIME":            TK_MAIN_OPEN,
		"YOU HAVE BEEN TERMINATED": TK_MAIN_CLOSE,
		"TALK TO THE HAND":         TK_PRINT,
	}
}
