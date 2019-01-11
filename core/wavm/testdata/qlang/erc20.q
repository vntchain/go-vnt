program PERC20 do
  contract ERC20 do
    store moneys: Ĝ Map, [I32, U64]
    #store approved: Ĝ Map, [I32 * I32, U64]
    api decimals do
      8_i32
    end
    api totalSupply do
      1234567890123456_u64
    end
    api balanceOf(u: I32) do
      @moneys[u]
    end
    api transfer(to: I32, val: U64) do
      if @moneys[0] < val
        false
      else
        @moneys[0] -= val
        @moneys[to] += val
        true
      end
    end
    #api transferFrom(from: I32, to: I32, val: I64) do
    #  x = @approved[ME]
    #  return false unless x && x[from] >= val && @moneys[from] >= val
    #  x[from] -= val
    #  @moneys[from] -= val
    #  @moneys[to] += val
    #  return true
    #end
    #api approve(spender: I32, val: I64) do
    #  (@approved[ME] ||= Map(I32, U64).new)[spender] = val
    #end
    #api allowance(owner: I32, spender: I32) do
    #  x = @approved[owner]
    #  x || x[spender]
    #end
  end
end
