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
	Leds         []Led
	Painter      ColorGenerator
	OffsetPar    int
	LengthPar    int
	ColorOpacity float64
	BlackOpacity float64
	myDisplay    Display
	mux          sync.Mutex
}

func NewEffect(disp Display, colorOp, blackOp float64) Effect {
	ef := Effect{
		Leds:         make([]Led, disp.NrOfLeds()),
		OffsetPar:    0,
		LengthPar:    disp.NrOfLeds(),
		ColorOpacity: colorOp,
		BlackOpacity: blackOp,
		myDisplay:    disp}

	return ef
}

// An Effector is able to generate (light-)Patterns. The Update() method gets
// called by its EffectHandler and has to call the AddPattern() method of its
// Display.
type Effector interface {
	Update()
}

// A Display is able to receive, merge and show the (light-)Patterns generated
// by various Effectors.
type Display interface {
	NrOfLeds() int
	AddPattern([]Led, float64, float64)
	Show()
}

// An EffectHandler creates, deletes, configures and updates Effectors. After
// updating all active Effectors the EffectHandler calls the show() method of
// its Display. The updateRate defines the frames per second.
type EffectHandler struct {
	activeEffects    []Effector
	effectParameter  []map[string]string
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
		myDisplay:  newDisplay,
		lastUpdate: time.Now(),
		updateRate: fps,
		loopTime:   time.Duration(float64(1000000000) / float64(fps)), // in nanoseconds
		//effectParameter: make([]map[string][]string),
		effectRequest: make(chan string),
		audio:         newAudio}

	eH.listAvailableEffects()
	eH.AddMultiRunningLight()
	go eH.handleRequests()
	return
}

func (eH *EffectHandler) render() {
	loopStart := time.Now()

	eH.updateAll()
	eH.myDisplay.Show()

	if time.Since(loopStart) < (eH.loopTime) {
		time.Sleep(eH.loopTime - time.Since(loopStart))
	}
}

func (eH *EffectHandler) updateAll() {
	for _, effect := range eH.activeEffects {
		effect.Update()
	}
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

//func (eH *EffectHandler) AddCellularAutomata() {
//	ca := NewCellularAutomata(eH.myDisplay)
//	eH.activeEffects = append(eH.activeEffects, ca)
//}

//func (eH *EffectHandler) AddWhiteSpectrum() {
//	ws := NewWhiteSpectrum(eH.myDisplay, eH.audio)
//	eH.activeEffects = append(eH.activeEffects, ws)
//}

//func (eH *EffectHandler) AddWhiteEdgeSpectrum() {
//	wes := NewWhiteEdgeSpectrum(eH.myDisplay, eH.audio)
//	eH.activeEffects = append(eH.activeEffects, wes)
//}

//func (eH *EffectHandler) AddEdgeVolume() {
//	ev := NewEdgeVolume(eH.myDisplay, eH.audio)
//	eH.activeEffects = append(eH.activeEffects, ev)
//}

//func (eH *EffectHandler) AddHsvFade() {
//	hsv := NewHsvFade(eH.myDisplay, eH.updateRate)
//	eH.activeEffects = append(eH.activeEffects, hsv)
//}

func (eH *EffectHandler) AddRunningLight() {
	rl := NewRunningLight(eH.myDisplay)
	eH.activeEffects = append(eH.activeEffects, rl)
}

func (eH *EffectHandler) AddMultiRunningLight() {
	grl := NewMultiRunningLight(eH.myDisplay, eH.updateRate)
	eH.activeEffects = append(eH.activeEffects, grl)
}

func (eH *EffectHandler) AddPrettyPrint() {
	fmt.Println("look at me I'm so pretty :-*")
}
