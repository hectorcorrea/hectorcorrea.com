_zeroPad = (number, zeroes) =>
  return ('000000' + number).slice(-zeroes)


toSortableDate = (date) ->
  day = date.getDate()
  month = date.getMonth() + 1
  return date.getFullYear() + '-' + _zeroPad(month, 2) + '-' + _zeroPad(day, 2) 


addDays = (date, days) ->
  date.setDate(date.getDate() + days)
  return date


addHours = (date, hours) ->
  date.setHours(date.getHours() + hours)
  return date


module.exports = {
  toSortableDate: toSortableDate
  addHours: addHours
  addDays: addDays
}
