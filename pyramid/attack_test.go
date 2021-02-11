package pyramid

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAttacks_Add(t *testing.T) {
	attacks := Attacks{}

	attacks.Add(Attack{
		Attacker: Player{Name: "P1"},
		Target:   Player{Name: "P2"},
	})

	assert.Equal(t, "P1", attacks.Attacks[0].Attacker.Name)
	assert.Equal(t, "P2", attacks.Attacks[0].Target.Name)
}

func TestAttacks_Len(t *testing.T) {
	attacks := Attacks{}

	attacks.Attacks = []Attack{
		{
			Attacker: Player{},
			Target:   Player{},
		},
		{
			Attacker: Player{},
			Target:   Player{},
		},
		{
			Attacker: Player{},
			Target:   Player{},
		},
	}

	assert.Equal(t, attacks.Len(), 3)
}

func TestAttacks_Remove(t *testing.T) {
	attacks := Attacks{}

	attacks.Attacks = []Attack{
		{
			Attacker: Player{Name: "P1A"},
			Target:   Player{Name: "P1T"},
		},
		{
			Attacker: Player{Name: "P2A"},
			Target:   Player{Name: "P2T"},
		},
		{
			Attacker: Player{Name: "P3A"},
			Target:   Player{Name: "P3T"},
		},
	}

	attackToRemove := Attack{
		Attacker: Player{Name: "P2A"},
		Target:   Player{Name: "P2T"},
	}
	attacks.Remove(attackToRemove)

	assert.Equal(t, attacks.Len(), 2)

	assert.Equal(t, attacks.Attacks[0].Attacker.Name, "P1A")
	assert.Equal(t, attacks.Attacks[0].Target.Name, "P1T")
	assert.Equal(t, attacks.Attacks[1].Attacker.Name, "P3A")
	assert.Equal(t, attacks.Attacks[1].Target.Name, "P3T")
}