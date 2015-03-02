$ = jQuery

window.onload = ->
  window.setInterval(update, 1000)
  $('#cube').css('margin', 1).css('width', 330).css('height', 430).load('fancycube.svg') #need 4 resizehandler

update = ->
  $.get('status', (data) ->
    updateCubeLED(data) #to be written
  , 'json')