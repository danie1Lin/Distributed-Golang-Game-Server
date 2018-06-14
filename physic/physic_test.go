package physic

import (
	//"github.com/gazed/vu/math/lin"
	"fmt"
	"testing"
)

func TestWorld(t *testing.T) {
	fmt.Printf("Bits:%b", SetBitExcept(SetAllBits(), Skill_Bit, Terrain_Bit))
	//
}
