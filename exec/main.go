package main

import (
	"fmt"
	"log"
)

//Place the current brick at chosen spot
func placeBrick(x1, y1, x2, y2 int, brickNo int16, second [][]int16) {
	second[x1][y1] = brickNo
	second[x2][y2] = brickNo
}

//Assigning -1 means that the spot is free
func removeBrick(x1, y1, x2, y2 int, second [][]int16) {
	second[x1][y1] = -1
	second[x2][y2] = -1
}

//Assert validity of the coordinates
func withinBounds(x1, y1, x2, y2 int, grid [][]int16) bool {
	return x1 < len(grid) && y1 < len(grid[x1]) &&
		x2 < len(grid) && y2 < len(grid[x2])
}

//When we are trying to place a brick on the second layer,
//this function ensures that we dont have a brick in the exact same position on layer one
func identicalPosition(x1, y1, x2, y2 int, bottomLayer [][]int16) bool {
	return withinBounds(x1, y1, x2, y2, bottomLayer) &&
		bottomLayer[x1][y1] == bottomLayer[x2][y2]
}

//Checks if no brick has been placed in this spot till now
func freeSpot(x1, y1, x2, y2 int, second [][]int16) bool {
	return second[x1][y1] == -1 && second[x2][y2] == -1
}

//Checks boundary conditions and assert there is free space to place a brick
func fits(x1, y1, x2, y2 int, second [][]int16) bool {
	return withinBounds(x1, y1, x2, y2, second) &&
		freeSpot(x1, y1, x2, y2, second)
}

//This indicates whether the current brick we're trying to place overlaps with an already placed brick.
//Such an overlap may happen when we are moving to the next row and there is a vertical brick already placed there
func overlap(x1, y1, x2, y2 int, second [][]int16) bool {
	return withinBounds(x1, y1, x2, y2, second) &&
		second[x1][y1] > 0
}

//Check if one half of the brick is in the bottom-right corner of the layer
func bottomRight(x, y int, grid[][]int16) bool {
	return x == len(grid) - 1 && y == len(grid[x - 1]) - 1
}

//The number of bricks needed for a properly formed grid is (n * m) / 2
func allBricksInPlace(brickNo int16, grid [][]int16) bool {
	return int(brickNo) -1 == (len(grid) * len(grid[0]) / 2)
}


func explore(x1, y1, x2, y2 int, brickNo int16, first [][]int16, second [][]int16) bool {

	//Detects the bottom-right corner of the grid. Here we check if we managed to place all bricks correctly.
	if bottomRight(x2, y2, second) {
		//Check if we can place a brick in this last spot
		if !identicalPosition(x1, y1, x2, y2, first) {
			placeBrick(x1, y1, x2, y2, brickNo, second)
			return true
		}else  {
			//If the last spot is taken and we have placed the number of bricks required -> return true
			return !freeSpot(x1, y1, x2, y2, second) && allBricksInPlace(brickNo, second)
		}
	}

	//The brick we are trying to place overlaps with a brick that is already placed, so we try moving to the right
	if overlap(x1, y1, x2, y2, second) && explore(x2, y2, x2, y2+1, brickNo, first, second) {
		return true
	}

	//Try to place a horizontal brick
	if fits(x1, y1, x2, y2, second) && !identicalPosition(x1, y1, x2, y2, first) {

		placeBrick(x1, y1, x2, y2, brickNo, second)

		//If we reach the end of the current row we move on to the next one, otherwise keep moving right
		if y2 == len(second[x2])-1 && explore(x1+1, 0, x2+1, 1, brickNo+1, first, second) {
			return true
		} else if explore(x1, y1+2, x2, y2+2, brickNo+1, first, second) {
			return true
		}

		removeBrick(x1, y1, x2, y2, second)
	}

	//If we happen to be unable to fit a horizontal brick we try to place the brick in vertical position
	if fits(x1, y1, x2+1, y2-1, second) && !identicalPosition(x1, y1, x2+1, y2-1, first) {
		placeBrick(x1, y1, x2+1, y2-1, brickNo, second)

		//If we have reached the last column, explore the next two rows
		if y2 - 1 == len(second[x1]) - 1  {
			if explore(x2+1, 0, x2+1, 1, brickNo+1, first, second) {
				return true
			} else if explore(x2+2, 0, x2+2, 1, brickNo+1, first, second) {
				return true
			}
		}

		if explore(x2, y2, x2, y2+1, brickNo+1, first, second) {
			return true
		}
		removeBrick(x1, y1, x2+1, y2-1, second)
	}


	return false
}

//Assert proper value of dimensions
func validDimension(x int) bool {
	return x > 1 && x < 100 && x%2 == 0
}

func main()  {
	//Dimensions of the first and second layers
	var n, m int

	for {
		fmt.Println("Please enter two even integers in range [2, 98]")

		_, err := fmt.Scanf("%d %d", &n, &m)
		if err != nil {
			log.Fatal(err)
		}

		if validDimension(n) && validDimension(m) {
			break
		}
	}

	firstLayer := make([][]int16, n)
	for i := range firstLayer {
		firstLayer[i] = make([]int16, m)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			_, err := fmt.Scan(&firstLayer[i][j])
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	var frequency [100]int8
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if frequency[firstLayer[i][j]] == 2  {
				log.Fatalln("Exceeded number of elements", firstLayer[i][j])
			}
			frequency[firstLayer[i][j]]++
		}
	}

	secondLayer := make([][]int16, n)
	for i := 0; i < n; i++ {
		secondLayer[i] = make([]int16, m)
		for j := 0; j < m; j++ {
			secondLayer[i][j] = -1
		}
	}

	fmt.Println("Layer 1")
	for i := range firstLayer {
		fmt.Println(firstLayer[i])
	}

	if explore(0, 0, 0, 1, 1, firstLayer, secondLayer){
		fmt.Println("After construction: Layer 2 ")
		for i := 0; i < n; i++ {
			fmt.Println(secondLayer[i])
		}
	}else {
		fmt.Println("-1")
	}

	fmt.Println("Layer 1")
	for i := range firstLayer {
		fmt.Println(firstLayer[i])
	}
}