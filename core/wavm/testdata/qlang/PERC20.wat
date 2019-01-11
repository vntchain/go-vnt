(module
  (type $t0 (func))
  (type $t1 (func (param i32)))
  (type $t2 (func (param i32) (result i32)))
  (type $t3 (func (param i32)))
  (type $t4 (func (result i32)))
  (type $t5 (func (param i32) (result i64)))
  (type $t6 (func (param i32 i32)))
  (type $t7 (func (param i32 i32) (result i32)))
  (type $t8 (func (param i32 i32) (result i32)))
  (type $t9 (func (param i32 i32) (result i32)))
  (type $t10 (func (param i32 i32)))
  (type $t11 (func (param i32 i32 i64) (result i64)))
  (type $t12 (func (param i32 i32 i32)))
  (type $t13 (func (param i32) (result i32)))
  (type $t14 (func (param i32 i32) (result i32)))
  (type $t15 (func (param i32 i32) (result i32)))
  (type $t16 (func (param i32) (result i32)))
  (type $t17 (func (result i32)))
  (type $t18 (func (result i64)))
  (type $t19 (func (param i32 i32) (result i64)))
  (type $t20 (func (param i32) (result i64)))
  (type $t21 (func (param i32 i32 i64) (result i32)))
  (type $t22 (func (param i32 i64) (result i32)))
  (type $t23 (func (param i32 i32 i32) (result i32)))
  (type $t24 (func (param i32 i32) (result i64)))
  (type $t25 (func (param i32 i32 i32 i64) (result i32)))
  (type $t26 (func (param i32 i32 i64) (result i32)))
  (type $t27 (func (param i32 i32 i64) (result i64)))
  (type $t28 (func (param i32 i64) (result i64)))
  (type $t29 (func (param i32 i32 i32) (result i64)))
  (type $t30 (func (param i32 i32) (result i64)))
  (import "env" "Sender" (func $env.Sender (type $t3)))
  (import "env" "Store" (func $env.Store (type $t10)))
  (import "env" "Load" (func $env.Load (type $t14)))
  (func $f3 (type $t0)
    (unreachable))
  (func $f4 (type $t1) (param $p0 i32)
    (local $l0 i32)
    (set_local $l0
      (get_local $p0))
    (i32.store offset=4
      (get_local $l0)
      (i32.const 0))
    (i32.store offset=8
      (get_local $l0)
      (i32.const 0)))
  (func $f5 (type $t2) (param $p0 i32) (result i32)
    (local $l0 i32) (local $l1 i32)
    (set_local $l0
      (get_local $p0))
    (set_global $heap_pointer
      (i32.add
        (i32.load offset=16
          (get_local $l0))
        (tee_local $l1
          (get_global $heap_pointer))))
    (i32.store
      (get_local $l1)
      (get_local $l0))
    (call_indirect (type $t1)
      (get_local $l1)
      (i32.add
        (i32.mul
          (i32.const 4)
          (i32.load offset=12
            (i32.load
              (get_local $l1))))
        (i32.const 0)))
    (get_local $l1))
  (func $f6 (type $t4) (result i32)
    (local $l0 i32)
    (set_global $heap_pointer
      (i32.add
        (i32.const 20)
        (tee_local $l0
          (get_global $heap_pointer))))
    (call $env.Sender
      (get_local $l0))
    (get_local $l0))
  (func $f7 (type $t5) (param $p0 i32) (result i64)
    (i64.const 1234567890123456))
  (func $f8 (type $t7) (param $p0 i32) (param $p1 i32) (result i32)
    (local $l0 i32) (local $l1 i32)
    (set_local $l0
      (get_local $p0))
    (set_global $heap_pointer
      (i32.add
        (i32.load offset=16
          (get_local $l0))
        (tee_local $l1
          (get_global $heap_pointer))))
    (i32.store
      (get_local $l1)
      (get_local $l0))
    (call_indirect (type $t6)
      (get_local $l1)
      (get_local $p1)
      (i32.add
        (i32.mul
          (i32.const 4)
          (i32.load offset=12
            (i32.load
              (get_local $l1))))
        (i32.const 1)))
    (get_local $l1))
  (func $f9 (type $t2) (param $p0 i32) (result i32)
    (local $l0 i32)
    (set_local $l0
      (get_local $p0))
    (if $I0
      (i32.eqz
        (i32.load offset=4
          (get_local $l0)))
      (then
        (i32.store offset=8
          (get_local $l0)
          (get_global $g0))
        (set_global $g0
          (get_local $l0))
        (i32.store offset=4
          (get_local $l0)
          (i32.add
            (get_global $g1)
            (i32.const 1)))))
    (i32.load offset=4
      (get_local $l0)))
  (func $f10 (type $t8) (param $p0 i32) (param $p1 i32) (result i32)
    (local $l0 i32) (local $l1 i32)
    (set_local $l0
      (get_local $p0))
    (set_local $l1
      (i32.load offset=16
        (get_local $l0)))
    (i32.store
      (i32.add
        (i32.load offset=20
          (get_local $l0))
        (get_local $l1))
      (get_local $p1))
    (i32.store offset=16
      (get_local $l0)
      (i32.add
        (get_local $l1)
        (i32.const 4)))
    (get_local $l0))
  (func $f11 (type $t9) (param $p0 i32) (param $p1 i32) (result i32)
    (local $l0 i32) (local $l1 i32) (local $l2 i32)
    (set_local $l0
      (get_local $p0))
    (set_local $l1
      (get_local $p1))
    (set_local $l2
      (i32.add
        (i32.load offset=20
          (get_local $l0))
        (i32.load offset=16
          (get_local $l0))))
    (i64.store
      (get_local $l2)
      (i64.load
        (get_local $l1)))
    (i64.store offset=8
      (get_local $l2)
      (i64.load offset=8
        (get_local $l1)))
    (i32.store offset=16
      (get_local $l2)
      (i32.load offset=16
        (get_local $l1)))
    (i32.store offset=16
      (get_local $l0)
      (i32.add
        (i32.load offset=16
          (get_local $l0))
        (i32.const 20)))
    (get_local $l0))
  (func $f12 (type $t11) (param $p0 i32) (param $p1 i32) (param $p2 i64) (result i64)
    (local $l0 i32) (local $l1 i64)
    (set_local $l1
      (get_local $p2))
    (i64.store
      (i32.const 392)
      (get_local $l1))
    (set_local $l0
      (call $f11
        (call $f10
          (call $f8
            (i32.const 180)
            (i32.const 64))
          (call $f9
            (get_local $p0)))
        (get_local $p1)))
    (i32.store
      (i32.const 360)
      (i32.load offset=16
        (get_local $l0)))
    (i32.store offset=4
      (i32.const 360)
      (i32.load offset=20
        (get_local $l0)))
    (i32.store
      (i32.const 368)
      (i32.const 8))
    (i32.store offset=4
      (i32.const 368)
      (i32.const 392))
    (get_local $l1))
  (func $f13 (type $t1) (param $p0 i32)
    (block $B0
      (i32.store offset=12
        (get_local $p0)
        (call $f5
          (i32.const 216)))
      (i32.store offset=16
        (get_local $p0)
        (call $f5
          (i32.const 288)))
      (drop
        (call $f12
          (i32.load offset=12
            (get_local $p0))
          (call $f6)
          (call $f7
            (get_local $p0))))
      (nop)))
  (func $f14 (type $t6) (param $p0 i32) (param $p1 i32)
    (local $l0 i32) (local $l1 i32)
    (set_local $l0
      (get_local $p0))
    (i32.store offset=16
      (get_local $l0)
      (i32.const 0))
    (i32.store offset=12
      (get_local $l0)
      (get_local $p1))
    (set_global $heap_pointer
      (i32.add
        (i32.const 32)
        (tee_local $l1
          (get_global $heap_pointer))))
    (i32.store offset=20
      (get_local $l0)
      (get_local $l1)))
  (func $f15 (type $t12) (param $p0 i32) (param $p1 i32) (param $p2 i32)
    (local $l0 i32)
    (set_local $l0
      (get_local $p0))
    (i32.store offset=12
      (get_local $l0)
      (get_local $p1))
    (i32.store offset=16
      (get_local $l0)
      (get_local $p2))
    (nop))
  (func $f16 (type $t8) (param $p0 i32) (param $p1 i32) (result i32)
    (call $f10
      (get_local $p1)
      (call $f9
        (get_local $p0))))
  (func $f17 (type $t8) (param $p0 i32) (param $p1 i32) (result i32)
    (local $l0 i32)
    (set_local $l0
      (get_local $p0))
    (call $f11
      (call $f11
        (get_local $p1)
        (i32.load offset=12
          (get_local $l0)))
      (i32.load offset=16
        (get_local $l0))))
  (func $malloc (export "malloc") (type $t13) (param $p0 i32) (result i32)
    (local $l0 i32)
    (set_global $heap_pointer
      (i32.add
        (get_local $p0)
        (tee_local $l0
          (get_global $heap_pointer))))
    (get_local $l0))
  (func $f19 (type $t2) (param $p0 i32) (result i32)
    (local $l0 i32) (local $l1 i32) (local $l2 i32) (local $l3 i32) (local $l4 i32) (local $l5 i32)
    (set_local $l0
      (get_local $p0))
    (if $I0
      (i32.load offset=32
        (tee_local $l2
          (i32.load
            (get_local $l0))))
      (then
        (set_local $l1
          (i32.add
            (get_local $l0)
            (i32.load offset=24
              (i32.load
                (get_local $l0)))))
        (loop $L1
          (set_local $l4
            (i32.add
              (get_local $l0)
              (i32.load offset=24
                (get_local $l2))))
          (set_local $l5
            (i32.add
              (get_local $l0)
              (i32.load offset=28
                (get_local $l2))))
          (set_local $l3
            (i32.add
              (get_local $l0)
              (i32.load offset=16
                (tee_local $l2
                  (i32.load offset=32
                    (get_local $l2))))))
          (loop $L2
            (block $B3
              (br_if $B3
                (i32.gt_u
                  (get_local $l3)
                  (tee_local $l5
                    (i32.sub
                      (get_local $l5)
                      (i32.const 4)))))
              (i32.store
                (get_local $l5)
                (call $f9
                  (i32.load
                    (get_local $l5))))
              (br $L2)))
          (if $I4
            (i32.eq
              (get_local $l4)
              (get_local $l1))
            (then
              (set_local $l1
                (tee_local $l4
                  (get_local $l3))))
            (else
              (loop $L5
                (block $B6
                  (br_if $B6
                    (i32.gt_u
                      (get_local $l3)
                      (tee_local $l4
                        (i32.sub
                          (get_local $l4)
                          (i32.const 8)))))
                  (i64.store
                    (tee_local $l1
                      (i32.sub
                        (get_local $l1)
                        (i32.const 8)))
                    (i64.load
                      (get_local $l4)))
                  (br $L5)))
              (if $I7
                (i32.le_u
                  (get_local $l3)
                  (tee_local $l4
                    (i32.add
                      (get_local $l4)
                      (i32.const 4))))
                (then
                  (i32.store
                    (tee_local $l1
                      (i32.sub
                        (get_local $l1)
                        (i32.const 4)))
                    (i32.load offset=4
                      (tee_local $l4
                        (i32.sub
                          (get_local $l4)
                          (i32.const 4)))))))
              (if $I8
                (i32.le_u
                  (get_local $l3)
                  (tee_local $l4
                    (i32.add
                      (get_local $l4)
                      (i32.const 2))))
                (then
                  (i32.store16
                    (tee_local $l1
                      (i32.sub
                        (get_local $l1)
                        (i32.const 2)))
                    (i32.load16_u offset=2
                      (tee_local $l4
                        (i32.sub
                          (get_local $l4)
                          (i32.const 2)))))))
              (if $I9
                (i32.le_u
                  (get_local $l3)
                  (tee_local $l4
                    (i32.add
                      (get_local $l4)
                      (i32.const 1))))
                (then
                  (i32.store8
                    (tee_local $l1
                      (i32.sub
                        (get_local $l1)
                        (i32.const 1)))
                    (i32.load8_u
                      (get_local $l4))))
                (else
                  (set_local $l4
                    (i32.add
                      (get_local $l4)
                      (i32.const 1)))))))
          (br_if $L1
            (get_local $l2))))
      (else
        (set_local $l1
          (i32.add
            (get_local $l0)
            (i32.const 0)))))
    (if $I10
      (i32.eqz
        (i32.load offset=4
          (get_local $l0)))
      (then
        (i32.store offset=8
          (get_local $l0)
          (get_global $g0))
        (set_global $g0
          (get_local $l0))
        (i32.store offset=4
          (get_local $l0)
          (i32.add
            (get_global $g1)
            (i32.const 1)))))
    (i32.store offset=4
      (i32.const 360)
      (i32.add
        (get_local $l0)
        (i32.const 4)))
    (i32.store
      (i32.const 368)
      (i32.sub
        (i32.add
          (get_local $l0)
          (i32.load offset=24
            (i32.load
              (get_local $l0))))
        (get_local $l1)))
    (i32.store offset=4
      (i32.const 368)
      (get_local $l1))
    (get_local $l0))
  (func $ERC20::init (export "ERC20::init") (type $t0)
    (local $l0 i32)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l0
      (call $f5
        (i32.const 324)))
    (i32.store offset=8
      (get_local $l0)
      (get_global $g0))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (nop))
  (func $f21 (type $t15) (param $p0 i32) (param $p1 i32) (result i32)
    (local $l0 i32) (local $l1 i32) (local $l2 i32) (local $l3 i32) (local $l4 i32) (local $l5 i32)
    (set_local $l0
      (get_local $p0))
    (set_local $l1
      (get_local $p1))
    (set_global $heap_pointer
      (i32.add
        (i32.load offset=16
          (get_local $l0))
        (tee_local $l3
          (get_global $heap_pointer))))
    (i32.store
      (i32.const 372)
      (get_local $l1))
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 372))
    (set_local $l2
      (i32.add
        (call $env.Load
          (i32.const 360)
          (get_local $l3))
        (get_local $l3)))
    (loop $L0
      (set_local $l4
        (i32.sub
          (get_local $l2)
          (i32.load offset=20
            (get_local $l0))))
      (set_local $l5
        (i32.add
          (get_local $l3)
          (i32.load offset=24
            (get_local $l0))))
      (if $I1
        (i32.eq
          (get_local $l2)
          (get_local $l5))
        (then
          (set_local $l5
            (tee_local $l2
              (get_local $l4))))
        (else
          (loop $L2
            (block $B3
              (br_if $B3
                (i32.gt_u
                  (get_local $l4)
                  (tee_local $l2
                    (i32.sub
                      (get_local $l2)
                      (i32.const 8)))))
              (i64.store
                (tee_local $l5
                  (i32.sub
                    (get_local $l5)
                    (i32.const 8)))
                (i64.load
                  (get_local $l2)))
              (br $L2)))
          (if $I4
            (i32.le_u
              (get_local $l4)
              (tee_local $l2
                (i32.add
                  (get_local $l2)
                  (i32.const 4))))
            (then
              (i32.store
                (tee_local $l5
                  (i32.sub
                    (get_local $l5)
                    (i32.const 4)))
                (i32.load offset=4
                  (tee_local $l2
                    (i32.sub
                      (get_local $l2)
                      (i32.const 4)))))))
          (if $I5
            (i32.le_u
              (get_local $l4)
              (tee_local $l2
                (i32.add
                  (get_local $l2)
                  (i32.const 2))))
            (then
              (i32.store16
                (tee_local $l5
                  (i32.sub
                    (get_local $l5)
                    (i32.const 2)))
                (i32.load16_u offset=2
                  (tee_local $l2
                    (i32.sub
                      (get_local $l2)
                      (i32.const 2)))))))
          (if $I6
            (i32.le_u
              (get_local $l4)
              (tee_local $l2
                (i32.add
                  (get_local $l2)
                  (i32.const 1))))
            (then
              (i32.store8
                (tee_local $l5
                  (i32.sub
                    (get_local $l5)
                    (i32.const 1)))
                (i32.load8_u
                  (get_local $l2))))
            (else
              (set_local $l2
                (i32.add
                  (get_local $l2)
                  (i32.const 1)))))))
      (br_if $L0
        (tee_local $l0
          (i32.load offset=32
            (get_local $l0)))))
    (i32.store offset=4
      (get_local $l3)
      (get_local $l1))
    (get_local $l3))
  (func $f22 (type $t16) (param $p0 i32) (result i32)
    (i32.const 8))
  (func $decimals (export "decimals") (type $t17) (result i32)
    (local $l0 i32) (local $l1 i32)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l1
      (call $f21
        (i32.const 324)
        (i32.const 1)))
    (set_local $l0
      (call $f22
        (get_local $l1)))
    (i32.store offset=8
      (get_local $l1)
      (get_global $g0))
    (set_global $g0
      (get_local $l1))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (get_local $l0))
  (func $totalSupply (export "totalSupply") (type $t18) (result i64)
    (local $l0 i32) (local $l1 i64)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l0
      (call $f21
        (i32.const 324)
        (i32.const 1)))
    (set_local $l1
      (call $f7
        (get_local $l0)))
    (i32.store offset=8
      (get_local $l0)
      (get_global $g0))
    (set_global $g0
      (get_local $l0))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (get_local $l1))
  (func $f25 (type $t19) (param $p0 i32) (param $p1 i32) (result i64)
    (local $l0 i32) (local $l1 i32)
    (set_local $l0
      (call $f11
        (call $f10
          (call $f8
            (i32.const 180)
            (i32.const 64))
          (call $f9
            (get_local $p0)))
        (get_local $p1)))
    (i32.store
      (i32.const 360)
      (i32.load offset=16
        (get_local $l0)))
    (i32.store offset=4
      (i32.const 360)
      (i32.load offset=20
        (get_local $l0)))
    (set_global $heap_pointer
      (i32.add
        (call $env.Load
          (i32.const 360)
          (tee_local $l1
            (get_global $heap_pointer)))
        (get_global $heap_pointer)))
    (i64.load
      (get_local $l1)))
  (func $f26 (type $t19) (param $p0 i32) (param $p1 i32) (result i64)
    (call $f25
      (i32.load offset=12
        (get_local $p0))
      (get_local $p1)))
  (func $balanceOf (export "balanceOf") (type $t20) (param $p0 i32) (result i64)
    (local $l0 i32) (local $l1 i64)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l0
      (call $f21
        (i32.const 324)
        (i32.const 1)))
    (set_local $l1
      (call $f26
        (get_local $l0)
        (get_local $p0)))
    (i32.store offset=8
      (get_local $l0)
      (get_global $g0))
    (set_global $g0
      (get_local $l0))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (get_local $l1))
  (func $f28 (type $t21) (param $p0 i32) (param $p1 i32) (param $p2 i64) (result i32)
    (local $l0 i32) (local $l1 i32)
    (if $I0 (result i32)
      (i64.lt_u
        (call $f25
          (i32.load offset=12
            (get_local $p0))
          (call $f6))
        (get_local $p2))
      (then
        (i32.const 0))
      (else
        (set_local $l0
          (call $f6))
        (drop
          (call $f12
            (i32.load offset=12
              (get_local $p0))
            (get_local $l0)
            (i64.sub
              (call $f25
                (i32.load offset=12
                  (get_local $p0))
                (get_local $l0))
              (get_local $p2))))
        (set_local $l1
          (get_local $p1))
        (drop
          (call $f12
            (i32.load offset=12
              (get_local $p0))
            (get_local $l1)
            (i64.add
              (call $f25
                (i32.load offset=12
                  (get_local $p0))
                (get_local $l1))
              (get_local $p2))))
        (i32.const 1))))
  (func $transfer (export "transfer") (type $t22) (param $p0 i32) (param $p1 i64) (result i32)
    (local $l0 i32) (local $l1 i32)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l1
      (call $f21
        (i32.const 324)
        (i32.const 1)))
    (set_local $l0
      (call $f28
        (get_local $l1)
        (get_local $p0)
        (get_local $p1)))
    (i32.store offset=8
      (get_local $l1)
      (get_global $g0))
    (set_global $g0
      (get_local $l1))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (get_local $l0))
  (func $f30 (type $t23) (param $p0 i32) (param $p1 i32) (param $p2 i32) (result i32)
    (local $l0 i32) (local $l1 i32)
    (set_local $l0
      (get_local $p0))
    (set_global $heap_pointer
      (i32.add
        (i32.load offset=16
          (get_local $l0))
        (tee_local $l1
          (get_global $heap_pointer))))
    (i32.store
      (get_local $l1)
      (get_local $l0))
    (call_indirect (type $t12)
      (get_local $l1)
      (get_local $p1)
      (get_local $p2)
      (i32.add
        (i32.mul
          (i32.const 4)
          (i32.load offset=12
            (i32.load
              (get_local $l1))))
        (i32.const 2)))
    (get_local $l1))
  (func $f31 (type $t24) (param $p0 i32) (param $p1 i32) (result i64)
    (local $l0 i32) (local $l1 i32)
    (set_local $l0
      (get_local $p1))
    (set_local $l0
      (call_indirect (type $t8)
        (get_local $l0)
        (call $f10
          (call $f8
            (i32.const 180)
            (i32.const 64))
          (call $f9
            (get_local $p0)))
        (i32.add
          (i32.mul
            (i32.const 4)
            (i32.load offset=12
              (i32.load
                (get_local $l0))))
          (i32.const 3))))
    (i32.store
      (i32.const 360)
      (i32.load offset=16
        (get_local $l0)))
    (i32.store offset=4
      (i32.const 360)
      (i32.load offset=20
        (get_local $l0)))
    (set_global $heap_pointer
      (i32.add
        (call $env.Load
          (i32.const 360)
          (tee_local $l1
            (get_global $heap_pointer)))
        (get_global $heap_pointer)))
    (i64.load
      (get_local $l1)))
  (func $f32 (type $t25) (param $p0 i32) (param $p1 i32) (param $p2 i32) (param $p3 i64) (result i32)
    (local $l0 i32) (local $l1 i32) (local $l2 i64)
    (block $B0 (result i32)
      (set_local $l2
        (call $f31
          (i32.load offset=16
            (get_local $p0))
          (call $f30
            (i32.const 252)
            (get_local $p1)
            (get_local $p2))))
      (if $I1 (result i32)
        (i64.lt_u
          (get_local $l2)
          (get_local $p3))
        (then
          (i32.const 0))
        (else
          (set_local $l0
            (get_local $p1))
          (drop
            (call $f12
              (i32.load offset=12
                (get_local $p0))
              (get_local $l0)
              (i64.sub
                (call $f25
                  (i32.load offset=12
                    (get_local $p0))
                  (get_local $l0))
                (get_local $p3))))
          (set_local $l1
            (get_local $p2))
          (drop
            (call $f12
              (i32.load offset=12
                (get_local $p0))
              (get_local $l1)
              (i64.add
                (call $f25
                  (i32.load offset=12
                    (get_local $p0))
                  (get_local $l1))
                (get_local $p3))))
          (i32.const 1)))))
  (func $transferFrom (export "transferFrom") (type $t26) (param $p0 i32) (param $p1 i32) (param $p2 i64) (result i32)
    (local $l0 i32) (local $l1 i32)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l1
      (call $f21
        (i32.const 324)
        (i32.const 1)))
    (set_local $l0
      (call $f32
        (get_local $l1)
        (get_local $p0)
        (get_local $p1)
        (get_local $p2)))
    (i32.store offset=8
      (get_local $l1)
      (get_global $g0))
    (set_global $g0
      (get_local $l1))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (get_local $l0))
  (func $f34 (type $t27) (param $p0 i32) (param $p1 i32) (param $p2 i64) (result i64)
    (local $l0 i32) (local $l1 i64)
    (set_local $l0
      (get_local $p1))
    (set_local $l1
      (get_local $p2))
    (i64.store
      (i32.const 392)
      (get_local $l1))
    (set_local $l0
      (call_indirect (type $t8)
        (get_local $l0)
        (call $f10
          (call $f8
            (i32.const 180)
            (i32.const 64))
          (call $f9
            (get_local $p0)))
        (i32.add
          (i32.mul
            (i32.const 4)
            (i32.load offset=12
              (i32.load
                (get_local $l0))))
          (i32.const 3))))
    (i32.store
      (i32.const 360)
      (i32.load offset=16
        (get_local $l0)))
    (i32.store offset=4
      (i32.const 360)
      (i32.load offset=20
        (get_local $l0)))
    (i32.store
      (i32.const 368)
      (i32.const 8))
    (i32.store offset=4
      (i32.const 368)
      (i32.const 392))
    (get_local $l1))
  (func $f35 (type $t11) (param $p0 i32) (param $p1 i32) (param $p2 i64) (result i64)
    (call $f34
      (i32.load offset=16
        (get_local $p0))
      (call $f30
        (i32.const 252)
        (call $f6)
        (get_local $p1))
      (get_local $p2)))
  (func $approve (export "approve") (type $t28) (param $p0 i32) (param $p1 i64) (result i64)
    (local $l0 i32) (local $l1 i64)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l0
      (call $f21
        (i32.const 324)
        (i32.const 1)))
    (set_local $l1
      (call $f35
        (get_local $l0)
        (get_local $p0)
        (get_local $p1)))
    (i32.store offset=8
      (get_local $l0)
      (get_global $g0))
    (set_global $g0
      (get_local $l0))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (get_local $l1))
  (func $f37 (type $t29) (param $p0 i32) (param $p1 i32) (param $p2 i32) (result i64)
    (call $f31
      (i32.load offset=16
        (get_local $p0))
      (call $f30
        (i32.const 252)
        (get_local $p1)
        (get_local $p2))))
  (func $allowance (export "allowance") (type $t30) (param $p0 i32) (param $p1 i32) (result i64)
    (local $l0 i32) (local $l1 i64)
    (i32.store
      (i32.const 360)
      (i32.const 4))
    (i32.store offset=4
      (i32.const 360)
      (i32.const 0))
    (drop
      (call $env.Load
        (i32.const 360)
        (i32.const 384)))
    (set_global $g1
      (i32.load
        (i32.const 384)))
    (set_local $l0
      (call $f21
        (i32.const 324)
        (i32.const 1)))
    (set_local $l1
      (call $f37
        (get_local $l0)
        (get_local $p0)
        (get_local $p1)))
    (i32.store offset=8
      (get_local $l0)
      (get_global $g0))
    (set_global $g0
      (get_local $l0))
    (if $I0
      (get_global $g0)
      (then
        (drop
          (call $f19
            (get_global $g0)))))
    (if $I1
      (i32.ne
        (get_global $g1)
        (i32.load
          (i32.const 384)))
      (then
        (i32.store
          (i32.const 384)
          (get_global $g1))
        (i32.store offset=4
          (i32.const 360)
          (i32.const 0))
        (i32.store
          (i32.const 368)
          (i32.const 4))
        (i32.store offset=4
          (i32.const 368)
          (i32.const 384))))
    (get_local $l1))
  (table $T0 40 40 anyfunc)
  (memory $default (export "default") 1 1)
  (global $g0 i32 (i32.const 0))
  (global $g1 i32 (i32.const 0))
  (global $heap_pointer (export "heap_pointer") i32 (i32.const 400))
  (elem (i32.const 0) $f3 $f3 $f3 $f3 $f4 $f3 $f3 $f16 $f4 $f3 $f3 $f16 $f4 $f3 $f3 $f16 $f4 $f3 $f3 $f16 $f4 $f14 $f3 $f16 $f4 $f3 $f3 $f16 $f4 $f3 $f15 $f17 $f4 $f3 $f3 $f16 $f13 $f3 $f3 $f16)
  (data (i32.const 0) "\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\01\00\00\00\0c\00\00\00\04\00\00\00\04\00\00\00\00\00\00\00\00\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\02\00\00\00$\00\00\00\00\00\00\00\0c\00\00\00\0c\00\00\00$\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\03\00\00\00\14\00\00\00\00\00\00\00\0c\00\00\00\0c\00\00\00$\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\04\00\00\00\0c\00\00\00\00\00\00\00\0c\00\00\00\0c\00\00\00$\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\05\00\00\00\18\00\00\00\00\00\00\00\0c\00\00\00\0c\00\00\00$\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\06\00\00\00\0c\00\00\00\00\00\00\00\0c\00\00\00\0c\00\00\00\90\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\07\00\00\00\14\00\00\00\08\00\00\00\14\00\00\00\0c\00\00\00$\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\08\00\00\00\0c\00\00\00\00\00\00\00\0c\00\00\00\0c\00\00\00\90\00\00\00H\00\00\00\00\00\00\00\00\00\00\00\09\00\00\00\14\00\00\00\08\00\00\00\14\00\00\00\14\00\00\00$\00\00\00")
  (data (i32.const 360) "\04\00\00\00\00\00\00\00")
  (data (i32.const 368) "\00\00\00\00\00\00\00\00")
  (data (i32.const 376) "\00\00\00\00\00\00\00\00")
  (data (i32.const 392) "\00\00\00\00\00\00\00\00"))
