package arnoldc

var (
	// This is a map of valid instructions for ArnoldC to their token.
	instructions map[string]int
	// This is a map of valid "macros" for ArnoldC to their token.
	macros map[string]int
)

func init() {
	macros = map[string]int{
		"NO PROBLEMO": TK_TRUE,
		"I LIED":      TK_FALSE,
	}

	instructions = map[string]int{
		// Functions
		"IT'S SHOWTIME":            TK_MAIN_OPEN,
		"YOU HAVE BEEN TERMINATED": TK_MAIN_CLOSE,

		// Simple built ins
		"TALK TO THE HAND": TK_PRINT,

		// Variable declaration
		"HEY CHRISTMAS TREE": TK_DECLARE,
		"YOU SET US UP":      TK_INITIALIZE,

		// Munging
		"GET TO THE CHOPPER":    TK_ASSIGNMENT,
		"ENOUGH TALK":           TK_ASSIGNMENT_END,
		"HERE IS MY INVITATION": TK_FIRST_OPERAND,

		// Arithmetic
		"GET UP":          TK_ADD,
		"GET DOWN":        TK_SUBTRACT,
		"YOU'RE FIRED":    TK_MULTIPLY,
		"HE HAD TO SPLIT": TK_DIVIDE,

		// Logic
		"YOU ARE NOT YOU YOU ARE ME": TK_EQUAL_TO,
		"LET OFF SOME STEAM BENNET":  TK_GREATER_THAN,
		"CONSIDER THAT A DIVORCE":    TK_OR,
		"KNOCK KNOCK":                TK_AND,

		// If/Else
		"BECAUSE I'M GOING TO SAY PLEASE": TK_IF,
		"BULLSHIT":                        TK_ELSE,
		"YOU HAVE NO RESPECT FOR LOGIC":   TK_END_IF,

		// While loops
		"STICK AROUND": TK_WHILE,
		"CHILL":        TK_END_WHILE,
	}
}
