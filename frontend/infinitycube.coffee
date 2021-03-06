$ = jQuery

document.ready = ->
  console.log("onLy")
  $.get('status', (data) ->
    fillEffectSelector(data)
    fillActiveEffectSelector(data)
    fillEffectParameter(data)
    console.log("test")
  , 'json')
  console.log("onLy")
  bindEffectSelector()
  bindClapSelector()
  bindRemoveButton()
  bindActiveSelector()

fillEffectSelector = (data) ->
  console.log("AvailableEffects")
  console.log(data['AvailableEffects'])
  opt = ""
  for i in data['AvailableEffects']
    console.log(i)
    opt += "<option value='" + i + "'>" + i + "</option>"

  $('#effects').html("")
  $('#effects').append(opt)

fillActiveEffectSelector = (data) ->
  opt = ""
  idx = 0
  for i in data['ActiveEffects']
    console.log(i)
    opt += "<option value='" + idx + "'>" + idx + ". " + i + "</option>"
    idx++

  $('#activeEffects').html("")
  $('#activeEffects').append(opt)

fillEffectParameter = (data) ->
  opt = "<h1>EffectParameter:</h1><UL>"
  console.log(data['EffectParameter'])
  for k,v of data['EffectParameter']
    console.log("first loop " + k + v)
    opt += "<LI>" + k
    if k.lastIndexOf("Int", 0) == 0
      opt += ": <input type='range' id='" + k + "' value='" + v + "' min ='" + 0 + "' max='" + (parseInt v*10) + "'></input>"
    else  if k.lastIndexOf("Float", 0) == 0
      opt += ": <input type='range' id='" + k + "' value='" + v + "' min ='" + 0 + "' max='" + (parseInt v*10) + " step='0.01'></input>"
    else
      opt += ": <input type='text' id='" + k + "' value='" + v + "' maxlength='5' size='5'></input>"
  opt += "</UL>"


  $('#Parameter').html("")
  $('#Parameter').append(opt)
  console.log("<<< Before bindLoop")
  for k,v of data['EffectParameter']
    do (k) ->
      obj = "input#" + k
      div = "#" + k
      effect = $("#activeEffects option:selected").val()

      $(obj).change ->
        value = 'set' + effect + 'Par' + $(obj).attr('id') + 'Val' + $(obj).val()
        $.post 'toggle', { act: value }, ((data, textStatus, jqXHR) ->
          console.log '+++ EffectParameter'
          console.log value
        ), 'json'




bindEffectSelector = ->
  $('#effects').bind 'change', ->
    effect = $("#effects option:selected").val()
    $.post('toggle', {t: effect}, (data, textStatus, jqXHR) ->
      console.log("--- bindEffectSelector")
      console.log(effect)
    , 'json')
    $.get('status', (data) ->
      fillActiveEffectSelector(data)
    , 'json')

bindActiveSelector = ->
  $('#activeEffects').bind 'change', ->
    effect = $("#activeEffects option:selected")
    parameter = $("#activeEffects option:selected").val()
    $.post('toggle', {act: parameter}, (data, textStatus, jqXHR) ->
      console.log("*** bindActiveSelector")
      console.log(effect)
      #console.log(parameter)
    , 'json')
    $.get('status', (data) ->
      #console.log("binded Active selector get()")
      fillEffectParameter(data)
    , 'json')


bindRemoveButton = ->
    $('#remover').bind 'click', ->
      console.log("remoooove!")
      value = $("#activeEffects option:selected").val()
      $.post('toggle', {r: value}, (data, textStatus, jqXHR) ->
        console.log(value)
      , 'json')
      $.get('status', (data) ->
        fillActiveEffectSelector(data)
      , 'json')


bindClapSelector = ->
  $('#clapSelect').bind 'change', ->
    value = $("#clapSelect").prop('checked')
    $.post('toggle', {c: value}, (data, textStatus, jqXHR) ->
      console.log(value)
    , 'json')
