# Fibonacci Sequences

A go package called fibonacci which generates slices of Fibonacci numbers.

- Step to run before each scenario

## Single Numbers
- The 1st item is 1
- The 3rd item is 2

## Numbers From the Start
- The first 3 items are 1, 1, 2
- The first 10 items are 1, 1, 2, 3, 5, 8, 13, 21, 34, 55

## Numbers From an Offset
- The 4th to the 8th items are 3, 5, 8, 13, 21
- The 5th to the 6th items are 5, 8

## Zeroth Item
- The 0th item is 0

## Negative Numbers
- The -1th item is 1
- The -2nd item is -1
- The -10th to the 0th items are -55, 34, -21, 13, -8, 5, -3, 2, -1, 1, 0
- The -10th to the 10th items are -55, 34, -21, 13, -8, 5, -3, 2, -1, 1, 0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55

## More Examples

Given the following test cases:

 From | To  | Result  
------|-----|---------
 2nd  | 5th | 1,2,3,5 
 3rd  | 6th | 2,3,5,8 

- The <From> to the <To> items are <Result>

## Table Parameters

- This step takes a table parameter

 A     | Table  | Demo  
-------|--------|------
 This  | Will   | Be
 Given | To     | Step

## Undefined Steps

- This will be logged as undefined rather than skipped

---

- Step to run after each scenario