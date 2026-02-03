Example 3: Loan Service ( system design and abstraction)
we are building a loan engine. A loan can a multiple state: proposed , approved,
invested, disbursed. the rule of state:
1. proposed is the initial state (when loan created it will has proposed state):
2. approved is once it approved by our staff.
1. a approval must contains several information:
1. the picture proof of the a field validator has visited the borrower
2. the employee id of field validator
3. date of approval
2. once approved it can not go back to proposed state
3. once approved loan is ready to be offered to investors/lender
3. invested is once total amount of invested is equal the loan principal
1. loan can have multiple investors, each with each their own amount

2. total of invested amount can not be bigger than the loan principal amount
3. once invested all investors will receive an email containing link to
agreement letter (pdf)

2. disbursed is when is loan is given to borrower.
1. a disbursement must contains several information:
1. the loan agreement letter signed by borrower (pdf/jpeg)
2. the employee id of the field officer that hands the money and/or
collect the agreement letter
3. date of disbursement

movement between state can only move forward, and a loan only need following
information:
● borrower id number
● principal amount
● rate, will define total interest that borrower will pay
● ROI return of investment, will define total profit received by investors
● link to the generated agreement letter
design a RESTFful api that satisfy above requirement.