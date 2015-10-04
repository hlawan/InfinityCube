$ = jQuery

window.onload = ->
  $(window).resize(resizeHandler)
  resizeHandler()
  window.setInterval(update, 1000/50) #set updaterate to __hz

update = ->
  $.get('status', (data) ->
    plotGraphs(data)
  , 'json')

plotGraphs = (data)->
    soundData = parseSoundData(data)
    freqData = parseFreqData(data)

    ylim = 1
    ymin = -1
    timeDomainData = [{label:"Time Domain", data: soundData , lines : {show: true}, curvedLines: { apply: true}}]
    freqDomainData = [{label:"Freq Domain", data: freqData , lines : {show: true}, curvedLines: { apply: true}}]
    optionsTime = {series:{curvedLines: {active: true}}, yaxis:{max: ylim, min: ymin}, legend:{show: true, position: "ne", backgroundColor: "black"}}
    optionsFreq = {series:{curvedLines: {active: true}}, yaxis:{max: 0.001, min: 0}, xaxis:{max: 5000, min: 0}, legend:{show: true, position: "ne", backgroundColor: "black"}}
    $.plot($("#timeDomainPlot"), timeDomainData, optionsTime)
    $.plot($("#freqDomainPlot"), freqDomainData, optionsFreq)

parseSoundData = (data) ->
    SoundData = []
    i = 0
    for v in data['SoundSignal']
        SoundData.push [i, v]
        i++
    return SoundData

parseFreqData = (data) ->
    FreqData = []
    i = 0
    for v in data['Freqs']
        FreqData.push [v, data['SpectralDensity'][i]]
        i++
    return FreqData

resizeHandler = ->
  w = window.innerWidth
  h = window.innerHeight
  chartWidth = 0.4*w
  console.log("resizing to #{w}, #{h}")
  $('#timeDomainPlot').css('margin', 1).css('width', chartWidth).css('height', chartWidth)
  $('#freqDomainPlot').css('margin', 1).css('width', chartWidth).css('height', chartWidth)
