package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// An Effect has to provide information about the current led-pattern,
// the blackOpacity and the colorOpacity of its leds and which Display
// it belongs to.
type Effect struct {
	Leds      []Color
	Painter   ColorGenerator
	Dimming   float64
	OffsetPar int
	LengthPar int
	myDisplay Display
	mux       sync.Mutex
}

func NewEffect(disp Display, colorOp, blackOp float64) Effect {
	ef := Effect{
		Leds:      make([]Color, disp.NrOfLeds()),
		Dimming:   1.0,
		OffsetPar: 0,
		LengthPar: disp.NrOfLeds(),
		myDisplay: disp}

	return ef
}

func (ef *Effect) SetDimming(d float64) {
	ef.Dimming = d
}

// An Effector is able to generate (light-)Patterns. The Update() method gets
// called by its EffectHandler and has to call the AddPattern() method of its
// Display.
type Effector interface {
	SetDimming(float64)
	Update()
}

// A Display is able to receive, merge and show the (light-)Patterns generated
// by various Effectors.
type Display interface {
	NrOfLeds() int
	AddEffect(Effect)
	Show()
}

// An EffectHandler creates, deletes, configures and updates Effectors. After
// updating all active Effectors the EffectHandler calls the show() method of
// its Display. The updateRate defines the frames per second.
type EffectHandler struct {
	activeEffects    []Effector
	effectParameter  []map[string]string
	currentPlaylist  *PlayList
	availableEffects []string
	effectRequest    chan string
	myDisplay        Display
	lastUpdate       time.Time
	updateRate       int
	loopTime         time.Duration
	audio            *ProcessedAudio
}

func NewEffectHandler(newDisplay Display, fps int, newAudio *ProcessedAudio) (eH *EffectHandler) {
	eH = &EffectHandler{
		myDisplay:     newDisplay,
		lastUpdate:    time.Now(),
		updateRate:    fps,
		loopTime:      time.Duration(float64(1000000000) / float64(fps)), // in nanoseconds
		effectRequest: make(chan string),
		audio:         newAudio}

	eH.listAvailableEffects()
	go eH.handleRequests()
	return
}

func (eH *EffectHandler) render() {
	loopStart := time.Now()

	eH.checkPlayList()
	eH.updateAll()
	eH.myDisplay.Show()

	if time.Since(loopStart) < (eH.loopTime) {
		time.Sleep(eH.loopTime - time.Since(loopStart))
	}
}

func (eH *EffectHandler) checkPlayList() {
	if eH.currentPlaylist != nil {
		newEffects := eH.currentPlaylist.SlotEffects()
		if newEffects != nil {
			eH.activeEffects = newEffects
		}
	}
}

func (eH *EffectHandler) stopPlayList() {
	if eH.currentPlaylist != nil {
		eH.currentPlaylist = nil
		eH.activeEffects = nil
	}
}

func (eH *EffectHandler) updateAll() {
	for _, effect := range eH.activeEffects {
		effect.Update()
	}
}

func (eH *EffectHandler) AddPlayAllEffects() {

	slots := []map[Effector]time.Duration{}
	slots = append(slots, RedSunsetStarDust(eH, 20*time.Second))
	slots = append(slots, CellularAutomatagGradient(eH, 10*time.Second))
	slots = append(slots, MultiRunningLightHSV(eH, 10*time.Second))
	slots = append(slots, MagmaPlasma(eH, 10*time.Second))
	slots = append(slots, GoldenStarDust(eH, 20*time.Second))
	slots = append(slots, CellularAutomataMonochrome(eH, 10*time.Second))
	slots = append(slots, EdgeVolumeRedGreen(eH, 5*time.Second))
	slots = append(slots, LinearEdgeSpectrumMonochrome(eH, 5*time.Second))
	slots = append(slots, LinearSpectrumMonochrome(eH, 5*time.Second))

	eH.currentPlaylist = NewPlayList("all Effects", slots)
}

func (eH *EffectHandler) handleRequests() (err error) {
	for {
		//receive request from webserver
		fmt.Println("Waiting for Request (Backend)...")
		req := <-eH.effectRequest

		if strings.HasPrefix(req, "par") {
			eH.sendEffectParameter(req)
		}

		if strings.HasPrefix(req, "set") {
			eH.setParameter(req)
		}

		if strings.HasPrefix(req, "del") {
			eH.removeEffect(req)
			eH.sendActiveEffects()
		}

		if strings.HasPrefix(req, "Add") {
			eH.appendEffect(req)
			eH.sendActiveEffects()
		}
	}
}

func (eH *EffectHandler) setParameter(req string) {
	req = strings.TrimPrefix(req, "set")
	in := strings.Split(req, "Par")
	effNr, err := strconv.Atoi(in[0])
	CheckErr(err)
	in = strings.Split(in[1], "Val")
	val := in[1]
	par := in[0]

	fmt.Println(effNr, par, val)

	if strings.HasPrefix(par, "Float") {
		eH.setFloat(effNr, par, val)
	}
	if strings.HasPrefix(par, "Int") {
		eH.setInt(effNr, par, val)
	}
}

func (eH *EffectHandler) setInt(effNr int, par, val string) {
	par = strings.TrimPrefix(par, "Int")
	par += "Par"
	valInt, err := strconv.ParseInt(val, 10, 64)
	CheckErr(err)

	eff := reflect.ValueOf(eH.activeEffects[effNr]).Elem().FieldByName(par)

	if eff.CanSet() {
		eff.SetInt(valInt)
		fmt.Println("set ", par, " to ", valInt, " --- now it is", eff.Int())
	} else {
		fmt.Println("Effect Parameter: ", par, " is not settable")
	}
}

func (eH *EffectHandler) setFloat(effNr int, par, val string) {
	par = strings.TrimPrefix(par, "Float")
	par += "Par"
	valFloat, err := strconv.ParseFloat(val, 64)
	CheckErr(err)

	eff := reflect.ValueOf(eH.activeEffects[effNr]).Elem().FieldByName(par)

	if eff.CanSet() {
		eff.SetFloat(valFloat)
		fmt.Println("set ", par, " to ", valFloat, " --- now it is", eff.Float())
	} else {
		fmt.Println("Effect Parameter: ", par, " is not settable")
	}
}

func (eH *EffectHandler) removeEffect(req string) {
	req = strings.TrimPrefix(req, "del")
	ele, err := strconv.Atoi(req)
	CheckErr(err)

	var newList []Effector
	for i, effect := range eH.activeEffects {
		if i != ele {
			newList = append(newList, effect)
		}
	}
	eH.activeEffects = newList
}

func (eH *EffectHandler) appendEffect(req string) {
	fmt.Println("append called with: ", req)
	eHValue := reflect.ValueOf(eH)
	m := eHValue.MethodByName(req)
	if !m.IsValid() {
		fmt.Errorf("Method not found \"%s\"", req)
	}
	in := make([]reflect.Value, 0)
	m.Call(in)
}

func (eH *EffectHandler) listAvailableEffects() {
	eHType := reflect.TypeOf(eH)
	for i := 0; i < eHType.NumMethod(); i++ {
		method := eHType.Method(i).Name
		if strings.Contains(method, "Add") {
			method = strings.TrimPrefix(method, "Add")
			eH.availableEffects = append(eH.availableEffects, method)
		}
	}
	//fmt.Println("*EffectHandler.listAvailableEffects()*: ", eH.availableEffects)
}

func (eH *EffectHandler) sendActiveEffects() {
	for _, effect := range eH.activeEffects {
		eType := reflect.TypeOf(effect)
		eH.effectRequest <- eType.Elem().Name()
	}
	eH.effectRequest <- "done"
}

func (eH *EffectHandler) sendEffectParameter(req string) {
	req = strings.TrimPrefix(req, "par")
	nr, err := strconv.Atoi(req)
	CheckErr(err)
	eH.listEffectParameter()

	for par, val := range eH.effectParameter[nr] {
		eH.effectRequest <- par
		fmt.Print(" sent: ", par, " as par and")
		eH.effectRequest <- val
		fmt.Println(" sent: ", val, " as val")
	}
	eH.effectRequest <- "done"
	fmt.Println(" sent: done (EffectParameter)")
}

func (eH *EffectHandler) listEffectParameter() {
	var propList []map[string]string

	for i, ele := range eH.activeEffects {
		propList = append(propList, make(map[string]string))
		propList[i] = analyseEffectParameter(reflect.ValueOf(ele).Elem())
	}
	fmt.Println(propList[0])
	eH.effectParameter = propList
}

func analyseEffectParameter(s reflect.Value) map[string]string {
	parList := make(map[string]string)
	//fmt.Println(reflect.ValueOf(o))
	//fmt.Println(reflect.ValueOf(o).Elem())

	//s := reflect.ValueOf(o).Elem()
	typ := s.Type()

	for p := 0; p < s.NumField(); p++ {
		prop := s.Field(p)
		id := typ.Field(p).Name
		fmt.Println(prop.Kind(), " -> ", id)
		if strings.HasSuffix(id, "Par") {
			id = strings.TrimSuffix(id, "Par")

			switch prop.Kind() {
			case reflect.Int:
				id = "Int" + id
				parList[id] = strconv.Itoa(int(prop.Int()))
			case reflect.Float64:
				id = "Float" + id
				parList[id] = FloatToString(prop.Float())
			default:
				parList[id] = " "
			}
		} else if prop.Kind() == reflect.Struct {
			deepFields := analyseEffectParameter(prop)
			fmt.Println("found struct")
			for k, v := range deepFields {
				parList[k] = v
			}
		}
	}
	fmt.Println(parList)
	return parList
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

// List of addable Effects
func (eH *EffectHandler) AddCellularAutomataMonochrome() {
	eH.stopPlayList()
	effectSet := CellularAutomataMonochrome(eH, 30*time.Second)
	for effect, _ := range effectSet {
		eH.activeEffects = append(eH.activeEffects, effect)
	}
}

func (eH *EffectHandler) AddLinearEdgeSpectrumMonochrome() {
	eH.stopPlayList()
	effectSet := LinearEdgeSpectrumMonochrome(eH, 30*time.Second)
	for effect, _ := range effectSet {
		eH.activeEffects = append(eH.activeEffects, effect)
	}
}

func (eH *EffectHandler) AddLinearSpectrumMonochrome() {
	eH.stopPlayList()
	effectSet := LinearSpectrumMonochrome(eH, 30*time.Second)
	for effect, _ := range effectSet {
		eH.activeEffects = append(eH.activeEffects, effect)
	}
}

func (eH *EffectHandler) AddEdgeVolumeRedGreen() {
	eH.stopPlayList()
	effectSet := EdgeVolumeRedGreen(eH, 30*time.Second)
	for effect, _ := range effectSet {
		eH.activeEffects = append(eH.activeEffects, effect)
	}
}

func (eH *EffectHandler) AddMultiRunningLightHSV() {
	eH.stopPlayList()
	effectSet := MultiRunningLightHSV(eH, 30*time.Second)
	for effect, _ := range effectSet {
		eH.activeEffects = append(eH.activeEffects, effect)
	}
}

func (eH *EffectHandler) AddMagmaPlasma() {
	eH.stopPlayList()
	effectSet := MagmaPlasma(eH, 30*time.Second)
	for effect, _ := range effectSet {
		eH.activeEffects = append(eH.activeEffects, effect)
	}
}

func (eH *EffectHandler) AddFullWhite() {
	eH.stopPlayList()
	cc1 := NewConstantColor(0, 0)
	effect := NewSolidBrightness(eH.myDisplay, cc1, 1.0)
	eH.activeEffects = append(eH.activeEffects, effect)
}

func (eH *EffectHandler) AddOff() {
	eH.activeEffects = nil
}

func (eH *EffectHandler) AddPrettyPrint() {
	fmt.Println("look at me I'm so pretty :-*")
}
