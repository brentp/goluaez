-- remove whitespace at the ends of a string.
function string:strip()
    return self:match'^%s*(.*%S)' or ''
end

-- split a string by a separator
function string:split(sep)
    local sep, fields = sep or "\t", {}
    local pattern = string.format("[^%s]+", sep)
    for tok in self:gmatch(pattern) do fields[#fields+1] = tok end
    return fields
end

----------------------------------------------------
-- Testing Functions Don't copy these to Go code  --
----------------------------------------------------
function test_strip()
    a = " aaa "
    assert(a:strip() == "aaa")
    a = " aaa\t"
    assert(a:strip() == "aaa")
    a = " aaa    "
    assert(a:strip() == "aaa")
end

function test_split()
    v = string.split("x xx x", "%s")
    assert(#v == 3)
    assert(v[1] == "x")
    assert(v[2] == "xx")

    v = string.split("x xx\tx", "%s")
    assert(v[1] == "x")
    assert(v[2] == "xx")
end


function test()
    test_strip()
    test_split()
end

test()
