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
 (type $FUNCSIG$vijii (func (param i32 i64 i32 i32)))
 (type $FUNCSIG$iiji (func (param i32 i64 i32) (result i32)))
 (type $FUNCSIG$v (func))
 (import "env" "GenerateKey" (func $GenerateKey (param i32 i32 i32) (result i32)))
 (import "env" "GenerateKeyWithType" (func $GenerateKeyWithType (param i32 i32 i32) (result i32)))
 (import "env" "ReadArrayWithType" (func $ReadArrayWithType (param i32 i64 i32) (result i32)))
 (import "env" "ReadWithType" (func $ReadWithType (param i32 i32) (result i32)))
 (import "env" "WriteArrayWithType" (func $WriteArrayWithType (param i32 i64 i32 i32)))
 (import "env" "WriteWithType" (func $WriteWithType (param i32 i32 i32)))
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
 (data (i32.const 64) "mapping_a\00")
 (data (i32.const 112) "mapping_c\00")
 (data (i32.const 160) "array_e\00")
 (data (i32.const 192) "1\00")
 (data (i32.const 224) "mapping_h\00")
 (data (i32.const 240) "mapping_h.key\00")
 (data (i32.const 256) "name\00")
 (data (i32.const 272) "mapping_h.value.key\00")
 (data (i32.const 304) "aaa\00")
 (data (i32.const 320) "address\00")
 (data (i32.const 336) "bbb\00")
 (data (i32.const 432) "mapping_i\00")
 (data (i32.const 448) "mapping_i.key\00")
 (data (i32.const 464) "mapping_i.value.name\00")
 (data (i32.const 496) "2\00")
 (data (i32.const 512) "mapping_i.value.address\00")
 (data (i32.const 544) "mapping_i.value.phone\00")
 (data (i32.const 596) "\02\00\00\00")
 (data (i32.const 600) "\01\00\00\00")
 (data (i32.const 604) "\01\00\00\00\01\00\00\00\00\00\00\00")
 (data (i32.const 616) "\01\00\00\00\00\00\00\00\01\00\00\00\00\00\00\00\01\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data (i32.const 648) "\01\00\00\00\01\00\00\00\01\00\00\00\01\00\00\00")
 (data (i32.const 664) "\c0\00\00\00\c0\00\00\00\c0\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
 (data (i32.const 728) "\01\00\00\00\01\00\00\00\01\00\00\00\01\00\00\00")
 (data (i32.const 744) "\00\00\00\00\00\00\00\00")
 (data (i32.const 752) "\00\00\00\00\00\00\00\00")
 (data (i32.const 760) "\00\00\00\00")
 (data (i32.const 768) "\00\00\00\00\00\00\00\00")
 (data (i32.const 776) "\00\00\00\00")
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
 (func $getmsg (; 30 ;) (param $0 i32)
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
 (func $getcontract (; 31 ;) (param $0 i32)
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
 (func $getblock (; 32 ;) (param $0 i32)
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
 (func $getutils (; 33 ;) (param $0 i32)
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
 (func $writemapping (; 34 ;) (param $0 i32) (param $1 i32)
  (i32.store offset=20
   (i32.const 0)
   (get_local $0)
  )
  (set_local $0
   (call $GenerateKeyWithType
    (i32.const 32)
    (i32.const 5)
    (i32.const 20)
   )
  )
  (i32.store offset=24
   (i32.const 0)
   (get_local $1)
  )
  (call $WriteWithType
   (get_local $0)
   (i32.const 5)
   (i32.const 24)
  )
 )
 (func $readmapping (; 35 ;) (param $0 i32) (result i32)
  (i32.store offset=20
   (i32.const 0)
   (get_local $0)
  )
  (i32.store offset=24
   (i32.const 0)
   (tee_local $0
    (i32.load
     (call $ReadWithType
      (call $GenerateKeyWithType
       (i32.const 32)
       (i32.const 5)
       (i32.const 20)
      )
      (i32.const 5)
     )
    )
   )
  )
  (get_local $0)
 )
 (func $writemapping_int32 (; 36 ;) (param $0 i32) (param $1 i32)
  (i32.store offset=44
   (i32.const 0)
   (get_local $0)
  )
  (set_local $0
   (call $GenerateKeyWithType
    (i32.const 64)
    (i32.const 1)
    (i32.const 44)
   )
  )
  (i32.store offset=48
   (i32.const 0)
   (get_local $1)
  )
  (call $WriteWithType
   (get_local $0)
   (i32.const 1)
   (i32.const 48)
  )
 )
 (func $readmapping_int32 (; 37 ;) (param $0 i32) (result i32)
  (i32.store offset=44
   (i32.const 0)
   (get_local $0)
  )
  (i32.store offset=48
   (i32.const 0)
   (tee_local $0
    (i32.load
     (call $ReadWithType
      (call $GenerateKeyWithType
       (i32.const 64)
       (i32.const 1)
       (i32.const 44)
      )
      (i32.const 1)
     )
    )
   )
  )
  (get_local $0)
 )
 (func $writemapping_int64 (; 38 ;) (param $0 i64) (param $1 i64)
  (local $2 i32)
  (i64.store offset=80
   (i32.const 0)
   (get_local $0)
  )
  (set_local $2
   (call $GenerateKeyWithType
    (i32.const 112)
    (i32.const 2)
    (i32.const 80)
   )
  )
  (i64.store offset=88
   (i32.const 0)
   (get_local $1)
  )
  (call $WriteWithType
   (get_local $2)
   (i32.const 2)
   (i32.const 88)
  )
 )
 (func $readmapping_int64 (; 39 ;) (param $0 i64) (result i64)
  (i64.store offset=80
   (i32.const 0)
   (get_local $0)
  )
  (i64.store offset=88
   (i32.const 0)
   (tee_local $0
    (i64.load
     (call $ReadWithType
      (call $GenerateKeyWithType
       (i32.const 112)
       (i32.const 2)
       (i32.const 80)
      )
      (i32.const 2)
     )
    )
   )
  )
  (get_local $0)
 )
 (func $writearray (; 40 ;) (param $0 i64) (param $1 i32) (param $2 i64)
  (local $3 i32)
  (i64.store offset=144
   (i32.const 0)
   (get_local $2)
  )
  (call $WriteWithType
   (tee_local $3
    (call $GenerateKeyWithType
     (i32.const 160)
     (i32.const 8)
     (i32.const 0)
    )
   )
   (i32.const 2)
   (i32.const 144)
  )
  (i32.store offset=136
   (i32.const 0)
   (get_local $1)
  )
  (i64.store offset=128
   (i32.const 0)
   (get_local $0)
  )
  (call $WriteArrayWithType
   (get_local $3)
   (get_local $0)
   (i32.const 5)
   (i32.const 136)
  )
 )
 (func $readarray (; 41 ;) (param $0 i64) (result i32)
  (local $1 i32)
  (i64.store offset=128
   (i32.const 0)
   (get_local $0)
  )
  (i32.store offset=136
   (i32.const 0)
   (tee_local $1
    (i32.load
     (call $ReadArrayWithType
      (call $GenerateKeyWithType
       (i32.const 160)
       (i32.const 8)
       (i32.const 0)
      )
      (i64.load offset=128
       (i32.const 0)
      )
      (i32.const 5)
     )
    )
   )
  )
  (get_local $1)
 )
 (func $readarraylength (; 42 ;) (param $0 i64) (result i64)
  (local $1 i64)
  (i64.store offset=144
   (i32.const 0)
   (tee_local $1
    (i64.load
     (call $ReadWithType
      (call $GenerateKeyWithType
       (i32.const 160)
       (i32.const 8)
       (i32.const 0)
      )
      (i32.const 2)
     )
    )
   )
  )
  (get_local $1)
 )
 (func $TestPop (; 43 ;) (result i32)
  (local $0 i64)
  (local $1 i32)
  (i64.store offset=144
   (i32.const 0)
   (tee_local $0
    (i64.load
     (call $read
      (i32.const 160)
      (i32.const 176)
     )
    )
   )
  )
  (i64.store offset=128
   (i32.const 0)
   (tee_local $0
    (i64.add
     (get_local $0)
     (i64.const -1)
    )
   )
  )
  (i32.store offset=136
   (i32.const 0)
   (tee_local $1
    (i32.load
     (call $read
      (i32.const 160)
      (i32.wrap/i64
       (get_local $0)
      )
     )
    )
   )
  )
  (i64.store offset=144
   (i32.const 0)
   (i64.add
    (i64.load offset=144
     (i32.const 0)
    )
    (i64.const -1)
   )
  )
  (call $write
   (i32.const 160)
   (i32.const 176)
   (i32.const 144)
  )
  (get_local $1)
 )
 (func $TestPush (; 44 ;) (param $0 i32)
  (local $1 i64)
  (i64.store offset=144
   (i32.const 0)
   (i64.add
    (i64.load
     (call $read
      (i32.const 160)
      (i32.const 176)
     )
    )
    (i64.const 1)
   )
  )
  (call $write
   (i32.const 160)
   (i32.const 176)
   (i32.const 144)
  )
  (i32.store offset=136
   (i32.const 0)
   (get_local $0)
  )
  (i64.store offset=128
   (i32.const 0)
   (tee_local $1
    (i64.add
     (i64.load offset=144
      (i32.const 0)
     )
     (i64.const -1)
    )
   )
  )
  (call $write
   (i32.const 160)
   (i32.wrap/i64
    (get_local $1)
   )
   (i32.const 136)
  )
 )
 (func $TestComplexMappingWithNesting (; 45 ;)
  (local $0 i32)
  (local $1 i32)
  (i32.store offset=196
   (i32.const 0)
   (i32.const 192)
  )
  (set_local $0
   (call $GenerateKey
    (i32.const 224)
    (i32.const 240)
    (i32.const 192)
   )
  )
  (i32.store offset=200
   (i32.const 0)
   (i32.const 256)
  )
  (set_local $1
   (call $GenerateKey
    (get_local $0)
    (i32.const 272)
    (i32.const 256)
   )
  )
  (i32.store offset=204
   (i32.const 0)
   (i32.const 304)
  )
  (call $write
   (i32.const 224)
   (get_local $1)
   (i32.const 204)
  )
  (i32.store offset=200
   (i32.const 0)
   (i32.const 320)
  )
  (set_local $0
   (call $GenerateKey
    (get_local $0)
    (i32.const 272)
    (i32.const 320)
   )
  )
  (i32.store offset=204
   (i32.const 0)
   (i32.const 336)
  )
  (call $write
   (i32.const 224)
   (get_local $0)
   (i32.const 204)
  )
 )
 (func $TestWriteComplexMappingWithStruct (; 46 ;)
  (local $0 i32)
  (local $1 i32)
  (i32.store offset=344
   (i32.const 0)
   (i32.const 192)
  )
  (set_local $1
   (call $GenerateKey
    (tee_local $0
     (call $GenerateKey
      (i32.const 432)
      (i32.const 448)
      (i32.const 192)
     )
    )
    (i32.const 464)
    (i32.const 176)
   )
  )
  (i32.store offset=352
   (i32.const 0)
   (i32.const 496)
  )
  (call $write
   (i32.const 432)
   (get_local $1)
   (i32.const 352)
  )
  (set_local $1
   (call $GenerateKey
    (get_local $0)
    (i32.const 512)
    (i32.const 176)
   )
  )
  
  (set_local $0
   (call $GenerateKey
    (get_local $0)
    (i32.const 544)
    (i32.const 176)
   )
  )
  (i32.store offset=360
   (i32.const 0)
   (i32.const 496)
  )
  (call $write
   (i32.const 432)
   (get_local $0)
   (i32.const 360)
  )
 )
 (func $TestReadComplexMappingWithStruct (; 47 ;)
  (local $0 i32)
  (i32.store offset=344
   (i32.const 0)
   (i32.const 192)
  )
  (i64.store32 offset=352
   (i32.const 0)
   (i64.load
    (call $read
     (i32.const 432)
     (call $GenerateKey
      (tee_local $0
       (call $GenerateKey
        (i32.const 432)
        (i32.const 448)
        (i32.const 192)
       )
      )
      (i32.const 464)
      (i32.const 176)
     )
    )
   )
  )
  (i64.store32 offset=356
   (i32.const 0)
   (i64.load
    (call $read
     (i32.const 432)
     (call $GenerateKey
      (get_local $0)
      (i32.const 512)
      (i32.const 176)
     )
    )
   )
  )
  (i64.store32 offset=360
   (i32.const 0)
   (i64.load
    (call $read402
     (i32.const 432)
     (call $GenerateKey
      (get_local $0)
      (i32.const 544)
      (i32.const 176)
     )
    )
   )
  )
 )
 (func $testcomplex (; 48 ;)
  (i64.store offset=568 align=4
   (i32.const 0)
   (i64.const 4294967297)
  )
  (i64.store offset=576 align=4
   (i32.const 0)
   (i64.const 4294967297)
  )
 )
 (func $__wasm_nullptr (; 49 ;) (type $FUNCSIG$v)
  (unreachable)
 )
 (func $__importThunk_GetOrigin (; 50 ;) (type $FUNCSIG$i) (result i32)
  (call $GetOrigin)
 )
 (func $__importThunk_GetSender (; 51 ;) (type $FUNCSIG$i) (result i32)
  (call $GetSender)
 )
 (func $__importThunk_GetValue (; 52 ;) (type $FUNCSIG$j) (result i64)
  (call $GetValue)
 )
 (func $__importThunk_GetContractValue (; 53 ;) (type $FUNCSIG$j) (result i64)
  (call $GetContractValue)
 )
 (func $__importThunk_GetContractAddress (; 54 ;) (type $FUNCSIG$i) (result i32)
  (call $GetContractAddress)
 )
 (func $__importThunk_GetBalanceFromAddress (; 55 ;) (type $FUNCSIG$ji) (param $0 i32) (result i64)
  (call $GetBalanceFromAddress
   (get_local $0)
  )
 )
 (func $__importThunk_SendFromContract (; 56 ;) (type $FUNCSIG$vij) (param $0 i32) (param $1 i64)
  (call $SendFromContract
   (get_local $0)
   (get_local $1)
  )
 )
 (func $__importThunk_GetBlockHash (; 57 ;) (type $FUNCSIG$ij) (param $0 i64) (result i32)
  (call $GetBlockHash
   (get_local $0)
  )
 )
 (func $__importThunk_GetDifficulty (; 58 ;) (type $FUNCSIG$j) (result i64)
  (call $GetDifficulty)
 )
 (func $__importThunk_GetBlockNumber (; 59 ;) (type $FUNCSIG$j) (result i64)
  (call $GetBlockNumber)
 )
 (func $__importThunk_GetTimestamp (; 60 ;) (type $FUNCSIG$j) (result i64)
  (call $GetTimestamp)
 )
 (func $__importThunk_GetCoinBase (; 61 ;) (type $FUNCSIG$i) (result i32)
  (call $GetCoinBase)
 )
 (func $__importThunk_GetGas (; 62 ;) (type $FUNCSIG$j) (result i64)
  (call $GetGas)
 )
 (func $__importThunk_GetGasLimit (; 63 ;) (type $FUNCSIG$j) (result i64)
  (call $GetGasLimit)
 )
 (func $__importThunk_StorageWrite (; 64 ;) (type $FUNCSIG$vii) (param $0 i32) (param $1 i32)
  (call $StorageWrite
   (get_local $0)
   (get_local $1)
  )
 )
 (func $__importThunk_StorageRead (; 65 ;) (type $FUNCSIG$ii) (param $0 i32) (result i32)
  (call $StorageRead
   (get_local $0)
  )
 )
 (func $__importThunk_SHA3 (; 66 ;) (type $FUNCSIG$ii) (param $0 i32) (result i32)
  (call $SHA3
   (get_local $0)
  )
 )
 (func $__importThunk_FromI64 (; 67 ;) (type $FUNCSIG$ij) (param $0 i64) (result i32)
  (call $FromI64
   (get_local $0)
  )
 )
 (func $__importThunk_FromU64 (; 68 ;) (type $FUNCSIG$ij) (param $0 i64) (result i32)
  (call $FromU64
   (get_local $0)
  )
 )
 (func $__importThunk_ToI64 (; 69 ;) (type $FUNCSIG$ji) (param $0 i32) (result i64)
  (call $ToI64
   (get_local $0)
  )
 )
 (func $__importThunk_ToU64 (; 70 ;) (type $FUNCSIG$ji) (param $0 i32) (result i64)
  (call $ToU64
   (get_local $0)
  )
 )
 (func $__importThunk_Concat (; 71 ;) (type $FUNCSIG$iii) (param $0 i32) (param $1 i32) (result i32)
  (call $Concat
   (get_local $0)
   (get_local $1)
  )
 )
)
