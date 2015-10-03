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
    plotData = parseSoundData(data)

    ylim = 2500000000
    ymin = -2500000000
    timeDomainData = [{label:"Time Domain", data: plotData , lines : {show: true}, curvedLines: { apply: true}}]
    options = {series:{curvedLines: {active: true}}, yaxis:{max: ylim, min: ymin}, legend:{show: true, position: "ne", backgroundColor: "black"}}
    $.plot($("#timeDomainPlot"), timeDomainData, options)

parseSoundData = (data) ->
    SoundData = []
    i = 0
    for v in data['SoundSignal']
        SoundData.push [i, v]
        i++
    return SoundData

resizeHandler = ->
  w = window.innerWidth
  h = window.innerHeight
  chartWidth = 0.4*w
  console.log("resizing to #{w}, #{h}")
  $('#timeDomainPlot').css('margin', 1).css('width', chartWidth).css('height', chartWidth)
