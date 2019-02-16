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
  expression Expression
  expressions []Expression
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

%type <value> value
%type <expression> expression
%type <expressions> expressions
%type <function> main
%type <program> program

%type <str> TK_PRINT

%start program

%%

program: main
         {
           setProgram(yylex, Program{Main: $1})
         }

main: TK_MAIN_OPEN '\n' expressions '\n' TK_MAIN_CLOSE
      {
        $$ = Function{Name: "", Arguments: []string{}, Expressions: $3}
      }

expressions: expression
             {
               $$ = []Expression{$1}
             }
           | expressions '\n' expression
             {
               $1 = append($1, $3)
               $$ = $1
             }

expression: TK_PRINT value
            {
              $$ = Expression{$1, []Value{$2}}
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
