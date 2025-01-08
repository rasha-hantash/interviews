package main

import (
    "fmt"
    "coderpad/echo"
)


// # A 2-D grid consisting of some blocked (represented as '#') and some unblocked (represented as '.') cells is given. The starting position of a pointer is in the top-left corner of the grid. It is guaranteed that the starting position is in an unblocked cell [.]. It is also guaranteed that the bottom-right cell is unblocked [.]. Each cell of the grid is connected with its right, left, top, and bottom cells (if those cells exist). It takes 1 second for a pointer to move from a cell to its adjacent cell. If the pointer can reach the bottom-right corner of the grid within maxTime seconds, return the string 'Yes'. Otherwise, return the string 'No'. 

 
// # Example:

// #     rows = 3
// #     grid = ['..##', 
// #             '..##', 
// #             '#...']
// #     maxTime = 5

// # It will take the pointer 5 seconds to reach the bottom right corner. As long as maxTime ≥ 5, return 'Yes'.

// # Returns:
// #    string: the final string; either 'Yes' or 'No'

 

// # Constraints

// # 1 ≤ rows ≤ 500
// # 0 ≤ maxTime ≤ 106


/*
- isVisited hashmaps [col][row] -> true 
- currMaxTime -> how many seconds has passed (each second represents 1 adjacent move)
- before moving on to next step just check if we are at MaxTIme && are not at bottom right 
- 
- keep track of column and row that we are at (along with currMaxTime)
*/


type CurPosition struct {
	Row int
	Col int 
}

func main() {
    for i := 0; i < 5; i++ {
        fmt.Println(echo.EchoRepeat("Hello, World!"))
    }

		isVisited := make(map[col][row]int)


		curMaxTime := 0 
		// col, row of where we're at 
		// meant to be a tuple 
		


}


func isReachedWithinMaxTime(rows int, maxTime int, grid [][]string, curPosition CurPosition) {

		topRow := 0 
		bottomRow := rows - 1  // 3
		leftMostColumn := 0 
		rightMostColumn := len(grid[0]) -1 // 2



		leftMove := [curPosition.Row][curPosition.Col -1]string 
		rightMove := [curPosition.Row][curPosition.Col +1]string 
		bottomMove := [curPosition.Row-1 ][curPosition.Col]string 
		topMove := [curPosition.Row+1 ][curPosition.Col]string 

		if leftMove[][] >= leftMostColumn   {
					/* move left */ 
		} 

		if rightMove <= rightMostColumn {
			/* move right */ 

		}

		if topMove == leftMostColumn {

		}

		if curPosition.Col == rightMostColumn { 

		}




		if curPosition.Row == 0 && curPosition.Col == 0 && ((curPosition.Row + 1) < rows) {

			/* Check right & check bottom & check rows to make sure we aren't out of bounds */
			

		}







}


