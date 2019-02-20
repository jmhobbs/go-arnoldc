package arnoldc

// This is a map of valid instructions for ArnoldC to their token.
var instructions map[string]int

func init() {
	instructions = map[string]int{
		"IT'S SHOWTIME":                   TK_MAIN_OPEN,
		"YOU HAVE BEEN TERMINATED":        TK_MAIN_CLOSE,
		"TALK TO THE HAND":                TK_PRINT,
		"HEY CHRISTMAS TREE":              TK_DECLARE,
		"YOU SET US UP":                   TK_INITIALIZE,
		"GET TO THE CHOPPER":              TK_ASSIGNMENT,
		"HERE IS MY INVITATION":           TK_FIRST_OPERAND,
		"GET UP":                          TK_ADD,
		"GET DOWN":                        TK_SUBTRACT,
		"YOU'RE FIRED":                    TK_MULTIPLY,
		"HE HAD TO SPLIT":                 TK_DIVIDE,
		"BECAUSE I'M GOING TO SAY PLEASE": TK_IF,
		"BULLSHIT":                        TK_ELSE,
		"YOU HAVE NO RESPECT FOR LOGIC":   TK_END_IF,
		"STICK AROUND":                    TK_WHILE,
		"CHILL":                           TK_END_WHILE,
	}
}
