package main

type Effector interface {
  Update()
}

type EffectHandler struct {
  Effects []Effector
}

func NewEffectHandler() (eH *EffectHandler) {
  eH = &EffectHandler{}
  return
}

func (eH *EffectHandler) updateAll() {
  for _,effect := range eH.Effects {
    effect.Update()
  }
}

func (eH *EffectHandler) addCellularAutomata(colorOpacity, blackOpacity, secsperGen float64, rule int) {
  cA := NewCellularAutomata(colorOpacity, blackOpacity, rule, secsperGen)
  //cA.Consumer/Display = cube....
  eH.Effects = append(eH.Effects, cA)
}
