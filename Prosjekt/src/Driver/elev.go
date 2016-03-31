package Driver

import (
	"fmt"
)

const MOTOR_SPEED = 2800
const NUMFLOORS = 4

var buttonChannels = [NUMFLOORS][3]int{
	[3]int{BUTTON_DOWN1, BUTTON_UP1, BUTTON_COMMAND1},
	[3]int{BUTTON_DOWN2, BUTTON_UP2, BUTTON_COMMAND2},
	[3]int{BUTTON_DOWN3, BUTTON_UP3, BUTTON_COMMAND3},
	[3]int{BUTTON_DOWN4, BUTTON_UP4, BUTTON_COMMAND4},
}

var lightChannels = [NUMFLOORS][3]int{
	[3]int{LIGHT_DOWN1, LIGHT_UP1, LIGHT_COMMAND1},
	[3]int{LIGHT_DOWN2, LIGHT_UP2, LIGHT_COMMAND2},
	[3]int{LIGHT_DOWN3, LIGHT_UP3, LIGHT_COMMAND3},
	[3]int{LIGHT_DOWN4, LIGHT_UP4, LIGHT_COMMAND4},
}

func ElevSetMotorDirection(direction int) {
	if direction == 0 {
		IoWriteAnalog(MOTOR, 0)
	} else if direction > 0 {
		IoClearBit(MOTORDIR)
		IoWriteAnalog(MOTOR, 2800)
	} else if direction < 0 {
		IoSetBit(MOTORDIR)
		IoWriteAnalog(MOTOR, 2800)
	}
}

func ElevInit() {
	initCheck := IoInit()
	if initCheck == 0 {
		fmt.Println("Unable to initialize elevator")
	}

	for floor := 0; floor < NUMFLOORS; floor++ {
		for button := 0; button < 3; button++ {
			ElevSetButtonLight(button, floor, 0)
		}
	}

}

func ElevSetButtonLight(button int, floor int, value int) { //We want assert?
	if floor >= 0 && floor < NUMFLOORS {
		if button >= 0 && button < 3 {
			if value > 0 {
				IoSetBit(lightChannels[floor][button])
			} else {
				IoClearBit(lightChannels[floor][button])
			}
		}
	}
}

func ElevSetFloorIndicator(floor int) {
	if floor >= 0 && floor < NUMFLOORS {
		switch floor {
		case 0:
			IoClearBit(LIGHT_FLOOR_IND1)
			IoClearBit(LIGHT_FLOOR_IND2)
		case 1:
			IoSetBit(LIGHT_FLOOR_IND2)
			IoClearBit(LIGHT_FLOOR_IND1)
		case 2:
			IoClearBit(LIGHT_FLOOR_IND2)
			IoSetBit(LIGHT_FLOOR_IND1)
		case 3:
			IoSetBit(LIGHT_FLOOR_IND1)
			IoSetBit(LIGHT_FLOOR_IND2)
		}
	}
}

func ElevSetDoorLight(value int) {
	if value > 0 {
		IoSetBit(LIGHT_DOOR_OPEN)
	} else {
		IoClearBit(LIGHT_DOOR_OPEN)
	}
}

func ElevGetButtonSignal(button int, floor int) int { //We want assert?
	if floor >= 0 && floor < NUMFLOORS {
		if button >= 0 && button < 3 {
			if IoReadBit(buttonChannels[floor][button]) > 0 {
				return 1
			} else {
				return 0
			}
		}
	}
	return 0
}

func ElevGetFloorSensorSignal() int {
	if IoReadBit(SENSOR_FLOOR1) > 0 {
		return 0
	} else if IoReadBit(SENSOR_FLOOR2) > 0 {
		return 1
	} else if IoReadBit(SENSOR_FLOOR3) > 0 {
		return 2
	} else if IoReadBit(SENSOR_FLOOR4) > 0 {
		return 3
	} else {
		return -1
	}
}
