$ = jQuery

window.onload = ->
  $.get('status', (data) ->
    fillEffectSelector(data)
  , 'json')
  bindEffectSelector()
  bindClapSelector()
  bindRemoveButton()
  bindActiveSelector()

fillEffectSelector = (data) ->
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
  opt = "<p><h1>EffectParameter:</h1><UL>"
  console.log(data['EffectParameter'])
  for k,v of data['EffectParameter']
    console.log(k)
    opt += "<LI>" + k + ": <input type='text' name='" + k + "' value='" + v + "' maxlength='5' size='5'><br></p>"
  opt += "</UL>"

  $('#Parameter').html("")
  $('#Parameter').append(opt)


bindEffectSelector = ->
  $('#effects').bind 'change', ->
    effect = $("#effects option:selected").val()
    $.post('toggle', {t: effect}, (data, textStatus, jqXHR) ->
      console.log(effect)
    , 'json')
    $.get('status', (data) ->
      fillActiveEffectSelector(data)
    , 'json')

bindActiveSelector = ->
  $('#activeEffects').bind 'change', ->
    effect = $("#activeEffects option:selected").val()
    $.post('toggle', {act: effect}, (data, textStatus, jqXHR) ->
      console.log(effect)
    , 'json')
    $.get('status', (data) ->
      fillEffectParameter(data)
    , 'json')


bindRemoveButton = ->
    $('#remover').bind 'click', ->
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
