%{
package arnoldc

func setProgram(l yyLexer, v Program) {
  l.(*ArnoldC).program = v
}
%}

%union{
  str         string
  integer     int
  value       Value
  statement   Statement
  statements  []Statement
  expression  Expression
  expressions []Expression
  block       Block
  blocks      []Block
  function    Function
  functions   []Function
  program     Program
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
%token TK_EQUAL_TO TK_GREATER_THAN TK_OR TK_AND
%token TK_IF TK_ELSE TK_END_IF
%token TK_WHILE TK_END_WHILE

%type <str> arithmetic_token
%type <value> value number
%type <statements> statements arithmetics
%type <expression> expression arithmetic
%type <block> block
%type <function> main
%type <program> program

%type <str> TK_PRINT
%type <str> TK_DECLARE TK_INITIALIZE
%type <str> TK_ASSIGNMENT TK_FIRST_OPERAND
%type <str> TK_ADD TK_SUBTRACT TK_MULTIPLY TK_DIVIDE
%type <str> TK_EQUAL_TO TK_GREATER_THAN TK_OR TK_AND
%type <str> TK_IF TK_ELSE TK_END_IF
%type <str> TK_WHILE TK_END_WHILE

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
          | TK_DECLARE Variable TK_INITIALIZE number
            {
              $$ = Expression{$1, []Value{VariableValue{$2}, $4}}
            }

block: TK_ASSIGNMENT Variable TK_FIRST_OPERAND number arithmetics TK_ASSIGNMENT_END
       {
         $$ = Block{$1, []Value{VariableValue{$2}}, append([]Statement{Expression{$3, []Value{$4}}}, $5...)}
       }
     | TK_IF number statements TK_END_IF
       {
         $$ = Block{$1, []Value{$2}, []Statement{Block{"__TRUE", []Value{}, $3}}}
       }
     | TK_IF number statements TK_ELSE statements TK_END_IF
       {
         $$ = Block{$1, []Value{$2}, []Statement{Block{"__TRUE", []Value{}, $3}, Block{"__FALSE", []Value{}, $5}}}
       }
     | TK_WHILE number statements TK_END_WHILE
       {
         $$ = Block{$1, []Value{$2}, $3}
       }


arithmetics: arithmetic
             {
               $$ = []Statement{$1}
             }
           | arithmetics arithmetic
             {
               $1 = append($1, $2)
               $$ = $1
             }

arithmetic: arithmetic_token number
            {
              $$ = Expression{$1, []Value{$2}}
            }

arithmetic_token: TK_SUBTRACT
                | TK_ADD
                | TK_MULTIPLY
                | TK_DIVIDE
                | TK_EQUAL_TO
                | TK_GREATER_THAN
                | TK_OR
                | TK_AND

number: Variable
         {
           $$ = VariableValue{$1}
         }
       | Integer
         {
           $$ = IntegerValue{$1}
         }
       | TK_TRUE
         {
           $$ = IntegerValue{1}
         }
       | TK_FALSE
         {
           $$ = IntegerValue{0}
         }

value: String
       {
         $$ = StringValue{$1}
       }
     | number
