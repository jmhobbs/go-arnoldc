%{
package arnoldc

func setProgram(l yyLexer, v Program) {
  l.(*ArnoldC).program = v
}
%}

%union{
  str string 
  integer int
  value Value
  statement Statement
  statements []Statement
  expression Expression 
  expressions []Expression
  block Block
  blocks []Block
  function Function
  functions []Function
  program Program
}

%token LexError

%token <str> Variable
%token <str> String
%token <integer> Integer

%token TK_MAIN_OPEN TK_MAIN_CLOSE
%token TK_PRINT
%token TK_TRUE TK_FALSE
%token TK_DECLARE TK_INITIALIZE
%token TK_ASSIGNMENT TK_FIRST_OPERAND TK_ASSIGNMENT_END
%token TK_ADD TK_SUBTRACT TK_MULTIPLY TK_DIVIDE
%token TK_IF TK_ELSE TK_END_IF
%token TK_WHILE TK_END_WHILE

%type <value> value
%type <statements> statements
%type <expression> expression
%type <block> block
%type <function> main
%type <program> program

%type <str> TK_PRINT TK_DECLARE TK_INITIALIZE TK_ASSIGNMENT

%start program

%%

program: main
         {
           setProgram(yylex, Program{Main: $1})
         }

main: TK_MAIN_OPEN statements TK_MAIN_CLOSE
      {
        $$ = Function{Name: "", Arguments: []string{}, Statements: $2}
      }

statements: expression
             {
               $$ = []Statement{$1}
             }
           | block
             {
               $$ = []Statement{$1}
             }
           | statements expression
             {
               $1 = append($1, $2)
               $$ = $1
             }
           | statements block
             {
               $1 = append($1, $2)
               $$ = $1
             }

expression: TK_PRINT value
            {
              $$ = Expression{$1, []Value{$2}}
            }
          | TK_DECLARE Variable TK_INITIALIZE Integer
            {
              $$ = Expression{$1, []Value{VariableValue{$2}, IntegerValue{$4}}}
            }

block: TK_ASSIGNMENT Variable TK_FIRST_OPERAND Integer TK_ASSIGNMENT_END
       {
         $$ = Block{$1, []Value{VariableValue{$1}}, []Statement{}}
       }

value: String
       {
         $$ = StringValue{$1}
       }
     | Variable
       {
         $$ = VariableValue{$1}
       }
     | Integer
       {
         $$ = IntegerValue{$1}
       }
     | TK_TRUE
       {
         $$ = BoolValue{true}
       }
     | TK_FALSE
       {
         $$ = BoolValue{false}
       }
