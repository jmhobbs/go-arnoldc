%{
package arnoldc

func setProgram(l yyLexer, v Program) {
  l.(*ArnoldC).program = v
}
%}

%union{
  str         string
  strs        []string
  integer     int
  value       Value
  values      []Value
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
%token TK_METHOD_OPEN TK_METHOD_CLOSE TK_DECLARE_PARAMETER TK_END_PARAMETER_DECLARATION TK_RETURN
%token TK_CALL_METHOD TK_ASSIGN_FROM_CALL
%token TK_PRINT
%token TK_TRUE TK_FALSE
%token TK_DECLARE TK_INITIALIZE
%token TK_ASSIGNMENT TK_FIRST_OPERAND TK_ASSIGNMENT_END
%token TK_ADD TK_SUBTRACT TK_MULTIPLY TK_DIVIDE
%token TK_EQUAL_TO TK_GREATER_THAN TK_OR TK_AND
%token TK_IF TK_ELSE TK_END_IF
%token TK_WHILE TK_END_WHILE

%type <str> arithmetic_token
%type <strs> parameters
%type <value> value number
%type <values> values
%type <statements> statements arithmetics
%type <expression> expression arithmetic invoke
%type <block> block
%type <function> main method
%type <functions> methods
%type <program> program

%type <str> TK_PRINT
%type <str> TK_DECLARE TK_INITIALIZE
%type <str> TK_ASSIGNMENT TK_FIRST_OPERAND
%type <str> TK_ADD TK_SUBTRACT TK_MULTIPLY TK_DIVIDE
%type <str> TK_EQUAL_TO TK_GREATER_THAN TK_OR TK_AND
%type <str> TK_IF TK_ELSE TK_END_IF
%type <str> TK_WHILE TK_END_WHILE
%type <str> TK_CALL_METHOD TK_ASSIGN_FROM_CALL
%type <str> TK_RETURN

%start program

%%

program: main
         {
           setProgram(yylex, Program{Main: $1})
         }
       | main methods
         {
           setProgram(yylex, Program{Main: $1, Methods: $2})
         }

main: TK_MAIN_OPEN statements TK_MAIN_CLOSE
      {
        $$ = Function{Name: "", Arguments: []string{}, Statements: $2}
      }

methods: method
         {
           $$ = []Function{$1}
         }
       | methods method
         {
           $1 = append($1, $2)
           $$ = $1
         }

method: TK_METHOD_OPEN Variable statements TK_METHOD_CLOSE
        {
          $$ = Function{$2, []string{}, $3}
        }
      | TK_METHOD_OPEN Variable parameters TK_END_PARAMETER_DECLARATION statements TK_METHOD_CLOSE
        {
          $$ = Function{$2, $3, $5}
        }

parameters: TK_DECLARE_PARAMETER Variable
            {
              $$ = []string{$2}
            }
          | parameters TK_DECLARE_PARAMETER Variable 
            {
              $1 = append($1, $3)
              $$ = $1
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
          | TK_RETURN value
            {
              $$ = Expression{$1, []Value{$2}}
            }
          | invoke

invoke: TK_CALL_METHOD Variable
        {
          $$ = Expression{$1, []Value{VariableValue{$2}}}
        }
      | TK_CALL_METHOD Variable values
        {
          $3 = append([]Value{VariableValue{$2}}, $3...)
          $$ = Expression{$1, $3}
        }
      | TK_ASSIGN_FROM_CALL Variable TK_CALL_METHOD Variable
        {
          $$ = Expression{$1, []Value{VariableValue{$2}, VariableValue{$4}}}
        }
      | TK_ASSIGN_FROM_CALL Variable TK_CALL_METHOD Variable values
        {
          $5 = append([]Value{VariableValue{$2}, VariableValue{$4}}, $5...)
          $$ = Expression{$1, $5}
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

values: value
        {
          $$ = []Value{$1}
        }
      | values value
        {
          $1 = append($1, $2)
          $$ = $1
        }

value: String
       {
         $$ = StringValue{$1}
       }
     | number
