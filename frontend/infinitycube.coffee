$ = jQuery

window.onload = ->
  $.get('status', (data) ->
    fillEffectSelector(data)
  , 'json')
  bindEffectSelector()
  bindClapSelector()
  bindRemoveButton()

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


bindEffectSelector = ->
  $('#effects').bind 'change', ->
    effect = $("#effects option:selected").val()
    $.post('toggle', {t: effect}, (data, textStatus, jqXHR) ->
      console.log(effect)
    , 'json')
    $.get('status', (data) ->
      fillActiveEffectSelector(data)
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
