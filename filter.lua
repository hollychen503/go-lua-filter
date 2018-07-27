function filter(line)
--    found = string.find(line, "if")
--    if found == nil then
--        return false
--    end

    local json = require("json")
    assert(type(json) == "table")
    assert(type(json.decode) == "function")
    assert(type(json.encode) == "function")
        
    local jsonObj = json.decode(line)
    if (jsonObj.name == "holly") then 
        return true
    end

    return false
end
