$ = jQuery

window.onload = ->
  window.setInterval(update, 1000/15) #set updaterate to __hz

update = ->
  $.get('status', (data) ->
    plotGraphs(data)
  , 'json')

plotGraphs = (data)->
  ylim = 2500000000
  ymin = -2500000000
  timeDomainData = [{label:"Time Domain", data: HitsPerMinuteLeft, lines : {show: true}, curvedLines: { apply: true}}]
  options = {series:{curvedLines: {active: true}}, yaxis:{max: ylim, min: ymin}, legend:{show: true, position: "ne", backgroundColor: "black"}}
  $.plot($("#timeDomainPlot"),timeDomainData,options)
