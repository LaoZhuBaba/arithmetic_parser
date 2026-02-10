### A recursive descent parser for infix arithment

* Supports addition (+), subtraction (-), multiplication (*) division (/) and exponentiation (^)  with correct order for evaluation
* All operators are left associative except for exponentiation which associates from right to left.  E.g., 2^2^3 is evaluated as 2^(2^3)
* In this version a lot of logic is implemented in types.go so you could modify the code for other types of evaluation that uses
  infix operators and the same concepts of precedence, associativity and parentheses.
* Parentheses are interpreted correctly and spaces between tokens are ignored.
* Currently only supports 32 bit integer arithmetic but it would be pretty trivial to modify to support floating point.  You would
  just need to update the lexer to interpret floating point strings and change the type in the parser.
* All the tests pass so it seems to be working :-)
