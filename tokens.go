package arnoldc

var (
	// This is a map of valid instructions for ArnoldC to their token.
	instructions map[string]int
	// This is a map of valid "macros" for ArnoldC to their token.
	macros map[string]int
)

func init() {
	macros = map[string]int{
		"NO PROBLEMO": TRUE,
		"I LIED":      FALSE,
	}

	instructions = map[string]int{
		// Functions
		"IT'S SHOWTIME":                                      MAIN_OPEN,
		"YOU HAVE BEEN TERMINATED":                           MAIN_CLOSE,
		"LISTEN TO ME VERY CAREFULLY":                        METHOD_OPEN,
		"HASTA LA VISTA, BABY":                               METHOD_CLOSE,
		"I NEED YOUR CLOTHES YOUR BOOTS AND YOUR MOTORCYCLE": DECLARE_PARAMETER,
		"GIVE THESE PEOPLE AIR":                              END_PARAMETER_DECLARATION,
		"I'LL BE BACK":                                       RETURN,
		"DO IT NOW":                                          CALL_METHOD,
		"GET YOUR ASS TO MARS":                               ASSIGN_FROM_CALL,

		// Simple built ins
		"TALK TO THE HAND": PRINT,

		// Variable declaration
		"HEY CHRISTMAS TREE": DECLARE,
		"YOU SET US UP":      INITIALIZE,

		// Munging
		"GET TO THE CHOPPER":    ASSIGNMENT,
		"ENOUGH TALK":           ASSIGNMENT_END,
		"HERE IS MY INVITATION": FIRST_OPERAND,

		// Arithmetic
		"GET UP":          ADD,
		"GET DOWN":        SUBTRACT,
		"YOU'RE FIRED":    MULTIPLY,
		"HE HAD TO SPLIT": DIVIDE,

		// Logic
		"YOU ARE NOT YOU YOU ARE ME": EQUAL_TO,
		"LET OFF SOME STEAM BENNET":  GREATER_THAN,
		"CONSIDER THAT A DIVORCE":    OR,
		"KNOCK KNOCK":                AND,

		// If/Else
		"BECAUSE I'M GOING TO SAY PLEASE": IF,
		"BULLSHIT":                        ELSE,
		"YOU HAVE NO RESPECT FOR LOGIC":   END_IF,

		// While loops
		"STICK AROUND": WHILE,
		"CHILL":        END_WHILE,
	}
}

func instructionToString(instruction int) string {
	for s, tk := range instructions {
		if tk == instruction {
			return s
		}
	}
	return "UNKNOWN"
}
