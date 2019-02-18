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
%token TK_DECLARE TK_INITIALIZE

%type <value> value
%type <expression> expression
%type <expressions> expressions
%type <function> main
%type <program> program

%type <str> TK_PRINT TK_DECLARE TK_INITIALIZE

%start program

%%

program: main
         {
           setProgram(yylex, Program{Main: $1})
         }

main: TK_MAIN_OPEN newlines expressions newlines TK_MAIN_CLOSE
      {
        $$ = Function{Name: "", Arguments: []string{}, Expressions: $3}
      }

expressions: expression
             {
               $$ = []Expression{$1}
             }
           | expressions newlines expression
             {
               $1 = append($1, $3)
               $$ = $1
             }

expression: TK_PRINT value
            {
              $$ = Expression{$1, []Value{$2}}
            }
          | TK_DECLARE Variable newlines TK_INITIALIZE Integer
            {
              $$ = Expression{$1, []Value{VariableValue{$2}, IntegerValue{$5}}}
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

newlines: '\n'
        | '\r' '\n'
        | newlines '\n'
        | newlines '\r' '\n'
