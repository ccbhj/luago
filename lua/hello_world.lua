local function sum(...)
  local l = {...}
  return l[1] + l[2]
end

return sum(1, 2) == 3
