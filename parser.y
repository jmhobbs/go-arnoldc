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
  function    Method
  functions   []Method
  program     Program
}

%token LexError

%token <str> Variable
%token <str> String
%token <integer> Integer

%token MAIN_OPEN MAIN_CLOSE
%token METHOD_OPEN METHOD_CLOSE DECLARE_PARAMETER END_PARAMETER_DECLARATION RETURN
%token CALL_METHOD ASSIGN_FROM_CALL
%token PRINT
%token TRUE FALSE
%token DECLARE INITIALIZE
%token ASSIGNMENT FIRST_OPERAND ASSIGNMENT_END
%token ADD SUBTRACT MULTIPLY DIVIDE
%token EQUAL_TO GREATER_THAN OR AND
%token IF ELSE END_IF IF_TRUE IF_FALSE
%token WHILE END_WHILE

%type <integer> arithmetic_token
%type <strs> parameters
%type <value> value number
%type <values> values
%type <statements> statements arithmetics
%type <expression> expression arithmetic invoke
%type <block> block
%type <function> main method
%type <functions> methods
%type <program> program

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

main: MAIN_OPEN statements MAIN_CLOSE
      {
        $$ = Method{Name: "", Statements: $2}
      }

methods: method
         {
           $$ = []Method{$1}
         }
       | methods method
         {
           $1 = append($1, $2)
           $$ = $1
         }

method: METHOD_OPEN Variable statements METHOD_CLOSE
        {
          $$ = Method{$2, []string{}, $3}
        }
      | METHOD_OPEN Variable parameters END_PARAMETER_DECLARATION statements METHOD_CLOSE
        {
          $$ = Method{$2, $3, $5}
        }

parameters: DECLARE_PARAMETER Variable
            {
              $$ = []string{$2}
            }
          | parameters DECLARE_PARAMETER Variable
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

expression: PRINT value
            {
              $$ = Expression{PRINT, []Value{$2}}
            }
          | DECLARE Variable INITIALIZE number
            {
              $$ = Expression{DECLARE, []Value{VariableValue{$2}, $4}}
            }
          | RETURN value
            {
              $$ = Expression{RETURN, []Value{$2}}
            }
          | invoke

invoke: CALL_METHOD Variable
        {
          $$ = Expression{CALL_METHOD, []Value{VariableValue{$2}}}
        }
      | CALL_METHOD Variable values
        {
          $3 = append([]Value{VariableValue{$2}}, $3...)
          $$ = Expression{CALL_METHOD, $3}
        }
      | ASSIGN_FROM_CALL Variable CALL_METHOD Variable
        {
          $$ = Expression{ASSIGN_FROM_CALL, []Value{VariableValue{$2}, VariableValue{$4}}}
        }
      | ASSIGN_FROM_CALL Variable CALL_METHOD Variable values
        {
          $5 = append([]Value{VariableValue{$2}, VariableValue{$4}}, $5...)
          $$ = Expression{ASSIGN_FROM_CALL, $5}
        }

block: ASSIGNMENT Variable FIRST_OPERAND number arithmetics ASSIGNMENT_END
       {
         $$ = Block{ASSIGNMENT, []Value{VariableValue{$2}}, append([]Statement{Expression{FIRST_OPERAND, []Value{$4}}}, $5...)}
       }
     | IF number statements END_IF
       {
         $$ = Block{IF, []Value{$2}, []Statement{Block{IF_TRUE, []Value{}, $3}}}
       }
     | IF number statements ELSE statements END_IF
       {
         $$ = Block{IF, []Value{$2}, []Statement{Block{IF_TRUE, []Value{}, $3}, Block{IF_FALSE, []Value{}, $5}}}
       }
     | WHILE number statements END_WHILE
       {
         $$ = Block{WHILE, []Value{$2}, $3}
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

arithmetic_token: SUBTRACT
                  {
                    $$ = SUBTRACT
                  }
                | ADD
                  {
                    $$ = ADD
                  }
                | MULTIPLY
                  {
                    $$ = MULTIPLY
                  }
                | DIVIDE
                  {
                    $$ = DIVIDE
                  }
                | EQUAL_TO
                  {
                    $$ = EQUAL_TO
                  }
                | GREATER_THAN
                  {
                    $$ = GREATER_THAN
                  }
                | OR
                  {
                    $$ = OR
                  }
                | AND
                  {
                    $$ = AND
                  }

number: Variable
         {
           $$ = VariableValue{$1}
         }
       | Integer
         {
           $$ = IntegerValue{$1}
         }
       | TRUE
         {
           $$ = IntegerValue{1}
         }
       | FALSE
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
