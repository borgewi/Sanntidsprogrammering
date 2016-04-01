package Driver

import (
	"fmt"
)

const(
	MOTOR_SPEED = 2800
	NUMFLOORS = 4
	NUMBUTTONS = 3
)

var buttonChannels = [NUMFLOORS][NUMBUTTONS]int{
	[NUMBUTTONS]int{BUTTON_DOWN1, BUTTON_UP1, BUTTON_COMMAND1},
	[NUMBUTTONS]int{BUTTON_DOWN2, BUTTON_UP2, BUTTON_COMMAND2},
	[NUMBUTTONS]int{BUTTON_DOWN3, BUTTON_UP3, BUTTON_COMMAND3},
	[NUMBUTTONS]int{BUTTON_DOWN4, BUTTON_UP4, BUTTON_COMMAND4},
}

var lightChannels = [NUMFLOORS][3]int{
	[NUMBUTTONS]int{LIGHT_DOWN1, LIGHT_UP1, LIGHT_COMMAND1},
	[NUMBUTTONS]int{LIGHT_DOWN2, LIGHT_UP2, LIGHT_COMMAND2},
	[NUMBUTTONS]int{LIGHT_DOWN3, LIGHT_UP3, LIGHT_COMMAND3},
	[NUMBUTTONS]int{LIGHT_DOWN4, LIGHT_UP4, LIGHT_COMMAND4},
}



func ElevSetMotorDirection(direction int) {
	if direction == 0 {
		io_write_analog(MOTOR, 0)
	} else if direction > 0 {
		io_clear_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	} else if direction < 0 {
		io_set_bit(MOTORDIR)
		io_write_analog(MOTOR, MOTOR_SPEED)
	}
}

func ElevInit() {
	initCheck := io_init()
	if initCheck == false {
		fmt.Println("Unable to initialize elevator")
	}

	for floor := 0; floor < NUMFLOORS; floor++ {
		for button := 0; button < 3; button++ {
			ElevSetButtonLight(button, floor, 0)
		}
	}
	ElevSetFloorIndicator(0)
	ElevSetDoorLight(1)
}

func ElevSetButtonLight(button int, floor int, value int) { //We want assert?
	if floor >= 0 && floor < NUMFLOORS {
		if button >= 0 && button < 3 {
			if value > 0 {
				io_set_bit(lightChannels[floor][button])
			} else {
				io_clear_bit(lightChannels[floor][button])
			}
		}
	}
}

func ElevSetFloorIndicator(floor int) {
	if floor >= 0 && floor < NUMFLOORS {
		switch floor {
		case 0:
			io_clear_bit(LIGHT_FLOOR_IND1)
			io_clear_bit(LIGHT_FLOOR_IND2)
		case 1:
			io_set_bit(LIGHT_FLOOR_IND2)
			io_clear_bit(LIGHT_FLOOR_IND1)
		case 2:
			io_clear_bit(LIGHT_FLOOR_IND2)
			io_set_bit(LIGHT_FLOOR_IND1)
		case 3:
			io_set_bit(LIGHT_FLOOR_IND1)
			io_set_bit(LIGHT_FLOOR_IND2)
		}
	}
}

func ElevSetDoorLight(value int) {
	if value > 0 {
		io_set_bit(LIGHT_DOOR_OPEN)
	} else {
		io_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func ElevGetButtonSignal(button int, floor int) int { //We want assert?
	if floor >= 0 && floor < NUMFLOORS {
		if button >= 0 && button < 3 {
			if io_read_bit(buttonChannels[floor][button]) > 0 {
				fmt.Println("yo")
				return 1
			} else {
				return 0
			}
		}
	}
	return 0
}

func ElevGetFloorSensorSignal() int {
	if io_read_bit(SENSOR_FLOOR1) > 0 {
		return 0
	} else if io_read_bit(SENSOR_FLOOR2) > 0 {
		return 1
	} else if io_read_bit(SENSOR_FLOOR3) > 0 {
		return 2
	} else if io_read_bit(SENSOR_FLOOR4) > 0 {
		return 3
	} else {
		return -1
	}
}
