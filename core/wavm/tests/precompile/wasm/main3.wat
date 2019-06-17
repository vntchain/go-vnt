(module
 (type $FUNCSIG$i (func (result i32)))
 (type $FUNCSIG$j (func (result i64)))
 (type $FUNCSIG$ji (func (param i32) (result i64)))
 (type $FUNCSIG$vij (func (param i32 i64)))
 (type $FUNCSIG$ij (func (param i64) (result i32)))
 (type $FUNCSIG$ii (func (param i32) (result i32)))
 (type $FUNCSIG$vii (func (param i32 i32)))
 (type $FUNCSIG$iii (func (param i32 i32) (result i32)))
 (type $FUNCSIG$iiii (func (param i32 i32 i32) (result i32)))
 (type $FUNCSIG$viii (func (param i32 i32 i32)))
 (type $FUNCSIG$v (func))
 (import "env" "GenerateKey" (func $GenerateKey (param i32 i32 i32) (result i32)))
 (import "env" "read" (func $read (param i32 i32) (result i32)))
 (import "env" "write" (func $write (param i32 i32 i32)))
 (import "env" "GetOrigin" (func $GetOrigin (result i32)))
 (import "env" "GetSender" (func $GetSender (result i32)))
 (import "env" "GetValue" (func $GetValue (result i64)))
 (import "env" "GetContractValue" (func $GetContractValue (result i64)))
 (import "env" "GetContractAddress" (func $GetContractAddress (result i32)))
 (import "env" "GetBalanceFromAddress" (func $GetBalanceFromAddress (param i32) (result i64)))
 (import "env" "SendFromContract" (func $SendFromContract (param i32 i64)))
 (import "env" "GetBlockHash" (func $GetBlockHash (param i64) (result i32)))
 (import "env" "GetDifficulty" (func $GetDifficulty (result i64)))
 (import "env" "GetBlockNumber" (func $GetBlockNumber (result i64)))
 (import "env" "GetTimestamp" (func $GetTimestamp (result i64)))
 (import "env" "GetCoinBase" (func $GetCoinBase (result i32)))
 (import "env" "GetGas" (func $GetGas (result i64)))
 (import "env" "GetGasLimit" (func $GetGasLimit (result i64)))
 (import "env" "StorageWrite" (func $StorageWrite (param i32 i32)))
 (import "env" "StorageRead" (func $StorageRead (param i32) (result i32)))
 (import "env" "SHA3" (func $SHA3 (param i32) (result i32)))
 (import "env" "FromI64" (func $FromI64 (param i64) (result i32)))
 (import "env" "FromU64" (func $FromU64 (param i64) (result i32)))
 (import "env" "ToI64" (func $ToI64 (param i32) (result i64)))
 (import "env" "ToU64" (func $ToU64 (param i32) (result i64)))
 (import "env" "Concat" (func $Concat (param i32 i32) (result i32)))
 (table 23 23 anyfunc)
 (elem (i32.const 0) $__wasm_nullptr $__importThunk_GetOrigin $__importThunk_GetSender $__importThunk_GetValue $__importThunk_GetContractValue $__importThunk_GetContractAddress $__importThunk_GetBalanceFromAddress $__importThunk_SendFromContract $__importThunk_GetBlockHash $__importThunk_GetDifficulty $__importThunk_GetBlockNumber $__importThunk_GetTimestamp $__importThunk_GetCoinBase $__importThunk_GetGas $__importThunk_GetGasLimit $__importThunk_StorageWrite $__importThunk_StorageRead $__importThunk_SHA3 $__importThunk_FromI64 $__importThunk_FromU64 $__importThunk_ToI64 $__importThunk_ToU64 $__importThunk_Concat)
 (memory $0 1)
 (data (i32.const 12) "\00\00\00\00")
 (data (i32.const 16) "\00\00\00\00")
 (data (i32.const 32) "mapping_e\00")
 (data (i32.const 48) "mapping_e.key\00")
 (data (i32.const 64) "mapping_e.value\00")
 (data (i32.const 96) "mapping_a\00")
 (data (i32.const 112) "mapping_a.key\00")
 (data (i32.const 144) "mapping_c\00")
 (data (i32.const 160) "mapping_c.key\00")
 (data (i32.const 208) "array_e\00")
 (data (i32.const 224) "array_e.index\00")
 (data (i32.const 240) "array_e.length\00")
 (data (i32.const 272) "1\00")
 (data (i32.const 288) "mapping_h\00")
 (data (i32.const 304) "mapping_h.key\00")
 (data (i32.const 320) "name\00")
 (data (i32.const 336) "mapping_h.value.key\00")
 (data (i32.const 368) "aaa\00")
 (data (i32.const 384) "address\00")
 (data (i32.const 400) "bbb\00")
 (data (i32.const 464) "mapping_i\00")
 (data (i32.const 480) "mapping_i.key\00")
 (data (i32.const 496) "mapping_i.value.name\00")
 (data (i32.const 528) "2\00")
 (data (i32.const 544) "mapping_i.value.address\00")
 (data (i32.const 576) "mapping_i.value.phone\00")
 (data (i32.const 616) "\02\00\00\00")
 (data (i32.const 620) "\01\00\00\00")
 (data (i32.const 624) "\01\00\00\00\01\00\00\00")
 (data (i32.const 632) "\01\00\00\00\00\00\00\00\01\00\00\00\00\00\00\00\01\00\00\00\00\00\00\00")
 (data (i32.const 656) "\01\00\00\00\01\00\00\00\01\00\00\00\01\00\00\00")
 (data (i32.const 672) "\10\01\00\00\10\01\00\00\10\01\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data (i32.const 720) "\01\00\00\00\01\00\00\00\01\00\00\00\01\00\00\00")
 (data (i32.const 736) "\00\00\00\00\00\00\00\00")
 (data (i32.const 744) "\00\00\00\00\00\00\00\00")
 (data (i32.const 752) "\00\00\00\00")
 (data (i32.const 760) "\00\00\00\00\00\00\00\00")
 (data (i32.const 768) "\00\00\00\00")
 (export "memory" (memory $0))
 (export "getmsg" (func $getmsg))
 (export "getcontract" (func $getcontract))
 (export "getblock" (func $getblock))
 (export "getutils" (func $getutils))
 (export "writemapping" (func $writemapping))
 (export "readmapping" (func $readmapping))
 (export "writemapping_int32" (func $writemapping_int32))
 (export "readmapping_int32" (func $readmapping_int32))
 (export "writemapping_int64" (func $writemapping_int64))
 (export "readmapping_int64" (func $readmapping_int64))
 (export "writearray" (func $writearray))
 (export "readarray" (func $readarray))
 (export "readarraylength" (func $readarraylength))
 (export "TestPop" (func $TestPop))
 (export "TestPush" (func $TestPush))
 (export "TestComplexMappingWithNesting" (func $TestComplexMappingWithNesting))
 (export "TestWriteComplexMappingWithStruct" (func $TestWriteComplexMappingWithStruct))
 (export "TestReadComplexMappingWithStruct" (func $TestReadComplexMappingWithStruct))
 (export "testcomplex" (func $testcomplex))
 (func $getmsg (; 25 ;) (param $0 i32)
  (i32.store offset=4
   (get_local $0)
   (i32.const 1)
  )
  (i32.store
   (get_local $0)
   (i32.const 2)
  )
  (i32.store offset=8
   (get_local $0)
   (i32.const 3)
  )
 )
 (func $getcontract (; 26 ;) (param $0 i32)
  (i32.store offset=4
   (get_local $0)
   (i32.const 4)
  )
  (i32.store
   (get_local $0)
   (i32.const 5)
  )
  (i32.store offset=8
   (get_local $0)
   (i32.const 6)
  )
  (i32.store offset=12
   (get_local $0)
   (i32.const 7)
  )
 )
 (func $getblock (; 27 ;) (param $0 i32)
  (i32.store offset=4
   (get_local $0)
   (i32.const 8)
  )
  (i32.store
   (get_local $0)
   (i32.const 9)
  )
  (i32.store offset=8
   (get_local $0)
   (i32.const 10)
  )
  (i32.store offset=12
   (get_local $0)
   (i32.const 11)
  )
  (i32.store offset=16
   (get_local $0)
   (i32.const 12)
  )
  (i32.store offset=20
   (get_local $0)
   (i32.const 13)
  )
  (i32.store offset=24
   (get_local $0)
   (i32.const 14)
  )
 )
 (func $getutils (; 28 ;) (param $0 i32)
  (i32.store offset=4
   (get_local $0)
   (i32.const 15)
  )
  (i32.store
   (get_local $0)
   (i32.const 16)
  )
  (i32.store offset=8
   (get_local $0)
   (i32.const 17)
  )
  (i32.store offset=12
   (get_local $0)
   (i32.const 18)
  )
  (i32.store offset=16
   (get_local $0)
   (i32.const 19)
  )
  (i32.store offset=20
   (get_local $0)
   (i32.const 20)
  )
  (i32.store offset=24
   (get_local $0)
   (i32.const 21)
  )
  (i32.store offset=28
   (get_local $0)
   (i32.const 22)
  )
 )
 (func $writemapping (; 29 ;) (param $0 i32) (param $1 i32)
  (i32.store offset=24
   (i32.const 0)
   (get_local $1)
  )
  (i32.store offset=20
   (i32.const 0)
   (get_local $0)
  )
  (call $write
   (i32.const 64)
   (call $GenerateKey
    (i32.const 32)
    (i32.const 48)
    (get_local $0)
   )
   (i32.const 24)
  )
 )
 (func $readmapping (; 30 ;) (param $0 i32) (result i32)
  (i32.store offset=20
   (i32.const 0)
   (get_local $0)
  )
  (i32.store offset=24
   (i32.const 0)
   (tee_local $0
    (i32.load
     (call $read
      (i32.const 64)
      (call $GenerateKey
       (i32.const 32)
       (i32.const 48)
       (get_local $0)
      )
     )
    )
   )
  )
  (get_local $0)
 )
 (func $writemapping_int32 (; 31 ;) (param $0 i32) (param $1 i32)
  (i32.store offset=84
   (i32.const 0)
   (get_local $1)
  )
  (i32.store offset=80
   (i32.const 0)
   (get_local $0)
  )
  (call $write
   (i32.const 96)
   (call $GenerateKey
    (i32.const 96)
    (i32.const 112)
    (get_local $0)
   )
   (i32.const 84)
  )
 )
 (func $readmapping_int32 (; 32 ;) (param $0 i32) (result i32)
  (i32.store offset=80
   (i32.const 0)
   (get_local $0)
  )
  (i32.store offset=84
   (i32.const 0)
   (tee_local $0
    (i32.load
     (call $read
      (i32.const 96)
      (call $GenerateKey
       (i32.const 96)
       (i32.const 112)
       (get_local $0)
      )
     )
    )
   )
  )
  (get_local $0)
 )
 (func $writemapping_int64 (; 33 ;) (param $0 i64) (param $1 i64)
  (i64.store offset=136
   (i32.const 0)
   (get_local $1)
  )
  (i64.store offset=128
   (i32.const 0)
   (get_local $0)
  )
  (call $write
   (i32.const 144)
   (call $GenerateKey
    (i32.const 144)
    (i32.const 160)
    (i32.wrap/i64
     (get_local $0)
    )
   )
   (i32.const 136)
  )
 )
 (func $readmapping_int64 (; 34 ;) (param $0 i64) (result i64)
  (i64.store offset=128
   (i32.const 0)
   (get_local $0)
  )
  (i64.store offset=136
   (i32.const 0)
   (tee_local $0
    (i64.load
     (call $read
      (i32.const 144)
      (call $GenerateKey
       (i32.const 144)
       (i32.const 160)
       (i32.wrap/i64
        (get_local $0)
       )
      )
     )
    )
   )
  )
  (get_local $0)
 )
 (func $writearray (; 35 ;) (param $0 i64) (param $1 i32) (param $2 i64)
  (i32.store offset=184
   (i32.const 0)
   (get_local $1)
  )
  (i64.store offset=176
   (i32.const 0)
   (get_local $0)
  )
  (i64.store offset=192
   (i32.const 0)
   (get_local $2)
  )
  (set_local $1
   (call $GenerateKey
    (i32.const 208)
    (i32.const 224)
    (i32.wrap/i64
     (get_local $0)
    )
   )
  )
  (call $write
   (i32.const 208)
   (call $GenerateKey
    (i32.const 208)
    (i32.const 240)
    (i32.const 256)
   )
   (i32.const 192)
  )
  (call $write
   (i32.const 208)
   (get_local $1)
   (i32.const 184)
  )
 )
 (func $readarray (; 36 ;) (param $0 i64) (result i32)
  (local $1 i32)
  (i64.store offset=176
   (i32.const 0)
   (get_local $0)
  )
  (i32.store offset=184
   (i32.const 0)
   (tee_local $1
    (i32.load
     (call $read
      (i32.const 208)
      (call $GenerateKey
       (i32.const 208)
       (i32.const 224)
       (i32.wrap/i64
        (get_local $0)
       )
      )
     )
    )
   )
  )
  (get_local $1)
 )
 (func $readarraylength (; 37 ;) (param $0 i64) (result i64)
  (local $1 i64)
  (i64.store offset=192
   (i32.const 0)
   (tee_local $1
    (i64.load
     (call $read
      (i32.const 208)
      (call $GenerateKey
       (i32.const 208)
       (i32.const 240)
       (i32.const 256)
      )
     )
    )
   )
  )
  (get_local $1)
 )
 (func $TestPop (; 38 ;) (result i32)
  (local $0 i64)
  (local $1 i32)
  (i64.store offset=192
   (i32.const 0)
   (tee_local $0
    (i64.load
     (call $read
      (i32.const 208)
      (i32.const 256)
     )
    )
   )
  )
  (i64.store offset=176
   (i32.const 0)
   (tee_local $0
    (i64.add
     (get_local $0)
     (i64.const -1)
    )
   )
  )
  (i32.store offset=184
   (i32.const 0)
   (tee_local $1
    (i32.load
     (call $read
      (i32.const 208)
      (i32.wrap/i64
       (get_local $0)
      )
     )
    )
   )
  )
  (i64.store offset=192
   (i32.const 0)
   (i64.add
    (i64.load offset=192
     (i32.const 0)
    )
    (i64.const -1)
   )
  )
  (call $write
   (i32.const 208)
   (i32.const 256)
   (i32.const 192)
  )
  (get_local $1)
 )
 (func $TestPush (; 39 ;) (param $0 i32)
  (local $1 i64)
  (i64.store offset=192
   (i32.const 0)
   (i64.add
    (i64.load
     (call $read
      (i32.const 208)
      (i32.const 256)
     )
    )
    (i64.const 1)
   )
  )
  (call $write
   (i32.const 208)
   (i32.const 256)
   (i32.const 192)
  )
  (i32.store offset=184
   (i32.const 0)
   (get_local $0)
  )
  (i64.store offset=176
   (i32.const 0)
   (tee_local $1
    (i64.add
     (i64.load offset=192
      (i32.const 0)
     )
     (i64.const -1)
    )
   )
  )
  (call $write
   (i32.const 208)
   (i32.wrap/i64
    (get_local $1)
   )
   (i32.const 184)
  )
 )
 (func $TestComplexMappingWithNesting (; 40 ;)
  (local $0 i32)
  (local $1 i32)
  (i32.store offset=276
   (i32.const 0)
   (i32.const 272)
  )
  (set_local $0
   (call $GenerateKey
    (i32.const 288)
    (i32.const 304)
    (i32.const 272)
   )
  )
  (i32.store offset=280
   (i32.const 0)
   (i32.const 320)
  )
  (set_local $1
   (call $GenerateKey
    (get_local $0)
    (i32.const 336)
    (i32.const 320)
   )
  )
  (i32.store offset=284
   (i32.const 0)
   (i32.const 368)
  )
  (call $write
   (i32.const 288)
   (get_local $1)
   (i32.const 284)
  )
  (i32.store offset=280
   (i32.const 0)
   (i32.const 384)
  )
  (set_local $0
   (call $GenerateKey
    (get_local $0)
    (i32.const 336)
    (i32.const 384)
   )
  )
  (i32.store offset=284
   (i32.const 0)
   (i32.const 400)
  )
  (call $write
   (i32.const 288)
   (get_local $0)
   (i32.const 284)
  )
 )
 (func $TestWriteComplexMappingWithStruct (; 41 ;)
  (local $0 i32)
  (local $1 i32)
  (i32.store offset=408
   (i32.const 0)
   (i32.const 272)
  )
  (set_local $1
   (call $GenerateKey
    (tee_local $0
     (call $GenerateKey
      (i32.const 464)
      (i32.const 480)
      (i32.const 272)
     )
    )
    (i32.const 496)
    (i32.const 256)
   )
  )
  (i32.store offset=416
   (i32.const 0)
   (i32.const 528)
  )
  (call $write
   (i32.const 464)
   (get_local $1)
   (i32.const 416)
  )
  (set_local $1
   (call $GenerateKey
    (get_local $0)
    (i32.const 544)
    (i32.const 256)
   )
  )
  (i32.store offset=420
   (i32.const 0)
   (i32.const 528)
  )
  (call $write
   (i32.const 464)
   (get_local $1)
   (i32.const 420)
  )
  (set_local $0
   (call $GenerateKey
    (get_local $0)
    (i32.const 576)
    (i32.const 256)
   )
  )
  (i32.store offset=424
   (i32.const 0)
   (i32.const 528)
  )
  (call $write
   (i32.const 464)
   (get_local $0)
   (i32.const 424)
  )
 )
 (func $TestReadComplexMappingWithStruct (; 42 ;)
  (local $0 i32)
  (i32.store offset=408
   (i32.const 0)
   (i32.const 272)
  )
  (i64.store32 offset=416
   (i32.const 0)
   (i64.load
    (call $read
     (i32.const 464)
     (call $GenerateKey
      (tee_local $0
       (call $GenerateKey
        (i32.const 464)
        (i32.const 480)
        (i32.const 272)
       )
      )
      (i32.const 496)
      (i32.const 256)
     )
    )
   )
  )
  (i64.store32 offset=420
   (i32.const 0)
   (i64.load
    (call $read
     (i32.const 464)
     (call $GenerateKey
      (get_local $0)
      (i32.const 544)
      (i32.const 256)
     )
    )
   )
  )
  (i64.store32 offset=424
   (i32.const 0)
   (i64.load
    (call $read
     (i32.const 464)
     (call $GenerateKey
      (get_local $0)
      (i32.const 576)
      (i32.const 256)
     )
    )
   )
  )
 )
 (func $testcomplex (; 43 ;)
  (i64.store offset=600 align=4
   (i32.const 0)
   (i64.const 4294967297)
  )
  (i64.store offset=608 align=4
   (i32.const 0)
   (i64.const 4294967297)
  )
 )
 (func $__wasm_nullptr (; 44 ;) (type $FUNCSIG$v)
  (unreachable)
 )
 (func $__importThunk_GetOrigin (; 45 ;) (type $FUNCSIG$i) (result i32)
  (call $GetOrigin)
 )
 (func $__importThunk_GetSender (; 46 ;) (type $FUNCSIG$i) (result i32)
  (call $GetSender)
 )
 (func $__importThunk_GetValue (; 47 ;) (type $FUNCSIG$j) (result i64)
  (call $GetValue)
 )
 (func $__importThunk_GetContractValue (; 48 ;) (type $FUNCSIG$j) (result i64)
  (call $GetContractValue)
 )
 (func $__importThunk_GetContractAddress (; 49 ;) (type $FUNCSIG$i) (result i32)
  (call $GetContractAddress)
 )
 (func $__importThunk_GetBalanceFromAddress (; 50 ;) (type $FUNCSIG$ji) (param $0 i32) (result i64)
  (call $GetBalanceFromAddress
   (get_local $0)
  )
 )
 (func $__importThunk_SendFromContract (; 51 ;) (type $FUNCSIG$vij) (param $0 i32) (param $1 i64)
  (call $SendFromContract
   (get_local $0)
   (get_local $1)
  )
 )
 (func $__importThunk_GetBlockHash (; 52 ;) (type $FUNCSIG$ij) (param $0 i64) (result i32)
  (call $GetBlockHash
   (get_local $0)
  )
 )
 (func $__importThunk_GetDifficulty (; 53 ;) (type $FUNCSIG$j) (result i64)
  (call $GetDifficulty)
 )
 (func $__importThunk_GetBlockNumber (; 54 ;) (type $FUNCSIG$j) (result i64)
  (call $GetBlockNumber)
 )
 (func $__importThunk_GetTimestamp (; 55 ;) (type $FUNCSIG$j) (result i64)
  (call $GetTimestamp)
 )
 (func $__importThunk_GetCoinBase (; 56 ;) (type $FUNCSIG$i) (result i32)
  (call $GetCoinBase)
 )
 (func $__importThunk_GetGas (; 57 ;) (type $FUNCSIG$j) (result i64)
  (call $GetGas)
 )
 (func $__importThunk_GetGasLimit (; 58 ;) (type $FUNCSIG$j) (result i64)
  (call $GetGasLimit)
 )
 (func $__importThunk_StorageWrite (; 59 ;) (type $FUNCSIG$vii) (param $0 i32) (param $1 i32)
  (call $StorageWrite
   (get_local $0)
   (get_local $1)
  )
 )
 (func $__importThunk_StorageRead (; 60 ;) (type $FUNCSIG$ii) (param $0 i32) (result i32)
  (call $StorageRead
   (get_local $0)
  )
 )
 (func $__importThunk_SHA3 (; 61 ;) (type $FUNCSIG$ii) (param $0 i32) (result i32)
  (call $SHA3
   (get_local $0)
  )
 )
 (func $__importThunk_FromI64 (; 62 ;) (type $FUNCSIG$ij) (param $0 i64) (result i32)
  (call $FromI64
   (get_local $0)
  )
 )
 (func $__importThunk_FromU64 (; 63 ;) (type $FUNCSIG$ij) (param $0 i64) (result i32)
  (call $FromU64
   (get_local $0)
  )
 )
 (func $__importThunk_ToI64 (; 64 ;) (type $FUNCSIG$ji) (param $0 i32) (result i64)
  (call $ToI64
   (get_local $0)
  )
 )
 (func $__importThunk_ToU64 (; 65 ;) (type $FUNCSIG$ji) (param $0 i32) (result i64)
  (call $ToU64
   (get_local $0)
  )
 )
 (func $__importThunk_Concat (; 66 ;) (type $FUNCSIG$iii) (param $0 i32) (param $1 i32) (result i32)
  (call $Concat
   (get_local $0)
   (get_local $1)
  )
 )
)
