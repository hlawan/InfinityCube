$ = jQuery

window.onload = ->
  window.setInterval(update, 1000/300) #set updaterate to 30hz
  $('#cube').css('margin', 1).css('width', 330).css('height', 430).svg({loadURL: 'fancycube.svg'}) #need 4 resizehandler

update = ->
  $.get('status', (data) ->
    updateCubeLED(data) 
  , 'json')

updateCubeLED = (data) ->
	for i in [0..5]
		for o in [0..3]
			for p in [0..13]
				color = Number(0x1000000 + data['LedR'][i][o][p]*0x10000 + data['LedG'][i][o][p]*0x100 + data['LedB'][i][o][p]).toString(16).substring(1);
				cube = "0"
				side = (i).toString()
				edge = (o).toString()
				led  = (p).toString()
				field = "#" + cube + side + edge + led #matching number format to svg path
				console.log(field)
				setColor(field, '#' + color) 

setColor = (field, color) ->
	#console.log("setColor called with field: ", field, " color: ", color)
	$(field, $('#cube').svg('get').root()).css('fill', color)