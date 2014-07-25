// Completed Sun Sep 22 14:36:10 CDT 2013
//
// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[3], respectively.)

// Put your code here.
  @R2
  M=0

  //If R0 is zero, end program
  @R0
  D=M
  @END
  D;JEQ

(ADDAGAIN)
  //Check if R1 is 0, meaning no more iterations
  @R1
  D=M
  @END
  D;JEQ
  
  //R2 += R0
  @R0
  D=M
  @R2
  M=D+M
  
  @R1
  M=M-1
  
  @ADDAGAIN
  0;JMP
(END)
  @END
  0;JMP
