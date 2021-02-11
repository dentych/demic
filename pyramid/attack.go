package pyramid

import "sync"

type Attack struct {
	Attacker Player
	Target   Player
}

func (a *Attack) EqualTo(other Attack) bool {
	if a.Attacker.Name == "" || a.Target.Name == "" {
		return false
	}
	return a.Attacker.Name == other.Attacker.Name && a.Target.Name == other.Target.Name
}

type Attacks struct {
	Attacks []Attack `json:"attacks"`
	mutex   sync.Mutex
}

func (a *Attacks) Add(attack Attack) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.Attacks = append(a.Attacks, attack)
}

func (a *Attacks) Remove(attack Attack) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if len(a.Attacks) == 0 {
		return
	}

	var attacks []Attack
	for _, v := range a.Attacks {
		if v.EqualTo(attack) {
			continue
		}
		attacks = append(attacks, v)
	}

	a.Attacks = attacks
}

func (a *Attacks) Len() int {
	return len(a.Attacks)
}